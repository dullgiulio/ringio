package client

import (
	"errors"
	"flag"
	"net/rpc"

	"github.com/dullgiulio/ringio/server"
	"github.com/dullgiulio/ringio/utils"
)

type CommandStart struct {
	client   *rpc.Client
	response *server.RpcResp
}

func NewCommandStart() *CommandStart {
	return &CommandStart{
		response: new(server.RpcResp),
	}
}

func (c *CommandStart) Help() string {
	return `Start a specific agent`
}

func (c *CommandStart) Init(fs *flag.FlagSet) bool {
	// nothing to do yet.
	return false
}

func (c *CommandStart) Run(cli *Cli) error {
	c.client = cli.GetClient()

	argsErr := errors.New("Start must be followed by an argument.")

	if cli.Filter == nil {
		utils.Fatal(argsErr)
	}

	in := cli.Filter.GetIn()

	if len(in) == 0 {
		utils.Fatal(argsErr)
	}

	for _, id := range in {
		if err := c.client.Call("RpcServer.Start", id, &c.response); err != nil {
			utils.Error(err)
		}
	}

	return nil
}
