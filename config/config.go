package config

import (
    "bytes"
    "crypto/md5"
    "encoding/json"
    "fmt"

    "io/ioutil"
    "log"
    "os"
    "path"
    "runtime"
    "strings"
    "time"

    "github.com/eoscanada/eos-go/btcsuite/btcutil/base58"
    "github.com/libp2p/go-libp2p-core/peer"
    "github.com/multiformats/go-multiaddr"
    "github.com/spf13/viper"

    ci "github.com/libp2p/go-libp2p-core/crypto"
    "github.com/yottachain/YTDataNode/util"
    ytfsOpts "github.com/yottachain/YTFS/opt"
)

type peerInfo struct {
	ID    string   `json:"ID"`
	Addrs []string `json:"Addrs"`
}

var isDebug = false

func init() {
	filename := path.Join(util.GetYTFSPath(), "debug.yaml")
	ok, err := util.PathExists(filename)
	if err == nil && ok {
		viper.SetConfigType("yaml")
		fl, err := os.OpenFile(filename, os.O_RDONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer fl.Close()
		viper.ReadConfig(fl)
		isDebug = true
		fmt.Println("---------DEBUG MODE------")
	}
}

// Config 配置
type Config struct {
	ID            string `json:"ID"`
	privKey       ci.PrivKey
	PubKey        string
	BPList        []peerInfo    `json:"BPList"`
	Adminacc      string        `json:"Adminacc"`
	Relay         bool          `json:"Relay"`
	ListenAddr    string        `json:"ListenAddr"`
	APIListen     string        `json:"APIListen"`
	IndexID       uint32        `json:"IndexID"`
	PoolID        string        `json:"PoolID"`
	MaxConn       int           `json:"MaxConn"`
	TokenInterval time.Duration `json:"TokenInterval"`
	*ytfsOpts.Options
	UpdateURL string `json:"update_url"`
	bpListMd5    []byte
	DisableWrite bool
	AllocSpace   uint64
}

// DefaultYTFSOptions default config
func DefaultYTFSOptions() *ytfsOpts.Options {
	yp := util.GetYTFSPath()
	opts := ytfsOpts.DefaultOptions()
	for index, storage := range opts.Storages {
		storage.StorageName = fmt.Sprintf("%s/storage-%d", yp, index)
		storage.StorageVolume = 2 << 40
		storage.DataBlockSize = 1 << 14
		opts.Storages[index] = storage
	}
	opts.DataBlockSize = 1 << 14
	opts.TotalVolumn = 2 << 41
	opts.IndexTableCols = 1 << 14
	opts.IndexTableRows = 1 << 28
	if runtime.GOOS == "linux" {
		opts.UseKvDb = true
	}
	return opts
}

func InitRowsCols(size uint64, n uint32, db string)(uint64, uint64, error){
	var d uint32 = 1 << 14
	var m uint64

	//expendRatioM = 1.2 = 12 / 10
	m = size / uint64(d) / uint64(n)
	if db == "rocksdb"{
		return m, uint64(n), nil
	}

	fmt.Println("InitRowsCols,M=",m)
	// IndexTableCols(M) = m * expendRatioM,  expendRatioM = 1.2  (512 <= IndexTableCols(M) <= 2048)
	if m < 420 || m > 1700{
		err := fmt.Errorf("IndexTableCols not suitable,M=",m)
		fmt.Println("[error]init failed, IndexTableCols not suitable,M=",m)
		return 0, 0, err
	}

	return m, uint64(n), nil
}

// GetYTFSOptionsByParams 通过参数生成YTFS配置
func GetYTFSOptionsByParams(size uint64, m uint32) *ytfsOpts.Options {
    yp := util.GetYTFSPath()
    var d uint32 = 16384
    n := size / uint64(d) / uint64(m)

    for {
        if (n & (n - 1)) == 0 {
            break
        }
        n = n + 1
    }
	m = uint32(size / uint64(d) / n)

    opts := &ytfsOpts.Options{
        YTFSTag: "ytfs",
        Storages: []ytfsOpts.StorageOptions{
            {
                StorageName:   path.Join(yp, "storage"),
                StorageType:   0,
                ReadOnly:      false,
                SyncPeriod:    1,
                StorageVolume: size,
                DataBlockSize: 16384,
            },
        },
        ReadOnly:       false,
        SyncPeriod:     1,
        IndexTableCols: m,
        IndexTableRows: uint32(n),
        DataBlockSize:  d,
        TotalVolumn:    size,
        UseKvDb:        true,
    }

    if runtime.GOOS == "windows" {
        opts.UseKvDb = false
    }

    return opts
}

// GetYTFSOptionsByParams2 通过参数生成YTFS配置, 多storage配置
func GetYTFSOptionsByParams2(totalSize uint64, storageSize uint64, m uint32) *ytfsOpts.Options {
	yp := util.GetYTFSPath()
	n := totalSize / uint64(m)
	opts := &ytfsOpts.Options{
		YTFSTag: "ytfs",
		Storages: []ytfsOpts.StorageOptions{
			{
				StorageName:   path.Join(yp, "storage"),
				StorageType:   0,
				ReadOnly:      false,
				SyncPeriod:    1,
				StorageVolume: storageSize,
				DataBlockSize: 1 << 14,
			},
		},
		ReadOnly:       false,
		SyncPeriod:     1,
		IndexTableCols: m,
		IndexTableRows: uint32(n),
		DataBlockSize:  1 << 14,
		TotalVolumn:    totalSize,
	}
	return opts
}

// NewConfig ..
func NewConfig() *Config {
	cfg := NewConfigByYTFSOptions(DefaultYTFSOptions())
	return cfg
}

func NewConfigByYTFSOptions(opts *ytfsOpts.Options) *Config {
	if opts == nil {
		fmt.Println("[error] new config opts is nil")
		return nil
	}
	cfg := new(Config)
	cfg.ListenAddr = "/ip4/0.0.0.0/tcp/9001"
	cfg.APIListen = ":9002"
	cfg.Options = opts
	cfg.AllocSpace = opts.TotalVolumn
	var err error
	cfg.privKey, cfg.PubKey, err = util.RandomIdentity2()
	if err != nil {
		fmt.Printf("config create privkey err:%s\n", err.Error())
		return nil
	}

	cfg.Relay = true
	cfg.BPList = getBPList()
	log.Println("系统", runtime.GOOS)
	if runtime.GOOS == "windows" {
		cfg.UseKvDb = false
	} else {
		cfg.UseKvDb = true
	}
	return cfg
}

func getBPList() []peerInfo {
	var bplist []peerInfo

	if isDebug {
		snmaddr := viper.GetStringSlice("snmaddrs")
		for _, v := range snmaddr {
			ma, err := multiaddr.NewMultiaddr(v)
			if err != nil {
				fmt.Println(err)
				continue
			}
			ai, err := peer.AddrInfoFromP2pAddr(ma)
			if err != nil {
				fmt.Println(err)
				continue
			}
			var addrs = []string{}
			for _, v := range ai.Addrs {
				addrs = append(addrs, v.String())
			}
			bplist = append(bplist, peerInfo{ai.ID.String(), addrs})
		}
		return bplist
	}

//	jsdata := `
//    [
//    {
//      "ID": "16Uiu2HAm7o24DSgWTrcu5sLCgSkf3D3DQqzpMz9W1Bi7F2Cc4SF6",
//      "Addrs": ["/dns4/sn00.yottachain.net/tcp/9999"]
//    },
//    {
//      "ID": "16Uiu2HAmNe1bZF2s7msxqy9tFT7WDfUaJa98h1KBhAmTTHvcZqpA",
//      "Addrs": ["/dns4/sn01.yottachain.net/tcp/9999"]
//    },
//    {
//      "ID": "16Uiu2HAkyZAuzcjmpFhk1pCLAZaYusV3wXmrEhnnNDfeJjkVoQc6",
//      "Addrs": ["/dns4/sn02.yottachain.net/tcp/9999"]
//    },
//    {
//      "ID": "16Uiu2HAmSFq7SbwcfYVn3NzWuuV7SizQEVjKEwty1knZuzTA7jDq",
//      "Addrs": ["/dns4/sn03.yottachain.net/tcp/9999"]
//    },
//    {
//      "ID": "16Uiu2HAmC7DSN4kNi64sB5N9aMgv9DjTTrtydf4YKS3Q56hYsDNS",
//      "Addrs": ["/dns4/sn04.yottachain.net/tcp/9999"]
//    },
//    {
//      "ID": "16Uiu2HAmSFgs5Pj6hFdAzCAvFGH78ew7egakT6VqL1xaLdvxnnSc",
//      "Addrs": ["/dns4/sn05.yottachain.net/tcp/9999"]
//    },
//    {
//      "ID": "16Uiu2HAmP2RuNAkXdtQDiFqVuBA8yERh91JV6b29rQpAGKkb3PiM",
//      "Addrs": ["/dns4/sn06.yottachain.net/tcp/9999"]
//    },
//    {
//      "ID": "16Uiu2HAmJQ7cjPzi7u4NdgYK7xWqDgEVBqAqNPrmxY2KVKGpND2W",
//      "Addrs": ["/dns4/sn07.yottachain.net/tcp/9999"]
//    },
//    {
//      "ID": "16Uiu2HAmDKteRvgXPtzz3pvhGn56HH7uo8WqoGqJWPArY4G1kuWP",
//      "Addrs": ["/dns4/sn08.yottachain.net/tcp/9999"]
//    },
//    {
//      "ID": "16Uiu2HAmDpQ6527dqtiv5fixTptQBtGa561BZeUTDuALiAZwQNGR",
//      "Addrs": ["/dns4/sn09.yottachain.net/tcp/9999"]
//    },
//    {
//      "ID": "16Uiu2HAmATZsCop9hkKDbmtyLbizLQU92jrCVpWvzRChKRQbwzy7",
//      "Addrs": ["/dns4/sn10.yottachain.net/tcp/9999"]
//    },
//    {
//      "ID": "16Uiu2HAmGheSFhwbpihhEnyZUxsVr6Rn9z5v2XDMeEyAfK2K4nwG",
//      "Addrs": ["/dns4/sn11.yottachain.net/tcp/9999"]
//    },
//    {
//      "ID": "16Uiu2HAm2JANbeeaXDa9JaDTU5Q1h2hmjJGJx91LpYd36pdoDWdx",
//      "Addrs": ["/dns4/sn12.yottachain.net/tcp/9999"]
//    },
//    {
//      "ID": "16Uiu2HAm3Cmzqg9TKR6FvEH5NSgzLZgDZb4xtPC9aYhqbc9p7WM5",
//      "Addrs": ["/dns4/sn13.yottachain.net/tcp/9999"]
//    },
//    {
//      "ID": "16Uiu2HAm2CzizQh2AU8NXK5z2bvJUaFuPiM9Z6R1uDEFKDvob4mJ",
//      "Addrs": ["/dns4/sn14.yottachain.net/tcp/9999"]
//    },
//    {
//      "ID": "16Uiu2HAmTd1jqEGLThwcrD9yYG1JsHHj7qsDJDBcdgMLMvaBnksU",
//      "Addrs": ["/dns4/sn15.yottachain.net/tcp/9999"]
//    },
//    {
//      "ID": "16Uiu2HAmHufUv4udcL1f1bNP4r6VqDBppmKH495iQKSgv6nWGoZA",
//      "Addrs": ["/dns4/sn16.yottachain.net/tcp/9999"]
//    },
//    {
//      "ID": "16Uiu2HAmKS35S4JQk8BDUvgWhjGLMJ1f9zWJhT3QeRRyFdReXeue",
//      "Addrs": ["/dns4/sn17.yottachain.net/tcp/9999"]
//    },
//    {
//      "ID": "16Uiu2HAmUyPbR4wcKtGi6n84CGkHsXsHZZ2sGrnhJPAqJmFCMfDW",
//      "Addrs": ["/dns4/sn18.yottachain.net/tcp/9999"]
//    },
//    {
//      "ID": "16Uiu2HAmLUCp92e25HXiZW8fMwpCUfhQRcNGL7PibTDtg51JTRCq",
//      "Addrs": ["/dns4/sn19.yottachain.net/tcp/9999"]
//    },
//    {
//      "ID": "16Uiu2HAmBG1d8HHBApLg9MrDqgUX4LoKcFCSCrq54QW3mkqRheo1",
//      "Addrs": ["/dns4/sn20.yottachain.net/tcp/9999"]
//    }
//  ]
//`

	jsdata := `
    [
    {
      "ID": "16Uiu2HAm7o24DSgWTrcu5sLCgSkf3D3DQqzpMz9W1Bi7F2Cc4SF6",
      "Addrs": ["/dns4/sn.yottachain.net/tcp/9999"]
    }
	]
`
	buf := bytes.NewBufferString(jsdata)
	_ = json.Unmarshal(buf.Bytes(), &bplist)
	return bplist
}

// Save ..
func (cfg *Config) Save() error {
	yp := util.GetYTFSPath()
	if ok, err := util.PathExists(yp); ok != true || err != nil {
		_ = os.Mkdir(yp, os.ModePerm)
	}

	cfgPath := util.GetConfigPath()

	keyBytes, err := cfg.privKey.Raw()
	if err != nil {
		log.Println("配置保存失败", err)
	}

	_ = ioutil.WriteFile(fmt.Sprintf("%s/swarm.key", yp), keyBytes, os.ModePerm)
	peerID, err := util.IdFromPublicKey(cfg.PubKey)
	if err != nil {
		return err
	}
	cfg.ID = peerID.Pretty()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cfgPath, data, os.ModePerm)
}

