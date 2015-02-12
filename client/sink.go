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

func addErrorsAgentPipe(client *rpc.Client, meta *agents.AgentMetadata, response *server.RPCResp, pipeName string) {
	_addSinkAgentPipe(client, meta, response, pipeName, agents.AgentRoleErrors)
}

func addSinkAgentPipe(client *rpc.Client, meta *agents.AgentMetadata, response *server.RPCResp, pipeName string) {
	_addSinkAgentPipe(client, meta, response, pipeName, agents.AgentRoleSink)
}

func addLogAgentPipe(client *rpc.Client, response *server.RPCResp, pipeName string) {
	_addSinkAgentPipe(client, &agents.AgentMetadata{}, response, pipeName, agents.AgentRoleLog)
}

func _addSinkAgentPipe(client *rpc.Client, meta *agents.AgentMetadata,
	response *server.RPCResp, pipeName string, role agents.AgentRole) {
	var id int

	p := pipe.New(pipeName)

	meta.Role = role

	if err := p.Create(); err != nil {
		utils.Fatal(fmt.Errorf("Couldn't create pipe: %s", err))
	}

	done := make(chan struct{})

	go func(done chan struct{}) {
		// Open will block until the pipe is opened on the other side.
		if err := p.OpenReadErr(); err != nil {
			_removePipe(p)
			utils.Fatal(fmt.Errorf("Couldn't open pipe for reading: %s", err))
		}

		r := bufio.NewReader(p)

		if _, err := r.WriteTo(os.Stdout); err != nil {
			_removePipe(p)
			utils.Fatal(err)
		}

		done <- struct{}{}
	}(done)

	if err := client.Call("RPCServer.Add", &server.RPCReq{
		Agent: &agents.AgentDescr{
			Args: []string{pipeName},
			Meta: *meta,
			Type: agents.AgentTypePipe,
		},
	}, &id); err != nil {
		_removePipe(p)
		utils.Fatal(err)
	}

	_removePipe(p)

	onexit.Defer(func() {
		if err := client.Call("RPCServer.Stop", id, &response); err != nil {
			utils.Fatal(err)
		}
	})

	// Wait for the pipe to be read.
	<-done

	if err := p.Remove(); err != nil {
		utils.Fatal(err)
	}
}

func addErrorsAgentCmd(client *rpc.Client, meta *agents.AgentMetadata, response *server.RPCResp, args []string) {
	_addSinkAgentCmd(client, meta, response, args, agents.AgentRoleErrors)
}

func addSinkAgentCmd(client *rpc.Client, meta *agents.AgentMetadata, response *server.RPCResp, args []string) {
	_addSinkAgentCmd(client, meta, response, args, agents.AgentRoleSink)
}

func _addSinkAgentCmd(client *rpc.Client, meta *agents.AgentMetadata,
	response *server.RPCResp, args []string, role agents.AgentRole) {
	var id int

	meta.Role = role

	if err := client.Call("RPCServer.Add", &server.RPCReq{
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
