package topsdk

import (
	"context"
	"encoding/json"
	"log"
	"toprelayer/sdk"
	"toprelayer/util"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

type TopSdk struct {
	*sdk.SDK
	url string
}

type ElectBlockType = uint

const (
	ElectBlock_Current ElectBlockType = iota
	ElectBlock_Next
)

const (
	GETLATESTETTOPELECTBLOCKHEADER = "getLatestTopElectBlockHeader"
	GETTOPELECTBLOCKHEADBYHEIGHT   = "top_getBlockByNumber"
	GETLATESTTOPELECTBLOCKHEIGHT   = "top_blockNumber"
)

func NewTopSdk(url string) (*TopSdk, error) {
	sdk, err := sdk.NewSDK(url)
	if err != nil {
		return nil, err
	}
	return &TopSdk{SDK: sdk, url: url}, nil
}

func (t *TopSdk) SendBlockHeadTransaction(ctx context.Context, btx *types.Transaction) error {
	return t.SendTransaction(ctx, btx)
}

func (t *TopSdk) GetLatestTopElectBlockHeader() (*types.Block, error) {
	return t.getLatestTopElectBlockHeader()
}

func (t *TopSdk) getLatestTopElectBlockHeader() (*types.Block, error) {
	var data json.RawMessage
	err := t.Rpc.CallContext(context.Background(), &data, GETLATESTETTOPELECTBLOCKHEADER)
	if err != nil {
		return nil, err
	} else if len(data) == 0 {
		return nil, ethereum.NotFound
	}
	// Decode header and transactions.
	var head *types.Block
	if err := json.Unmarshal(data, &head); err != nil {
		return nil, err
	}
	return head, nil
}

func (t *TopSdk) GetTransactionByHash(hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	return t.TransactionByHash(context.Background(), hash)
}

func (t *TopSdk) GetTransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	return t.TransactionReceipt(context.Background(), hash)
}

func (t *TopSdk) GetTopElectBlockHeadByHeight(height uint64) (string, error) {
	return t.getTopElectBlockHeadByHeight(height)
}

func (t *TopSdk) getTopElectBlockHeadByHeight(height uint64) (string, error) {
	var result hexutil.Bytes
	err := t.Rpc.CallContext(context.Background(), &result, GETTOPELECTBLOCKHEADBYHEIGHT, util.Uint64ToHexString(height))
	if err != nil {
		return "", err
	} else if len(result) == 0 {
		return "", ethereum.NotFound
	}
	return result.String(), nil
}

func (t *TopSdk) GetLatestTopElectBlockHeight() (uint64, error) {
	return t.getLatestTopElectBlockHeight()
}

func (t *TopSdk) getLatestTopElectBlockHeight() (uint64, error) {
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
