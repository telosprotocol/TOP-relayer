package util

import (
	"encoding/json"
	"log"
	"os"
	"toprelayer/config"

	"github.com/urfave/cli/v2"
)

var ( // relayer config
	ConfigFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "configuration file",
		Value: "./config/relayerconfig.json",
	}
	// user config
	PasswordFileFlag = cli.StringFlag{
		Name:  "password",
		Usage: "Password file to use for non-interactive password input",
		Value: "",
	}
)

func MakePasswordList(ctx *cli.Context, cfg *config.Config) (map[string]string, error) {
	path := ctx.String(PasswordFileFlag.Name)
	if path == "" {
		return GetPasses(cfg)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Failed to read password file:", err)
		return nil, err
	}
	passes := make(map[string]string)
	err = json.Unmarshal(data, &passes)
	if err != nil {
		log.Fatal("Umarshal password file failed:", err)
		return nil, err
	}
	// lines := strings.Split(string(text), "\n")
	// // Sanitise DOS line endings.
	// for i := range lines {
	// 	lines[i] = strings.TrimRight(lines[i], "\r")
	// }
	return passes, err
}
