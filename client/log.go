package client

import (
	"flag"
	"net/rpc"

	"github.com/dullgiulio/ringio/server"
	"github.com/dullgiulio/ringio/utils"
)

type CommandLog struct {
	client   *rpc.Client
	response *server.RpcResp
}

func NewCommandLog() *CommandLog {
	return &CommandLog{
		response: new(server.RpcResp),
	}
}

func (c *CommandLog) Help() string {
	return `Print internal log`
}

func (c *CommandLog) Init(fs *flag.FlagSet) bool {
	// nothing to do yet.
	return false
}

func (c *CommandLog) Run(cli *Cli) error {
	c.client = cli.GetClient()

	addLogAgentPipe(c.client, c.response, utils.GetRandomDotfile())

	return nil
}
