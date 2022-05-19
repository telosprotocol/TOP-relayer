package eth2top

import (
	"math/big"
	"testing"
	"toprelayer/base"
	"toprelayer/msg"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const SUBMITTERURL string = "http://0.0.0.0:37399"

var DEFAULTPATH = "../../.relayer/wallet/top"
var CONTRACT common.Address = common.HexToAddress("0xa6D2b331B03fdDB8c6A8830A63fE47E42c4bDF4E")

func TestSubmitHeader(t *testing.T) {
	sub := &Eth2TopRelayer{}
	err := sub.Init(SUBMITTERURL, SUBMITTERURL, DEFAULTPATH, "", base.TOP, CONTRACT, 10, 0, false)
	if err != nil {
		t.Fatal(err)
	}

	var headers []*types.Header
	for i := 1; i <= 200; i++ {
		headers = append(headers, &types.Header{Number: big.NewInt(int64(i))})
	}

	data, err := msg.EncodeHeaders(&headers)
	if err != nil {
		t.Fatal("EncodeToBytes:", err)
	}

	if sub.wallet == nil {
		t.Fatal("nil wallet!!!")
	}
	tx, err := sub.submitEthHeader(data, 1)
	if err != nil {
		t.Fatal("SubmitHeader error:", err)
	}
	t.Log("SubmitHeader hash:", tx.Hash())
}
