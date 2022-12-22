package toprelayer

import (
	"context"
	"math/big"
	"testing"
	"toprelayer/config"
	"toprelayer/relayer/toprelayer/congress"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/wonderivan/logger"
)

const hecoUrl = "https://http-mainnet.hecochain.com"

func TestGetHecoInitData(t *testing.T) {
	ethclient, err := ethclient.Dial(hecoUrl)
	if err != nil {
		t.Fatal(err)
	}
	destHeight, err := ethclient.BlockNumber(context.Background())
	height := (destHeight - 11) / 200 * 200
	logger.Info("heco init with height: %v - %v", height, height+11)
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

func TestGetHecoSyncData(t *testing.T) {
	var start_height uint64 = 17276022
	var sync_num uint64 = 1

	ethsdk, err := ethclient.Dial(hecoUrl)
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

func TestHecoInitContract(t *testing.T) {
	var topUrl string = "http://192.168.30.200:8080"
	var keyPath = "../../.relayer/wallet/top"

	cfg := &config.Relayer{
		Url:     []string{topUrl},
		KeyPath: keyPath,
	}
	relayer := &Heco2TopRelayer{}
	err := relayer.Init(cfg, []string{hecoUrl}, defaultPass)
	if err != nil {
		t.Fatal(err)
	}
	nonce, err := relayer.wallet.NonceAt(context.Background(), relayer.wallet.Address(), nil)
	if err != nil {
		t.Error(err)
	}
	gaspric, err := relayer.wallet.SuggestGasPrice(context.Background())
	if err != nil {
		t.Error(err)
	}
	//must init ops as bellow
	ops := &bind.TransactOpts{
		From:      relayer.wallet.Address(),
		Nonce:     big.NewInt(0).SetUint64(nonce),
		GasLimit:  50000,
		GasFeeCap: gaspric,
		GasTipCap: big.NewInt(0),
		Signer:    relayer.signTransaction,
		Context:   context.Background(),
		NoSend:    false,
	}
	var initHeaders []byte
	tx, err := relayer.transactor.Init(ops, initHeaders, "")
	if err != nil {
		t.Error(err)
	}
	t.Log(tx.Hash())
}

func TestHecoInit(t *testing.T) {
	// changable
	var start_height uint64 = 17276000
	var sync_num uint64 = 12
	var topUrl string = "http://192.168.30.200:8080"
	var keyPath = "../../.relayer/wallet/top"

	cfg := &config.Relayer{
		Url:     []string{topUrl},
		KeyPath: keyPath,
	}
	relayer := &Heco2TopRelayer{}
	err := relayer.Init(cfg, []string{ethUrl}, defaultPass)
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

	nonce, err := relayer.wallet.NonceAt(context.Background(), relayer.wallet.Address(), nil)
	if err != nil {
		t.Fatal(err)
	}
	gaspric, err := relayer.wallet.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	ops := &bind.TransactOpts{
		From:      relayer.wallet.Address(),
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
