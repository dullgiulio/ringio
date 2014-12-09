package client

import (
	"errors"
	"flag"
	"net/rpc"

	"github.com/dullgiulio/ringio/server"
	"github.com/dullgiulio/ringio/utils"
)

type CommandStop struct {
	client   *rpc.Client
	response *server.RpcResp
}

func NewCommandStop() *CommandStop {
	return &CommandStop{
		response: new(server.RpcResp),
	}
}

func (c *CommandStop) Help() string {
	return `List all agents`
}

func (c *CommandStop) Init(fs *flag.FlagSet) error {
	// nothing to do yet.
	return nil
}

func (c *CommandStop) Run(cli *Cli) error {
	c.client = cli.GetClient()

	argsErr := errors.New("Stop must be followed by an argument.")

	if cli.Filter == nil {
		utils.Fatal(argsErr)
	}

	in := cli.Filter.GetIn()

	if len(in) == 0 {
		utils.Fatal(argsErr)
	}

	for _, id := range in {
		if err := c.client.Call("RpcServer.Stop", id, &c.response); err != nil {
			utils.Error(err)
		}
	}

	return nil
}
