package wallet

import (
	"context"
	"math/big"
	"testing"
	"toprelayer/util"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	defaultPass = "asd123"
	defaultPath = "../.relayer/wallet/top"
	url         = "http://192.168.50.204:19086"
)

func TestGetBalance(t *testing.T) {
	w, err := NewWallet(url, defaultPath, defaultPass)
	if err != nil {
		t.Fatalf("new wallet error:%v", err)
	}

	addr := w.Address()
	b, err := w.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		t.Fatalf("get[%v] balance error:%v", addr, err)
	}
	
	t.Logf("addr[%v] balance:%v", addr, b.Uint64())
}

func TestGetNonce(t *testing.T) {
	w, err := NewWallet(url, defaultPath, defaultPass)
	if err != nil {
		t.Fatalf("new wallet error:%v", err)
	}

	addr := w.Address()
	//addr := common.HexToAddress("0xd8aE0197425C0eA651264b06978580DcB62f3c91")
	nonce, err := w.NonceAt(context.Background(), addr, nil)
	if err != nil {
		t.Fatalf("get[%v] balance error:%v", addr, err)
	}
	t.Logf("addr[%v] nonce:%v", addr, nonce)
}

func TestGasPrice(t *testing.T) {
	w, err := NewWallet(url, defaultPath, defaultPass)
	if err != nil {
		t.Fatalf("new wallet error:%v", err)
	}

	pric, err := w.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatal("gas price error:", err)
	}
	t.Log("gas price:", pric.Uint64())
}

func TestGasTip(t *testing.T) {
	w, err := NewWallet(url, defaultPath, defaultPass)
	if err != nil {
		t.Fatalf("new wallet error:%v", err)
	}
	tip, err := w.SuggestGasTipCap(context.Background())
	if err != nil {
		t.Fatal("GasTip error:", err)
	}
	t.Log("gas tip:", tip.Uint64())
}

func TestSendTransaction(t *testing.T) {
	w, err := NewWallet(url, defaultPath, defaultPass)
	if err != nil {
		t.Fatalf("new wallet error:%v", err)
	}

	contract := common.HexToAddress("0xa3e165d80c949833C5c82550D697Ab31Fd3BB446")
	gasprice, err := w.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatal("gas price error:", err)
	}
	amount := big.NewInt(1000)
	var gaslimit uint64 = 500000

	t.Log("current account:", w.Address().Hex())

	nonce, err := w.NonceAt(context.Background(), w.Address(), nil)
	if err != nil {
		t.Fatal("GetNonce error:", err)
	}

	balance, err := w.BalanceAt(context.Background(), w.Address(), nil)
	if err != nil {
		t.Fatalf("GetBalance error:%v", err)
	}

	t.Log("account:", w.Address(), "balance:", balance, "nonce:", nonce, "gaslimit:", gaslimit, "gasprice:", gasprice)

	etx := types.NewTransaction(nonce, contract, amount, gaslimit, gasprice, []byte{2, 4, 44, 63, 22, 120})
	stx, err := w.SignTx(etx)
	if err != nil {
		t.Fatalf("Failed to sign with unlocked account: %v", err)
	}

	err = util.VerifyEthSignature(stx)
	if err != nil {
		t.Fatalf("VerifyEthSignature error: %v", err)
	}

	data, err := stx.MarshalBinary()
	if err != nil {
		t.Fatal("MarshalBinary error:", err)
	}
	rawtx := hexutil.Encode(data)

	t.Log("To:", stx.To())

	t.Log("rawtx:", rawtx)

	err = w.SendTransaction(context.Background(), stx)
	if err != nil {
		t.Fatal("SendTransaction error:", err)
	}
	t.Log("SendTransaction success hash:", stx.Hash())
}
