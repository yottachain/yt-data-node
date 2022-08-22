package recover

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	//"github.com/yottachain/YTDataNode/activeNodeList"
	"sync/atomic"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/yottachain/YTDataNode/config"
	"github.com/yottachain/YTDataNode/recover/actuator"
	"github.com/yottachain/YTDataNode/recover/shardDownloader"
	"github.com/yottachain/YTDataNode/statistics"
	"github.com/yottachain/YTElkProducer"
	"github.com/yottachain/YTElkProducer/conf"
	"github.com/yottachain/YTHost/client"

	//"github.com/docker/docker/pkg/locker"
	_ "net/http/pprof"
	"strings"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/klauspost/reedsolomon"
	"github.com/mr-tron/base58/base58"
	"github.com/yottachain/YTDataNode/TokenPool"
	log "github.com/yottachain/YTDataNode/logger"
	"github.com/yottachain/YTDataNode/message"
	node "github.com/yottachain/YTDataNode/storageNodeInterface"
	"github.com/yottachain/YTFS/common"
	lrcpkg "github.com/yottachain/YTLRC"
	//"io"
)

const (
	max_reply_num       = 1000
	max_task_num        = 1000
	max_reply_wait_time = time.Millisecond* 10
)

// var elkClt = util.NewElkClient("rebuild_reply", &config.Gconfig.ElkReport2)

type elkErrorLog struct {
	ErrorMsg  string
	RetryTime int
}

type Engine struct {
	sn                node.StorageNode
	waitQueue         *TaskWaitQueue
	replyQueue        chan *TaskMsgResult
	le                *LRCEngine
	Upt               *TokenPool.TokenPool
	startTskTmCtl     uint8
	DefaultDownloader shardDownloader.ShardDownloader
	lck 	*sync.Mutex		//引擎的全局锁  cgo调用的时候使用
}

func New(sn node.StorageNode) (*Engine, error) {

	var re = new(Engine)
	re.waitQueue = NewTaskWaitQueue()
	re.replyQueue = make(chan *TaskMsgResult, max_reply_num)
	re.sn = sn
	re.DefaultDownloader = shardDownloader.New(sn.Host().ClientStore(), 20)
	re.le = NewLRCEngine(statistics.DefaultRebuildCount.IncRbdSucc)
	re.Upt = TokenPool.Utp()
	re.lck = &sync.Mutex{}

	return re, nil
}

func (re *Engine) Len() uint32 {
	return uint32(re.waitQueue.Len())
}

func (re *Engine) recoverShard(description *message.TaskDescription) error {
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
			shard, err := re.getShard2(ctx, v.NodeId, base58.Encode(description.Id), v.Addrs, description.Hashs[k], &number)
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
	_, err = re.sn.YTFS().BatchPut(map[common.IndexTableKey][]byte{common.IndexTableKey{Hsh:vhf, Id:0}: shards[int(description.RecoverId)]})
	if err != nil && (err.Error() != "YTFS: hash key conflict happens" || err.Error() == "YTFS: conflict hash value") {
		log.Printf("[recover:%s]YTFS Put error %s\n", base58.Encode(description.Id), err.Error())
		return err
	}
	return nil
}

func (re *Engine) getShard2(ctx context.Context, id string, taskID string, addrs []string, hash []byte, n *int) ([]byte, error) {
	clt, err := re.sn.Host().ClientStore().GetByAddrString(ctx, id, addrs)
	if err != nil {
		log.Println("[recover][debug] getShardcnn  err=", err)
		return nil, err
	}

	var msg message.DownloadShardRequest
	msg.VHF = hash

	buf, err := proto.Marshal(&msg)
	if err != nil {
		return nil, err
	}
	ctx2, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	shardBuf, err := clt.SendMsg(ctx2, message.MsgIDDownloadShardRequest.Value(), buf)

	if err != nil {
		if strings.Contains(err.Error(), "Get data Slice fail") {
			log.Printf("[recover:%s]  error %s addr %v\n",
				base58.Encode(hash), err.Error(), addrs)
		} else {
			log.Printf("[recover:%s]  error %s addr %v\n",
				base58.Encode(hash), err.Error(), addrs)
		}
		return nil, err
	}

	if len(shardBuf) < 3 {
		log.Printf("[recover:%s] error: shard empty! addr %v\n",
			base58.Encode(hash), addrs)
		return nil, fmt.Errorf("error: shard less then 3, len=", len(shardBuf))
	}

	var resMsg message.DownloadShardResponse
	err = proto.Unmarshal(shardBuf[2:], &resMsg)
	if err != nil {
		return nil, err
	}

	return resMsg.Data, err
}

