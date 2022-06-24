package eth2top

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"toprelayer/base"
	"toprelayer/relayer/eth2top/ethashapp"
	"toprelayer/sdk/ethsdk"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func TestGetHeaderRlp(t *testing.T) {
	var height uint64 = 12969999

	const url string = "https://api.mycryptoapi.com/eth"
	ethsdk, err := ethsdk.NewEthSdk(url)
	if err != nil {
		t.Fatal("NewEthSdk: ", err)
	}
	header, err := ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(height))
	if err != nil {
		t.Fatal("HeaderByNumber: ", err)
	}
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		t.Fatal("EncodeToBytes: ", err)
	}
	t.Log("headers hex data:", common.Bytes2Hex(data))
}

func TestGetBatchHeadersRlp(t *testing.T) {
	const url string = "https://api.mycryptoapi.com/eth"
	ethsdk, err := ethsdk.NewEthSdk(url)
	if err != nil {
		t.Fatal("NewEthSdk: ", err)
	}

	var currH uint64 = 12970000
	var num uint64 = 5
	var batchHeaders []*types.Header

	for i := uint64(1); i <= num; i++ {
		header, err := ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(currH+i))
		if err != nil {
			t.Fatal("HeaderByNumber", err)
		}
		batchHeaders = append(batchHeaders, header)
	}

	data, err := base.EncodeHeaders(batchHeaders)
	if err != nil {
		t.Fatal("EncodeToBytes:", err)
	}
	t.Log("headers hex data:", common.Bytes2Hex(data))
}

func TestGetHeaderWithProofsRlp(t *testing.T) {
	var height uint64 = 12970001

	const url string = "https://api.mycryptoapi.com/eth"
	ethsdk, err := ethsdk.NewEthSdk(url)
	if err != nil {
		t.Fatal("NewEthSdk: ", err)
	}

	header, err := ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(height))
	if err != nil {
		t.Fatal("HeaderByNumber: ", err)
	}
	out, err := ethashapp.EthashWithProofs(height, header)
	if err != nil {
		t.Fatal("HeaderByNumber: ", err)
	}
	rlp_bytes, err := rlp.EncodeToBytes(out)
	if err != nil {
		t.Fatal("rlp encode error: ", err)
	}
	fmt.Println("rlp output: ", common.Bytes2Hex(rlp_bytes))
}

func TestSyncHeaderWithProofsRlp(t *testing.T) {
	// changable
	var height uint64 = 12970000
	var contract common.Address = common.HexToAddress("0x0eD0BA13032aDD72398042B931aecCEFCc66A826")
	var submitUrl string = "http://192.168.30.200:8080"
	var accountPath = "../../.relayer/wallet/top"
	// fix
	var listenUrl string = "https://api.mycryptoapi.com/eth"
	var abiPath string = "../../contract/topbridge/topbridge.abi"

	sub := &Eth2TopRelayer{}
	err := sub.Init(submitUrl, listenUrl, accountPath, "", abiPath, base.TOP, contract, 5, 0, false)
	if err != nil {
		t.Fatal(err)
	}
	ethsdk, err := ethsdk.NewEthSdk(listenUrl)
	if err != nil {
		t.Fatal("NewEthSdk: ", err)
	}
	header, err := ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(height))
	if err != nil {
		t.Fatal("HeaderByNumber: ", err)
	}
	out, err := ethashapp.EthashWithProofs(height, header)
	if err != nil {
		t.Fatal("HeaderByNumber: ", err)
	}
	rlp_bytes, err := rlp.EncodeToBytes(out)
	if err != nil {
		t.Fatal("rlp encode error: ", err)
	}
	nonce, err := sub.wallet.GetNonce(sub.wallet.CurrentAccount().Address)
	if err != nil {
		t.Fatal("GasPrice:", err)
	}
	tx, err := sub.submitEthHeader(rlp_bytes, nonce)
	if err != nil {
		t.Fatal("submitEthHeader:", err)
	}
	t.Log("hash:", tx.Hash())
}

func TestSyncHeaderWithProofsRlpGas(t *testing.T) {
	// changable
	var contract common.Address = common.HexToAddress("0x0eD0BA13032aDD72398042B931aecCEFCc66A826")
	var submitUrl string = "http://192.168.30.200:8080"
	var accountPath = "../../.relayer/wallet/top"
	// fix
	var listenUrl string = "https://api.mycryptoapi.com/eth"
	var abiPath string = "../../contract/topbridge/topbridge.abi"

	sub := &Eth2TopRelayer{}
	err := sub.Init(submitUrl, listenUrl, accountPath, "", abiPath, base.TOP, contract, 5, 0, false)
	if err != nil {
		t.Fatal(err)
	}
	ethsdk, err := ethsdk.NewEthSdk(listenUrl)
	if err != nil {
		t.Fatal("NewEthSdk: ", err)
	}
	header, err := ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(12970000))
	if err != nil {
		t.Fatal("HeaderByNumber: ", err)
	}
	out, err := ethashapp.EthashWithProofs(12970000, header)
	if err != nil {
		t.Fatal("HeaderByNumber: ", err)
	}
	rlp_bytes, err := rlp.EncodeToBytes(out)
	if err != nil {
		t.Fatal("rlp encode error: ", err)
	}
	gaspric, err := sub.wallet.GasPrice(context.Background())
	if err != nil {
		t.Fatal("GasPrice error: ", err)
	}
	fmt.Println("data_len: ", len(rlp_bytes))
	fmt.Println("price: ", gaspric)
	gaslimit, err := sub.estimateGas(gaspric, rlp_bytes)
	if err != nil {
		t.Fatal("GasPrice error: ", err)
	}
	fmt.Println("limit: ", gaslimit)
}

func TestGetTopBridgeHeight(t *testing.T) {
	// changable
	var contract common.Address = common.HexToAddress("0x0eD0BA13032aDD72398042B931aecCEFCc66A826")
	var submitUrl string = "http://192.168.30.200:8080"
	var accountPath = "../../.relayer/wallet/top"
	// fix
	var listenUrl string = "https://api.mycryptoapi.com/eth"
	var abiPath string = "../../contract/topbridge/topbridge.abi"

	sub := &Eth2TopRelayer{}
	err := sub.Init(submitUrl, listenUrl, accountPath, "", abiPath, base.TOP, contract, 90, 0, false)
	if err != nil {
		t.Fatal(err)
	}

	height, err := sub.getTopBridgeCurrentHeight()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("current height:", height)
}
