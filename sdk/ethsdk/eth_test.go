package ethsdk

import (
	"context"
	"testing"
)

func TestNewethsdk(t *testing.T) {
	url := "http://0.0.0.0:37399"
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
