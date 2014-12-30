package agents

import (
	"testing"

	"github.com/dullgiulio/ringio/log"
)

func readLogs(t *testing.T) {
	nw := new(log.NullWriter)
	log.AddWriter(nw)

	if !log.Run(log.LevelError) {
		t.Log("Did not expect logger to return an error")
		t.Fail()
	}
}

func TestAgentAdding(t *testing.T) {
	go readLogs(t)

	proceed := make(chan struct{})
	ac := NewCollection()
	go func() {
		ac.Run(true)
		proceed <- struct{}{}
	}()

	a0 := NewAgentNull(0, AgentRoleSource, nil, &AgentOptions{})
	resp := NewAgentMessageResponseBool()

	ac.Add(a0, &resp)
	if _, err := resp.Get(); err != nil {
		t.Error(err)
	}

	meta := a0.Meta()

	if meta.Id == 0 {
		t.Error("Did not expect ID to be still zero after adding")
	}
	firstId := meta.Id

	if meta.Role != AgentRoleSource {
		t.Error("Adding an agent changed its role")
	}

	if meta.Status == AgentStatusRunning {
		t.Error("Did not expected agent to be in running state")
	}

	if !meta.Started.IsZero() {
		t.Error("Did not expected starting time to be set")
	}

	if !meta.Finished.IsZero() {
		t.Error("Expected finish time to be still undefined")
	}

	a1 := NewAgentNull(0, AgentRoleSource, nil, &AgentOptions{})
	resp = NewAgentMessageResponseBool()

	ac.Add(a1, &resp)
	resp.Get()

	ac.Start(a1, &resp)
	resp.Get()

	meta = a1.Meta()

	if meta.Id == 0 {
		t.Error("Did not expect ID to be still zero after adding")
	}

	if meta.Id <= firstId {
		t.Error("Expected second ID to be greater than the first")
	}

	if meta.Role != AgentRoleSource {
		t.Error("Adding an agent changed its role")
	}

	if meta.Status != AgentStatusRunning {
		t.Error("Expected agent to be in running state")
	}

	if meta.Started.IsZero() {
		t.Error("Expected starting time to be set")
	}

	if !meta.Finished.IsZero() {
		t.Error("Expected finish time to be still undefined")
	}

	ac.Start(a0, &resp)
	resp.Get()

	ac.Cancel(&resp)
	resp.Get()

	<-proceed

	meta = a0.Meta()
	if meta.Id == 0 {
		t.Error("Did not expect ID to be still zero after cancel")
	}

	if meta.Role != AgentRoleSource {
		t.Error("Stopping an agent changed its role")
	}

	if meta.Status != AgentStatusFinished {
		t.Error("Expected agent to be in finished state")
	}

	if meta.Started.IsZero() {
		t.Error("Expected starting time to be set")
	}

	if meta.Finished.IsZero() {
		t.Error("Expected finish time to be still undefined")
	}

	meta = a1.Meta()
	if meta.Id == 0 {
		t.Error("Did not expect ID to be still zero after cancel")
	}

	if meta.Role != AgentRoleSource {
		t.Error("Stopping an agent changed its role")
	}

	if meta.Status != AgentStatusFinished {
		t.Error("Expected agent to be in finished state")
	}

	if meta.Started.IsZero() {
		t.Error("Expected starting time to be set")
	}

	if meta.Finished.IsZero() {
		t.Error("Expected finish time to be still undefined")
	}

	log.Cancel()
}

func TestInvalidActions(t *testing.T) {
	go readLogs(t)

	proceed := make(chan struct{})
	ac := NewCollection()
	go func() {
		ac.Run(false)
		proceed <- struct{}{}
	}()

	resp := NewAgentMessageResponseBool()

	// Test that cancel exits correctly.
	ac.Cancel(&resp)
	resp.Get()

	<-proceed

	log.Cancel()
}
