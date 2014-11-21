package agents

import (
	"bitbucket.org/dullgiulio/ringbuf"
	"bitbucket.org/dullgiulio/ringio/pipe"
)

type AgentPipe struct {
	pipe   *pipe.Pipe
	meta   *AgentMetadata
	cancel chan bool
}

func NewAgentPipe(pipeName string, role AgentRole) *AgentPipe {
	return &AgentPipe{
		pipe:   pipe.New(pipeName),
		meta:   &AgentMetadata{Role: role},
		cancel: make(chan bool),
	}
}

func (a *AgentPipe) Init() {
	if a.meta.Role == AgentRoleSink || a.meta.Role == AgentRoleErrors {
		if ok := a.pipe.OpenWrite(); !ok {
			return
		}
	}
}

func (a *AgentPipe) Meta() *AgentMetadata {
	return a.meta
}

func (a *AgentPipe) Descr() AgentDescr {
	return AgentDescr{
		Args: []string{a.String()},
		Meta: *a.meta,
		Type: AgentTypePipe,
	}
}

func (a *AgentPipe) String() string {
	return "|" + a.pipe.String()
}

func (a *AgentPipe) Cancel() error {
	a.cancel <- true
	return nil
}

func (a *AgentPipe) OutputToRingbuf(rErrors, rOutput *ringbuf.Ringbuf) {
	if ok := a.pipe.OpenRead(); !ok {
		return
	}

	a.pipe.Remove()

	id := a.meta.Id

	cancelled := writeToRingbuf(id, a.pipe, rOutput, a.cancel, nil)

	if !cancelled {
		<-a.cancel
	}

	close(a.cancel)

	a.pipe.Close()
}

func (a *AgentPipe) InputFromRingbuf(rStdout, rErrors, rOutput *ringbuf.Ringbuf) {
	cancelled := readFromRingbuf(a.pipe, rOutput, a.cancel, nil)

	if !cancelled {
		<-a.cancel
	}

	close(a.cancel)

	a.pipe.Close()
}
