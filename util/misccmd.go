package util

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"

	"toprelayer/config"
	"toprelayer/relayer"
	"toprelayer/version"
)

const (
	clientIdentifier = "xrelayer"
)

func versionPrint(ctx *cli.Context) error {
	fmt.Println(clientIdentifier)
	fmt.Println("Version:", version.VersionWithMeta)
	fmt.Println("Architecture:", runtime.GOARCH)
	fmt.Println("Operating System:", runtime.GOOS)
	return nil
}

func getInitData(ctx *cli.Context) error {
	if ctx.Args().Len() != 1 {
		return errors.New("need chain_name as the only argument")
	}
	chainName := ctx.Args().First()
	if len(chainName) == 0 {
		return errors.New("invalid chain_name")
	}

	cfg, err := config.LoadRelayerConfig(ctx.String("config"))
	if err != nil {
		return err
	}
	pass, err := MakePassword(ctx, cfg)
	if err != nil {
		return err
	}
	bytes, err := relayer.GetInitData(cfg, pass, chainName)
	if err != nil {
		return err
	}
	fmt.Println(common.Bytes2Hex(bytes))
	return nil
}

var (
	VersionCommand = &cli.Command{
		Action:    versionPrint,
		Name:      "version",
		Usage:     "Print version numbers",
		ArgsUsage: " ",
		Category:  "MISCELLANEOUS COMMANDS",
		Description: `
The output of this command is supposed to be machine-readable.
`,
	}
	GetInitDataCommand = &cli.Command{
		Action:    getInitData,
		Name:      "get_init_data",
		Usage:     "Print init hex data",
		ArgsUsage: "<chain_name>",
		Category:  "MISCELLANEOUS COMMANDS",
		Description: `
The output of this command is hex data.
`,
	}
)