func (cfg *Config) ReloadBPList() {
	cfg.BPList = getBPList()
	cfg.Save()
}

// NewKey 创建新的key
func (cfg *Config) NewKey() error {
	cfg.privKey, _ = util.RandomIdentity()
	id, err := util.IdFromPublicKey(cfg.PubKey)
	if err != nil {
		return err
	}
	cfg.ID = id.Pretty()
	return nil
}

// GetBPIndex 返回bpindex
func (cfg *Config) GetBPIndex() int {
	id := cfg.IndexID
	bpnum := len(cfg.BPList)
	if bpnum == 0 {
		return 0
	}
	bpindex := id % uint32(bpnum)

	//log.Printf("len bplist:%d ,id %d, bpindex %d\n", bpnum, id, bpindex)
	return int(bpindex)
}

// ReadConfig 读配置
func ReadConfig() (*Config, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("读取配置失败部分功能不可用")
		}
	}()
	var cfg Config
	data, err := ioutil.ReadFile(util.GetConfigPath())
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("配置文件格式错误:%s", err)
	}
	keyBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/swarm.key", util.GetYTFSPath()))
	if err != nil {
		return nil, fmt.Errorf("read swarm.key error:%s", err)
	}
	privk, err := ci.UnmarshalSecp256k1PrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("private key error:%s", err)
	}

	cfg.privKey = privk

	return &cfg, nil
}

