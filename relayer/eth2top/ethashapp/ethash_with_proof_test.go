package ethashapp

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"toprelayer/sdk/ethsdk"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
// METHOD_GETBRIDGESTATE = "getCurrentBlockHeight"
// SYNCHEADERS           = "syncBlockHeader"

// SUCCESSDELAY int64 = 15 //mainnet 120
// FATALTIMEOUT int64 = 24 //hours
// FORKDELAY    int64 = 5  //mainnet 10000 seconds
// ERRDELAY     int64 = 10
// CONFIRMDELAY int64 = 5

// BLOCKS_PER_EPOCH uint64 = 30000
// BLOCKS_TO_END_OF_EPOCH uint64 = 5000
)

func TestGetHeaderWithProof(t *testing.T) {
	var url string = "https://api.mycryptoapi.com/eth"
	ethsdk, err := ethsdk.NewEthSdk(url)
	if err != nil {
		t.Fatal("NewEthSdk: ", err)
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
