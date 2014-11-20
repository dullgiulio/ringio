package client

import (
	"flag"

	"bitbucket.org/dullgiulio/ringio/server"
)

type CommandOpen struct {
}

func NewCommandOpen() *CommandOpen {
	return &CommandOpen{}
}

func (c *CommandOpen) Help() string {
	return `Run all processes`
}

func (c *CommandOpen) Init(fs *flag.FlagSet) error {
	// nothing to do yet.
	return nil
}

func (c *CommandOpen) Run(cli *Cli) error {
	server.Init()
	server.Run(cli.Session)

	return nil
}
