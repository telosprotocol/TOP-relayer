package openalliance_rpc

import (
	"context"
	"math/big"
	"toprelayer/util"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

type OpenAllianceRpc struct {
	rpc *rpc.Client
}

func NewOpenAllianceRpc(url string) (*OpenAllianceRpc, error) {
	rpcclient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	return &OpenAllianceRpc{
		rpc: rpcclient,
	}, nil
}

func (c *OpenAllianceRpc) BlockNumber(ctx context.Context) (uint64, error) {
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

func (c *OpenAllianceRpc) HeaderByNumber(ctx context.Context, number *big.Int) (*util.TopHeader, error) {
	var head *util.TopHeader
	err := c.rpc.CallContext(ctx, &head, "topRelay_getBlockByNumber", toBlockNumArg(number))
	if err == nil && head == nil {
		err = ethereum.NotFound
	}
	return head, err
}
