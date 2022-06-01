package top2eth

import (
	"bytes"
	"context"
	"math/big"
	"testing"
	"toprelayer/base"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const SUBMITTERURL string = "http://192.168.50.235:8545"
const LISTENURL string = "http://192.168.50.204:19086"

var DEFAULTPATH = "../../.relayer/wallet/eth"

var CONTRACT common.Address = common.HexToAddress("0xd287F92c8cB8Cd54DC4C93a1619b04481E4a62F9")
var abipath string = "../../contract/ethbridge/ethbridge.abi"

func TestSubmitHeader(t *testing.T) {
	sub := new(Top2EthRelayer)
	err := sub.Init(SUBMITTERURL, LISTENURL, DEFAULTPATH, "", abipath, base.ETH, CONTRACT, 100, 0, false)
	if err != nil {
		t.Fatal(err)
	}

	/* h, err := sub.topsdk.GetLatestTopElectBlockHeight()
	if err != nil {
		t.Fatalf("GetLatestTopElectBlockHeight failed,error:%v", err)
	}

	header, err := sub.topsdk.GetTopElectBlockHeadByHeight(h)
	if err != nil {
		t.Fatalf("GetTopElectBlockHeadByHeight failed,error:%v", err)
	} */

	var batchHeaders [][]byte
	header := "0xf902faf8a980808080a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000008011c0f90228f843a01f49cece388273d5fbf8248d023143428ac035ff2edb57d2c3ad80f9de24e40fa0334e7866658e01fd2f23d0cf8a2a56ef849a05a056ef9618089db89ed4b9b59480f843a05fcaf45a46bace3b25acec4ed00e34ef65a28f354dfdbf20010f59d231b3f9bfa059d1b64da30c75c05d10e16b3c6716537b1c398e4e6c2ec0b3a8d16beafc0b5980f843a024c107a15840045be1176e681252000cd49357ac1eb3dfb91da2f5380aa6c67fa02586dab5dae8030b93076cc397129d5319e431a42cd0ce8bf37ba3fb6ffdf91180f843a07e114de9907a9623f7a33a9cab75d7d7f97194eb6c569c0d2c9d1f52a2b6a4bca0290aee80ae525e79b01e47ccef4dfa73a88762c27181704a9901df10c045646f80f843a0f90e2bc1e2e74e606b294f9d6eb1b2631efcd47fbbeb228b46cdc8c1ed281f71a027963972089d6355bea43c39835141d06f43975b254ce7cbe97ebd32c801d82c01f843a0c1e8fada4ae896ffe973b913c81e50834a0ed7d4355667841fcb1d514c216b01a064efb5e3126d95cac81ae1d87924fd0408a83e1b33f4f4b941b1c7c256f3567c80f843a0da65a9248d67db1d7ec33609600eb8fe68e9f91d89e3629c0f946d7dc3781518a03f42be47ac7176b9bb3eec681d557cefa009fd0eca223b8e6f98c898f43f68e980f843a0191e3e4fc0e05827601d746723482d6675ff2b3080388ca9c09b1444a7cd0bada00cd3ec900d51eab003000aa5ec71e9e3a8508e028c7a98955ecac2d16c6554b080"

	batchHeaders = append(batchHeaders, common.Hex2Bytes(header[2:]))
	batchHeaders = append(batchHeaders, common.Hex2Bytes(header[2:]))
	batchHeaders = append(batchHeaders, common.Hex2Bytes(header[2:]))

	data := bytes.Join(batchHeaders, []byte{})

	t.Log("header encode data:", data, "len:", len(data))
	t.Log("hex data:", common.Bytes2Hex(data))

	nonce, err := sub.wallet.GetNonce(sub.wallet.CurrentAccount().Address)
	if err != nil {
		t.Fatal("GasPrice:", err)
	}

	tx, err := sub.submitTopHeader(data, nonce)
	if err != nil {
		t.Fatal("submitEthHeader:", err)
	}
	t.Log("hash:", tx.Hash())

	/* err = sub.StartRelayer(&sync.WaitGroup{})
	if err != nil {
		t.Fatal(err)
	} */
}

func TestEstimateGas(t *testing.T) {
	sub := &Top2EthRelayer{}
	err := sub.Init(SUBMITTERURL, LISTENURL, DEFAULTPATH, "", abipath, base.ETH, CONTRACT, 90, 0, false)
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
	t.Log("gaslimit:", gaslimit)

	_, err = sub.getEthBridgeCurrentHeight()
	if err != nil {
		t.Fatal("getEthBridgeCurrentHeight:", err)
	}
}

func TestGetETHBridgeState(t *testing.T) {
	sub := &Top2EthRelayer{}
	err := sub.Init(SUBMITTERURL, LISTENURL, DEFAULTPATH, "", abipath, base.ETH, CONTRACT, 90, 0, false)
	if err != nil {
		t.Fatal(err)
	}
	curr, err := sub.getEthBridgeCurrentHeight()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("current height:", curr)
}
