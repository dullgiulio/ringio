package client

import (
	"flag"
	"net/rpc"

	"bitbucket.org/dullgiulio/ringio/server"
	"bitbucket.org/dullgiulio/ringio/utils"
)

type CommandClose struct {
	client   *rpc.Client
	response *server.RpcResp
}

func NewCommandClose() *CommandClose {
	return &CommandClose{
		response: new(server.RpcResp),
	}
}

func (c *CommandClose) Help() string {
	return `Close a session`
}

func (c *CommandClose) Init(fs *flag.FlagSet) error {
	// nothing to do yet.
	return nil
}

func (c *CommandClose) Run(cli *Cli) error {
	c.client = cli.GetClient()

	if err := c.client.Call("RpcServer.Close", &server.RpcReq{}, &c.response); err != nil {
		utils.Fatal(err)
	}

	return nil
}
