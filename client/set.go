package client

import (
	"flag"
	"net/rpc"

	"bitbucket.org/dullgiulio/ringio/server"
	"bitbucket.org/dullgiulio/ringio/utils"
)

type CommandSet struct {
	client   *rpc.Client
	response *server.RpcResp
}

func NewCommandSet() *CommandSet {
	return &CommandSet{
		response: new(server.RpcResp),
	}
}

func (c *CommandSet) Help() string {
	return `List all agents`
}

func (c *CommandSet) Init(fs *flag.FlagSet) error {
	// nothing to do yet.
	return nil
}

func (c *CommandSet) Run(cli *Cli) error {
	if client, err := rpc.Dial("unix", utils.FileInDotpath(cli.Session)); err != nil {
		utils.Fatal(err)
	} else {
		c.client = client
	}

	var action server.ServerAction

	// TODO: set minimum logging level here.
	switch cli.Args[0] {
	case "verbose":
		action = server.ActionSetVerbose
	case "quiet":
		action = server.ActionUnsetVerbose
	case "locked":
		action = server.ActionSetLocked
	}

	if err := c.client.Call("RpcServer.Set", &server.RpcReq{
		Action: action,
	}, &c.response); err != nil {
		utils.Fatal(err)
	}

	return nil
}
