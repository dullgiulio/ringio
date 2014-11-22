package client

import (
	"flag"
	"net/rpc"

	"bitbucket.org/dullgiulio/ringio/server"
	"bitbucket.org/dullgiulio/ringio/utils"
)

type CommandRun struct {
	client   *rpc.Client
	response *server.RpcResp
}

func NewCommandRun() *CommandRun {
	return &CommandRun{
		response: new(server.RpcResp),
	}
}

func (c *CommandRun) Help() string {
	return `Run all processes`
}

func (c *CommandRun) Init(fs *flag.FlagSet) error {
	// nothing to do yet.
	return nil
}

func (c *CommandRun) Run(cli *Cli) error {
	c.client = cli.GetClient()

	if err := c.client.Call("RpcServer.Run", &server.RpcReq{}, &c.response); err != nil {
		utils.Fatal(err)
	}

	return nil
}
