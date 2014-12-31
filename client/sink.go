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

func addErrorsAgentPipe(client *rpc.Client, meta *agents.AgentMetadata, response *server.RpcResp, pipeName string) {
	_addSinkAgentPipe(client, meta, response, pipeName, agents.AgentRoleErrors)
}

func addSinkAgentPipe(client *rpc.Client, meta *agents.AgentMetadata, response *server.RpcResp, pipeName string) {
	_addSinkAgentPipe(client, meta, response, pipeName, agents.AgentRoleSink)
}

func addLogAgentPipe(client *rpc.Client, response *server.RpcResp, pipeName string) {
	_addSinkAgentPipe(client, &agents.AgentMetadata{}, response, pipeName, agents.AgentRoleLog)
}

func _addSinkAgentPipe(client *rpc.Client, meta *agents.AgentMetadata,
	response *server.RpcResp, pipeName string, role agents.AgentRole) {
	var id int

	p := pipe.New(pipeName)

	meta.Role = role

	if err := p.Create(); err != nil {
		utils.Fatal(fmt.Errorf("Couldn't create pipe: %s", err))
	}

	defer p.Remove()

	done := make(chan struct{})

	go func(done chan struct{}) {
		// Open will block until the pipe is opened on the other side.
		if err := p.OpenReadErr(); err != nil {
			utils.Fatal(fmt.Errorf("Couldn't open pipe for reading: %s", err))
		}

		r := bufio.NewReader(p)

		if _, err := r.WriteTo(os.Stdout); err != nil {
			utils.Fatal(err)
		}

		done <- struct{}{}
	}(done)

	if err := client.Call("RpcServer.Add", &server.RpcReq{
		Agent: &agents.AgentDescr{
			Args: []string{pipeName},
			Meta: *meta,
			Type: agents.AgentTypePipe,
		},
	}, &id); err != nil {
		utils.Fatal(err)
	}

	onexit.Defer(func() {
		if err := client.Call("RpcServer.Stop", id, &response); err != nil {
			utils.Fatal(err)
		}
	})

	// Wait for the pipe to be read.
	<-done
}

func addErrorsAgentCmd(client *rpc.Client, meta *agents.AgentMetadata, response *server.RpcResp, args []string) {
	_addSinkAgentCmd(client, meta, response, args, agents.AgentRoleErrors)
}

func addSinkAgentCmd(client *rpc.Client, meta *agents.AgentMetadata, response *server.RpcResp, args []string) {
	_addSinkAgentCmd(client, meta, response, args, agents.AgentRoleSink)
}

func _addSinkAgentCmd(client *rpc.Client, meta *agents.AgentMetadata,
	response *server.RpcResp, args []string, role agents.AgentRole) {
	var id int

	meta.Role = role

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
