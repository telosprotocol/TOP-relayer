package eth2top

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"toprelayer/base"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const SUBMITTERURL string = "http://192.168.50.204:19086"

const LISTENURL string = "http://192.168.50.235:8545"

var DEFAULTPATH = "../../.relayer/wallet/top"
var CONTRACT common.Address = common.HexToAddress("0xa3e165d80c949833C5c82550D697Ab31Fd3BB446")
var abipath string = "../../contract/topbridge/topbridge.abi"

func TestSubmitHeader(t *testing.T) {
	sub := &Eth2TopRelayer{}
	err := sub.Init(SUBMITTERURL, LISTENURL, DEFAULTPATH, "", abipath, base.TOP, CONTRACT, 90, 0, false)
	if err != nil {
		t.Fatal(err)
	}
	var batchHeaders []*types.Header

	currH, err := sub.getTopBridgeCurrentHeight()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("bridge contract current height:", currH)

	header, err := sub.ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(currH+1))
	if err != nil {
		t.Fatal(err)
	}
	{
		bt, _ := header.MarshalJSON()
		t.Log("header:", string(bt))
	}

	batchHeaders = append(batchHeaders, header)
	header2, err := sub.ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(currH+2))
	if err != nil {
		t.Fatal(err)
	}
	{
		bt, _ := header2.MarshalJSON()
		t.Log("header2:", string(bt))
	}
	batchHeaders = append(batchHeaders, header2)

	data, err := base.EncodeHeaders(batchHeaders)
	if err != nil {
		t.Fatal("EncodeToBytes:", err)
	}

	t.Log("header data:", data)

	nonce, err := sub.wallet.GetNonce(sub.wallet.CurrentAccount().Address)
	if err != nil {
		t.Fatal("GasPrice:", err)
	}

	tx, err := sub.submitEthHeader(data, nonce)
	if err != nil {
		t.Fatal("submitEthHeader:", err)
	}
	t.Log("hash:", tx.Hash())

	/* 	hashes, err := sub.signAndSendTransactions(1, 10)
	   	if err != nil {
	   		t.Fatal("signAndSendTransactions error:", err)
	   	}
	   	t.Log("hashes:", hashes) */

	/* nonce, err := sub.wallet.GetNonce(sub.wallet.CurrentAccount().Address)
	if err != nil {
		t.Fatal("GetNonce error:", err)
	}
	balance, err := sub.wallet.GetBalance(sub.wallet.CurrentAccount().Address)
	if err != nil {
		t.Fatal("GetBalance error:", err)
	}
	t.Log("balance:", balance, "nonce:", nonce)

	var headers []*types.Header
	for i := 1; i <= 2; i++ {
		headers = append(headers, &types.Header{Number: big.NewInt(int64(i))})
	}
	hash, err := sub.batch(headers, nonce)
	if err != nil {
		t.Fatal("batch error:", err)
	}
	t.Log("stx hash:", hash) */

	/* data, err := base.EncodeHeaders(&headers)
	if err != nil {
		t.Fatal("EncodeToBytes:", err)
	}

	if sub.wallet == nil {
		t.Fatal("nil wallet!!!")
	}

	stx, err := sub.submitEthHeader(data, nonce)
	if err != nil {
		t.Fatal("SubmitHeader error:", err)
	}
	t.Log("stx hash:", stx.Hash(), "type:", stx.Type())

	byt, err := stx.MarshalBinary()
	if err != nil {
		t.Fatal("MarshalBinary error:", err)
	}
	t.Log("rawtx:", hexutil.Encode(byt)) */
}

func TestEstimateGas(t *testing.T) {
	sub := &Eth2TopRelayer{}
	err := sub.Init(SUBMITTERURL, LISTENURL, DEFAULTPATH, "", abipath, base.TOP, CONTRACT, 90, 0, false)
	if err != nil {
		t.Fatal(err)
	}

	header := &types.Header{Number: big.NewInt(int64(1))}
	data, err := base.EncodeHeaders(header)
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
	t.Log("gasprice", pric, "gaslimit:", gaslimit)
}

func TestGetTopBridgeState(t *testing.T) {
	sub := &Eth2TopRelayer{}
	err := sub.Init(SUBMITTERURL, LISTENURL, DEFAULTPATH, "", abipath, base.TOP, CONTRACT, 90, 0, false)
	if err != nil {
		t.Fatal(err)
	}

	curr, err := sub.getTopBridgeCurrentHeight()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("current height:", curr)
}

func TestStartRelayer(t *testing.T) {
	sub := &Eth2TopRelayer{}
	err := sub.Init(SUBMITTERURL, LISTENURL, DEFAULTPATH, "", abipath, base.TOP, CONTRACT, 90, 0, false)
	if err != nil {
		t.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	err = sub.StartRelayer(wg)
	if err != nil {
		t.Fatal(err)
	}
}
