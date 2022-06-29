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
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/wonderivan/logger"
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

func TestGetHeaderTxData(t *testing.T) {
	var height uint64 = 12969999
	var url string = "https://api.mycryptoapi.com/eth"
	var abiPath string = "../../contract/topbridge/topbridge.abi"

	ethsdk, err := ethsdk.NewEthSdk(url)
	if err != nil {
		t.Fatal("NewEthSdk: ", err)
	}
	header, err := ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(height))
	if err != nil {
		t.Fatal("HeaderByNumber: ", err)
	}
	rlp_bytes, err := rlp.EncodeToBytes(header)
	if err != nil {
		t.Fatal("EncodeToBytes: ", err)
	}
	abi, err := initABI(abiPath)
	if err != nil {
		t.Fatal(err)
	}
	input, err := abi.Pack("sync", rlp_bytes)
	if err != nil {
		t.Fatal(err)
	}
	logger.Debug("data:", common.Bytes2Hex(input))
}

func TestGetHeaderWithProofsRlpTxData(t *testing.T) {
	// changable
	var height uint64 = 12970000
	var listenUrl string = "https://api.mycryptoapi.com/eth"
	var abiPath string = "../../contract/topbridge/topbridge.abi"

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
	abi, err := initABI(abiPath)
	if err != nil {
		t.Fatal(err)
	}
	input, err := abi.Pack("sync", rlp_bytes)
	if err != nil {
		t.Fatal(err)
	}
	logger.Debug("data:", common.Bytes2Hex(input))
}

func TestGetHeightTxData(t *testing.T) {
	var abiPath string = "../../contract/topbridge/topbridge.abi"

	abi, err := initABI(abiPath)
	if err != nil {
		t.Fatal(err)
	}
	input, err := abi.Pack(METHOD_GETHEIGHT)
	if err != nil {
		t.Fatal(err)
	}
	logger.Debug("data:", common.Bytes2Hex(input))
}

func TestIsConfirmedTxData(t *testing.T) {
	hash := common.Hex2Bytes("13049bb8cfd97fe2333829f06df37c569db68d42c23097fbac64f2c61471f281")
	var abiPath string = "../../contract/topbridge/topbridge.abi"

	abi, err := initABI(abiPath)
	if err != nil {
		t.Fatal(err)
	}
	input, err := abi.Pack("is_confirmed", hash)
	if err != nil {
		t.Fatal(err)
	}
	logger.Debug("data:", common.Bytes2Hex(input))
}

func TestSyncHeaderWithProofsRlp(t *testing.T) {
	// changable
	var height uint64 = 12970000
	var contract common.Address = common.HexToAddress("0x0eD0BA13032aDD72398042B931aecCEFCc66A826")
	var submitUrl string = "http://192.168.50.204:8080"
	var accountPath = "../../.relayer/wallet/top"
	// fix
	var listenUrl string = "https://api.mycryptoapi.com/eth"

	sub := &Eth2TopRelayer{}
	err := sub.Init(submitUrl, listenUrl, accountPath, "", base.TOP, contract, 5)
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
	err = sub.submitEthHeader(rlp_bytes, nonce)
	if err != nil {
		t.Fatal("submitEthHeader:", err)
	}
}

func TestSyncHeaderWithProofsRlpGas(t *testing.T) {
	// changable
	var contract common.Address = common.HexToAddress("0x0eD0BA13032aDD72398042B931aecCEFCc66A826")
	var submitUrl string = "http://192.168.30.200:8080"
	var accountPath = "../../.relayer/wallet/top"
	// fix
	var listenUrl string = "https://api.mycryptoapi.com/eth"

	sub := &Eth2TopRelayer{}
	err := sub.Init(submitUrl, listenUrl, accountPath, "", base.TOP, contract, 5)
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
	gaslimit, err := sub.estimateSyncGas(gaspric, rlp_bytes)
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

	sub := &Eth2TopRelayer{}
	err := sub.Init(submitUrl, listenUrl, accountPath, "", base.TOP, contract, 5)
	if err != nil {
		t.Fatal(err)
	}

	height, err := sub.getTopBridgeCurrentHeight()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("current height:", height)
}
