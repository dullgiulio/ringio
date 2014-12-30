package client

import (
	"flag"
	"net/rpc"

	"github.com/dullgiulio/ringio/agents"
	"github.com/dullgiulio/ringio/server"
	"github.com/dullgiulio/ringio/utils"
)

type CommandOutput struct {
	client   *rpc.Client
	response *server.RpcResp
}

func NewCommandOutput() *CommandOutput {
	return &CommandOutput{
		response: new(server.RpcResp),
	}
}

func (c *CommandOutput) Help() string {
	return `Output data from the ringbuf`
}

func (c *CommandOutput) Init(fs *flag.FlagSet) bool {
	// nothing to do yet.
	return false
}

func (c *CommandOutput) Run(cli *Cli) error {
	c.client = cli.GetClient()

	meta := &agents.AgentMetadata{Filter: cli.Filter}

	if len(cli.Args) == 0 {
		addSinkAgentPipe(c.client, meta, c.response, utils.GetRandomDotfile())
	} else {
		addSinkAgentCmd(c.client, meta, c.response, cli.Args)
	}

	return nil
}
