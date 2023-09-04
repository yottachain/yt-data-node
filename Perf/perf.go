package Perf

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/yottachain/YTDataNode/TokenPool"
	"github.com/yottachain/YTDataNode/config"
	log "github.com/yottachain/YTDataNode/logger"
	"github.com/yottachain/YTDataNode/message"
	"github.com/yottachain/YTDataNode/statistics"
	"github.com/yottachain/YTDataNode/storageNodeInterface"
	"github.com/yottachain/YTHost/client"
)

const MSG_DOWNLOAD = "download"
const MSG_CHECKOUT = "checkout"
const MSG_UPLOAD = "upload"
const MSG_DOWNLOAD_TK = "dtk"
const MSG_UPLOAD_TK = "utk"

var Sn storageNodeInterface.StorageNode

func TestMinerPerfHandler(data []byte) (res []byte, err error) {
	var successCount int64
	var errorCount int64
	var successLatency int64
	var errorLatency int64

	var task message.TestMinerPerfTask
	err = proto.UnmarshalMerge(data, &task)
	if err != nil {
		return
	}

	var pi = &peer.AddrInfo{}
	// 解系地址
	for _, addr := range task.TargetMa {
		ma, err2 := multiaddr.NewMultiaddr(addr)
		if err2 != nil {
			continue
		}
		i, err2 := peer.AddrInfoFromP2pAddr(ma)
		if err2 != nil {
			continue
		}
		pi.ID = i.ID
		pi.Addrs = append(pi.Addrs, i.Addrs...)
	}

	if len(pi.Addrs) <= 0 {
		err = fmt.Errorf("no addr")
		return
	}

	// 建立连接
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(task.TimeOut))
	defer cancel()
	clt, err := Sn.Host().ClientStore().Get(ctx, pi.ID, pi.Addrs)
	if err != nil {
		return
	}

	outtime := time.Now().Add(time.Duration(task.TestTime) * time.Second)
	for {
		if outtime.Unix() < time.Now().Unix() {
			break
		}

		timeStart := time.Now()
		testerr := testOne(clt, &task, task.TimeOut)
		timeEnd := time.Now()

		if testerr == nil {
			successLatency += timeEnd.Sub(timeStart).Milliseconds()
			successCount += 1
		} else {
			errorLatency += timeEnd.Sub(timeStart).Milliseconds()
			errorCount += 1
		}
	}

	// 构造返回消息
	var minerPerfResMsg message.TestMinerPerfTaskRes
	minerPerfResMsg.TargetMa = task.TargetMa
	minerPerfResMsg.TestType = task.TestType
	minerPerfResMsg.SuccessCount = successCount
	minerPerfResMsg.ErrorCount = errorCount
	minerPerfResMsg.SuccessLatency = successLatency
	minerPerfResMsg.ErrorLatency = errorLatency

	res, err = proto.Marshal(&minerPerfResMsg)
	log.Println("[test] test task return", minerPerfResMsg)
	return
}

func testOne(clt *client.YTHostClient, task *message.TestMinerPerfTask, timeOut int64) (err error) {
	// 构造请求
	var requestMsg message.TestGetBlock
	if task.TestType == 0 {
		requestMsg.Msg = MSG_DOWNLOAD
	}

	requestbuf, err := proto.Marshal(&requestMsg)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeOut))
	defer cancel()
	// 发送消息
	resbuf, err := clt.SendMsg(ctx, message.MsgIDTestGetBlock.Value(), requestbuf)
	if err != nil {
		return
	}

	// 解析回复消息
	var resMsg message.TestGetBlockRes
	err = proto.Unmarshal(resbuf, &resMsg)
	if err != nil {
		return
	}
	return nil
}

func GetBlock(data []byte) (res []byte, err error) {

	var msg message.TestGetBlock
	err = proto.Unmarshal(data, &msg)
	if err != nil {
		log.Println("[perf]", err)
		return
	}
	var resMsg message.TestGetBlockRes
	switch msg.Msg {
	case MSG_DOWNLOAD:
		resMsg.Msg = make([]byte, config.Global_Shard_Size * 1024)
		rand.Read(resMsg.Msg)
	case MSG_CHECKOUT:
		statistics.DefaultStat.TXTest.AddSuccess()
	case MSG_UPLOAD:
		tk := &TokenPool.Token{}
		err := tk.FillFromString(msg.AllocID)
		if err == nil && TokenPool.Utp().Check(tk) {

			statistics.DefaultStat.RXTest.AddSuccess()
		} else {
			log.Println("[perf]check token", msg.AllocID, err)
		}

	case MSG_DOWNLOAD_TK:
	case MSG_UPLOAD_TK:
	}
	res, err = proto.Marshal(&resMsg)
	return
}
