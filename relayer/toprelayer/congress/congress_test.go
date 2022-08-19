package congress

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wonderivan/logger"
)

const hecoUrl = "https://http-mainnet.hecochain.com"

func TestCheckValidatorsNum(t *testing.T) {
	var height uint64 = 17276000

	ethsdk, err := ethclient.Dial(hecoUrl)
	if err != nil {
		t.Fatal(err)
	}
	header, err := ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(height))
	if err != nil {
		t.Fatal("HeaderByNumber: ", err)
	}
	t.Log(common.Bytes2Hex(header.Extra))
	validators := make([]common.Address, (len(header.Extra)-extraVanity-extraSeal)/common.AddressLength)
	for i := 0; i < len(validators); i++ {
		copy(validators[i][:], header.Extra[extraVanity+i*common.AddressLength:])
	}
	for _, v := range validators {
		t.Log("validator:", v)
	}
	t.Log("validator num:", len(validators))
}

func TestInit(t *testing.T) {
	var height uint64 = 17276012

	ethsdk, err := ethclient.Dial(hecoUrl)
	if err != nil {
		t.Fatal(err)
	}

	con := New(ethsdk)
	err = con.Init(height)
	if err != nil {
		t.Fatal(err)
	}
	header, err := ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(height+1))
	if err != nil {
		t.Fatal("HeaderByNumber: ", err)
	}
	bytes, err := con.GetLastSnapBytes(header)
	if err != nil {
		t.Fatal(err)
	}
	logger.Debug(bytes)
}

func TestEncodeSnapOut(t *testing.T) {
	snap := new(Snapshot)
	snap.Number = 17276004
	snap.Hash = common.HexToHash("0x6fd4c6243c79e0a1d424e59a33abfc495decd7a57a9e571a6c52eeba1274435a")
	snap.Validators = make(map[common.Address]struct{})
	snap.Recents = make(map[uint64]common.Address)
	snap.Validators[common.HexToAddress("0x09Dc2AbA0419dd915e675b08C385Ff565783978F")] = struct{}{}
	snap.Recents[17276004] = common.HexToAddress("0xE5FBAB23d0117b7f0f32ea01bfCabe38eFA0bFD4")

	var height uint64 = 17276000

	ethclient, err := ethclient.Dial(hecoUrl)
	if err != nil {
		t.Fatal(err)
	}
	header, err := ethclient.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(height+5))
	if err != nil {
		t.Fatal("HeaderByNumber: ", err)
	}
	out, err := encodeSnapshot(header, snap)
	if err != nil {
		t.Fatal(err)
	}
	logger.Debug(common.Bytes2Hex(out))
}
