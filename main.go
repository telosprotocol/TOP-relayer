package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"toprelayer/config"
	"toprelayer/relayer"
	"toprelayer/util"

	"github.com/urfave/cli/v2"
)

var (
	app = cli.NewApp()

	nodeFlags = []cli.Flag{
		&util.PasswordFileFlag,
		&util.ConfigFileFlag,
	}
)

func init() {
	app.Name = filepath.Base(os.Args[0])
	app.Usage = "the TOP-relayer command line interface"
	app.Copyright = "2017-present Telos Foundation & contributors"
	app.Action = start
	app.Flags = nodeFlags
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal("Run relayer error:", err)
		os.Exit(1)
	}
}

func start(ctx *cli.Context) error {
	cfg, err := config.LoadRelayerConfig(ctx.String("config"))
	if err != nil {
		return err
	}

	pass, err := util.MakePassword(ctx, cfg)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("starting relayer...")

	err = config.InitLogConfig()
	if err != nil {
		return err
	}

	wg := new(sync.WaitGroup)
	err = relayer.StartRelayer(cfg, pass, wg)
	if err != nil {
		return err
	}

	wg.Wait()
	return nil
}
