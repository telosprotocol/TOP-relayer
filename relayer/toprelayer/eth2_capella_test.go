package toprelayer

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"
	"os"
	"testing"
	"time"
	"toprelayer/config"
)

var relayer *Eth2TopRelayerV2

func init() {
	var topUrl string = "http://192.168.2.104:8080"
	var keyPath = "../../.relayer/wallet/top"

	cfg := &config.Relayer{
		Url:     []string{topUrl},
		KeyPath: keyPath,
	}
	relayer = &Eth2TopRelayerV2{}
	if err := relayer.Init(cfg, []string{config.ETHAddr, config.ETHPrysm}, defaultPass); err != nil {
		panic(err.Error())
	}
}

func TestEthInitContract(t *testing.T) {
	nonce, err := relayer.wallet.NonceAt(context.Background(), relayer.wallet.Address(), nil)
	if err != nil {
		t.Fatal(err)
	}
	gaspric, err := relayer.wallet.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	//must init ops as bellow
	ops := &bind.TransactOpts{
		From:      relayer.wallet.Address(),
		Nonce:     big.NewInt(0).SetUint64(nonce),
		GasLimit:  5000000,
		GasFeeCap: gaspric,
		GasTipCap: big.NewInt(0),
		Signer:    relayer.signTransaction,
		Context:   context.Background(),
		NoSend:    false,
	}
	data, err := os.ReadFile("../../util/eth_init_data")
	if err != nil {
		t.Fatal(err)
	}
	if tx, err := relayer.transactor.Init(ops, data); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(tx.Hash())
	}
}

func TestReset(t *testing.T) {
	nonce, err := relayer.wallet.NonceAt(context.Background(), relayer.wallet.Address(), nil)
	if err != nil {
		t.Fatal(err)
	}
	gaspric, err := relayer.wallet.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	//must init ops as bellow
	ops := &bind.TransactOpts{
		From:      relayer.wallet.Address(),
		Nonce:     big.NewInt(0).SetUint64(nonce),
		GasLimit:  5000000,
		GasFeeCap: gaspric,
		GasTipCap: big.NewInt(0),
		Signer:    relayer.signTransaction,
		Context:   context.Background(),
		NoSend:    false,
	}
	if tx, err := relayer.transactor.Reset(ops); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(tx.Hash())
	}
}

func TestEthClient(t *testing.T) {
	b, err := relayer.beaconrpcclient.GetLastSlotNumber()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(b)

	//ethClient, err := ethclient.Dial(config.ETHAddr)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//relayer
}

func TestInitialized(t *testing.T) {
	if ok, err := relayer.callerSession.Initialized(); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(ok)
	}
}

func TestGetClientMode(t *testing.T) {
	mode, err := relayer.callerSession.GetClientMode()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(mode)
}

func TestLCU(t *testing.T) {
	v2, err := relayer.beaconrpcclient.GetLastFinalizedLightClientUpdateV2()
	if err != nil {
		t.Fatal(err)
	}
	encode, err := v2.Encode()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(encode)
	fmt.Println(v2)

	fmt.Println(time.Now())
}

//func TestSubmit(t *testing.T) {
//	relayer := newRelayerClient()
//	headers, err := relayer.submitHeaders()
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println(headers)
//}
