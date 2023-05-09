package util

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"

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
	if argsNum < 1 {
		return errors.New("invalid args")
	}
	chainName := ctx.Args().First()

	var bytes []byte
	var err error
	switch chainName {
	case config.ETH_CHAIN:
		switch argsNum {
		case 1: // get_init_data ETH
			bytes, err = getEthInitData(config.ETHAddr, config.ETHPrysm)
		case 2: // get_init_data ETH  123
			if false == isNumber(ctx.Args().Get(1)) {
				return errors.New("the height needs to be a number")
			}
			bytes, err = getEthInitDataWithHeight(config.ETHAddr, config.ETHPrysm, ctx.Args().Get(1))
		default:
			return errors.New("invalid arg nums")
		}
	case config.BSC_CHAIN:
		switch argsNum {
		case 1: // get_init_data BSC
			bytes, err = getBscInitData(config.BSCAddr)
		case 2: // get_init_data BSC 123
			if false == isNumber(ctx.Args().Get(1)) {
				return errors.New("the height needs to be a number")
			}
			bytes, err = getBscInitDataWithHeight(config.BSCAddr, ctx.Args().Get(1))
		default:
			return errors.New("invalid arg nums")
		}
	case config.HECO_CHAIN:
		switch argsNum {
		case 1: // get_init_data HECO
			bytes, err = getHecoInitData(config.HECOAddr)
		case 2: // get_init_data HECO 123
			if false == isNumber(ctx.Args().Get(1)) {
				return errors.New("the height needs to be a number")
			}
			bytes, err = getHecoInitDataWithHeight(config.HECOAddr, ctx.Args().Get(1))
		default:
			return errors.New("invalid arg nums")
		}
	default:
		return fmt.Errorf("the %s chain is not supported", chainName)
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

func isNumber(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}
