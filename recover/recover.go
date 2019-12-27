package recover

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/klauspost/reedsolomon"
	"github.com/mr-tron/base58/base58"
	log "github.com/yottachain/YTDataNode/logger"
	"github.com/yottachain/YTDataNode/message"
	node "github.com/yottachain/YTDataNode/storageNodeInterface"
	"github.com/yottachain/YTFS/common"
	lrcpkg "github.com/yottachain/YTLRC"
	_ "net/http/pprof"
	"sync"
	"time"
)

const (
	max_reply_num       = 1000
	max_task_num        = 1000
	max_reply_wait_time = time.Second * 60
)

type RecoverEngine struct {
	sn         node.StorageNode
	queue      chan []byte
	replyQueue chan *TaskMsgResult
	le         *LRCEngine
}

func New(sn node.StorageNode) (*RecoverEngine, error) {
	var re = new(RecoverEngine)
	re.queue = make(chan []byte, max_task_num)
	re.replyQueue = make(chan *TaskMsgResult, max_reply_num)
	re.sn = sn
	re.le = NewLRCEngine(re.getShard)

	return re, nil
}

func (re *RecoverEngine) recoverShard(description *message.TaskDescription) error {
	defer func() {
		err := recover()
		fmt.Println("err:", err)
	}()
	var size = len(description.Hashs)
	var shards [][]byte = make([][]byte, size)
	encoder, err := reedsolomon.New(size-int(description.ParityShardCount), int(description.ParityShardCount))
	if err != nil {
		return err
	}
	var wg = sync.WaitGroup{}
	var number int
	wg.Add(len(description.Locations))
	log.Printf("[recover:%s]recover start %d\n", base58.Encode(description.Id), size)
	for k, v := range description.Locations {
		go func(k int, v *message.P2PLocation) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			shard, err := re.getShard(ctx, v.NodeId, base58.Encode(description.Id), v.Addrs, description.Hashs[k], &number)
			if err == nil {
				shards[k] = shard
			} else {
				log.Printf("[recover:%s]error:%s, %v, %s\n", base58.Encode(description.Id), err.Error(), v.Addrs, v.NodeId)
			}
		}(k, v)
	}
	wg.Wait()
	shards[description.RecoverId] = nil
	err = encoder.Reconstruct(shards)
	if err != nil {
		log.Printf("[recover:%s]datas recover error:%s\n", base58.Encode(description.Id), err.Error())
		return err
	}
	log.Printf("[recover:%s]datas recover success\n", base58.Encode(description.Id))
	var vhf [16]byte
	copy(vhf[:], description.Hashs[description.RecoverId])
	err = re.sn.YTFS().Put(common.IndexTableKey(vhf), shards[int(description.RecoverId)])
	if err != nil && err.Error() != "YTFS: hash key conflict happens" || err.Error() == "YTFS: conflict hash value" {
		log.Printf("[recover:%s]YTFS Put error %s\n", base58.Encode(description.Id), err.Error())
		return err
	}
	return nil
}

func (re *RecoverEngine) getShard(ctx context.Context, id string, taskID string, addrs []string, hash []byte, n *int) ([]byte, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	clt, err := re.sn.Host().ClientStore().GetByAddrString(ctx, id, addrs)
	if err != nil {
		return nil, err
	}
	// todo: 这里需要单独处理，连接自己失败的错误
	//if err != nil {
	//	if err.Error() == "new stream error:dial to self attempted" {
	//		var vhf [16]byte
	//		copy(vhf[:], hash)
	//		return re.sn.YTFS().Get(common.IndexTableKey(vhf))
	//	}
	//	return nil, err
	//}

	var msg message.DownloadShardRequest
	var res message.DownloadShardResponse
	msg.VHF = hash
	buf, err := proto.Marshal(&msg)
	if err != nil {
		return nil, err
	}
	shardBuf, err := clt.SendMsgClose(ctx, message.MsgIDDownloadShardRequest.Value(), buf)

	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(shardBuf[2:], &res)
	if err != nil {
		log.Printf("[recover:%s]get shard [%s] error[%d] %s\n", taskID, base58.Encode(hash), *n, err.Error())
		return nil, err
	}
	*n = *n + 1
	log.Printf("[recover:%s]get shard [%s] success[%d]\n", taskID, base58.Encode(hash), *n)
	return res.Data, nil
}

type TaskMsgResult struct {
	ID   []byte
	RES  int32
	BPID byte
}

func (re *RecoverEngine) PutTask(task []byte) error {
	select {
	case re.queue <- task:
	default:
	}
	return nil
}

// 多重建任务消息处理
func (re *RecoverEngine) HandleMuilteTaskMsg(msgData []byte) error {
	var mtdMsg message.MultiTaskDescription
	if err := proto.Unmarshal(msgData, &mtdMsg); err != nil {
		return err
	}
	log.Printf("[recover]multi recover task start, pack size %d\n", len(mtdMsg.Tasklist))
	for _, task := range mtdMsg.Tasklist {
		if err := re.PutTask(task); err != nil {
			log.Printf("[recover]put recover task error: %s\n", err.Error())
		}
	}
	return nil
}

func (re *RecoverEngine) Run() {
	go func() {
		for {
			msg := <-re.queue
			if bytes.Equal(msg[0:2], message.MsgIDTaskDescript.Bytes()) {
				res := re.execRCTask(msg[2:])
				re.PutReplyQueue(res)
			} else if bytes.Equal(msg[0:2], message.MsgIDLRCTaskDescription.Bytes()) {
				log.Printf("[recover]LRC start\n")
				res := re.execLRCTask(msg[2:])
				re.PutReplyQueue(res)
			} else {
				res := re.execCPTask(msg[2:])
				re.PutReplyQueue(res)
			}
		}
	}()

	for {
		re.MultiReply()
	}
}

