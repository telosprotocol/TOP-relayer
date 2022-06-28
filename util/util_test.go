package util

import (
	"fmt"
	"os"
	"testing"
	"toprelayer/config"
)

func TestGetchainpass(t *testing.T) {
	configPath := "../config/relayerconfig.json"
	handlercfg, err := config.InitHeaderSyncConfig(configPath)
	if err != nil {
		t.Fatal("init config error: ", err)
		return
	}
	chainpass, err := Getchainpass(handlercfg)
	if err != nil {
		t.Fatal("get chainpass error: ", err)
		return
	}
	for chainid := range chainpass {
		fmt.Println("id: ", chainid, ", pass: ", chainpass[chainid])
	}
}
