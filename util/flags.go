package util

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"syscall"
	"toprelayer/config"

	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	// relayer config
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

func MakePassword(ctx *cli.Context, cfg *config.Config) (string, error) {
	path := ctx.String(PasswordFileFlag.Name)
	if path == "" {
		return ReadPassword(cfg)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Failed to read password file:", err)
		return string(""), err
	}
	// passes := make(map[string]string)
	// err = json.Unmarshal(data, &passes)
	// if err != nil {
	// 	log.Fatal("Umarshal password file failed:", err)
	// 	return nil, err
	// }
	lines := strings.Split(string(data), "\n")
	// Sanitise DOS line endings.
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], "\r")
	}
	return lines[0], err
}

func ReadPassword(cfg *config.Config) (string, error) {
	fmt.Print(">>> Please Enter " + cfg.RelayerToRun + " pasword:\n>>> ")

	var passwd string
	if terminal.IsTerminal(int(syscall.Stdin)) {
		pass, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return string(pass), err
		}
		passwd = string(pass)
	} else {
		var b [1]byte
		var pw []byte
		for {
			n, err := os.Stdin.Read(b[:])
			// terminal.ReadPassword discards any '\r', so we do the same
			if n > 0 && b[0] != '\r' {
				if b[0] == '\n' {
					return string(pw), nil
				}
				pw = append(pw, b[0])
				// limit size, so that a wrong input won't fill up the memory
				if len(pw) > 1024 {
					err = errors.New("password too long")
				}
			}
			if err != nil {
				// terminal.ReadPassword accepts EOF-terminated passwords
				// if non-empty, so we do the same
				if err == io.EOF && len(pw) > 0 {
					err = nil
				}
				return string(pw), err
			}
		}
	}
	fmt.Println()
	return passwd, nil
}
