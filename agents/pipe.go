package agents

import (
	"fmt"

	"github.com/dullgiulio/ringbuf"
	"github.com/dullgiulio/ringio/log"
	"github.com/dullgiulio/ringio/msg"
	"github.com/dullgiulio/ringio/pipe"
)

type AgentPipe struct {
	pipe     *pipe.Pipe
	meta     *AgentMetadata
	cancelCh chan bool
}

func NewAgentPipe(pipeName string, role AgentRole, filter *msg.Filter, options *AgentOptions) *AgentPipe {
	return &AgentPipe{
		pipe:     pipe.New(pipeName),
		meta:     &AgentMetadata{Role: role, Filter: filter, Options: options},
		cancelCh: make(chan bool),
	}
}

func (a *AgentPipe) Init() {
	if a.meta.Options == nil {
		a.meta.Options = &AgentOptions{}
	}

	if a.meta.Role == AgentRoleSink ||
		a.meta.Role == AgentRoleErrors ||
		a.meta.Role == AgentRoleLog {
		if ok := a.pipe.OpenWrite(); !ok {
			// TODO:XXX: Must return the error to the user.
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

func (a *AgentPipe) cancel() error {
	if a.meta.Status.IsRunning() {
		a.cancelCh <- true
	}

	a.pipe.Close()
	return nil
}

func (a *AgentPipe) Stop() error {
	if err := a.cancel(); err != nil {
		return err
	}

	log.Info(log.FacilityAgent, fmt.Sprintf("PipeAgent %d has been stopped", a.meta.Id))
	return nil
}

func (a *AgentPipe) Kill() error {
	if err := a.Stop(); err != nil {
		return err
	}

	log.Info(log.FacilityAgent, fmt.Sprintf("PipeAgent %d has been killed", a.meta.Id))
	return nil
}

func (a *AgentPipe) WaitFinish() error {
	return nil
}

func (a *AgentPipe) InputToRingbuf(rErrors, rOutput *ringbuf.Ringbuf) {
	if ok := a.pipe.OpenRead(); !ok {
		return
	}

	a.pipe.Remove()

	id := a.meta.Id

	cancelled := writeToRingbuf(id, a.pipe, rOutput, a.cancelCh, nil)

	if !cancelled {
		<-a.cancelCh
	}

	close(a.cancelCh)
}

func (a *AgentPipe) OutputFromRingbuf(rStdout, rErrors, rOutput *ringbuf.Ringbuf, filter *msg.Filter) {
	cancelled := readFromRingbuf(a.pipe, filter, rOutput,
		makeReaderOptions(a.meta.Options), a.cancelCh, nil)

	if !cancelled {
		<-a.cancelCh
	}

	close(a.cancelCh)
}
