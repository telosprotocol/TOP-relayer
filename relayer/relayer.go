package relayer

import (
	"sync"
	"time"

	"toprelayer/config"
	"toprelayer/relayer/eth2top"
	"toprelayer/relayer/top2eth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wonderivan/logger"
)

var relayerMap = map[string]IChainRelayer{config.TOP_CHAIN: new(eth2top.Eth2TopRelayer), config.ETH_CHAIN: new(top2eth.Top2EthRelayer)}

type IChainRelayer interface {
	Init(fromUrl, toUrl, keypath, pass string, chainid uint64, contract common.Address) error
	StartRelayer(*sync.WaitGroup) error
	ChainId() uint64
}

func StartRelayer(cfg *config.Config, chainpass map[uint64]string, wg *sync.WaitGroup) (err error) {
	for _, chain := range cfg.RelayerConfig {
		chainName := chain.Chain
		_, exist := relayerMap[chainName]
		if !exist {
			logger.Warn("unknown chain config: %v", chainName)
			continue
		}
		if chain.Start {
			err := relayerMap[chainName].Init(
				chain.SubmitUrl,
				chain.ListenUrl,
				chain.KeyPath,
				chainpass[chain.ChainId],
				chain.ChainId,
				common.HexToAddress(chain.Contract),
			)
			if err != nil {
				return err
			}

			wg.Add(1)
			go func() {
				err = relayerMap[chainName].StartRelayer(wg)
			}()
			if err != nil {
				return err
			}
			time.Sleep(time.Second * 1)
		}
	}

	return nil
}
