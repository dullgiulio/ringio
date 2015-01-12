package server

import (
	"testing"

	"github.com/dullgiulio/ringio/agents"
	"github.com/dullgiulio/ringio/log"
	"github.com/dullgiulio/ringio/onexit"
)

func TestRPCServer(t *testing.T) {
	proceed := make(chan struct{})

	go func() {
		if !log.Run(log.LevelWarn) {
			t.Error("Expected success after cancel")
		}
		proceed <- struct{}{}
	}()

	s := NewRPCServer("test-socket", true)

	resp := 0

	req := &RPCReq{
		Agent: &agents.AgentDescr{
			Args: []string{},
			Meta: agents.AgentMetadata{Role: agents.AgentRoleSource},
			Type: agents.AgentTypeNull,
		},
	}

	if err := s.Ping(req, &resp); err != nil || resp != 1 {
		t.Error(err)
	}

	if err := s.Add(req, &resp); err != nil {
		t.Error(err)
	}

	if resp == 0 {
		t.Error("Expected an ID as response")
	}

	var done RPCResp

	if err := s.Kill(resp, &done); err != nil {
		t.Error(err)
	}

	if !done {
		t.Error("Expected Kill to return success")
	}

	done = false

	if err := s.Run(nil, &done); err != nil {
		t.Error(err)
	}

	if !done {
		t.Error("Expected Run to return success")
	}

	// TODO: Test List
	list := agents.NewAgentMessageResponseList()

	if err := s.List(nil, &list); err != nil {
		t.Error(err)
	}

	if !done {
		t.Error("Expected List to return success")
	}

	if list.Agents == nil {
		t.Error("Did not expect retrieved list for List command to be null")
	}

	exit := make(chan int)

	onexit.SetFunc(func(i int) {
		exit <- i
	})
	onexit.HandleInterrupt()

	go func() {
		if err := s.Close(nil, &done); err != nil {
			t.Error(err)
		}
	}()

	if i := <-exit; i != 0 {
		t.Error("Successful Close did not exit with zero")
	}

	// Check that the server doesn't accept commands after Close.
	if err := s.Add(req, &resp); err == nil {
		t.Error("Expected error on Add on closed server")
	}

	if err := s.Stop(resp, &done); err == nil {
		t.Error("Expected error on Stop on closed server")
	}

	if err := s.List(nil, &list); err == nil {
		t.Error("Expected error on List on closed server")
	}

	if err := s.Run(nil, &done); err == nil {
		t.Error("Expected error on Run on closed server")
	}

	if err := s.Close(nil, &done); err == nil {
		t.Error("Expected error on Close on closed server")
	}

	<-proceed
}
