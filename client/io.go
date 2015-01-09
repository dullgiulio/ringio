package client

import (
	"flag"
	"net/rpc"
	"os"

	"github.com/dullgiulio/ringio/agents"
	"github.com/dullgiulio/ringio/server"
	"github.com/dullgiulio/ringio/utils"
)

type CommandIO struct {
	client   *rpc.Client
	response *server.RpcResp
	name     string
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
	fs.StringVar(&c.name, "name", "", "Name or comment of this stream")
	return true
}

func (c *CommandIO) Run(cli *Cli) error {
	c.client = cli.GetClient()

	metaSource := agents.AgentMetadata{
		User: os.Getenv("USER"),
		Name: c.name,
	}
	metaSink := metaSource
	metaSink.Filter = cli.Filter

	go addSourceAgentPipe(c.client, c.response, &metaSource, utils.GetRandomDotfile())
	addSinkAgentPipe(c.client, &metaSink, c.response, utils.GetRandomDotfile())

	return nil
}
