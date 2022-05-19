package wallet

import (
	"context"
	"math/big"
	"testing"
	"toprelayer/sdk"
	"toprelayer/util"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const CHAINID uint64 = 1337

var url string = "http://192.168.30.32:8545"

func newsdk() (*sdk.SDK, error) {
	return sdk.NewSDK(url)
}

func newWallet() (*Wallet, error) {
	sdk, err := newsdk()
	if err != nil {
		return nil, err
	}

	w := new(Wallet)
	w.chainId = CHAINID
	w.sdk = sdk

	err = w.initWallet(DEFAULTPATH, pass)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func newsdkWallet() (*Wallet, error) {
	sdk, err := newsdk()
	if err != nil {
		return nil, err
	}
	return &Wallet{sdk: sdk}, nil
}

func TestGetBalance(t *testing.T) {
	iw, err := NewWallet(url, DEFAULTPATH, pass, CHAINID)
	if err != nil {
		t.Fatalf("new wallet error:%v", err)
	}
	w := iw.(IWallet_ETH)
	acc := w.CurrentAccount()
	b, err := w.GetBalance()
	if err != nil {
		t.Fatalf("get[%v] balance error:%v", acc.Address, err)
	}
	//acc1 := "0xd7c8e1e98a4985c4be72f605f91ae31bea074b24"
	t.Logf("addr[%v] balance:%v", acc.Address, b.Uint64())
}

func TestGasPrice(t *testing.T) {
	/* w, err := newsdkWallet()
	if err != nil {
		t.Fatalf("new wallet error:%v", err)
	} */

	iw, err := NewWallet(url, DEFAULTPATH, pass, CHAINID)
	if err != nil {
		t.Fatalf("new wallet error:%v", err)
	}
	w := iw.(IWallet_ETH)

	pric, err := w.GasPrice(context.Background())
	if err != nil {
		t.Fatal("gas price error:", err)
	}
	t.Log("gas price:", pric.Uint64())
}

func TestGasTip(t *testing.T) {
	w, err := newsdkWallet()
	if err != nil {
		t.Fatalf("new wallet error:%v", err)
	}
	if w.sdk == nil {
		t.Fatal("fatal error: nil sdk!")
	}
	tip, err := w.GasTipCap(context.Background())
	if err != nil {
		t.Fatal("GasTip error:", err)
	}
	t.Log("gas tip:", tip.Uint64())
}

func TestSendTransaction(t *testing.T) {
	w, err := newWallet()
	if err != nil {
		t.Fatalf("new wallet error:%v", err)
	}

	to := "0xd8aE0197425C0eA651264b06978580DcB62f3c91"
	gasprice, err := w.GasPrice(context.Background())
	if err != nil {
		t.Fatal("gas price error:", err)
	}
	amount := big.NewInt(1000)
	var gaslimit uint64 = 30000

	acc := w.currentAccount()
	t.Log("current account:", acc.Address.Hex())

	nonce, err := w.GetNonce(acc.Address)
	if err != nil {
		t.Fatal("GetNonce error:", err)
	}

	t.Log("nonce:", nonce, "gaslimit:", gaslimit)

	etx := types.NewTransaction(nonce, common.HexToAddress(to), amount, gaslimit, gasprice, nil)
	stx, err := w.SignTx(etx)
	if err != nil {
		t.Fatalf("Failed to sign with unlocked account: %v", err)
	}

	err = util.VerifyEthSignature(stx)
	if err != nil {
		t.Fatalf("VerifyEthSignature error: %v", err)
	}
	err = w.sdk.SendTransaction(context.Background(), stx)
	if err != nil {
		t.Fatal("SendTransaction error:", err)
	}

	t.Log("SendTransaction success!!! hash:", stx.Hash(), "nonce:", nonce, "gaslimit:", gaslimit)
}
