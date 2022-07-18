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
	topRelayer         = new(eth2top.Eth2TopRelayer)
	crossChainRelayers = map[string]IChainRelayer{
		config.ETH_CHAIN:  new(top2eth.Top2EthRelayer),
		config.BSC_CHAIN:  new(top2eth.Top2EthRelayer),
		config.HECO_CHAIN: new(top2eth.Top2EthRelayer)}
)

type IChainRelayer interface {
	Init(cfg *config.Relayer, toUrl map[string]string, pass string) error
	StartRelayer(*sync.WaitGroup) error
	ChainId() uint64
}

func startOneRelayer(relayer IChainRelayer, cfg *config.Relayer, listenUrls map[string]string, pass string, wg *sync.WaitGroup) error {
	err := relayer.Init(cfg, listenUrls, pass)
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

func StartRelayer(cfg *config.Config, chainpass map[string]string, wg *sync.WaitGroup) (err error) {
	topConfig, exist := cfg.RelayerConfig[config.TOP_CHAIN]
	if !exist {
		return errors.New("not found TOP chain config")
	}
	topListenUrls := make(map[string]string)
	for name, cfg := range cfg.RelayerConfig {
		if name == config.TOP_CHAIN {
			continue
		}
		crossChainRelayer, exist := crossChainRelayers[name]
		if !exist {
			logger.Warn("unknown chain config:", name)
			continue
		}
		if cfg.Start {
			err := startOneRelayer(crossChainRelayer, cfg, topListenUrls, chainpass[cfg.Chain], wg)
			if err != nil {
				logger.Error("StartRelayer error:", err)
				return err
			}
		}
		topListenUrls[name] = cfg.SubmitUrl
	}
	if topConfig.Start {
		err := startOneRelayer(topRelayer, topConfig, topListenUrls, chainpass[topConfig.Chain], wg)
		if err != nil {
			logger.Error("StartRelayer error:", err)
			return err
		}
	}

	return nil
}
