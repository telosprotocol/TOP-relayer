package relayer

import (
	"sync"

	"toprelayer/base"
	"toprelayer/config"
	"toprelayer/relayer/eth2top"
	"toprelayer/relayer/top2eth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wonderivan/logger"
)

type IChainRelayer interface {
	Init(fromUrl, toUrl, keypath, pass string, chainid uint64, contract common.Address, batch, cert int, verify bool) error
	StartRelayer(*sync.WaitGroup) error
	ChainId() uint64
}

func StartRelayer(wg *sync.WaitGroup, handlercfg *config.HeaderSyncConfig, chainpass map[uint64]string) (err error) {
	handler := NewHeaderSyncHandler(handlercfg)
	err = handler.Init(wg, chainpass)
	if err != nil {
		return err
	}
	return handler.StartRelayer()
}

func GetRelayer(chain uint64) (relayer IChainRelayer) {
	switch chain {
	case base.ETH:
		relayer = new(top2eth.Top2EthRelayer)
	case base.TOP:
		relayer = new(eth2top.Eth2TopRelayer)
	default:
		logger.Error("Unsupport chain id:", chain)
	}
	return
}
