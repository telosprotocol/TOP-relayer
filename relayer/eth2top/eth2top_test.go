package eth2top

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"toprelayer/relayer/eth2top/ethashapp"
	"toprelayer/sdk/ethsdk"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/wonderivan/logger"
)

// main net
// https://api.mycryptoapi.com/eth
// https://web3.1inch.exchange/
// https://eth-mainnet.gateway.pokt.network/v1/5f3453978e354ab992c4da79
// https://eth-mainnet.token.im
const ethUrl string = "https://eth-mainnet.token.im"
const topChainId uint64 = 1023

func TestGetHeaderRlp(t *testing.T) {
	var height uint64 = 12969999

	ethsdk, err := ethsdk.NewEthSdk(ethUrl)
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

func TestGetHeadersWithProofsRlp(t *testing.T) {
	var start_height uint64 = 12970000
	var sync_num uint64 = 1

	ethsdk, err := ethsdk.NewEthSdk(ethUrl)
	if err != nil {
		t.Fatal("NewEthSdk: ", err)
	}

	var batch []byte
	for h := start_height; h <= start_height+sync_num-1; h++ {
		header, err := ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(h))
		if err != nil {
			t.Fatal("HeaderByNumber: ", err)
		}
		out, err := ethashapp.EthashWithProofs(h, header)
		if err != nil {
			t.Fatal("HeaderByNumber: ", err)
		}
		rlp_bytes, err := rlp.EncodeToBytes(out)
		if err != nil {
			t.Fatal("rlp encode error: ", err)
		}
		batch = append(batch, rlp_bytes...)
	}
	fmt.Println("rlp output: ", common.Bytes2Hex(batch))
}

func TestGetInitTxData(t *testing.T) {
	var height uint64 = 12969999
	var abiPath string = "../../contract/topbridge/topbridge.abi"

	ethsdk, err := ethsdk.NewEthSdk(ethUrl)
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

func TestGetSyncTxData(t *testing.T) {
	// changable
	var start_height uint64 = 12970000
	var sync_num uint64 = 1
	var abiPath string = "../../contract/topbridge/topbridge.abi"

	ethsdk, err := ethsdk.NewEthSdk(ethUrl)
	if err != nil {
		t.Fatal("NewEthSdk: ", err)
	}
	var batch []byte
	for h := start_height; h <= start_height+sync_num-1; h++ {
		header, err := ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(h))
		if err != nil {
			t.Fatal("HeaderByNumber: ", err)
		}
		out, err := ethashapp.EthashWithProofs(h, header)
		if err != nil {
			t.Fatal("HeaderByNumber: ", err)
		}
		rlp_bytes, err := rlp.EncodeToBytes(out)
		if err != nil {
			t.Fatal("rlp encode error: ", err)
		}
		batch = append(batch, rlp_bytes...)
	}
	abi, err := initABI(abiPath)
	if err != nil {
		t.Fatal(err)
	}
	input, err := abi.Pack("sync", batch)
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

func TestGetIsConfirmedTxData(t *testing.T) {
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

func TestSync(t *testing.T) {
	// changable
	var height uint64 = 12970000
	var contract common.Address = common.HexToAddress("0x0eD0BA13032aDD72398042B931aecCEFCc66A826")
	var submitUrl string = "http://192.168.50.204:8080"
	var accountPath = "../../.relayer/wallet/top"
	// fix

	sub := &Eth2TopRelayer{}
	err := sub.Init(submitUrl, ethUrl, accountPath, "", topChainId, contract, 5)
	if err != nil {
		t.Fatal(err)
	}
	ethsdk, err := ethsdk.NewEthSdk(ethUrl)
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

	sub := &Eth2TopRelayer{}
	err := sub.Init(submitUrl, ethUrl, accountPath, "", topChainId, contract, 5)
	if err != nil {
		t.Fatal(err)
	}
	ethsdk, err := ethsdk.NewEthSdk(ethUrl)
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

	sub := &Eth2TopRelayer{}
	err := sub.Init(submitUrl, ethUrl, accountPath, "", topChainId, contract, 5)
	if err != nil {
		t.Fatal(err)
	}

	height, err := sub.getTopBridgeCurrentHeight()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("current height:", height)
}
