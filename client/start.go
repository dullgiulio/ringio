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
	startAll bool
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
	fs.BoolVar(&c.startAll, "all", false, "Start all that has never been started")
	return true
}

func (c *CommandStart) agent(cli *Cli) {
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
}

func (c *CommandStart) all(cli *Cli) {
	if cli.Filter != nil {
		utils.Fatal(errors.New("Filtering when --all is specified makes no sense."))
	}

	if err := c.client.Call("RpcServer.StartAll", &server.RpcReq{}, &c.response); err != nil {
		utils.Fatal(err)
	}
}

func (c *CommandStart) Run(cli *Cli) error {
	c.client = cli.GetClient()

	if c.startAll {
		c.all(cli)
	} else {
		c.agent(cli)
	}

	return nil
}
