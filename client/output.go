package client

import (
	"flag"
	"net/rpc"

	"bitbucket.org/dullgiulio/ringio/server"
	"bitbucket.org/dullgiulio/ringio/utils"
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
	return `Output data from the ringbuf.`
}

func (c *CommandOutput) Init(fs *flag.FlagSet) error {
	// nothing to do yet.
	return nil
}

func (c *CommandOutput) Run(cli *Cli) error {
	if client, err := rpc.Dial("unix", utils.FileInDotpath(cli.Session)); err != nil {
		utils.Fatal(err)
	} else {
		c.client = client
	}

	if len(cli.Args) == 0 {
		addSinkAgentPipe(c.client, c.response, utils.GetRandomDotfile())
	} else {
		addSinkAgentCmd(c.client, c.response, cli.Args)
	}

	return nil
}
