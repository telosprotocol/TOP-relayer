package top_rpc

import (
	"context"
	"math/big"
	"toprelayer/util"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

type TopRpc struct {
	rpc *rpc.Client
}

func NewTopRpc(url string) (*TopRpc, error) {
	rpcclient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	return &TopRpc{
		rpc: rpcclient,
	}, nil
}

func (w *TopRpc) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (balance *big.Int, err error) {
	var result hexutil.Big
	err = w.rpc.CallContext(ctx, &result, "top_getBalance", account)
	return (*big.Int)(&result), err
}

func (c *TopRpc) BlockNumber(ctx context.Context) (uint64, error) {
	var result hexutil.Uint64
	err := c.rpc.CallContext(ctx, &result, "topRelay_blockNumber")
	return uint64(result), err
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}

func (c *TopRpc) HeaderByNumber(ctx context.Context, number *big.Int) (*util.TopHeader, error) {
	var head *util.TopHeader
	err := c.rpc.CallContext(ctx, &head, "topRelay_getBlockByNumber", toBlockNumArg(number))
	if err == nil && head == nil {
		err = ethereum.NotFound
	}
	return head, err
}