func NewElkClient(tbstr string) *YTElkProducer.Client {
	elkConf := elasticsearch.Config{
		Addresses: []string{"https://c1-bj-elk.yottachain.net/"},
		Username:  "dnreporter",
		Password:  "dnreporter@yottachain",
	}

	ytESConfig := conf.YTESConfig{
		ESConf:      elkConf,
		DebugMode:   false,
		IndexPrefix: tbstr,
		IndexType:   "log",
	}

	client, _ := YTElkProducer.NewClient(ytESConfig)
	return &client
}

func (re *Engine) reportLog(body interface{}) {
	//if ! config.Gconfig.ElkReport{
	//	return
	//}

	time.Sleep(time.Second * 10)
}

func (re *Engine) MakeReportLog(nodeid string, hash []byte, errtype string, err error) *RcvDbgLog {
	//if ! config.Gconfig.ElkReport{
	//	return nil
	//}

	ShardId := base64.StdEncoding.EncodeToString(hash)
	NowTm := time.Now().Format("2006/01/02 15:04:05")
	localNodeId := re.sn.Config().ID
	localNdVersion := re.sn.Config().Version()
	return &RcvDbgLog{
		nodeid,
		ShardId,
		localNodeId,
		localNdVersion,
		NowTm,
		errtype,
		err.Error(),
	}
}

func (re *Engine) parmCheck(id string, taskID string, addrs []string, hash []byte, n *int, sw *Switchcnt) ([]byte, error) {
	if 0 == sw.swget {
		statistics.DefaultRebuildCount.IncShardForRbd()
		sw.swget++
	}

	statistics.DefaultRebuildCount.IncGetShardWK()
	btid, err := base58.Decode(taskID)
	if err != nil {
		statistics.DefaultRebuildCount.IncFailDcdTask()
		return btid, err
	}

	if 0 == len(id) {
		err = fmt.Errorf("zero length id")
		statistics.DefaultRebuildCount.IncFailDcdTask()
		return btid, err
	}

	if 0 == len(addrs) {
		err = fmt.Errorf("zero length addrs")
		statistics.DefaultRebuildCount.IncFailDcdTask()
		return btid, err
	}

	if 0 == len(hash) {
		err = fmt.Errorf("zero length hash")
		statistics.DefaultRebuildCount.IncFailDcdTask()
		return btid, err
	}
	return btid, nil
}

func (re *Engine) getRdToken(clt *client.YTHostClient, sw *Switchcnt) ([]byte, error) {
	var getToken message.NodeCapacityRequest
	getToken.RequestMsgID = message.MsgIDMultiTaskDescription.Value() + 1
	getTokenData, _ := proto.Marshal(&getToken)

	ctxto, cancels := context.WithTimeout(context.Background(), time.Second*10)
	defer cancels()
	tok, err := clt.SendMsg(ctxto, message.MsgIDNodeCapacityRequest.Value(), getTokenData)

	if err != nil {
		if config.Gconfig.ElkReport {
			//logelk:=re.MakeReportLog(id,hash,"failToken",err)
			//go re.reportLog(logelk)
		}

		return nil, err
	}

	if len(tok) < 3 {
		err = fmt.Errorf("the length of token less 3 byte")
		if config.Gconfig.ElkReport {
			//logelk:=re.MakeReportLog(id,hash,"failToken",err)
			//go re.reportLog(logelk)
		}
		return nil, err
	}
	return tok, err
}

