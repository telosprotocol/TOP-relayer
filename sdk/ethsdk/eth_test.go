package ethsdk

import (
	"testing"
	"toprelayer/sdk"
)

func newethsdk() (*EthSdk, error) {
	url := "http://0.0.0.0:37389"
	sdk, err := sdk.NewSDK(url)
	if err != nil {
		return nil, err
	}
	return NewEthSdk(url, sdk)
}

func TestGetTopElectBlockHeadByHeight(t *testing.T) {
	sdk, err := newethsdk()
	if err != nil {
		t.Fatal("new eth sdk error:", err)
	}

	b, err := sdk.GetTopElectBlockHeadByHeight(1, 1)
	if err != nil {
		t.Fatalf("GetTopElectBlockHeadByHeight failed,error:%v", err)
	}
	t.Logf("GetTopElectBlockHeadByHeight ok,top block hash:%v", b.Hash)
}

func TestGetLatestTopElectBlockHeight(t *testing.T) {
	sdk, err := newethsdk()
	if err != nil {
		t.Fatal("new eth sdk error:", err)
	}

	h, err := sdk.GetLatestTopElectBlockHeight()
	if err != nil {
		t.Fatalf("GetLatestTopElectBlockHeight failed,error:%v", err)
	}
	t.Logf("GetLatestTopElectBlockHeight ok,height:%v", h)
}
