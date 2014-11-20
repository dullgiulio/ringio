package client

import (
	"errors"
	"flag"
	"fmt"
	"net/rpc"

	"bitbucket.org/dullgiulio/ringio/agents"
	"bitbucket.org/dullgiulio/ringio/server"
	"bitbucket.org/dullgiulio/ringio/utils"
)

type CommandList struct {
	client *rpc.Client
}

func NewCommandList() *CommandList {
	return &CommandList{}
}

func (c *CommandList) Help() string {
	return `List all agents`
}

func (c *CommandList) Init(fs *flag.FlagSet) error {
	// nothing to do yet.
	return nil
}

func (c *CommandList) Run(cli *Cli) error {
	if client, err := rpc.Dial("unix", utils.FileInDotpath(cli.Session)); err != nil {
		utils.Fatal(err)
	} else {
		c.client = client
	}

	responseList := agents.NewAgentMessageResponseList()
	if err := c.client.Call("RpcServer.List", &server.RpcReq{}, &responseList); err != nil {
		utils.Fatal(err)
	} else {
		if responseList.Agents != nil {
			for _, a := range *responseList.Agents {
				fmt.Printf("%s\n", &a)
			}
		} else {
			utils.Fatal(errors.New("This session is currently empty."))
		}
	}
	return nil
}
