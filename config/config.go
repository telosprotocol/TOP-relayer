package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/wonderivan/logger"
)

const (
	TOPAddr = ""
)

const (
	ETHAddr     = ""
	ETHPrysm    = ""
	ETHContract = ""
)

const (
	BSCAddr     = "https://bsc-dataseed4.binance.org"
	BSCContract = ""
)

const (
	HECOAddr     = "https://http-mainnet.hecochain.com"
	HECOContract = ""
)

const (
	OAAddr          = ""
	OAContract      = ""
	OAContractOnTop = ""
)

const (
	TOP_CHAIN  string = "TOP"
	ETH_CHAIN  string = "ETH"
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
	Url      []string
	Contract string
	KeyPath  string `json:"keypath"`
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
	if err = json.Unmarshal(data, config); err != nil {
		log.Fatal("Unmarshal config file failed:", err)
		return nil, err
	}
	err = fullConfigInfo(config)
	return config, err
}

func fullConfigInfo(cfg *Config) error {
	// reset to
	if len(cfg.RelayerConfig) < 2 {
		return fmt.Errorf("config information error")
	}
	v, ok := cfg.RelayerConfig[cfg.RelayerToRun]
	if !ok {
		return fmt.Errorf("config infomation of relayerToRun(%s) is not found", cfg.RelayerToRun)
	}
	switch cfg.RelayerToRun {
	case TOP_CHAIN:
		v.Url = make([]string, 1)
		v.Url[0] = TOPAddr
	case ETH_CHAIN:
		v.Url = make([]string, 1)
		v.Url[0] = ETHAddr
		v.Contract = ETHContract
	default:
		return fmt.Errorf("invalid RelayerToRun(%s)", cfg.RelayerToRun)
	}
	// 2. reset from
	for chinaName, v := range cfg.RelayerConfig {
		if chinaName == cfg.RelayerToRun {
			continue
		}
		switch chinaName {
		case TOP_CHAIN:
			v.Url = make([]string, 1)
			v.Url[0] = TOPAddr
		case ETH_CHAIN:
			v.Url = make([]string, 3)
			v.Url[0] = ETHAddr
			v.Url[1] = ETHPrysm
		default:
			return fmt.Errorf("invalid RelayerToRun(%s)", cfg.RelayerToRun)
		}
	}
	return nil
}

func InitLogConfig() error {
	os.Mkdir(LOG_DIR, os.ModePerm)
	logger.SetLogger(LOG_CONFIG)
	return nil
}
