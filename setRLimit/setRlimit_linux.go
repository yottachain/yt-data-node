package setRLimit

import (
	log "github.com/yottachain/YTDataNode/logger"
	"syscall"
)

func SetRLimit() {
	//var rLimit syscall.Rlimit
	//
	//err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	//
	//if err != nil {
	//	fmt.Println("Error Getting Rlimit ", err)
	//}
	//fmt.Println(rLimit)
	//
	//rLimit.Max = 60000
	//rLimit.Cur = 60000
	//
	//err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	//if err != nil {
	//	fmt.Println("Error Setting Rlimit ", err)
	//}
	//err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	//
	//if err != nil {
	//	fmt.Println("Error Getting Rlimit ", err)
	//}
	//fmt.Println("Rlimit Final", rLimit)

	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return
	}
	rLimit.Max = 999999
	rLimit.Cur = 999999
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return
	}
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return
	}
	log.Printf("[SetRLimit] Ulimit -a, return %d\n", rLimit.Cur)
}
