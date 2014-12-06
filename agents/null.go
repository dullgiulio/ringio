package agents

import (
	"github.com/dullgiulio/ringbuf"
	"github.com/dullgiulio/ringio/msg"
)

type AgentNull struct {
	meta *AgentMetadata
	Kill chan bool
}

func NewAgentNull(id int, role AgentRole, filter *msg.Filter) *AgentNull {
	return &AgentNull{
		meta: &AgentMetadata{Id: id, Role: role, Filter: filter},
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

func (a *AgentNull) InputToRingbuf(errors, output *ringbuf.Ringbuf) {
	<-a.Kill
}

func (a *AgentNull) OutputFromRingbuf(rStdout, rErrors, rOutput *ringbuf.Ringbuf, filter *msg.Filter) {
	<-a.Kill
}
