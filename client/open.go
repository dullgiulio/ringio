package client

import (
	"flag"

	"github.com/dullgiulio/ringio/config"
	"github.com/dullgiulio/ringio/server"
)

type CommandOpen struct {
}

func NewCommandOpen() *CommandOpen {
	return &CommandOpen{}
}

func (c *CommandOpen) Help() string {
	return `Open a new session`
}

func (c *CommandOpen) Init(fs *flag.FlagSet) bool {
	fs.Int64Var(&config.C.RingbufSize, "ringbuf-size", config.C.RingbufSize, "Max number of lines contained in a ringbuffer")
	fs.Int64Var(&config.C.MaxLineSize, "line-size", config.C.MaxLineSize, "Max size in bytes of a single line read from stdin")
	fs.Int64Var(&config.C.RingbufLogSize, "ringbuf-log-size", config.C.RingbufSize, "Max number of lines for log scrollback")
	fs.BoolVar(&config.C.AutoRun, "autorun", config.C.AutoRun, "Automatically run new commands when added")
	fs.BoolVar(&config.C.PrintLog, "verbose", config.C.PrintLog, "Print logging on the command line")
	return true
}

func (c *CommandOpen) Run(cli *Cli) error {
	server.Init()
	server.Run(cli.Session)

	return nil
}
