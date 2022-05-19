package topsdk

import "testing"

func newtopsdk() (*TopSdk, error) {
	url := "http://0.0.0.0:37389"
	return NewTopSdk(url)
}

func TestSaveEthBlockHead(t *testing.T) {
	sdk, err := newtopsdk()
	if err != nil {
		t.Fatalf("NewSDK failed,error:%v", err)
	}

	rawTx := "0xf86e0385174876e80082520894a6d2b331b03fddb8c6a8830a63fe47e42c4bdf4e881bc16d674ec8000080820a93a0e17d1040b31459b966e9fc9815761e8c940b21ec7910ba6796262cada6dd7ffca0704c3b023dc136ae5980e2218078d0ad3df0778a356036209fd96f9f66016a7c"
	hash, err := sdk.SaveEthBlockHead(rawTx)
	if err != nil {
		t.Fatalf("SaveEthBlockHead failed,error:%v", err)
	}
	t.Logf("SaveEthBlockHead ok,hash:%v", hash)
}

func TestGetLatestEthBlockHead(t *testing.T) {
	sdk, err := newtopsdk()
	if err != nil {
		t.Fatalf("NewSDK failed,error:%v", err)
	}

	ehd, err := sdk.GetLatestEthBlockHead()
	if err != nil {
		t.Fatalf("GetLatestEthBlockHead failed,error:%v", err)
	}
	t.Logf("GetLatestEthBlockHead ok,eth block:%v", ehd.ReceivedFrom)
}
