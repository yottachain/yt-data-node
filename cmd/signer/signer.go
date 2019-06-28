package main

import (
	"encoding/json"
	"flag"
	"log"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"os"
	//"github.com/rocket049/gocui"
)

var key string
var keys []string
var tx string
var baseNodeUrl = "http://152.136.11.202:8888"
var api = eos.New(baseNodeUrl)
var kb = eos.NewKeyBag()

func main() {
	var signedTx eos.SignedTransaction
	flag.StringVar(&key, "k", "", "签名私钥")
	flag.StringVar(&tx, "t", "", "签名交易")
	flag.Parse()

	if key == "" && len(os.Args) > 1 {
		keys = os.Args[1:]
		for _, v := range keys {
			kb.ImportPrivateKey(v)
		}
	}
	kb.ImportPrivateKey(key)
	if tx == "" {
		log.Println("请输入待签名交易：")
		fmt.Scanf("%s\n", &tx)
	}
	err := json.Unmarshal([]byte(tx), &signedTx)
	if err != nil {
		log.Println("签名失败:", err)
	}
	txopts := &eos.TxOptions{}
	txopts.FillFromChain(api)

	res, err := kb.Sign(&signedTx, txopts.ChainID, getPubkey()...)
	if err != nil {
		log.Println("签名失败:", err)
	}
	log.Println("交易签名：")
	buf, _ := json.Marshal(res)
	log.Println("-----------签名结果-----------")
	log.Println(string(buf))
	log.Println("-----------------------------")
}

func getPubkey() []ecc.PublicKey {
	var pkeys = make([]ecc.PublicKey, len(kb.Keys))
	for k, v := range kb.Keys {
		pkeys[k] = v.PublicKey()
	}
	return pkeys
}
