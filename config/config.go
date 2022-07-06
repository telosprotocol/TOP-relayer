package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/wonderivan/logger"
)

var (
	TOP_CHAIN string = "TOP"
	ETH_CHAIN string = "ETH"

	LOG_DIR    string = "log"
	LOG_CONFIG string = `{
		"TimeFormat":"2006-01-02 15:04:05",
		"Console": {
		  "level": "DEBG",
		  "color": true
		},
		"File": {
		  "filename": "./log/relayer.log",
		  "level": "DEBG",
		  "daily": true,
		  "maxlines": 1000000,
		  "maxsize": 1,
		  "maxdays": -1,
		  "append": true,
		  "permit": "0660"
		}
	}`
)

type Relayer struct {
	// chain symbol
	Chain   string `json:"chainName"`
	ChainId uint64 `json:"chainId"`

	//listen config
	ListenUrl string `json:"listenurl"`

	//submit config
	SubmitUrl string `json:"submiturl"`
	Contract  string `json:"contract"`
	KeyPath   string `json:"keypath"`
	Start     bool   `json:"start"`
}

type Config struct {
	RelayerConfig []*Relayer `json:"relayerconfig"`
}

func NewConfig(path string) (*Config, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		return nil, fmt.Errorf("[%v] not a json file,want a json config file.", path)
	}

	jsonFile, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()

	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("Read config file error %v", err)
	}

	config := &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("Parse config file error %v", err)
	}
	return config, nil
}

func InitLogConfig() error {
	os.Mkdir(LOG_DIR, os.ModePerm)
	logger.SetLogger(LOG_CONFIG)
	return nil
}
