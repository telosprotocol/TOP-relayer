package topsdk

import (
	"context"
	"encoding/json"
	"log"
	"toprelayer/sdk"
	"toprelayer/util"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/wonderivan/logger"
)

type TopSdk struct {
	*sdk.SDK
	url string
}

type TopBlock struct {
	BlockType string `json:"blockType"`
	Number    string `json:"number"`
	Header    string `json:"header"`
}

const (
	GETTOPELECTBLOCKHEADBYHEIGHT = "topRelay_getBlockByNumber"
	GETLATESTTOPELECTBLOCKHEIGHT = "topRelay_blockNumber"

	ELECTION_BLOCK    = "election"
	AGGREGATE_BLOCK   = "aggregate"
	TRANSACTION_BLOCK = "transactions"
)

func NewTopSdk(url string) (*TopSdk, error) {
	sdk, err := sdk.NewSDK(url)
	if err != nil {
		return nil, err
	}
	return &TopSdk{SDK: sdk, url: url}, nil
}

func (t *TopSdk) GetTopElectBlockHeadByHeight(height uint64) (bytes []byte, flag bool, err error) {
	var data json.RawMessage
	err = t.Rpc.CallContext(context.Background(), &data, GETTOPELECTBLOCKHEADBYHEIGHT, util.Uint64ToHexString(height))
	if err != nil {
		return []byte{}, false, err
	} else if len(data) == 0 {
		return []byte{}, false, ethereum.NotFound
	}

	var block TopBlock
	if err := json.Unmarshal(data, &block); err != nil {
		log.Printf("Unmarshal GetTopElectBlockHeadByHeight data: %v,error:%v", data, err)
		return []byte{}, false, err
	}
	logger.Debug("Top block: %v, type: %v", block.Number, block.BlockType)

	bytes = common.Hex2Bytes(block.Header[2:])
	if block.BlockType == ELECTION_BLOCK || block.BlockType == AGGREGATE_BLOCK {
		flag = true
	} else {
		flag = false
	}

	return bytes, flag, nil
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
	return util.HexToUint64(res)
}
