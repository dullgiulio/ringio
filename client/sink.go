package client

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"

	"bitbucket.org/dullgiulio/ringio/agents"
	"bitbucket.org/dullgiulio/ringio/onexit"
	"bitbucket.org/dullgiulio/ringio/pipe"
	"bitbucket.org/dullgiulio/ringio/server"
	"bitbucket.org/dullgiulio/ringio/utils"
)

func addSourceAgentPipe(client *rpc.Client, response *server.RpcResp, pipeName string) {
	var id int

	p := pipe.New(pipeName)

	if err := p.OpenWriteErr(); err != nil {
		utils.Fatal(fmt.Errorf("Couldn't open pipe for writing: %s", err))
	}

	if err := client.Call("RpcServer.Add", &server.RpcReq{
		Agent: &agents.AgentDescr{
			Args: []string{pipeName},
			Meta: agents.AgentMetadata{Role: agents.AgentRoleSource},
			Type: agents.AgentTypePipe,
		},
	}, &id); err != nil {
		utils.Fatal(err)
	}

	// Write to pipe from stdin.
	r := bufio.NewReader(os.Stdin)

	if _, err := r.WriteTo(p); err != nil {
		utils.Fatal(err)
	}
}

func addErrorsAgentPipe(client *rpc.Client, response *server.RpcResp, pipeName string) {
	_addSinkAgentPipe(client, response, pipeName, agents.AgentRoleErrors)
}

func addSinkAgentPipe(client *rpc.Client, response *server.RpcResp, pipeName string) {
	_addSinkAgentPipe(client, response, pipeName, agents.AgentRoleSink)
}

func _addSinkAgentPipe(client *rpc.Client, response *server.RpcResp, pipeName string, role agents.AgentRole) {
	var id int

	p := pipe.New(pipeName)

	if err := client.Call("RpcServer.Add", &server.RpcReq{
		Agent: &agents.AgentDescr{
			Args: []string{pipeName},
			Meta: agents.AgentMetadata{Role: role},
			Type: agents.AgentTypePipe,
		},
	}, &id); err != nil {
		utils.Fatal(err)
	}

	if err := p.OpenReadErr(); err != nil {
		utils.Fatal(fmt.Errorf("Couldn't open pipe for reading: %s", err))
	}

	p.Remove()

	onexit.Defer(func() {
		if err := client.Call("RpcServer.Stop", id, &response); err != nil {
			utils.Fatal(err)
		}
	})

	r := bufio.NewReader(p)

	if _, err := r.WriteTo(os.Stdout); err != nil {
		utils.Fatal(err)
	}
}

func addSourceAgentCmd(client *rpc.Client, response *server.RpcResp, args []string) {
	var id int

	if err := client.Call("RpcServer.Add", &server.RpcReq{
		Agent: &agents.AgentDescr{
			Args: args,
			Meta: agents.AgentMetadata{Role: agents.AgentRoleSource},
			Type: agents.AgentTypeCmd,
		},
	}, &id); err != nil {
		utils.Fatal(err)
	}
}

func addErrorsAgentCmd(client *rpc.Client, response *server.RpcResp, args []string) {
	_addSinkAgentCmd(client, response, args, agents.AgentRoleErrors)
}

func addSinkAgentCmd(client *rpc.Client, response *server.RpcResp, args []string) {
	_addSinkAgentCmd(client, response, args, agents.AgentRoleSource)
}

func _addSinkAgentCmd(client *rpc.Client, response *server.RpcResp, args []string, role agents.AgentRole) {
	var id int

	if err := client.Call("RpcServer.Add", &server.RpcReq{
		Agent: &agents.AgentDescr{
			Args: args,
			Meta: agents.AgentMetadata{Role: role},
			Type: agents.AgentTypeCmd,
		},
	}, &id); err != nil {
		utils.Fatal(err)
	}
}
