package ethsdk

import (
	"context"
	"testing"
)

func TestNewethsdk(t *testing.T) {
	url := "http://192.168.30.200:19086"
	sdk, err := NewEthSdk(url)
	if err != nil {
		t.Fatal(err)
	}
	id, err := sdk.ChainID(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log("new eth sdk success.chain id:", id.Uint64())
}
