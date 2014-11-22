package client

import (
	"flag"
	"net/rpc"

	"bitbucket.org/dullgiulio/ringio/server"
	"bitbucket.org/dullgiulio/ringio/utils"
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
	return `Output data from the ringbuf.`
}

func (c *CommandLog) Init(fs *flag.FlagSet) error {
	// nothing to do yet.
	return nil
}

func (c *CommandLog) Run(cli *Cli) error {
	c.client = cli.GetClient()

	addLogAgentPipe(c.client, c.response, utils.GetRandomDotfile())

	return nil
}