// TaskMsgResult 重建任务结果对象
type TaskMsgResult struct {
	ID          []byte // 重建任务ID
	RES         int32  // 重建任务结果 0：success 1：error
	BPID        int32  // 需要回复的BP的ID
	ExpiredTime int64  // 任务过期时间
	SrcNodeID   int32  // 来源节点ID
	ErrorMsg    error
}

// 多重建任务消息处理
func (re *Engine) HandleMuilteTaskMsg(msgData []byte) error {
	var mtdMsg message.MultiTaskDescription
	if err := proto.Unmarshal(msgData, &mtdMsg); err != nil {
		return err
	}

	//要先判断一下队列的剩余长度是否能容纳当前任务不能的话不接收，返回错误
	queueLen := re.waitQueue.Len()
	if re.waitQueue.Max - queueLen < len(mtdMsg.Tasklist) {
		log.Printf("[recover] queue space is not enough, max len is %d, " +
			"current len is %d, tasks is %d\n",
			re.waitQueue.Max, queueLen, len(mtdMsg.Tasklist))
		return fmt.Errorf("queue space is not enough\n")
	}

	log.Printf("[recover] queue max len is %d, current len is %d, tasks is %d\n",
		re.waitQueue.Max, queueLen, len(mtdMsg.Tasklist))

	for _, task := range mtdMsg.Tasklist {
		bys := task[12:14]
		bytebuff := bytes.NewBuffer(bys)
		var snID int16
		_ = binary.Read(bytebuff, binary.BigEndian, &snID)

		if err := re.waitQueue.PutTask(task, int32(snID), mtdMsg.ExpiredTime,
				mtdMsg.SrcNodeID, mtdMsg.ExpiredTimeGap, time.Now(), 0); err != nil {
			log.Printf("[recover] put recover task error: %s\n", err.Error())
		}
	}
	return nil
}

/**
 * @Description: 分发不同类型的任务给不同执行器，目前只考虑LRC任务
 * @receiver re
 * @param ts 任务
 * @param pkgstart 任务包开始时间
 */

var tskcnt uint64
func (re *Engine) dispatchTask(ts *Task) {
	var msgID uint16
	err := binary.Read(bytes.NewBuffer(ts.Data[:2]), binary.BigEndian, &msgID)
	if err != nil {
		log.Println("[recover] dispatchTask get msgId err:", err.Error())
		return
	}
	var res *TaskMsgResult

	taskC := atomic.AddUint64(&tskcnt, 1)

	switch int32(msgID) {
	case message.MsgIDLRCTaskDescription.Value():
		atomic.AddUint64(&statistics.DefaultStatusCount.Total, 1)
		ts.ExecTimes++
		log.Printf("[recover] execLRCTask, msgId: %d exec times %d\n", msgID, ts.ExecTimes)
		if ts.ExecTimes == 1 {
			log.Printf("[recover] execLRCTask exec_task %d, msgId: %d\n", taskC, msgID)
		}

		startTime := time.Now()
		taskPf := new(statistics.PerformanceStat)
		res = re.execLRCTask(ts.Data[2:], ts.ExpiredTime, ts.StartTime,
				ts.TaskLife, ts.SrcNodeID, taskPf)
		useTime := time.Now().Sub(startTime).Milliseconds()
		log.Printf("[recover] execLRCTask use time is %d ms\n", useTime)
		taskPf.ExecTimes = useTime
		taskPf.TaskPfStat()

		select {
		case statistics.PfStatChan <- taskPf:
		default:
			log.Printf("[recover] task=%d performance stat info is discarded\n", taskPf.TaskId)
		}

		if int32(time.Now().Sub(ts.StartTime)) < ts.TaskLife &&
			res.RES == 1 &&  ts.ExecTimes < 2 {
			err := re.waitQueue.PutTask(ts.Data, ts.SnID, ts.ExpiredTime, ts.SrcNodeID,
				ts.TaskLife, ts.StartTime, ts.ExecTimes)
			if err != nil {
				goto Reply
			}else {
				break
			}
		}
	Reply:
		if res.ErrorMsg != nil {
			log.Println("[recover] error:", res.ErrorMsg,)
			res.RES = 1
		}
		res.BPID = ts.SnID
		res.SrcNodeID = ts.SrcNodeID

		if res.RES == 0 {
			atomic.AddUint64(&statistics.DefaultStatusCount.Error, 1)
			log.Printf("[recover] execLRCTask success %d, msgId: %d\n", taskC, msgID)
		}else {
			log.Printf("[recover] execLRCTask fail %d, msgId: %d\n", taskC, msgID)
		}

		re.PutReplyQueue(res)
	case message.MsgIDTaskDescriptCP.Value():
		log.Println("[recover] execCPTask, msgId:", msgID)
		log.Printf("[recover] execCPTask exec_task %d, msgId: %d\n", taskC, msgID)

		res = re.execCPTask(ts.Data[2:], ts.ExpiredTime)

		if res.RES == 0 {
			log.Printf("[recover] execCPTask success %d, msgId: %d\n", taskC, msgID)
		}else {
			log.Printf("[recover] execCPTask fail %d, msgId: %d\n", taskC, msgID)
		}

		res.BPID = ts.SnID
		res.SrcNodeID = ts.SrcNodeID
		re.PutReplyQueue(res)
	default:
		log.Println("[recover] unknown msgID:", msgID)
	}
}