// PrivKey ..
func (cfg *Config) PrivKey() ci.PrivKey {
	return cfg.privKey
}
func (cfg *Config) PrivKeyString() string {
	buf, _ := cfg.privKey.Raw()
	return base58.Encode(buf)
}

func (cfg *Config) Version() uint32 {
	return 285
}

func Version() uint32 {
	return new(Config).Version()
}

func (cfg Config) ResetYTFSOptions(opts *ytfsOpts.Options) Config {
	cfg.Options = opts
	return cfg
}

func (cfg Config) GetAPIAddr() string {
	bpIndex := cfg.GetBPIndex()
	//ma, err := multiaddr.NewMultiaddr(cfg.BPList[bpIndex].Addrs[0])
	//if err != nil {
	//	return ""
	//}
	addrs := strings.Split(cfg.BPList[bpIndex].Addrs[0], "/")
	//addr, err := manet.ToNetAddr(ma)
	//if err != nil {
	//	return ""
	//}

	return fmt.Sprintf("http://%s:%s", addrs[2], "8082")
}

func (cfg *Config) BPMd5() []byte {
	if cfg.bpListMd5 == nil {
		m5 := md5.New()
		for _, v := range cfg.BPList {
			m5.Write([]byte(v.ID))
		}
		cfg.bpListMd5 = m5.Sum(nil)
	}
	return cfg.bpListMd5
}

var DefaultConfig, _ = ReadConfig()
