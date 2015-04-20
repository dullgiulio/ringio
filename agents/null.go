package agents

import (
	"github.com/dullgiulio/ringbuf"
	"github.com/dullgiulio/ringio/msg"
)

type AgentNull struct {
	meta   *AgentMetadata
	cancel chan bool
}

func NewAgentNull(id int, meta *AgentMetadata) *AgentNull {
	meta.ID = id
	return &AgentNull{
		meta:   meta,
		cancel: make(chan bool),
	}
}

func (a *AgentNull) Meta() *AgentMetadata {
	return a.meta
}

func (a *AgentNull) Init() {
	if a.meta.Options == nil {
		a.meta.Options = &AgentOptions{}
	}
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

func (a *AgentNull) Stop() error {
	return nil
}

func (a *AgentNull) Kill() error {
	return nil
}

func (a *AgentNull) WaitFinish() error {
	a.cancel <- true
	return nil
}

func (a *AgentNull) StartWrite() {}

func (a *AgentNull) InputToRingbuf(errors, output *ringbuf.Ringbuf) {
	<-a.cancel
}

func (a *AgentNull) OutputFromRingbuf(rStdout, rErrors, rOutput *ringbuf.Ringbuf, filter *msg.Filter) {
	<-a.cancel
}