func (re *Engine) PutReplyQueue(res *TaskMsgResult) {
	select {
		case re.replyQueue <- res:
			log.Printf("[recover] reply queue enqueue task=%d\n", binary.BigEndian.Uint64(res.ID[:8]))
		//default:
	}
}

func (re *Engine) reply(res *TaskMsgResult) error {
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

var replycnt uint64
//
func (re *Engine) MultiReply() error {
	var resmsg = make(map[int32]*message.MultiTaskOpResult)

	func() {
		for i := 0; i < max_reply_num; i++ {
			if len(re.replyQueue) > 0 {
				//这个时候可能不大于0， 所以这个打印不准确
				log.Printf("[recover] reply queue length is %d\n", len(re.replyQueue))
			}
			select {
			//这个时候可能就有数据了
			case res := <-re.replyQueue:
				if resmsg[res.BPID] == nil {
					resmsg[res.BPID] = &message.MultiTaskOpResult{}
				}
				_r := resmsg[res.BPID]
				log.Println("[recover_debugtime]  reply task=", binary.BigEndian.Uint64(res.ID[:8]))
				_r.Id = append(_r.Id, res.ID)
				_r.RES = append(_r.RES, res.RES)
				_r.ExpiredTime = res.ExpiredTime
				_r.SrcNodeID = res.SrcNodeID
				resmsg[res.BPID] = _r
				statistics.DefaultRebuildCount.IncReportRbdTask()
			case <-time.After(max_reply_wait_time):
				continue
			}
		}
	}()

	for k, v := range resmsg {
		replycnt++
		v.NodeID = int32(re.sn.Config().IndexID)
		data, err := proto.Marshal(v)
		if err != nil {
			log.Printf("[recover][report] marsnal failed %s\n", err.Error())
			continue
		}

		for reportTms := 0; reportTms < 5; reportTms++ {
			if isReturn, err := re.tryReply(int(k), data); err != nil {
				log.Printf("[recover][report] tryReply error: %s\n", err.Error())
				if !isReturn {
					// 如果报错且sn没有返回继续循环
					continue
				}
			} else {
				log.Printf("[recover][report] tryReply success, report nums %d\n", len(v.RES))
			}
			// 如果不报错退出循环
			break
		}
		if replycnt % 100 == 0 {
			log.Println("[recover][report]MultiReply,NodeID=",v.NodeID,"srcnodeid=",
				v.SrcNodeID,"id=",v.Id,"res=",v.RES,"replycnt", replycnt)
		}

		log.Printf("[recover] [report] reply count %d\n", replycnt)
	}



	return nil
}

/**
 * @Description: 尝试回报重建结果
 * @receiver re
 * @param index
 * @param data
 * @return bool
 * @return error
 */
func (re *Engine) tryReply(index int, data []byte) (bool, error) {
	resp, err := re.sn.SendBPMsg(index, message.MsgIDMultiTaskOPResult.Value(), data)
	if err != nil {
		return false, err
	}

	if len(resp) < 3 {
		return false, fmt.Errorf("response too short %d %s", len(resp), hex.EncodeToString(resp[0:2]))
	}

	var res message.MultiTaskOpResultRes
	err = proto.Unmarshal(resp[2:], &res)
	if err != nil {
		return false, err
	}
	log.Println("[recover][report] tryreply resErrcode=",res.ErrCode,"res.SuccNum",res.SuccNum)
	if 0 == res.ErrCode {
		statistics.DefaultRebuildCount.IncAckSuccRebuild(uint64(res.SuccNum))
	} else {
		return true, fmt.Errorf("sn return error %d", res.ErrCode)
	}
	return false, nil
}

func (re *Engine) execRCTask(msgData []byte, expired int64) *TaskMsgResult {
	var res TaskMsgResult
	res.ExpiredTime = expired
	var msg message.TaskDescription
	if err := proto.Unmarshal(msgData, &msg); err != nil {
		log.Printf("[recover]proto解析错误%s", err)
		res.RES = 1
	}
	res.ID = msg.Id
	if err := re.recoverShard(&msg); err != nil {
		res.RES = 1
	} else {
		res.RES = 0
	}
	return &res
}

type PreJudgeReport struct {
	LocalNdID  string
	LostHash   string
	LostIndex  uint16
	FailType   string
	ShardExist string
}

func (re *Engine) MakeJudgeElkReport(lrcShd *lrcpkg.Shardsinfo, msg message.TaskDescription) *PreJudgeReport {
	//if ! config.Gconfig.ElkReport{
	//	return nil
	//}
	localid := re.sn.Config().ID
	lostidx := lrcShd.Lostindex
	losthash := base64.StdEncoding.EncodeToString(msg.Hashs[msg.RecoverId])
	failtype := "failJudge"
	shardExist := lrcShd.ShardExist[:164]
	shdExistStr := make([]string, len(shardExist))
	for k, v := range shardExist {
		shdExistStr[k] = fmt.Sprintf("%d", v)
	}
	strExist := strings.Join(shdExistStr, "")
	return &PreJudgeReport{
		LocalNdID:  localid,
		LostHash:   losthash,
		LostIndex:  lostidx,
		FailType:   failtype,
		ShardExist: strExist,
	}
}

/**
 * @Description: 用任务消息初始化LRC任务句柄
 * @receiver re
 * @param msg
 * @return *LRCHandler
 * @return error
 */
func (re *Engine) initLRCHandlerByMsg(msg message.TaskDescription) (*LRCHandler, error) {
	lrc := &lrcpkg.Shardsinfo{}

	lrc.OriginalCount = uint16(len(msg.Hashs) - int(msg.ParityShardCount))
	lrc.RecoverNum = 13
	lrc.Lostindex = uint16(msg.RecoverId)
	return re.le.GetLRCHandler(lrc)
}

/**
 * @Description: 验证重建后的数据并保存
 * @receiver re
 * @param recoverData
 * @param msg
 * @param res
 * @return *TaskMsgResult
 */
func (re *Engine) verifyLRCRecoveredDataAndSave(recoverData []byte, msg message.TaskDescription, res *TaskMsgResult) error {
	hashBytes := md5.Sum(recoverData)
	hash := hashBytes[:]

	if !bytes.Equal(hash, msg.Hashs[msg.RecoverId]) {
		statistics.DefaultRebuildCount.IncFailRbd()

		return fmt.Errorf(
			"[recover]fail shard saved %s recoverID %x hash %s\n",
			BytesToInt64(msg.Id[0:8]),
			msg.RecoverId,
			base58.Encode(msg.Hashs[msg.RecoverId]),
		)
	}

	var key [common.HashLength]byte
	copy(key[:], hash)

	if _, err := re.sn.YTFS().BatchPut(map[common.IndexTableKey][]byte{common.IndexTableKey{Hsh:key, Id:0}: recoverData}); err != nil && err.Error() != "YTFS: hash key conflict happens" {
		statistics.DefaultRebuildCount.IncFailRbd()
		return fmt.Errorf("[recover]LRC recover shard saved failed%s\n", err)
	}
	return nil
}

/**
 * @Description: 执行lrc 重建任务
 * @receiver re
 * @param msgData 单个重建消息
 * @param expired 过期时间
 * @param pkgstart 任务包开始时间
 * @param tasklife 任务存活周期
 * @return *TaskMsgResult 任务执行结果
 */
func (re *Engine) execLRCTask(msgData []byte, expired int64, StartTime time.Time,
			taskLife int32, srcNodeid int32, taskPf *statistics.PerformanceStat) (res *TaskMsgResult) {
	// @TODO 初始化返回
	res = &TaskMsgResult{}

	res.ExpiredTime = expired
	res.RES = 1

	taskActuator := actuator.New(re.DefaultDownloader)
	taskActuator.SetPfStat(taskPf)

	defer taskActuator.Free()

	var recoverData []byte
	var realHash []byte

	// @TODO 执行恢复任务
	for _, opts := range []actuator.Options{
		{
			Expired: taskLife,
			STime:   StartTime,
			Stage:   actuator.RECOVER_STAGE_CP,
		},
		{
			Expired: taskLife,
			STime:   StartTime,
			Stage:   actuator.RECOVER_STAGE_ROW,
		},
		{
			Expired: taskLife,
			STime:   StartTime,
			Stage:   actuator.RECOVER_STAGE_COL,
		},
		{
			Expired: taskLife,
			STime:   StartTime,
			Stage:   actuator.RECOVER_STAGE_FULL,
		},
	} {
		switch opts.Stage {
		case 1:
			atomic.AddUint64(&statistics.DefaultRebuildCount.RowRebuildCount, 1)
		case 2:
			atomic.AddUint64(&statistics.DefaultRebuildCount.ColRebuildCount, 1)
		case 3:
			atomic.AddUint64(&statistics.DefaultRebuildCount.GlobalRebuildCount, 1)
		}

		st := time.Now()
		data, resID, srcHash, err := taskActuator.ExecTask(
			msgData,
			opts,
		)
		res.ID = resID
		log.Printf("[recover]  task=%d stage=%d ExecTask used time %d ms\n",
			binary.BigEndian.Uint64(res.ID[:8]), opts.Stage, time.Now().Sub(st).Milliseconds())

		var downloads = 0
		sMap := taskActuator.GetdownloadShards()
		if sMap != nil {
			for _, v := range sMap.GetMap() {
				if v.Data != nil {
					downloads++
				}
			}

			if downloads != 0 {
				log.Printf("[recover] task=%d stage=%d real_downloads=%d\n",
					binary.BigEndian.Uint64(res.ID[:8]), opts.Stage, downloads)
			}
		}

		realHash = srcHash

		if err != nil {
			log.Println("[recover_debugtime] ExecTask error:", err.Error())
		}

		//log.Println("[recover_debugtime] ExecTask end, taskid=", base58.Encode(resID))
		// @TODO 如果重建成功退出循环
		if err == nil && data != nil {
			recoverData = data
			switch opts.Stage {
			case 0:
				statistics.DefaultRebuildCount.IncBackupRbdSucc()
			case 1:
				statistics.DefaultRebuildCount.IncRowRbdSucc()
			case 2:
				statistics.DefaultRebuildCount.IncColRbdSucc()
			case 3:
				statistics.DefaultRebuildCount.IncGlobalRbdSucc()
			}
			break
		}
	}

	if taskPf != nil {
		taskPf.TaskId = binary.BigEndian.Uint64(res.ID[:8])
	}

	if recoverData == nil {
		log.Printf("[recover] all rebuild stage fail task=%d src node id %d source hash key is %s\n",
			binary.BigEndian.Uint64(res.ID[:8]), srcNodeid, base58.Encode(realHash))

		//test
		var downloads = 0
		sMap := taskActuator.GetdownloadShards()
		if sMap != nil {
			for _, v := range sMap.GetMap() {
				if v.Data != nil {
					downloads++
				}
			}

			if downloads != 0 {
				log.Printf("[recover] task=%d all stage fail real downloads %d\n",
					binary.BigEndian.Uint64(res.ID[:8]), downloads)
			}
		}
		//test------

		res.ErrorMsg = fmt.Errorf("all rebuild stage fail")
		res.RES = 1
		statistics.DefaultRebuildCount.IncFailLessShard()
		return
	}

	// @TODO 存储重建完成的分片
	hashBytes := md5.Sum(recoverData)
	hash := hashBytes[:]
	var key [common.HashLength]byte
	copy(key[:], hash)
	_, err := re.sn.YTFS().BatchPut(map[common.IndexTableKey][]byte{common.IndexTableKey{Hsh:key, Id:0}: recoverData})
	if err != nil {
		log.Printf("[recover] fail task=%d src node id %d recover hash key %s, source hash key is %s\n",
			binary.BigEndian.Uint64(res.ID[:8]), srcNodeid, base58.Encode(key[:]), base58.Encode(realHash))

		if err.Error() != "YTFS: hash key conflict happens" {
			res.ErrorMsg = fmt.Errorf("[recover] LRC recover shard saved failed %s", err.Error())
		}

		return
	}else {
		log.Printf("[recover] success task=%d src node id %d recover hash key %s, source hash key is %s\n",
			binary.BigEndian.Uint64(res.ID[:8]), srcNodeid, base58.Encode(key[:]), base58.Encode(realHash))
	}

	//test
	var downloads = 0
	sMap := taskActuator.GetdownloadShards()
	if sMap != nil {
		for _, v := range sMap.GetMap() {
			if v.Data != nil {
				downloads++
			}
		}
	}

	if downloads != 0 {
		log.Printf("[recover] task=%d all stage success real downloads %d\n",
			binary.BigEndian.Uint64(res.ID[:8]), downloads)
	}
	//test------

	res.RES = 0
	//log.Println("恢复成功")
	return
}

// 副本集任务
func (re *Engine) execCPTask(msgData []byte, expired int64) *TaskMsgResult {
	var msg message.TaskDescriptionCP
	var result TaskMsgResult
	result.ExpiredTime = expired
	err := proto.UnmarshalMerge(msgData, &msg)
	if err != nil {
		log.Printf("[recover] execCPTask 解析错误%s\n", err.Error())
		return nil
	}

	result.ID = msg.Id
	result.RES = 1
	var number int	//这他妈是干啥的？？？？？？服了
	// 循环从副本节点获取分片，只要有一个成功就返回
	for _, v := range msg.Locations {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		shard, err := re.getShard2(ctx, v.NodeId, base58.Encode(msg.Id), v.Addrs, msg.DataHash, &number)
		cancel()

		// 如果没有发生错误，分片下载成功，就存储分片
		if err == nil {
			hashBytes := md5.Sum(shard)
			var key [common.HashLength]byte
			copy(key[:], hashBytes[:])

			var vhf [16]byte
			copy(vhf[:], msg.DataHash)
			log.Printf("[recover:%s] execCPTask--, get shard DataHash %s shard len %d, remote miner NodeId:%s Addr:%s\n",
				base58.Encode(msg.DataHash), base58.Encode(key[:]), len(shard), v.NodeId, v.Addrs)

			_, err := re.sn.YTFS().BatchPut(map[common.IndexTableKey][]byte{common.IndexTableKey{Hsh:vhf, Id:0}:shard})
			// 存储分片没有错误，或者分片已存在返回0，代表成功
			if err != nil && (err.Error() != "YTFS: hash key conflict happens" ||
				err.Error() == "YTFS: conflict hash value") {
				log.Printf("[recover:%s] execCPTask, YTFS Put error %s\n",
					base58.Encode(vhf[:]), err.Error())
				result.RES = 1
			} else {
				log.Printf("[recover:%s] execCPTask success\n",
					base58.Encode(msg.DataHash))
				result.RES = 0
				break
			}
		}else{
			log.Printf("[recover:%s] execCPTask error %s\n",
				base58.Encode(msg.DataHash), err.Error())
		}
	}
	return &result
}

//BytesToInt64 convet byte slice to int64
func BytesToInt64(bys []byte) int64 {
	bytebuff := bytes.NewBuffer(bys)
	var data int64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return data
}
