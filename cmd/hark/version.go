package main

import (
	"fmt"
	"github.com/maruel/subcommands"
)

var Version = "master"

var cmdVersion = &subcommands.Command{
	UsageLine: "version",
	ShortDesc: "shows the version and exits",
	LongDesc:  "Shows the version and exits",
	CommandRun: func() subcommands.CommandRun {
		return &versionRun{}
	},
}

type versionRun struct {
	subcommands.CommandRunBase
}

func (vr *versionRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	fmt.Printf("%s\n", Version)
	return 0
}
