package util

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"

	"toprelayer/config"
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
	argsNum := ctx.Args().Len()
	if argsNum <= 1 {
		return errors.New("invalid args")
	}
	chainName := ctx.Args().First()

	var bytes []byte
	var err error
	if chainName == config.ETH_CHAIN {
		if argsNum == 4 {
			bytes, err = getEthInitData(ctx.Args().Get(1), ctx.Args().Get(2), ctx.Args().Get(3))
		} else if argsNum == 5 {
			bytes, err = getEthInitDataWithHeight(ctx.Args().Get(1), ctx.Args().Get(2), ctx.Args().Get(3), ctx.Args().Get(4))
		} else {
			return errors.New("invalid arg nums")
		}
	} else if chainName == config.BSC_CHAIN {
		if argsNum == 2 {
			bytes, err = getBscInitData(ctx.Args().Get(1))
		} else if argsNum == 3 {
			bytes, err = getBscInitDataWithHeight(ctx.Args().Get(1), ctx.Args().Get(2))
		} else {
			return errors.New("invalid arg nums")
		}
	} else if chainName == config.HECO_CHAIN {
		if argsNum == 2 {
			bytes, err = getHecoInitData(ctx.Args().Get(1))
		} else if argsNum == 3 {
			bytes, err = getHecoInitDataWithHeight(ctx.Args().Get(1), ctx.Args().Get(2))
		} else {
			return errors.New("invalid arg nums")
		}
	} else {
		return errors.New("invalid chain_name")
	}
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
		ArgsUsage: "<chain_name> <urls...> [height]",
		Category:  "MISCELLANEOUS COMMANDS",
		Description: `
The output of this command is hex data.
`,
	}
)
