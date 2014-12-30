package client

import (
	"flag"
	"net/rpc"

	"github.com/dullgiulio/ringio/agents"
	"github.com/dullgiulio/ringio/server"
	"github.com/dullgiulio/ringio/utils"
)

type CommandIO struct {
	client   *rpc.Client
	response *server.RpcResp
}

func NewCommandIO() *CommandIO {
	return &CommandIO{
		response: new(server.RpcResp),
	}
}

func (c *CommandIO) Help() string {
	return `Open an input and an ouput session to current terminal`
}

func (c *CommandIO) Init(fs *flag.FlagSet) bool {
	// nothing to do yet.
	return false
}

func (c *CommandIO) Run(cli *Cli) error {
	c.client = cli.GetClient()

	go addSourceAgentPipe(c.client, c.response, utils.GetRandomDotfile())
	addSinkAgentPipe(c.client, &agents.AgentMetadata{Filter: cli.Filter}, c.response, utils.GetRandomDotfile())

	return nil
}
