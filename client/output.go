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
	meta     *agents.AgentMetadata
}

func NewCommandOutput() *CommandOutput {
	return &CommandOutput{
		response: new(server.RpcResp),
		meta:     agents.NewAgentMetadata(),
	}
}

func (c *CommandOutput) Help() string {
	return `Output data from the ringbuf`
}

func (c *CommandOutput) Init(fs *flag.FlagSet) bool {
	fs.BoolVar(&c.meta.Options.NoWait, "no-wait", false, "Don't wait for future output, exit when finished dumping past data")
	fs.StringVar(&c.meta.Name, "name", "", "Name or comment of this sink")
	return true
}

func (c *CommandOutput) Run(cli *Cli) error {
	c.client = cli.GetClient()

	c.meta.Filter = cli.Filter

	if len(cli.Args) == 0 {
		addSinkAgentPipe(c.client, c.meta, c.response, utils.GetRandomDotfile())
	} else {
		addSinkAgentCmd(c.client, c.meta, c.response, cli.Args)
	}

	return nil
}
