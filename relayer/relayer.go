package relayer

import (
	"errors"
	"sync"

	"toprelayer/config"
	"toprelayer/relayer/eth2top"
	"toprelayer/relayer/top2eth"

	"github.com/wonderivan/logger"
)

var (
	topRelayers = map[string]IChainRelayer{
		config.ETH_CHAIN:  new(eth2top.Eth2TopRelayer),
		config.BSC_CHAIN:  new(eth2top.Eth2TopRelayer),
		config.HECO_CHAIN: new(eth2top.Eth2TopRelayer)}

	crossChainRelayer = new(top2eth.Top2EthRelayer)
)

type IChainRelayer interface {
	Init(chainName string, cfg *config.Relayer, listenUrl string, pass string) error
	StartRelayer(*sync.WaitGroup) error
	ChainId() uint64
}

func startOneRelayer(chainName string, relayer IChainRelayer, cfg *config.Relayer, listenUrl string, pass string, wg *sync.WaitGroup) error {
	err := relayer.Init(chainName, cfg, listenUrl, pass)
	if err != nil {
		logger.Error("startOneRelayer error:", err)
		return err
	}

	wg.Add(1)
	go func() {
		err = relayer.StartRelayer(wg)
	}()
	if err != nil {
		logger.Error("relayer.StartRelayer error:", err)
		return err
	}
	return nil
}

func StartRelayer(cfg *config.Config, pass string, wg *sync.WaitGroup) (err error) {
	topConfig, exist := cfg.RelayerConfig[config.TOP_CHAIN]
	if !exist {
		return errors.New("not found TOP chain config")
	}
	RelayerConfig, exist := cfg.RelayerConfig[cfg.RelayerToRun]
	if !exist {
		return errors.New("not found config of RelayerToRun")
	}
	if cfg.RelayerToRun == config.TOP_CHAIN {
		for name, c := range cfg.RelayerConfig {
			if name == config.TOP_CHAIN {
				continue
			}
			if name != config.ETH_CHAIN {
				logger.Warn("TopRelayer not support:", name)
				continue
			}
			topRelayer, exist := topRelayers[name]
			if !exist {
				logger.Warn("unknown chain config:", name)
				continue
			}
			err := startOneRelayer(name, topRelayer, topConfig, c.Url, pass, wg)
			if err != nil {
				logger.Error("StartRelayer error:", err)
				return err
			}
		}
	} else {
		err := startOneRelayer(cfg.RelayerToRun, crossChainRelayer, RelayerConfig, topConfig.Url, pass, wg)
		if err != nil {
			logger.Error("StartRelayer error:", err)
			return err
		}
	}

	return nil
}
