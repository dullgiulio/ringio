package agents

import (
	"fmt"
	"time"

	"github.com/dullgiulio/ringbuf"
	"github.com/dullgiulio/ringio/config"
	"github.com/dullgiulio/ringio/msg"
)

type AgentMetadata struct {
	ID       int
	Role     AgentRole
	Status   AgentStatus
	Started  time.Time
	Finished time.Time
	User     string
	Name     string
	Filter   *msg.Filter
	Options  *AgentOptions
}

func NewAgentMetadata() *AgentMetadata {
	return &AgentMetadata{
		Filter:  new(msg.Filter),
		Options: new(AgentOptions),
	}
}

type AgentOptions struct {
	NoWait bool
}

type Agent interface {
	Init()
	// Gracefully stop an Agent.
	Stop() error
	// Forced stop of an Agent. Can be force onto stopped agents.
	Kill() error
	// Returns when the Agent is definitely finished and can be marked as such.
	WaitFinish() error
	String() string
	Meta() *AgentMetadata
	Descr() AgentDescr
	InputToRingbuf(errors, output *ringbuf.Ringbuf)
	OutputFromRingbuf(stdout, errors, output *ringbuf.Ringbuf, filter *msg.Filter)
}

func (ac *Collection) isFilteringSinkAgents(filter *msg.Filter) bool {
	if filter == nil {
		return false
	}

	fin := filter.GetIn()
	haveSinks := false

	if len(fin) == 0 {
		return false
	}

	for _, a := range ac.agents {
		meta := a.Meta()

		if meta.Role == AgentRoleSink {
			for _, id := range fin {
				if meta.ID == id {
					return true
				}
			}

			haveSinks = true
		}
	}

	if haveSinks {
		return false
	}

	return true
}

func (ac *Collection) validateFilter(self int, filter *msg.Filter) error {
	if filter == nil {
		return nil
	}

	fin := filter.GetIn()
	fout := filter.GetOut()

	agentIDs := make(map[int]struct{})

	for i := 0; i < len(ac.agents); i++ {
		a := ac.agents[i]
		meta := a.Meta()

		agentIDs[meta.ID] = struct{}{}
	}

	for _, id := range fin {
		if _, ok := agentIDs[id]; !ok {
			return fmt.Errorf("%%%d is not a valid agent", id)
		}

		if id == self {
			return fmt.Errorf("Cannot filter in the agent itself")
		}
	}

	for _, id := range fout {
		if _, ok := agentIDs[id]; !ok {
			return fmt.Errorf("%%%d is not a valid agent", id)
		}

		if id == self {
			return fmt.Errorf("Cannot filter out the agent itself")
		}
	}

	return nil
}

func (ac *Collection) inputToRingbuf(a Agent) {
	a.InputToRingbuf(ac.errors, ac.output)
	ac.waitFinish(a)
}

func (ac *Collection) outputFromRingbuf(a Agent, filter *msg.Filter, filtersSinks bool) {
	if filtersSinks {
		a.OutputFromRingbuf(ac.stdout, ac.errors, ac.stdout, filter)
	} else {
		a.OutputFromRingbuf(ac.stdout, ac.errors, ac.output, filter)
	}

	// XXX: This may cause Wait() to be called twice. We get an error, but
	//      does no harm.
	ac.waitFinish(a)
}

func (ac *Collection) errorsFromRingbuf(a Agent, filter *msg.Filter) {
	// We both read and write on errors.
	a.OutputFromRingbuf(ac.errors, ac.errors, ac.errors, filter)

	// XXX: This may cause Wait() to be called twice. We get an error, but
	//      does no harm.
	ac.waitFinish(a)
}

func (ac *Collection) logFromRingbuf(a Agent) {
	logring := config.GetLogRingbuf()

	a.OutputFromRingbuf(logring, logring, logring, nil)

	// XXX: This may cause Wait() to be called twice. We get an error, but
	//      does no harm.
	ac.waitFinish(a)
}

func (ac *Collection) runAgent(a Agent) error {
	meta := a.Meta()

	// Notice that this validation doesn't make safe access to ac:
	// it is safe to call it here, but not an a go routine.
	if err := ac.validateFilter(meta.ID, meta.Filter); err != nil {
		return err
	}

	filtersSinks := ac.isFilteringSinkAgents(meta.Filter)

	meta.Started = time.Now()
	meta.Status = AgentStatusRunning

	go func(ac *Collection, a Agent) {
		switch meta.Role {
		case AgentRoleSource:
			ac.inputToRingbuf(a)
		case AgentRoleErrors:
			ac.errorsFromRingbuf(a, meta.Filter)
		case AgentRoleSink:
			ac.outputFromRingbuf(a, meta.Filter, filtersSinks)
		case AgentRoleLog:
			ac.logFromRingbuf(a)
		}
	}(ac, a)

	return nil
}
