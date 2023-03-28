package config

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestFullConfigInfo(t *testing.T) {
	data := `{
    "relayer_config": {
        "TOP": {
            "keypath": ".relayer/wallet/top"
        },
        "ETH": {
            "keypath": ".relayer/wallet/eth"
        }
    },
    "relayer_to_run": "TOP"
}`

	config := &Config{}
	if err := json.Unmarshal([]byte(data), config); err != nil {
		t.Fatal(err.Error())
	}
	if err := fullConfigInfo(config); err != nil {
		t.Fatal(err.Error())
	}
	for k, v := range config.RelayerConfig {
		fmt.Printf("key:%s, value:%+v \n", k, *v)
	}

	if config.RelayerToRun != "TOP" {
		t.Fail()
	}
	if relayer, ok := config.RelayerConfig["TOP"]; ok {
		if len(relayer.Url) != 1 {
			t.Fail()
		}
		if relayer.KeyPath != `.relayer/wallet/top` {
			t.Fail()
		}
		if len(relayer.Contract) != 0 {
			t.Fail()
		}
	} else {
		t.Fatal("top cant found")
	}
	if relayer, ok := config.RelayerConfig["ETH"]; ok {
		if len(relayer.Url) != 3 {
			t.Fail()
		}
		if relayer.KeyPath != `.relayer/wallet/eth` {
			t.Fail()
		}
		if len(relayer.Contract) != 0 {
			t.Fail()
		}
	} else {
		t.Fatal("eth cant found")
	}
}
