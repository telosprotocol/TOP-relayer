package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/wonderivan/logger"
)

const (
	TOP_CHAIN     string = "TOP"
	ETH_CHAIN     string = "ETH"
	BSC_CHAIN     string = "BSC"
	HECO_CHAIN    string = "HECO"
	OPEN_ALLIANCE string = "OPEN_ALLIANCE"

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
	//submit config
	Url      []string `json:"url"`
	Contract string   `json:"contract"`
	KeyPath  string   `json:"keypath"`
}

type Config struct {
	RelayerConfig map[string]*Relayer `json:"relayer_config"`
	RelayerToRun  string              `json:"relayer_to_run"`
}

func LoadRelayerConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Read config password file failed:", err)
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatal("Unmarshal config file failed:", err)
		return nil, err
	}
	return config, nil
}

func InitLogConfig() error {
	os.Mkdir(LOG_DIR, os.ModePerm)
	logger.SetLogger(LOG_CONFIG)
	return nil
}
