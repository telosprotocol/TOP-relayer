package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	CONFIG      *Config
	WALLET_PATH string
	CONFIG_PATH string
)

type Config struct {
	Chains []*ChainsConfig `josn:"chains"`
}

func newConfig(path string) (*Config, error) {
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

type ChainsConfig struct {
	//listen config
	ListenChainId uint64 `json:"listenchainid"`
	ListenUrl     string `json:"listenurl"`

	//submit config
	SubmitChainId  uint64 `json:"submitchainid"`
	SubmitUrl      string `json:"submiturl"`
	Contract       string `json:"contract"`
	KeyPath        string `json:"keypath"`
	BlockCertainty int    `json:"blockcertainty"`
	SubBatch       int    `json:"subBatch"`
	VerifyBlock    bool   `json:"verifyblock"`
	Start          bool   `json:"start"`
}

type HeaderSyncConfig struct {
	Config *Config
}

func InitHeaderSyncConfig(path string) (*HeaderSyncConfig, error) {
	cfg, err := newConfig(path)
	if err != nil {
		return nil, err
	}
	return &HeaderSyncConfig{Config: cfg}, nil
}
