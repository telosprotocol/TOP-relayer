package top2eth

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"toprelayer/base"
	"toprelayer/msg"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const SUBMITTERURL string = "http://192.168.50.235:8545"
const LISTENURL string = "http://192.168.50.204:19086"

var DEFAULTPATH = "../../.relayer/wallet/eth"

var CONTRACT common.Address = common.HexToAddress("0xC78f29Abb15c016cf562821E988f6C1431C5469A")

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

func TestEstimateGas(t *testing.T) {
	sub := &Top2EthRelayer{}
	err := sub.Init(SUBMITTERURL, LISTENURL, DEFAULTPATH, "", base.ETH, CONTRACT, 90, 0, false)
	if err != nil {
		t.Fatal(err)
	}

	header := &types.Header{Number: big.NewInt(int64(1))}
	data, err := msg.EncodeHeaders(header)
	if err != nil {
		t.Fatal("EncodeToBytes:", err)
	}

	pric, err := sub.wallet.GasPrice(context.Background())
	if err != nil {
		t.Fatal("GasPrice:", err)
	}

	gaslimit, err := sub.estimateGas(pric, data)
	if err != nil {
		t.Fatal("estimateGas:", err)
	}
	t.Log("gaslimit:", gaslimit)

	_, err = sub.getEthBridgeState()
	if err != nil {
		t.Fatal("getEthBridgeState:", err)
	}
}
