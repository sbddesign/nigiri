package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

var rpc = cli.Command{
	Name:   "rpc",
	Usage:  "invoke bitcoin-cli or elements-cli",
	Action: rpcAction,
	Flags: []cli.Flag{
		&liquidFlag,
		&cli.StringFlag{
			Name:  "rpcwallet",
			Usage: "rpcwallet to be used for node JSONRPC commands",
			Value: "",
		},
		&cli.IntFlag{
			Name:  "generate",
			Usage: "generate block",
			Value: 0,
		},
	},
}

func rpcAction(ctx *cli.Context) error {

	if isRunning, _ := nigiriState.GetBool("running"); !isRunning {
		return errors.New("nigiri is not running")
	}

	isLiquid := ctx.Bool("liquid")
	rpcWallet := ctx.String("rpcwallet")
	generate := ctx.Int("generate")

	rpcArgs := []string{"exec", "bitcoin", "bitcoin-cli", "-datadir=config", "-rpcwallet=" + rpcWallet}
	if isLiquid {
		rpcArgs = []string{"exec", "liquid", "elements-cli", "-datadir=config", "-rpcwallet=" + rpcWallet}
	}
	if generate > 0 {
		rpcArgs = append(rpcArgs, "-generate="+fmt.Sprint(generate))
	}
	cmdArgs := append(rpcArgs, ctx.Args().Slice()...)
	bashCmd := exec.Command("docker", cmdArgs...)
	bashCmd.Stdout = os.Stdout
	bashCmd.Stderr = os.Stderr

	if err := bashCmd.Run(); err != nil {
		return err
	}

	return nil
}
