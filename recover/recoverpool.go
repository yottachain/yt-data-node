package recover

import (
	"github.com/yottachain/YTDataNode/config"
	log "github.com/yottachain/YTDataNode/logger"
	"time"
	//"sync"
)

var poolG chan int
var totalCap int = 2000
var realConCurrent uint16 = 1     //can be changed by write-weight and config

//type Request struct {
//	 Tsk     *Task
//	 //Response chan []string    //用于存放请求结果的通道
//}
//
//var requestChannelG chan Request
//
//func (re *RecoverEngine) sendRequest() {
//	task := <- re.queue
//	requestChannelG <- Request{Tsk:task}
//	//responseChanT := make(chan []string, 1)
//	//defer close(responseChanT)
//	//requestChannelG <- Request{Tsk:task,Response:responseChanT}
//}

func (re *RecoverEngine) doRequest(task *Task){
    re.IncConTask()
    defer re.DecConTask()
    re.processTask(task)
    poolG <- 0
}

func (re *RecoverEngine)processRequests(){
	for {
		if len(poolG) > 0 {
			<- poolG
			requestT :=<- re.queue
            re.IncRbdTask()
            log.Println("[recover] create_gorutine, len_poolG=",len(poolG))
			go re.doRequest(requestT)
		} else {
			log.Println("[recover] create_gorutine pool is full, len_poolG=",len(poolG))
			//requestT.Response <- []string{"goroutine pool is full"}
			<- time.After(time.Second * 3)
		}
	}
}

func (re *RecoverEngine)modifyPoolSize(){
	utp := re.Upt
	//configweight := re.sn.Config().ShardRbdConcurrent
	configweight := config.Gconfig.ShardRbdConcurrent
    for{
    	<-time.After(time.Second * 600)
		tokenweight := time.Second/utp.FillTokenInterval
        realConCurrent_N := configweight
        if uint16(tokenweight) < realConCurrent_N {
           realConCurrent_N = uint16(tokenweight)
		}

        if realConCurrent_N > 2000 {
        	realConCurrent_N = 2000
		}

		if realConCurrent_N == 0 {
			realConCurrent_N = 10
		}

        if realConCurrent < realConCurrent_N {
        	for k := uint16(0); k < realConCurrent_N - realConCurrent; k++{
        		poolG <- 0
			}
			realConCurrent = realConCurrent_N
		}

		if realConCurrent > realConCurrent_N {
			for k := uint16(0); k < realConCurrent - realConCurrent_N; k++{
				<- poolG
			}
			realConCurrent = realConCurrent_N
		}
		log.Println("[recover] realConCurent=",realConCurrent)
	}
}

func (re *RecoverEngine)RunPool(){
	poolG = make(chan int, totalCap)
	defer close(poolG)
	//requestChannelG = make(chan Request, 100)
	//defer close(requestChannelG)

	go re.processRequests()

	go re.modifyPoolSize()

	for i := uint16(0); i < realConCurrent; i++ {
		poolG <- 0
	}

	for {
		re.MultiReply()
	}
}

