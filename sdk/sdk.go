package sdk

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type SDK struct {
	Rpc *rpc.Client
	*ethclient.Client
}

func NewSDK(url string) (*SDK, error) {
	rpcClient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}

	rawClient, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}

	return &SDK{
		Rpc:    rpcClient,
		Client: rawClient,
	}, nil
}
