package ethashapp

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

func TestGetHeaderWithProof(t *testing.T) {
	var url string = "https://api.mycryptoapi.com/eth"
	ethsdk, err := ethclient.Dial(url)
	if err != nil {
		t.Fatal(err)
	}
	header, err := ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(12970000))
	if err != nil {
		t.Fatal("HeaderByNumber: ", err)
	}
	out, err := EthashWithProofs(12970000, header)
	if err != nil {
		t.Fatal("HeaderByNumber: ", err)
	}
	rlp_bytes, err := rlp.EncodeToBytes(out)
	if err != nil {
		t.Fatal("rlp encode error: ", err)
	}
	fmt.Println("rlp output: ", common.Bytes2Hex(rlp_bytes))
}
