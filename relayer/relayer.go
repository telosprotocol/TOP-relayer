package relayer

import (
	"errors"
	"fmt"
	"sync"

	"toprelayer/config"
	"toprelayer/relayer/crosschainrelayer"
	"toprelayer/relayer/monitor"
	"toprelayer/relayer/toprelayer"

	"github.com/wonderivan/logger"
)

var (
	topRelayers = map[string]IChainRelayer{
		config.ETH_CHAIN:     new(toprelayer.Eth2TopRelayerV2),
		config.BSC_CHAIN:     new(toprelayer.Bsc2TopRelayer),
		config.HECO_CHAIN:    new(toprelayer.Heco2TopRelayer),
		config.OPEN_ALLIANCE: new(toprelayer.OpenAlliance2TopRelayer)}

	crossChainRelayer = new(crosschainrelayer.CrossChainRelayer)
)

type IChainRelayer interface {
	Init(cfg *config.Relayer, listenUrl []string, pass string) error
	StartRelayer(*sync.WaitGroup) error
	GetInitData() ([]byte, error)
}

type ICrossChainRelayer interface {
	Init(chainName string, cfg *config.Relayer, listenUrl []string, pass string, server config.Server) error
	StartRelayer(*sync.WaitGroup) error
}

func startTopRelayer(relayer IChainRelayer, cfg *config.Relayer, listenUrl []string, pass string, wg *sync.WaitGroup) error {
	err := relayer.Init(cfg, listenUrl, pass)
	if err != nil {
		logger.Error("startTopRelayer error:", err)
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

func startCrossChainRelayer(relayer ICrossChainRelayer, chainName string, cfg *config.Relayer, listenUrl []string, pass string, server config.Server, wg *sync.WaitGroup) error {
	err := relayer.Init(chainName, cfg, listenUrl, pass, server)
	if err != nil {
		logger.Error("startCrossChainRelayer error:", err)
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

func StartRelayer(cfg *config.Config, pass string, wg *sync.WaitGroup) error {
	// start monitor
	err := monitor.MonitorMsgInit(cfg.RelayerToRun)
	if err != nil {
		logger.Error("MonitorMsgInit fail:", err)
		return err
	}

	// start relayer
	topConfig, exist := cfg.RelayerConfig[config.TOP_CHAIN]
	if !exist {
		return fmt.Errorf("not found TOP chain config")
	}
	RelayerConfig, exist := cfg.RelayerConfig[cfg.RelayerToRun]
	if !exist {
		return fmt.Errorf("not found config of RelayerToRun")
	}
	if cfg.RelayerToRun == config.TOP_CHAIN {
		for name, c := range cfg.RelayerConfig {
			logger.Info("name: ", name)
			if name == config.TOP_CHAIN {
				continue
			}
			if name != config.ETH_CHAIN && name != config.BSC_CHAIN && name != config.HECO_CHAIN && name != config.OPEN_ALLIANCE {
				logger.Warn("TopRelayer not support:", name)
				continue
			}
			topRelayer, exist := topRelayers[name]
			if !exist {
				logger.Warn("unknown chain config:", name)
				continue
			}
			err := startTopRelayer(topRelayer, topConfig, c.Url, pass, wg)
			if err != nil {
				logger.Error("StartRelayer %v error: %v", name, err)
				continue
			}
		}
	} else {
		err := startCrossChainRelayer(crossChainRelayer, cfg.RelayerToRun, RelayerConfig, topConfig.Url, pass, cfg.ServerConfig, wg)
		if err != nil {
			logger.Error("StartRelayer error:", err)
			return err
		}
	}

	return nil
}

func GetInitData(cfg *config.Config, pass, chainName string) ([]byte, error) {
	if cfg.RelayerToRun != config.TOP_CHAIN {
		err := errors.New("RelayerToRun error")
		logger.Error(err)
		return nil, err
	}
	if chainName != config.ETH_CHAIN {
		err := errors.New("chain not support init data")
		logger.Error(err)
		return nil, err
	}
	c, exist := cfg.RelayerConfig[chainName]
	if !exist {
		err := errors.New("not found chain config")
		logger.Error(err)
		return nil, err
	}
	topRelayer, exist := topRelayers[chainName]
	if !exist {
		err := errors.New("not found chain relayer")
		logger.Error(err)
		return nil, err
	}
	err := topRelayer.Init(c, c.Url, pass)
	if err != nil {
		logger.Error("Init error:", err)
		return nil, err
	}
	return topRelayer.GetInitData()
}
