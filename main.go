package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"toprelayer/config"
	"toprelayer/relayer"
	"toprelayer/util"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "xrelayer",
		Usage:  "block chain relayer",
		Action: start,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "./config/relayerconfig.json",
				Usage: "configuration file",
			},
			&cli.StringFlag{
				Name:  "logconfig",
				Value: "./config/logconfig.json",
				Usage: "log config path",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal("Run relayer error:", err)
		os.Exit(1)
	}
}

func start(c *cli.Context) error {
	wg := new(sync.WaitGroup)
	handlercfg, err := config.InitHeaderSyncConfig(c.String("config"))
	if err != nil {
		return err
	}
	passes, err := util.Getchainpass(handlercfg)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("starting relayer...")

	err = config.InitLogConfig()
	if err != nil {
		return err
	}
	err = relayer.StartRelayer(wg, handlercfg, passes)
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}
