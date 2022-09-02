package util

import (
	"fmt"
	"runtime"

	"github.com/urfave/cli/v2"

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
)
