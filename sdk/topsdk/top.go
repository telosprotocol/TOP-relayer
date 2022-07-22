package topsdk

import (
	"context"
	"encoding/json"
	"log"
	"toprelayer/sdk"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/wonderivan/logger"
)

type TopSdk struct {
	*sdk.SDK
	url string
}

type TopBlock struct {
	Number    string `json:"number"`
	Header    string `json:"header"`
	BlockType string `json:"blockType"`
	ChainBits string `json:"chainBits"`
}

const (
	GETTOPELECTBLOCKHEADBYHEIGHT = "topRelay_getBlockByNumber"
	GETLATESTTOPELECTBLOCKHEIGHT = "topRelay_blockNumber"
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
	err := t.Rpc.CallContext(context.Background(), &data, GETTOPELECTBLOCKHEADBYHEIGHT, hexutil.EncodeUint64(height))
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
	err := t.Rpc.CallContext(context.Background(), &data, GETLATESTTOPELECTBLOCKHEIGHT)
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
