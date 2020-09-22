package TaskPool

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	"testing"
	"time"
)

func TestUploadTaskPool_Check(t *testing.T) {
	go Utp().FillToken()
	go Dtp().FillToken()
	//go func() {
	//	for {
	//		time.Sleep(time.Second * 5)
	//		Utp().MakeTokenQueue()
	//	}
	//}()
	func() {
		for {
			go func() {
				ctx, _ := context.WithTimeout(context.Background(), time.Second)
				tk, err := Utp().Get(ctx, peer.ID("111"), 0)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println("000", tk.String())
				}
			}()
			time.Sleep(time.Millisecond * 1)
		}
	}()
	time.Sleep(time.Second)
	//for {
	//	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	//	_, err := Dtp().Get(ctx, peer.ID("222"), 1)
	//
	//	if err != nil {
	//		fmt.Println(err.Error())
	//	} else {
	//		fmt.Println("111")
	//	}
	//}
}
