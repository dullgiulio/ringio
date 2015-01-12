package client

import (
	"flag"
	"net/rpc"

	"github.com/dullgiulio/ringio/server"
	"github.com/dullgiulio/ringio/utils"
)

type CommandClose struct {
	client   *rpc.Client
	response *server.RPCResp
}

func NewCommandClose() *CommandClose {
	return &CommandClose{
		response: new(server.RPCResp),
	}
}

func (c *CommandClose) Help() string {
	return `Close a session`
}

func (c *CommandClose) Init(fs *flag.FlagSet) bool {
	// nothing to do yet.
	return false
}

func (c *CommandClose) Run(cli *Cli) error {
	c.client = cli.GetClient()

	if err := c.client.Call("RpcServer.Close", &server.RPCReq{}, &c.response); err != nil {
		utils.Fatal(err)
	}

	return nil
}
