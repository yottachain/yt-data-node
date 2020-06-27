package selfcheck

import(
	"fmt"
	"github.com/yottachain/YTDataNode/util"
	"github.com/yottachain/YTFS/storage"
	"strconv"
	"github.com/yottachain/YTDataNode/config"
	"github.com/yottachain/YTDataNode/message"
	"time"
	"path"
	"io/ioutil"
	sni "github.com/yottachain/YTDataNode/storageNodeInterface"
	"github.com/mr-tron/base58/base58"
)

var resp message.SelfVarifyResp
var ti *storage.TableIterator
var nTabVarifyedFile string = "/gc/tab_index"

type Scker struct {
	sni.StorageNode
}

func init(){
	StartTime := time.Now()
	fmt.Println("StartTime:",StartTime)
	cfg,_ := config.ReadConfig()
	resp.Id = strconv.FormatUint(uint64(cfg.IndexID),10)
	pathname := path.Join(util.GetYTFSPath(),"index.db")
	ti,_ = storage.GetTableIterator(pathname,cfg.Options)
	pathTabIdxfile := path.Join(util.GetYTFSPath(),nTabVarifyedFile)
	strVal,_ := GetValueFromFile(pathTabIdxfile)
    val,_ := strconv.ParseUint(strVal,10,64)
    SetValuetoTableIter(uint32(val),ti)
}

func SetValuetoTableIter(value uint32,ti *storage.TableIterator){
	ti.SetTableIdx(value)
}

func GetValueFromFile(filePath string) (string ,error){
	status_exist,_ := util.PathExists(filePath)
	if status_exist == false {
		return strconv.FormatUint(uint64(0),10),nil
	}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("read file err=%v\r\n", err)
	}
	return string(content),err
}

func SetValuetoFile(filePath string, value string) error{
	err := ioutil.WriteFile(filePath,[]byte(value),0666)
	if err != nil{
		fmt.Println(err)
	}
	return err
}

func (Sck *Scker)SelfCheck() message.SelfVarifyResp {
	errNum := 0
	varifyedNum := 0
	beginTab := ti.GetBeginTab()
	nowTab := ti.GetBeginTab()
	pathTabIdxfile := path.Join(util.GetYTFSPath(),nTabVarifyedFile)

	for{
		tab,_ := ti.GetNoNilTableBytes()
		for key,_ := range tab {
			varifyedNum++
			resData, err := Sck.YTFS().Get(key)
			if err != nil{
				fmt.Println("[selfCheck] err:",err)
				errNum++
				continue
			}

			if ! message.VerifyVHF(resData,key[:]){
				fmt.Println("[selfCheck] err:",err," key:",base58.Encode(key[:]))
				errNum++
			}
		}

		nowTab = ti.GetBeginTab()
		if varifyedNum > 200000{
			SetValuetoFile(strconv.FormatUint(uint64(nowTab),10),pathTabIdxfile)
			break
		}

		if nowTab - beginTab >= 10 {
			SetValuetoFile(strconv.FormatUint(uint64(nowTab),10),pathTabIdxfile)
			break
		}
	}

	resp.Numth = strconv.FormatUint(uint64(nowTab),10)
	resp.ErrNum = strconv.FormatUint(uint64(varifyedNum),10)
	return resp
}