func (re *RecoverEngine) PutReplyQueue(res *TaskMsgResult) {
	select {
	case re.replyQueue <- res:
	default:
	}
}

func (re *RecoverEngine) reply(res *TaskMsgResult) error {
	var msgData message.TaskOpResult
	msgData.Id = res.ID
	msgData.RES = res.RES
	data, err := proto.Marshal(&msgData)
	if err != nil {
		return err
	}
	_, err = re.sn.SendBPMsg(int(res.BPID), message.MsgIDTaskOPResult.Value(), data)
	log.Println("[recover] reply to", int(res.BPID))
	return err
}

//
func (re *RecoverEngine) MultiReply() error {
	var resmsg = make(map[byte]message.MultiTaskOpResult)

	func() {
		for i := 0; i < max_reply_num; i++ {
			select {
			case res := <-re.replyQueue:
				_r := resmsg[res.BPID]
				_r.Id = append(_r.Id, res.ID)
				_r.RES = append(_r.RES, res.RES)
				resmsg[res.BPID] = _r

			case <-time.After(max_reply_wait_time):
				return
			}
		}
	}()
	if l := len(resmsg); l > 0 {
		fmt.Println("待上报重建消息：", len(resmsg))
	}
	for k, v := range resmsg {
		if data, err := proto.Marshal(&v); err != nil {
			log.Printf("[recover]marsnal failed %s\n", err.Error())
			continue
		} else {
			re.sn.SendBPMsg(int(k), message.MsgIDMultiTaskOPResult.Value(), data)
			log.Println("[recover] multi reply success")
		}
	}

	return nil
}

func (re *RecoverEngine) execRCTask(msgData []byte) *TaskMsgResult {
	var res TaskMsgResult
	var msg message.TaskDescription
	if err := proto.Unmarshal(msgData, &msg); err != nil {
		log.Printf("[recover]proto解析错误%s", err)
		res.RES = 1
	}
	res.ID = msg.Id
	res.BPID = msg.Id[12]
	if err := re.recoverShard(&msg); err != nil {
		res.RES = 1
	} else {
		res.RES = 0
	}
	return &res
}

func (re *RecoverEngine) execLRCTask(msgData []byte) *TaskMsgResult {

	var res TaskMsgResult
	var msg message.TaskDescription

	if err := proto.Unmarshal(msgData, &msg); err != nil {
		log.Printf("[recover]proto解析错误%s", err)
		res.RES = 1
	}

	res.ID = msg.Id
	res.BPID = msg.Id[12]
	res.RES = 1
	log.Printf("[recover]LRC 分片恢复开始%s", base58.Encode(msg.Id))
	defer log.Printf("[recover]LRC 分片恢复结束%s", base58.Encode(msg.Id))

	lrc := lrcpkg.Shardsinfo{}

	lrc.OriginalCount = uint16(len(msg.Hashs) - int(msg.ParityShardCount))
	lrc.RecoverNum = 13
	lrc.Lostindex = uint16(msg.RecoverId)
	h, err := re.le.GetLRCHandler(&lrc)
	if err != nil {
		log.Printf("[recover]LRC 获取Handler失败%s", err)
		return &res
	}

	recoverData, err := h.Recover(msg)
	if err != nil {
		log.Printf("[recover]LRC 恢复失败%s", err)
		return &res
	}

	m5 := md5.New()
	m5.Write(recoverData)
	hash := m5.Sum(nil)
	// 校验hash失败
	if !bytes.Equal(hash, msg.Hashs[msg.RecoverId]) {
		log.Printf("[recover]LRC 校验HASH失败%s %s\n", base58.Encode(hash), base58.Encode(msg.Hashs[msg.RecoverId]))
		return &res
	}

	var key [common.HashLength]byte
	copy(key[:], hash)
	if err := re.sn.YTFS().Put(common.IndexTableKey(key), recoverData); err != nil && err.Error() != "YTFS: hash key conflict happens" {
		log.Printf("[recover]LRC 保存已恢复分片失败%s\n", err)
		return &res
	}

	log.Printf("[recover]LRC 分片恢复成功\n")
	res.RES = 0
	return &res
}

// 副本集任务
func (re *RecoverEngine) execCPTask(msgData []byte) *TaskMsgResult {
	var msg message.TaskDescriptionCP
	var result TaskMsgResult
	err := proto.UnmarshalMerge(msgData, &msg)
	if err != nil {
		log.Printf("[recover]解析错误%s\n", err.Error())
	}
	result.ID = msg.Id
	result.BPID = msg.Id[12]
	result.RES = 1
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var number int
	// 循环从副本节点获取分片，只要有一个成功就返回
	for _, v := range msg.Locations {
		shard, err := re.getShard(ctx, v.NodeId, base58.Encode(msg.Id), v.Addrs, msg.DataHash, &number)
		// 如果没有发生错误，分片下载成功，就存储分片
		if err == nil {
			var vhf [16]byte
			copy(vhf[:], msg.DataHash)
			err := re.sn.YTFS().Put(common.IndexTableKey(vhf), shard)
			// 存储分片没有错误，或者分片已存在返回0，代表成功
			if err != nil && err.Error() != "YTFS: hash key conflict happens" || err.Error() == "YTFS: conflict hash value" {
				log.Printf("[recover:%s]YTFS Put error %s\n", base58.Encode(vhf[:]), err.Error())
				result.RES = 1
			} else {
				result.RES = 0
			}
			break
		}
	}
	return &result
}
