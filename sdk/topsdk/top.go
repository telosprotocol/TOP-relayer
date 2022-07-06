package topsdk

import (
	"context"
	"encoding/json"
	"log"
	"toprelayer/sdk"
	"toprelayer/util"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

type TopSdk struct {
	*sdk.SDK
	url string
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

func (t *TopSdk) GetTopElectBlockHeadByHeight(height uint64) ([]byte, error) {
	var data []string
	err := t.Rpc.CallContext(context.Background(), &data, GETTOPELECTBLOCKHEADBYHEIGHT, util.Uint64ToHexString(height))
	if err != nil {
		return []byte{}, err
	} else if len(data) == 0 {
		return []byte{}, ethereum.NotFound
	}
	bytes := common.Hex2Bytes(data[0][2:])
	return bytes, nil
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
		log.Printf("sdk getLatestTopElectBlockHeight data: %v,error:%v", string(data), err)
		return 0, err
	}
	return util.HexToUint64(res)
}
