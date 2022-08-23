package crosschainrelayer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"toprelayer/config"

	"github.com/wonderivan/logger"
)

var defaultPass = "asd123"
var defaultKeyPath = "../../.relayer/wallet/eth"

func TestVerifySeverConnection(t *testing.T) {
	url := "http://192.168.20.10:8080/v1/aggragation/bridge/manage/checkRelayBlockHashs"
	var list []string
	list = append(list, "0xa7d87ac5e42f67ada5f170e8e4b7de053c3dad580c6c2d33cbe841f9d1462c44")
	list = append(list, "0x138a2ede7a7f49d138f8c5c65678281df5f7d82aa04c8b503289191fbf876328")

	data := make(map[string]interface{})
	data["relayBlockHashs"] = list
	bytesData, _ := json.Marshal(data)
	resp, err := http.Post(url, "application/json", bytes.NewReader(bytesData))
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}

func TestSeverVerify(t *testing.T) {
	listenUrl := "http://192.168.50.31:8080"
	ethUrl := "https://ropsten.infura.io/v3/fb2a09e82a234971ad84203e6f75990e"
	ethContract := "0x0E62B4A03D615ae0B088345FD26eB8782F2a861F"

	cfg := &config.Relayer{
		Url:      ethUrl,
		KeyPath:  defaultKeyPath,
		Contract: ethContract,
	}
	relayer := &CrossChainRelayer{}
	err := relayer.Init(config.ETH_CHAIN, cfg, listenUrl, defaultPass, config.Server{})
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = relayer.queryBlocks(0x12, 0x49)
	if err != nil {
		t.Fatal(err)
	}
	element := relayer.verifyList.Front()
	if element == nil {
		logger.Error("txList get front nil")
		return
	}
	info, ok := element.Value.(VerifyInfo)
	if !ok {
		logger.Error("txList get front error")
		return
	}
	if !relayer.serverVerify(info.VerifyList) {
		logger.Info("%v verify not pass", info.Block.Hash)
		return
	}
}
