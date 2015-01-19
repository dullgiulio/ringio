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
	response *server.RPCResp
}

func NewCommandStop() *CommandStop {
	return &CommandStop{
		response: new(server.RPCResp),
	}
}

func (c *CommandStop) Help() string {
	return `Stop a specified agent by stop sending data to it`
}

func (c *CommandStop) Init(fs *flag.FlagSet) bool {
	// nothing to do yet.
	return false
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
		if err := c.client.Call("RPCServer.Stop", id, &c.response); err != nil {
			utils.Error(err)
		}
	}

	return nil
}
