package client

import (
	"errors"
	"flag"
	"net/rpc"
	"time"

	"github.com/dullgiulio/ringio/server"
	"github.com/dullgiulio/ringio/utils"
)

type CommandPing struct {
	client  *rpc.Client
	maxWait int
}

func NewCommandPing() *CommandPing {
	return &CommandPing{}
}

func (c *CommandPing) Help() string {
	return `Ping a session`
}

func (c *CommandPing) Init(fs *flag.FlagSet) bool {
	fs.IntVar(&c.maxWait, "wait", 2, "Number of seconds to wait before giving up")
	return true
}

func (c *CommandPing) Run(cli *Cli) error {
	c.client = cli.GetClient()

	resp := 0

	callRes := c.client.Go("RpcServer.Ping", &server.RPCReq{}, &resp, nil)

	select {
	case <-callRes.Done:
	case <-time.After(time.Duration(c.maxWait) * time.Second):
		break
	}

	if resp != 1 {
		utils.Fatal(errors.New("Looks like the session is not responding properly."))
	}

	return nil
}
