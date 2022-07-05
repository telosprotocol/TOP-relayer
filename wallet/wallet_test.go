package wallet

import (
	"context"
	"math/big"
	"testing"
	"toprelayer/util"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	defaultPass = "asd123"
)

var (
	DEFAULTPATH = "../.relayer/wallet/top"
	chainid     = uint64(1023)

	//URL string = "http://192.168.50.235:8545"
	URL string = "http://192.168.50.204:19086"
	//URL string = "http://127.0.0.1:37399"

)

func TestGetBalance(t *testing.T) {
	w, err := NewWallet(URL, DEFAULTPATH, defaultPass, chainid)
	if err != nil {
		t.Fatalf("new wallet error:%v", err)
	}

	addr := w.CurrentAccount().Address
	//addr := common.HexToAddress("0xd8aE0197425C0eA651264b06978580DcB62f3c91")
	b, err := w.GetBalance(addr)
	if err != nil {
		t.Fatalf("get[%v] balance error:%v", addr, err)
	}
	t.Logf("addr[%v] balance:%v", addr, b.Uint64())
}

func TestGetNonce(t *testing.T) {
	w, err := NewWallet(URL, DEFAULTPATH, defaultPass, chainid)
	if err != nil {
		t.Fatalf("new wallet error:%v", err)
	}

	addr := w.CurrentAccount().Address
	//addr := common.HexToAddress("0xd8aE0197425C0eA651264b06978580DcB62f3c91")
	nonce, err := w.GetNonce(addr)
	if err != nil {
		t.Fatalf("get[%v] balance error:%v", addr, err)
	}
	t.Logf("addr[%v] nonce:%v", addr, nonce)
}

func TestGasPrice(t *testing.T) {
	w, err := NewWallet(URL, DEFAULTPATH, defaultPass, chainid)
	if err != nil {
		t.Fatalf("new wallet error:%v", err)
	}

	pric, err := w.GasPrice(context.Background())
	if err != nil {
		t.Fatal("gas price error:", err)
	}
	t.Log("gas price:", pric.Uint64())
}

func TestGasTip(t *testing.T) {
	w, err := NewWallet(URL, DEFAULTPATH, defaultPass, chainid)
	if err != nil {
		t.Fatalf("new wallet error:%v", err)
	}
	tip, err := w.GasTipCap(context.Background())
	if err != nil {
		t.Fatal("GasTip error:", err)
	}
	t.Log("gas tip:", tip.Uint64())
}

func TestSendTransaction(t *testing.T) {
	w, err := NewWallet(URL, DEFAULTPATH, defaultPass, chainid)
	if err != nil {
		t.Fatalf("new wallet error:%v", err)
	}

	contract := common.HexToAddress("0xa3e165d80c949833C5c82550D697Ab31Fd3BB446")
	gasprice, err := w.GasPrice(context.Background())
	if err != nil {
		t.Fatal("gas price error:", err)
	}
	amount := big.NewInt(1000)
	var gaslimit uint64 = 500000

	acc := w.CurrentAccount()
	t.Log("current account:", acc.Address.Hex())

	nonce, err := w.GetNonce(acc.Address)
	if err != nil {
		t.Fatal("GetNonce error:", err)
	}

	balance, err := w.GetBalance(acc.Address)
	if err != nil {
		t.Fatalf("GetBalance error:%v", err)
	}

	t.Log("account:", acc.Address, "balance:", balance, "nonce:", nonce, "gaslimit:", gaslimit, "gasprice:", gasprice)

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

func TestSendDynamicTx(t *testing.T) {
	contract := common.HexToAddress("0xa3e165d80c949833C5c82550D697Ab31Fd3BB446")

	w, err := NewWallet(URL, DEFAULTPATH, defaultPass, chainid)
	if err != nil {
		t.Fatalf("new wallet error:%v", err)
	}
	acc := w.CurrentAccount()

	balance, err := w.GetBalance(acc.Address)
	if err != nil {
		t.Fatalf("GetBalance error:%v", err)
	}

	/* nonce, err := w.GetNonce(acc.Address)
	if err != nil {
		t.Fatal("GetNonce error:", err)
	}
	*/
	nonce := uint64(1)
	// gastip, err := w.GasTipCap(context.Background())
	// if err != nil {
	// 	t.Fatal("GasTipCap error:", err)
	// }
	gastip := big.NewInt(0).SetUint64(1000000000)
	capfee := big.NewInt(0).SetUint64(2000000000)

	// gaspric, err := w.GasPrice(context.Background())
	// if err != nil {
	// 	t.Fatal("GasPrice error:", err)
	// }

	// msg := ethereum.CallMsg{
	// 	From:     acc.Address,
	// 	To:       &contract,
	// 	GasPrice: gaspric,
	// 	Value:    big.NewInt(1000),
	// 	Data:     []byte{2, 4, 44, 63, 22, 120},
	// }

	// gaslimit, err := w.EstimateGas(context.Background(), msg)
	// if err != nil {
	// 	t.Fatal("EstimateGas error:", err)
	// }
	var gaslimit uint64 = 500000

	t.Log("account:", acc.Address, "balance:", balance, "nonce:", nonce, "gastipfee:", gastip, "gaslimit:", gaslimit)

	var headers []*types.Header
	for i := 1; i <= 2; i++ {
		headers = append(headers, &types.Header{Number: big.NewInt(int64(i))})
	}

	data, err := rlp.EncodeToBytes(&headers)
	if err != nil {
		t.Fatal("EncodeToBytes:", err)
	}

	baseTx := &types.DynamicFeeTx{
		To:        &contract,
		Nonce:     nonce,
		GasFeeCap: capfee,
		GasTipCap: gastip,
		Gas:       gaslimit,
		Value:     big.NewInt(1000),
		Data:      data,
	}

	tx := types.NewTx(baseTx)

	stx, err := w.SignTx(tx)
	if err != nil {
		t.Fatal("SignTx error:", err)
	}
	// err = util.VerifyEthSignature(stx)
	// if err != nil {
	// 	t.Fatalf("VerifyEthSignature error: %v", err)
	// }

	t.Log("transaction:", stx)
	err = w.SendTransaction(context.Background(), stx)
	if err != nil {
		t.Fatal("SendTransaction error:", err)
	}

	byt, err := stx.MarshalBinary()
	if err != nil {
		t.Fatal("MarshalBinary error:", err)
	}
	rawtx := hexutil.Encode(byt)
	t.Log("stx hash:", stx.Hash(), "type:", stx.Type())
	t.Log("rawtx:", rawtx)

}
