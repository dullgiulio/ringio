package client

import (
	"errors"
	"flag"
	"net/rpc"

	"github.com/dullgiulio/ringio/server"
	"github.com/dullgiulio/ringio/utils"
)

type CommandKill struct {
	client   *rpc.Client
	response *server.RpcResp
}

func NewCommandKill() *CommandKill {
	return &CommandKill{
		response: new(server.RpcResp),
	}
}

func (c *CommandKill) Help() string {
	return `Kill a specified agent`
}

func (c *CommandKill) Init(fs *flag.FlagSet) error {
	// nothing to do yet.
	return nil
}

func (c *CommandKill) Run(cli *Cli) error {
	c.client = cli.GetClient()

	argsErr := errors.New("Kill must be followed by an argument.")

	if cli.Filter == nil {
		utils.Fatal(argsErr)
	}

	in := cli.Filter.GetIn()

	if len(in) == 0 {
		utils.Fatal(argsErr)
	}

	for _, id := range in {
		if err := c.client.Call("RpcServer.Kill", id, &c.response); err != nil {
			utils.Error(err)
		}
	}

	return nil
}
