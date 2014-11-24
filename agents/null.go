package agents

import (
	"bitbucket.org/dullgiulio/ringbuf"
)

type AgentNull struct {
	meta *AgentMetadata
	Kill chan bool
}

func NewAgentNull(id int, role AgentRole) *AgentNull {
	return &AgentNull{
		meta: &AgentMetadata{Id: id, Role: role},
		Kill: make(chan bool),
	}
}

func (a *AgentNull) Meta() *AgentMetadata {
	return a.meta
}

func (a *AgentNull) Init() {
}

func (a *AgentNull) Descr() AgentDescr {
	return AgentDescr{
		Args: []string{"nullagent"},
		Meta: *a.meta,
		Type: AgentTypeNull,
	}
}

func (a *AgentNull) String() string {
	return "AgentNull"
}

func (a *AgentNull) Cancel() error {
	a.Kill <- true
	return nil
}

func (a *AgentNull) Stop() {
}

func (a *AgentNull) OutputToRingbuf(errors, output *ringbuf.Ringbuf) {
	<-a.Kill
}

func (a *AgentNull) InputFromRingbuf(rStdout, rErrors, rOutput *ringbuf.Ringbuf) {
	<-a.Kill
}
