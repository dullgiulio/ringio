package client

import (
	"flag"
	"net/rpc"

	"bitbucket.org/dullgiulio/ringio/server"
	"bitbucket.org/dullgiulio/ringio/utils"
)

type CommandInput struct {
	client   *rpc.Client
	response *server.RpcResp
}

func NewCommandInput() *CommandInput {
	return &CommandInput{
		response: new(server.RpcResp),
	}
}

func (c *CommandInput) Help() string {
	return `Take data and writes it into the ringbuf.`
}

func (c *CommandInput) Init(fs *flag.FlagSet) error {
	// nothing to do yet.
	return nil
}

func (c *CommandInput) Run(cli *Cli) error {
	if client, err := rpc.Dial("unix", utils.FileInDotpath(cli.Session)); err != nil {
		utils.Fatal(err)
	} else {
		c.client = client
	}

	if len(cli.Args) == 0 {
		addSourceAgentPipe(c.client, c.response, utils.GetRandomDotfile())
	} else {
		addSourceAgentCmd(c.client, c.response, cli.Args)
	}
	return nil
}
