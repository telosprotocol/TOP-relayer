package topsdk

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"toprelayer/sdk"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/wonderivan/logger"
)

type TopSdk struct {
	*sdk.SDK
	url string
}

type BlockList struct {
	Hash  string `json:"blockHash"`
	Index string `json:"blockIndex"`
}
type TopBlock struct {
	Number      string      `json:"number"`
	Hash        string      `json:"hash"`
	Header      string      `json:"header"`
	BlockType   string      `json:"blockType"`
	ChainBits   string      `json:"chainBits"`
	RelatedList []BlockList `json:"blockList"`
}

const (
	getTopRelayBlockByNumber = "topRelay_getBlockByNumber"
	getTopRelayBlockNumber   = "topRelay_blockNumber"
	getTopBalance            = "top_getBalance"
)

func NewTopSdk(url string) (*TopSdk, error) {
	sdk, err := sdk.NewSDK(url)
	if err != nil {
		return nil, err
	}
	return &TopSdk{SDK: sdk, url: url}, nil
}

func (t *TopSdk) GetTopElectBlockHeadByHeight(height uint64) (*TopBlock, error) {
	var data json.RawMessage
	err := t.Rpc.CallContext(context.Background(), &data, getTopRelayBlockByNumber, hexutil.EncodeUint64(height), false, "transaction")
	if err != nil {
		return &TopBlock{}, err
	} else if len(data) == 0 {
		return &TopBlock{}, ethereum.NotFound
	}

	block := new(TopBlock)
	if err := json.Unmarshal(data, &block); err != nil {
		log.Printf("Unmarshal GetTopElectBlockHeadByHeight data: %v,error:%v", data, err)
		return &TopBlock{}, err
	}
	logger.Debug("Top block: %v, type: %v", block.Number, block.BlockType)
	return block, nil
}

func (t *TopSdk) GetLatestTopElectBlockHeight() (uint64, error) {
	var data json.RawMessage
	err := t.Rpc.CallContext(context.Background(), &data, getTopRelayBlockNumber)
	if err != nil {
		return 0, err
	} else if len(data) == 0 {
		return 0, ethereum.NotFound
	}

	//var res string
	var res string
	if err := json.Unmarshal(data, &res); err != nil {
		logger.Error("sdk getLatestTopElectBlockHeight data: %v,error:%v", string(data), err)
		return 0, err
	}
	return hexutil.DecodeUint64(res)
}

func (t *TopSdk) TopBalanceAt(ctx context.Context, account common.Address) (*big.Int, error) {
	var result hexutil.Big
	err := t.Rpc.CallContext(ctx, &result, getTopBalance, account)
	return (*big.Int)(&result), err
}
