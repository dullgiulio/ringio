package client

import (
	"flag"
	"net/rpc"

	"github.com/dullgiulio/ringio/server"
	"github.com/dullgiulio/ringio/utils"
)

type CommandError struct {
	client   *rpc.Client
	response *server.RpcResp
}

func NewCommandError() *CommandError {
	return &CommandError{
		response: new(server.RpcResp),
	}
}

func (c *CommandError) Help() string {
	return `Output data gathered from the standard error (stderr)`
}

func (c *CommandError) Init(fs *flag.FlagSet) bool {
	// nothing to do yet.
	return false
}

func (c *CommandError) Run(cli *Cli) error {
	c.client = cli.GetClient()

	addErrorsAgentPipe(c.client, cli.Filter, c.response, utils.GetRandomDotfile())

	return nil
}
