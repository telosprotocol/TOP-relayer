package toprelayer

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"toprelayer/config"
	"toprelayer/contract/top/ethclient"
	"toprelayer/relayer/toprelayer/ethashapp"
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

// testnet
// https://ropsten.infura.io/v3/fb2a09e82a234971ad84203e6f75990e

// const ethUrl string = "https://eth-mainnet.token.im"
const ethUrl = "https://ropsten.infura.io/v3/fb2a09e82a234971ad84203e6f75990e"
const topChainId uint64 = 1023
const defaultPass = "asd123"

func TestGetHeaderRlp(t *testing.T) {
	var height uint64 = 12989998

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
	var height uint64 = 12622433

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
	input, err := ethclient.PackSyncParam(rlp_bytes)
	if err != nil {
		t.Fatal(err)
	}
	logger.Debug("data:", common.Bytes2Hex(input))
}

func TestGetSyncTxData(t *testing.T) {
	// changable
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
	input, err := ethclient.PackSyncParam(batch)
	if err != nil {
		t.Fatal(err)
	}
	logger.Debug("data:", common.Bytes2Hex(input))
}

func TestGetHeightTxData(t *testing.T) {
	input, err := ethclient.PackGetHeightParam()
	if err != nil {
		t.Fatal(err)
	}
	logger.Debug("data:", common.Bytes2Hex(input))
}

func TestGetIsConfirmedTxData(t *testing.T) {
	height := big.NewInt(12970000)
	hash := common.HexToHash("13049bb8cfd97fe2333829f06df37c569db68d42c23097fbac64f2c61471f281")
	input, err := ethclient.PackIsKnownParam(height, hash)
	if err != nil {
		t.Fatal(err)
	}
	logger.Debug("data:", common.Bytes2Hex(input))
}

func TestSync(t *testing.T) {
	// changable
	var height uint64 = 12970000
	var topUrl string = "http://192.168.30.200:8080"
	var keyPath = "../../.relayer/wallet/top"

	cfg := &config.Relayer{
		Url:     topUrl,
		ChainId: topChainId,
		KeyPath: keyPath,
	}
	topRelayer := &TopRelayer{}
	err := topRelayer.Init(config.ETH_CHAIN, cfg, ethUrl, defaultPass)
	if err != nil {
		t.Fatal(err)
	}
	err = topRelayer.signAndSendTransactions(height, height+1)
	if err != nil {
		t.Fatal("submitEthHeader:", err)
	}
}

func TestSyncHeaderWithProofsRlpGas(t *testing.T) {
	// changable
	var height uint64 = 12970000
	var topUrl string = "http://192.168.30.200:8080"
	var keyPath = "../../.relayer/wallet/top"

	cfg := &config.Relayer{
		Url:     topUrl,
		ChainId: topChainId,
		KeyPath: keyPath,
	}
	topRelayer := &TopRelayer{}
	err := topRelayer.Init(config.ETH_CHAIN, cfg, ethUrl, defaultPass)
	if err != nil {
		t.Fatal(err)
	}
	header, err := topRelayer.ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(height))
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
	gaspric, err := topRelayer.wallet.GasPrice(context.Background())
	if err != nil {
		logger.Fatal(err)
	}
	packHeader, err := ethclient.PackSyncParam(rlp_bytes)
	if err != nil {
		logger.Fatal(err)
	}
	gaslimit, err := topRelayer.wallet.EstimateGas(context.Background(), &topRelayer.contract, gaspric, packHeader)
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println("gaslimit: ", gaslimit)
}

func TestGetEthClientHeight(t *testing.T) {
	// changable
	var topUrl string = "http://192.168.30.200:8080"
	var keyPath = "../../.relayer/wallet/top"

	cfg := &config.Relayer{
		Url:     topUrl,
		ChainId: topChainId,
		KeyPath: keyPath,
	}
	topRelayer := &TopRelayer{}
	err := topRelayer.Init(config.ETH_CHAIN, cfg, ethUrl, defaultPass)
	if err != nil {
		t.Fatal(err)
	}
	destHeight, err := topRelayer.callerSession.GetHeight()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("current height:", destHeight)
}
