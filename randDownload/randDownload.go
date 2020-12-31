package randDownload

import (
	"context"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/yottachain/YTDataNode/TaskPool"
	"github.com/yottachain/YTDataNode/activeNodeList"
	"github.com/yottachain/YTDataNode/config"
	log "github.com/yottachain/YTDataNode/logger"
	"github.com/yottachain/YTDataNode/message"
	"github.com/yottachain/YTDataNode/statistics"
	"github.com/yottachain/YTDataNode/storageNodeInterface"
	"math/rand"
	"sync/atomic"
	"time"
)

var Sn storageNodeInterface.StorageNode

func GetRandNode() (*peer.AddrInfo, error) {
	nodeList := activeNodeList.GetNodeListByTimeAndGroupSize(time.Minute*time.Duration(config.Gconfig.NodeListUpdateTime), config.Gconfig.RandDownloadGroupSize)

	pi := &peer.AddrInfo{}
	nl := len(nodeList)

	var randNode activeNodeList.Data

	if nl <= 0 {
		return nil, fmt.Errorf("no node")
	}
	randIndex := rand.Intn(nl)
	randNode = nodeList[randIndex]

	id, err := peer.IDB58Decode(randNode.NodeID)
	if err != nil {
		return nil, err
	}
	pi.ID = id
	for _, v := range randNode.IP {
		ma, err := multiaddr.NewMultiaddr(v)
		if err != nil {
			continue
		}
		pi.Addrs = append(pi.Addrs, ma)
	}
	return pi, nil
}

func DownloadFromRandNode() error {
	var err error
	atomic.AddInt64(&statistics.DefaultStat.RandDownloadCount, 1)

	pi, err := GetRandNode()
	if err != nil {
		return err
	}

	if Sn == nil {
		return fmt.Errorf("no storage-node")
	}

	ctx, cancle := context.WithTimeout(context.Background(), time.Second*30)
	defer cancle()
	clt, err := Sn.Host().ClientStore().Get(ctx, pi.ID, pi.Addrs)
	if err != nil {
		return err
	}

	utk, err := TaskPool.Utp().Get(ctx, Sn.Host().Config().ID, 0)
	if err != nil {
		return err
	}

	var getTokenMsg message.NodeCapacityRequest
	getTokenMsg.RequestMsgID = message.MsgIDDownloadShardRequest.Value() + 1
	getTKMsgBuf, err := proto.Marshal(&getTokenMsg)
	if err != nil {
		TaskPool.Utp().Delete(utk)
		return err
	}

	getTKResBuf, err := clt.SendMsg(ctx, message.MsgIDNodeCapacityRequest.Value(), getTKMsgBuf)
	if err != nil {
		return err
	}
	var tokenMsg message.NodeCapacityResponse
	err = proto.Unmarshal(getTKResBuf[2:], &tokenMsg)
	if err != nil {
		return err
	}

	var downloadMsg message.TestGetBlock
	downloadBuf, err := proto.Marshal(&downloadMsg)
	if err != nil {
		return err
	}
	_, err = clt.SendMsg(ctx, message.MsgIDTestGetBlock.Value(), downloadBuf)
	if err != nil {
		return err
	}
	var checkTKMsg message.DownloadTKCheck
	checkTKMsg.Tk = tokenMsg.AllocId
	checkTKBuf, err := proto.Marshal(&checkTKMsg)
	if err != nil {
		return err
	}

	_, err = clt.SendMsgClose(ctx, message.MsgIDDownloadTKCheck.Value(), checkTKBuf)
	if err != nil {
		return err
	}
	TaskPool.Utp().Delete(utk)

	atomic.AddInt64(&statistics.DefaultStat.RandDownloadSuccess, 1)
	return nil
}

func Run() {
	var successCount uint64
	var errorCount uint64
	var execCount int64

	go func() {
		for {
			<-time.After(time.Minute)
			log.Println("[randDownload] success", successCount, "error", errorCount, "exec", atomic.LoadUint64(&execCount))
		}
	}()
	for {
		ec := atomic.LoadInt64(&execCount)
		if ec < int64(TaskPool.Utp().GetTFillTKSpeed())/2 || ec < int64(config.Gconfig.RandDownloadNum) {
			go func() {
				atomic.AddInt64(&execCount, 1)
				defer atomic.AddInt64(&execCount, -1)

				err := DownloadFromRandNode()
				if err != nil {
					atomic.AddUint64(&errorCount, 1)
				} else {
					atomic.AddUint64(&successCount, 1)
				}
			}()
		} else {
			time.Sleep(time.Millisecond * 100)
		}
	}
}
