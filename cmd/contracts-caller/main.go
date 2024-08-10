package main

import (
	"fmt"

	"os"

	"github.com/urfave/cli"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"

	contracts_caller "github.com/the-web3/contracts-caller"
	"github.com/the-web3/contracts-caller/flags"
)

var (
	GitVersion = ""
	GitCommit  = ""
	GitDate    = ""
)

func main() {
	log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stderr, log.LevelInfo, true)))
	app := cli.NewApp()
	app.Flags = flags.Flags
	app.Version = fmt.Sprintf("%s-%s", GitVersion, params.VersionWithCommit(GitCommit, GitDate))
	app.Name = "contracts-caller"
	app.Usage = "Contracts caller template project"
	app.Description = "Contracts caller template service, every one can develop self contracts caller base this project"
	app.Action = contracts_caller.Main(GitVersion)
	err := app.Run(os.Args)
	if err != nil {
		log.Crit("Contracts Caller Application failed", "message", err)
	}
}
