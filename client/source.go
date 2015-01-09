package client

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"

	"github.com/dullgiulio/ringio/agents"
	"github.com/dullgiulio/ringio/onexit"
	"github.com/dullgiulio/ringio/pipe"
	"github.com/dullgiulio/ringio/server"
	"github.com/dullgiulio/ringio/utils"
)

func addSourceAgentPipe(client *rpc.Client, response *server.RpcResp, meta *agents.AgentMetadata, pipeName string) {
	var id int

	p := pipe.New(pipeName)

	if err := p.Create(); err != nil {
		utils.Fatal(err)
	}

	if err := p.OpenWriteErr(); err != nil {
		_removePipe(p)
		utils.Fatal(err)
	}

	meta.Role = agents.AgentRoleSource

	if err := client.Call("RpcServer.Add", &server.RpcReq{
		Agent: &agents.AgentDescr{
			Args: []string{pipeName},
			Meta: *meta,
			Type: agents.AgentTypePipe,
		},
	}, &id); err != nil {
		_removePipe(p)
		utils.Fatal(err)
	}

	p.Remove()

	onexit.Defer(func() {
		if err := client.Call("RpcServer.Stop", id, &response); err != nil {
			utils.Fatal(err)
		}
	})

	// Write to pipe from stdin.
	r := bufio.NewReader(os.Stdin)

	if _, err := r.WriteTo(p); err != nil {
		utils.Fatal(err)
	}
}

func addSourceAgentCmd(client *rpc.Client, response *server.RpcResp, meta *agents.AgentMetadata, args []string) {
	var id int

	meta.Role = agents.AgentRoleSource

	if err := client.Call("RpcServer.Add", &server.RpcReq{
		Agent: &agents.AgentDescr{
			Args: args,
			Meta: *meta,
			Type: agents.AgentTypeCmd,
		},
	}, &id); err != nil {
		utils.Fatal(err)
	}

	fmt.Printf("Added agent %%%d\n", id)
}
