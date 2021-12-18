package main

import (
	"github.com/maruel/subcommands"
	"os"
)

var application = &subcommands.DefaultApplication{
	Name:  "hark",
	Title: "Multi-protocol load testing tool",
	Commands: []*subcommands.Command{
		cmdInit,
		subcommands.CmdHelp,
		cmdVersion,
	},
}

func main() {
	os.Exit(subcommands.Run(application, nil))
}
