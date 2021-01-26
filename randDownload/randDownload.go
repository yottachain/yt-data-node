package randDownload

import (
	"context"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/yottachain/YTDataNode/Perf"
	"github.com/yottachain/YTDataNode/TokenPool"
	"github.com/yottachain/YTDataNode/activeNodeList"
	"github.com/yottachain/YTDataNode/config"
	"github.com/yottachain/YTDataNode/logBuffer"
	log "github.com/yottachain/YTDataNode/logger"
	"github.com/yottachain/YTDataNode/message"
	"github.com/yottachain/YTDataNode/statistics"
	"github.com/yottachain/YTDataNode/storageNodeInterface"
	"github.com/yottachain/YTDataNode/util"
	"math"
	"math/rand"
	"os"
	"sync/atomic"
	"time"
)

var stop = false

var errNoTK = fmt.Errorf("notk")

var Sn storageNodeInterface.StorageNode
var wl = activeNodeList.NewWeightNodeList(
	time.Minute*5,
	time.Minute*time.Duration(config.Gconfig.NodeListUpdateTime),
	config.Gconfig.RandDownloadGroupSize,
	util.GetSelfIP())

func GetRandNode() (*peer.AddrInfo, error) {
	nodeList := wl.Get()

	pi := &peer.AddrInfo{}
	nl := len(nodeList)

	var randNode *activeNodeList.Data

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

func UploadFromRandNode(ctx context.Context) error {
	tk, err := TokenPool.Utp().Get(ctx, Sn.Host().Config().ID, 0)
	if err != nil {
		return errNoTK
	}
	defer TokenPool.Utp().Delete(tk)

	pi, err := GetRandNode()
	if err != nil {
		return err
	}
	if Sn == nil {
		return fmt.Errorf("no storage-node")
	}

	statistics.DefaultStat.RXTestConnectRate.AddCount()
	clt, err := Sn.Host().ClientStore().Get(ctx, pi.ID, pi.Addrs)
	if err != nil {
		return err
	}
	statistics.DefaultStat.RXTestConnectRate.AddSuccess()
	defer clt.Close()

	var testMsg message.TestGetBlock

	// 第一次发送消息模拟下载
	testMsg.Msg = Perf.MSG_DOWNLOAD
	testMsgBuf, err := proto.Marshal(&testMsg)
	if err != nil {
		return err
	}
	statistics.DefaultStat.RXTest.AddCount()
	_, err = clt.SendMsg(ctx, message.MsgIDTestGetBlock.Value(), testMsgBuf)
	if err != nil {
		return err
	}
	statistics.DefaultStat.RXTest.AddSuccess()
	return nil
}

func DownloadFromRandNode(ctx context.Context) error {
	var err error

	pi, err := GetRandNode()
	if err != nil {
		return err
	}

	if Sn == nil {
		return fmt.Errorf("no storage-node")
	}

	statistics.DefaultStat.RXTestConnectRate.AddCount()
	clt, err := Sn.Host().ClientStore().Get(ctx, pi.ID, pi.Addrs)
	if err != nil {
		return err
	}
	statistics.DefaultStat.RXTestConnectRate.AddSuccess()
	defer clt.Close()

	var getTokenMsg message.NodeCapacityRequest
	getTokenMsg.RequestMsgID = message.MsgIDTestGetBlock.Value()
	getTKMsgBuf, err := proto.Marshal(&getTokenMsg)
	if err != nil {
		return errNoTK
	}

	getTKResBuf, err := clt.SendMsg(ctx, message.MsgIDNodeCapacityRequest.Value(), getTKMsgBuf)
	if err != nil {
		return errNoTK
	}
	var tokenMsg message.NodeCapacityResponse
	err = proto.Unmarshal(getTKResBuf[2:], &tokenMsg)
	if err != nil {
		return errNoTK
	}

	var testMsg message.TestGetBlock

	// 第一次发送消息模拟下载
	testMsg.Msg = Perf.MSG_DOWNLOAD
	testMsgBuf, err := proto.Marshal(&testMsg)
	if err != nil {
		return err
	}
	_, err = clt.SendMsg(ctx, message.MsgIDTestGetBlock.Value(), testMsgBuf)
	if err != nil {
		return err
	}
	// 第二次发送消息消耗token
	testMsg.Msg = Perf.MSG_CHECKOUT
	testMsgBuf, err = proto.Marshal(&testMsg)
	if err != nil {
		return err
	}
	_, err = clt.SendMsg(ctx, message.MsgIDTestGetBlock.Value(), testMsgBuf)
	if err != nil {
		return err
	}
	return nil
}
func Stop() {
	stop = true
}
func RunRX() {
	var successCount uint64
	var errorCount uint64
	var execChan *chan struct{}
	rand.Seed(int64(os.Getpid()))

	go func() {
		for {
			<-time.After(time.Minute)
			log.Println("[randUpload] success", successCount, "error", errorCount, "exec", len(*execChan))
		}
	}()

	c := make(chan struct{}, int(math.Min(float64(TokenPool.Utp().GetTFillTKSpeed())/2, float64(config.Gconfig.RandDownloadNum))))
	execChan = &c

	go func() {
		for {
			c := make(chan struct{}, int(math.Min(float64(TokenPool.Utp().GetTFillTKSpeed())/2, float64(config.Gconfig.RandDownloadNum))))
			execChan = &c
			<-time.After(5 * time.Minute)
		}
	}()

	// rx
	for {
		if stop {
			<-time.After(time.Minute)
			continue
		}
		if execChan == nil {
			continue
		}
		ec := *execChan
		ec <- struct{}{}
		go func(ec chan struct{}) {
			defer func() {
				<-ec
			}()

			ctx, cancle := context.WithTimeout(context.Background(), time.Second*time.Duration(config.Gconfig.TTL))
			defer cancle()

			err := UploadFromRandNode(ctx)
			if err != nil && err.Error() != errNoTK.Error() {
				logBuffer.ErrorLogger.Println(err.Error())
				atomic.AddUint64(&errorCount, 1)
			} else if err == nil {
				atomic.AddUint64(&successCount, 1)
			}

		}(ec)
		<-time.After(time.Millisecond * time.Duration(config.Gconfig.RandDownloadSleepTime))
	}
}

func RunTX() {
	var successCount uint64
	var errorCount uint64
	var execChan *chan struct{}
	rand.Seed(int64(os.Getpid()))

	go func() {
		for {
			<-time.After(time.Minute)
			log.Println("[randDownload] success", successCount, "error", errorCount, "exec", len(*execChan))
		}
	}()

	c := make(chan struct{}, int(math.Min(float64(TokenPool.Utp().GetTFillTKSpeed())/2, float64(config.Gconfig.RandDownloadNum))))
	execChan = &c

	go func() {
		for {
			c := make(chan struct{}, int(math.Min(float64(TokenPool.Utp().GetTFillTKSpeed())/2, float64(config.Gconfig.RandDownloadNum))))
			execChan = &c
			<-time.After(5 * time.Minute)
		}
	}()

	// tx
	for {
		if execChan == nil {
			continue
		}
		ec := *execChan
		ec <- struct{}{}
		go func(ec chan struct{}) {
			defer func() {
				<-ec
			}()

			ctx, cancle := context.WithTimeout(context.Background(), time.Second*time.Duration(config.Gconfig.TTL))
			defer cancle()

			err := DownloadFromRandNode(ctx)
			if err != nil && err.Error() != errNoTK.Error() {
				logBuffer.ErrorLogger.Println(err.Error())
				atomic.AddUint64(&errorCount, 1)
			} else if err == nil {
				atomic.AddUint64(&successCount, 1)
			}

		}(ec)
		<-time.After(time.Millisecond * time.Duration(config.Gconfig.RandDownloadSleepTime))
	}
}

func Run() {
	go RunRX()
	go RunTX()
}
