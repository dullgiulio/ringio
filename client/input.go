package client

import (
	"flag"
	"net/rpc"
	"os"

	"github.com/dullgiulio/ringio/agents"
	"github.com/dullgiulio/ringio/server"
	"github.com/dullgiulio/ringio/utils"
)

type CommandInput struct {
	client   *rpc.Client
	response *server.RpcResp
	name     string
}

func NewCommandInput() *CommandInput {
	return &CommandInput{
		response: new(server.RpcResp),
	}
}

func (c *CommandInput) Help() string {
	return `Reads data from stdin (and stderr for processes) and writes it into the ringbuf`
}

func (c *CommandInput) Init(fs *flag.FlagSet) bool {
	// nothing to do yet.
	return false
}

func (c *CommandInput) Run(cli *Cli) error {
	c.client = cli.GetClient()

	meta := agents.AgentMetadata{
		User: os.Getenv("USER"),
		Name: c.name,
	}

	if len(cli.Args) == 0 {
		addSourceAgentPipe(c.client, c.response, &meta, utils.GetRandomDotfile())
	} else {
		addSourceAgentCmd(c.client, c.response, &meta, cli.Args)
	}
	return nil
}
