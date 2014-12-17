package client

import (
	"errors"
	"flag"
	"fmt"
	"net/rpc"

	"github.com/dullgiulio/ringio/agents"
	"github.com/dullgiulio/ringio/server"
	"github.com/dullgiulio/ringio/utils"
)

type CommandList struct {
	client  *rpc.Client
	verbose bool
}

func NewCommandList() *CommandList {
	return &CommandList{}
}

func (c *CommandList) Help() string {
	return `List all agents`
}

func (c *CommandList) Init(fs *flag.FlagSet) bool {
	fs.BoolVar(&c.verbose, "verbose", false, "Long version of the list")
	return true
}

func (c *CommandList) Run(cli *Cli) error {
	c.client = cli.GetClient()

	responseList := agents.NewAgentMessageResponseList()
	if err := c.client.Call("RpcServer.List", &server.RpcReq{}, &responseList); err != nil {
		utils.Fatal(err)
	} else {
		if responseList.Agents != nil {
			for _, a := range *responseList.Agents {
				if c.verbose {
					fmt.Printf("%s\n", a.Text())
				} else {
					fmt.Printf("%s\n", &a)
				}
			}
		} else {
			utils.Fatal(errors.New("This session is currently empty."))
		}
	}
	return nil
}
