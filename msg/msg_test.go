package msg

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func TestEncodeAndDecodeTopHeader(t *testing.T) {
	var headers []*TopElectBlockHeader
	for i := 1; i < 5; i++ {
		headers = append(headers, &TopElectBlockHeader{BlockNumber: uint64(i)})
	}
	t.Log("before encode:", headers[2].BlockNumber)
	buff, err := EncodeHeaders(&headers)
	if err != nil {
		t.Fatal("EncodeToBytes:", err)
	}

	var deHeaders []*TopElectBlockHeader
	err = rlp.DecodeBytes(buff, &deHeaders)
	if err != nil {
		t.Fatal("DecodeBytes:", deHeaders)
	}
	t.Log("after decode:", deHeaders[2].BlockNumber)
}

func TestEncodeAndDecodeEthHeader(t *testing.T) {
	var headers []*types.Header
	for i := 1; i < 5; i++ {
		headers = append(headers, &types.Header{Number: big.NewInt(int64(i))})
	}
	t.Log("before encode:", headers[2].Number)
	buff, err := EncodeHeaders(&headers)
	if err != nil {
		t.Fatal("EncodeToBytes:", err)
	}

	var deHeaders []*types.Header
	err = rlp.DecodeBytes(buff, &deHeaders)
	if err != nil {
		t.Fatal("DecodeBytes:", deHeaders)
	}
	t.Log("after decode:", deHeaders[2].Number)
}
