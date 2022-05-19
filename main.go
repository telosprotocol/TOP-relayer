package main

import (
	"os"
	"sync"
	"toprelayer/base"

	"toprelayer/config"
	"toprelayer/relayer"

	"github.com/urfave/cli/v2"
	"github.com/wonderivan/logger"
)

func main() {
	logger.SetLogger("./log/logconfig.json")

	app := &cli.App{
		Name:   "xrelayer",
		Usage:  "block chain relayer",
		Action: start,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "./config/config.json",
				Usage: "configuration file",
			},
			&cli.StringFlag{
				Name:  "ethpass",
				Value: "",
				Usage: "eth relayer keystore pass word",
			},
			&cli.StringFlag{
				Name:  "toppass",
				Value: "",
				Usage: "top relayer keystore pass word",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal("Run relayer error:", err)
	}
}

func start(c *cli.Context) error {
	wg := new(sync.WaitGroup)
	handlercfg, err := config.InitHeaderSyncConfig(c.String("config"))
	if err != nil {
		return err
	}

	err = relayer.StartRelayer(wg, handlercfg, getchainpass(c, handlercfg))
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}

func getchainpass(c *cli.Context, handlercfg *config.HeaderSyncConfig) map[uint64]string {
	chainpass := make(map[uint64]string)
	for _, chain := range handlercfg.Config.Chains {
		switch chain.SubmitChainId {
		case base.ETH:
			chainpass[base.ETH] = c.String("ethpass")
		case base.TOP:
			chainpass[base.TOP] = c.String("toppass")
		}
	}
	return chainpass
}
