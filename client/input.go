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
	client    *rpc.Client
	response  *server.RPCResp
	name      string
	ignoreEnv bool
}

func NewCommandInput() *CommandInput {
	return &CommandInput{
		response: new(server.RPCResp),
	}
}

func (c *CommandInput) Help() string {
	return `Reads data from stdin (and stderr for processes) and writes it into the ringbuf`
}

func (c *CommandInput) Init(fs *flag.FlagSet) bool {
	fs.StringVar(&c.name, "name", "", "Name or comment of this source")
	fs.BoolVar(&c.ignoreEnv, "ignore-env", false, "Ignore current environment when running a subprocess")
	return true
}

func (c *CommandInput) Run(cli *Cli) error {
	c.client = cli.GetClient()

	meta := agents.AgentMetadata{
		User: os.Getenv("USER"),
		Name: c.name,
	}

	if !c.ignoreEnv {
		meta.Env = os.Environ()
	}

	if len(cli.Args) == 0 {
		addSourceAgentPipe(c.client, c.response, &meta, utils.GetRandomDotfile())
	} else {
		addSourceAgentCmd(c.client, c.response, &meta, cli.Args)
	}
	return nil
}
