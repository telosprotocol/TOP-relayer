package toprelayer

import (
	"context"
	"math/big"
	"testing"
	"toprelayer/config"
	"toprelayer/relayer/toprelayer/congress"
	"toprelayer/sdk/ethsdk"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/wonderivan/logger"
)

const bscUrl = "https://bsc-dataseed4.binance.org"

func TestGetBscInitData(t *testing.T) {
	var height uint64 = 20250000

	ethclient, err := ethclient.Dial(hecoUrl)
	if err != nil {
		t.Fatal(err)
	}
	var batch []byte
	for i := height; i <= height+11; i++ {
		header, err := ethclient.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(i))
		if err != nil {
			t.Fatal(err)
		}
		rlp_bytes, err := rlp.EncodeToBytes(header)
		if err != nil {
			t.Fatal(err)
		}
		batch = append(batch, rlp_bytes...)
	}
	logger.Debug(common.Bytes2Hex(batch))
}

func TestGetBscSyncData(t *testing.T) {
	var start_height uint64 = 17276022
	var sync_num uint64 = 1

	ethsdk, err := ethsdk.NewEthSdk(hecoUrl)
	if err != nil {
		t.Fatal(err)
	}
	con := congress.New(ethsdk)
	err = con.Init(start_height - 1)
	if err != nil {
		t.Fatal(err)
	}

	var batch []byte
	for h := start_height; h <= start_height+sync_num-1; h++ {
		header, err := ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(17276012))
		if err != nil {
			t.Fatal(err)
		}
		out, err := con.GetLastSnapBytes(header)
		if err != nil {
			t.Fatal(err)
		}
		batch = append(batch, out...)
	}
	logger.Debug(common.Bytes2Hex(batch))
}

func TestBscInit(t *testing.T) {
	// changable
	var start_height uint64 = 20250000
	var sync_num uint64 = 12
	var bscUrl = "https://bsc-dataseed4.binance.org"
	var topUrl string = "http://192.168.30.200:8080"
	var keyPath = "../../.relayer/wallet/top"

	cfg := &config.Relayer{
		Url:     topUrl,
		ChainId: topChainId,
		KeyPath: keyPath,
	}
	relayer := &Heco2TopRelayer{}
	err := relayer.Init(config.BSC_CHAIN, cfg, bscUrl, defaultPass)
	if err != nil {
		t.Fatal(err)
	}
	var batch []byte
	for h := start_height; h <= start_height+sync_num-1; h++ {
		header, err := relayer.ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(h))
		if err != nil {
			t.Fatal(err)
		}
		rlp_bytes, err := rlp.EncodeToBytes(header)
		if err != nil {
			t.Fatal("EncodeToBytes: ", err)
		}
		batch = append(batch, rlp_bytes...)
	}

	nonce, err := relayer.wallet.GetNonce(relayer.wallet.CurrentAccount().Address)
	if err != nil {
		t.Fatal(err)
	}
	gaspric, err := relayer.wallet.GasPrice(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	ops := &bind.TransactOpts{
		From:      relayer.wallet.CurrentAccount().Address,
		Nonce:     big.NewInt(0).SetUint64(nonce),
		GasLimit:  500000,
		GasFeeCap: gaspric,
		GasTipCap: big.NewInt(0),
		Signer:    relayer.signTransaction,
		Context:   context.Background(),
		NoSend:    false,
	}
	tx, err := relayer.transactor.Init(ops, batch, string(""))
	if err != nil {
		t.Fatal(err)
	}
	logger.Debug(tx.Hash())
}
