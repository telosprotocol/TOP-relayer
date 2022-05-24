package top2eth

import (
	"sync"
	"testing"
	"toprelayer/base"

	"github.com/ethereum/go-ethereum/common"
)

const SUBMITTERURL string = "http://192.168.50.235:8545"
const LISTENURL string = "http://192.168.50.204:19086"

var DEFAULTPATH = "../../.relayer/wallet/eth"

var CONTRACT common.Address = common.HexToAddress("0xe8b713aee3e241831649a993f04c9f2026d99d55")

func TestSubmitHeader(t *testing.T) {
	sub := new(Top2EthRelayer)
	err := sub.Init(SUBMITTERURL, LISTENURL, DEFAULTPATH, "", base.ETH, CONTRACT, 100, 0, false)
	if err != nil {
		t.Fatal(err)
	}

	err = sub.StartRelayer(&sync.WaitGroup{})
	if err != nil {
		t.Fatal(err)
	}

	/* var headers []*types.Header //[]*msg.TopElectBlockHeader
	for i := 1; i <= 150; i++ {
		headers = append(headers, &types.Header{GasLimit: uint64(i)})
	}
	data, _ := msg.EncodeHeader(headers)
	nonce, err := sub.wallet.GetNonce(sub.wallet.CurrentAccount().Address)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("account:%v,nonce:%v\n", sub.wallet.CurrentAccount().Address, nonce)
	tx, err := sub.submitTopHeader(data, nonce)
	if err != nil {
		t.Fatal("SubmitHeader error:", err)
	}
	t.Log("SubmitHeader hash:", tx.Hash()) */
}
