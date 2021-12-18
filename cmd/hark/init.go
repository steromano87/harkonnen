package main

import (
	"fmt"
	"github.com/maruel/subcommands"
)

var cmdInit = &subcommands.Command{
	UsageLine: "init",
	ShortDesc: "initializes a Harkonnen project folder",
	LongDesc: "Initializes a Harkonnen project folder. " +
		"If no folder is specified, the current folder will be used as project root",
	Advanced: false,
	CommandRun: func() subcommands.CommandRun {
		return &initRun{}
	},
}

type initRun struct {
	subcommands.CommandRunBase
}

func (ir initRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	fmt.Println("To be implemented!")
	return 1
}
