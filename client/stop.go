package client

import (
	"errors"
	"flag"
	"net/rpc"
	"strconv"

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

	if len(cli.Args) < 1 {
		utils.Fatal(errors.New("Stop must be followed by an argument."))
	}

	// Stop program by id
	if id, err := strconv.Atoi(cli.Args[0]); err != nil {
		utils.Fatal(err)
	} else if err = c.client.Call("RpcServer.Stop", id, &c.response); err != nil {
		utils.Fatal(err)
	}

	return nil
}
