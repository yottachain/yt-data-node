package recover

import (
	"context"
	"github.com/gogo/protobuf/proto"
	"github.com/klauspost/reedsolomon"
	"github.com/mr-tron/base58/base58"
	"github.com/yottachain/YTDataNode/host"
	log "github.com/yottachain/YTDataNode/logger"
	"github.com/yottachain/YTDataNode/message"
	"github.com/yottachain/YTFS"
	"github.com/yottachain/YTFS/common"
	"sync"
	"time"
)

type RecoverEngine struct {
	reedsolomon.Encoder
	host *host.Host
	ytfs *ytfs.YTFS
}

func New(hst *host.Host, yt *ytfs.YTFS) (*RecoverEngine, error) {
	var re = new(RecoverEngine)
	encoder, err := reedsolomon.New(128, 32)
	if err != nil {
		return nil, err
	}
	re.host = hst
	re.ytfs = yt
	re.Encoder = encoder
	return re, nil
}

func (re *RecoverEngine) ExecRecoverTask(description *message.TaskDescription) error {
	var shards [][]byte = make([][]byte, 160)
	var wg = sync.WaitGroup{}
	var number int
	wg.Add(160)
	for k, v := range description.Locations {
		go func(k int, v *message.P2PLocation) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			shard, err := re.getShard(ctx, v.NodeId, v.Addrs, description.Hashs[k], &number)
			if err == nil {
				shards[k] = shard
			} else {
				log.Printf("[recover:%s]error:%s\n", base58.Encode(description.Id), err.Error())
			}
		}(k, v)
	}
	wg.Wait()
	shards[description.RecoverId] = nil
	err := re.Encoder.Reconstruct(shards)
	if err != nil {
		log.Printf("[recover:%s]datas recover error:%s\n", base58.Encode(description.Id), err.Error())
		return err
	}
	log.Printf("[recover:%s]datas recover success\n", base58.Encode(description.Id))
	var vhf [32]byte
	err = re.ytfs.Put(common.IndexTableKey(vhf), shards[int(description.RecoverId)])
	if err != nil {
		return err
	}
	return nil
}

func (re *RecoverEngine) getShard(ctx context.Context, id string, addrs []string, hash []byte, n *int) ([]byte, error) {
	err := re.host.ConnectAddrStrings(id, addrs)
	if err != nil {
		return nil, err
	}
	stm, err := re.host.NewMsgStream(ctx, id, "/node/0.0.2")
	if err != nil {
		return nil, err
	}
	var msg message.DownloadShardRequest
	var res message.DownloadShardResponse
	msg.VHF = hash
	buf, err := proto.Marshal(&msg)
	if err != nil {
		return nil, err
	}
	shardBuf, err := stm.SendMsgGetResponse(append(message.MsgIDDownloadShardRequest.Bytes(), buf...))

	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(shardBuf[2:], &res)
	if err != nil {
		return nil, err
	}
	*n = *n + 1
	log.Printf("[recover:%s]get shard success[%d]\n", base58.Encode(hash), *n)
	return res.Data, nil
}
func (re *RecoverEngine) HandleMsg(msgData []byte, stm *host.MsgStream) error {
	var msg message.TaskDescription
	var res message.TaskOpResult
	if err := proto.Unmarshal(msgData, &msg); err != nil {
		return err
	}
	res.Id = msg.Id
	if err := re.ExecRecoverTask(&msg); err != nil {
		res.RES = 1
	} else {
		res.RES = 0
	}
	buf, err := proto.Marshal(&res)
	if err != nil {
		return err
	}
	if err := re.replay(buf, stm); err != nil {
		return err
	}
	log.Printf("[recover:%s]success\n", base58.Encode(msg.Id))
	return nil
}

func (re *RecoverEngine) replay(data []byte, stm *host.MsgStream) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := re.host.ConnectAddrStrings(stm.Conn().RemotePeer().Pretty(), []string{stm.Conn().RemoteMultiaddr().String()}); err != nil {
		return nil
	}
	stm, err := re.host.NewMsgStream(ctx, stm.Conn().RemotePeer().Pretty(), "/bp/0.0.2")
	if err != nil {
		return err
	}
	stm.SendMsgClose(append(message.MsgIDTaskOPResult.Bytes(), data...))
	return nil
}
