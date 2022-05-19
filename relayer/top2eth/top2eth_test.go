package top2eth

import (
	"testing"
	"toprelayer/base"
	"toprelayer/msg"

	"github.com/ethereum/go-ethereum/common"
)

const SUBMITTERURL string = "http://192.168.30.32:8545"

var DEFAULTPATH = "../../.relayer/wallet/eth"

var CONTRACT common.Address = common.HexToAddress("0xf04e16dae6887e8cf235f0179b50e0cd1860647c")

func TestSubmitHeader(t *testing.T) {
	sub := new(Top2EthRelayer)
	err := sub.Init(SUBMITTERURL, SUBMITTERURL, DEFAULTPATH, "", base.ETH, CONTRACT, 10, 0, false)
	if err != nil {
		t.Fatal(err)
	}

	var headers []*msg.TopElectBlockHeader
	for i := 1; i <= 100; i++ {
		headers = append(headers, &msg.TopElectBlockHeader{BlockNumber: uint64(i)})
	}
	data, _ := msg.EncodeHeader(headers)
	nonce, _ := sub.wallet.GetNonce(sub.wallet.CurrentAccount().Address)

	t.Logf("account:%v,nonce:%v\n", sub.wallet.CurrentAccount().Address, nonce)
	tx, err := sub.submitTopHeader(data, nonce)
	if err != nil {
		t.Fatal("SubmitHeader error:", err)
	}
	t.Log("SubmitHeader hash:", tx.Hash())
}
