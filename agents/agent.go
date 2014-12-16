package agents

import (
	"fmt"
	"strings"
	"time"

	"github.com/dullgiulio/ringbuf"
	"github.com/dullgiulio/ringio/config"
	"github.com/dullgiulio/ringio/msg"
)

type AgentType int

const (
	AgentTypeNull AgentType = iota
	AgentTypeCmd
	AgentTypePipe
)

type AgentStatus int

const (
	AgentStatusNone AgentStatus = iota
	AgentStatusRunning
	AgentStatusKilled
	AgentStatusStopped
	AgentStatusFinished
)

type AgentRole int

const (
	AgentRoleSink AgentRole = iota
	AgentRoleSource
	AgentRoleErrors
	AgentRoleLog
)

func (s AgentStatus) IsRunning() bool {
	return s == AgentStatusRunning
}

func (s AgentStatus) String() string {
	switch s {
	case AgentStatusRunning:
		return "R"
	case AgentStatusStopped:
		return "S"
	case AgentStatusKilled:
		return "K"
	case AgentStatusFinished:
		return "F"
	}

	return "?"
}

type AgentMetadata struct {
	Id       int
	Role     AgentRole
	Status   AgentStatus
	Started  time.Time
	Finished time.Time
	Filter   *msg.Filter
}

type AgentDescr struct {
	Args []string
	Meta AgentMetadata
	Type AgentType
}

func (a *AgentDescr) String() string {
	var args string

	flow := "->"

	if a.Type == AgentTypeCmd {
		args = strings.Join(a.Args, " ")
	} else if a.Type == AgentTypePipe {
		args = "[pipe]"
	}

	if a.Meta.Role == AgentRoleSource {
		flow = "<-"
	}

	if a.Meta.Role == AgentRoleSink &&
		a.Meta.Filter != nil {
		flow = fmt.Sprintf("-> [%s]", a.Meta.Filter.String())
	}

	return fmt.Sprintf("%d %s %s %s",
		a.Meta.Id, a.Meta.Status.String(), flow, args)
}

func (a *AgentDescr) Text() string {
	var started, finished string

	str := a.String()

	if a.Meta.Status != AgentStatusNone {
		started = fmt.Sprintf("  Started: %s\n", a.Meta.Started.Format("2006-01-02 15:04:05 -0700 MST"))
	}

	if a.Meta.Status == AgentStatusFinished {
		finished = fmt.Sprintf("  Finished: %s\n", a.Meta.Finished.Format("2006-01-02 15:04:05 -0700 MST"))
	}

	return fmt.Sprintf("%s\n%s%s", str, started, finished)
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
				if meta.Id == id {
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

	agentIds := make(map[int]struct{})

	for i := 0; i < len(ac.agents); i++ {
		a := ac.agents[i]
		meta := a.Meta()

		agentIds[meta.Id] = struct{}{}
	}

	for _, id := range fin {
		if _, ok := agentIds[id]; !ok {
			return fmt.Errorf("%%%d is not a valid agent", id)
		}

		if id == self {
			return fmt.Errorf("Cannot filter in the agent itself")
		}
	}

	for _, id := range fout {
		if _, ok := agentIds[id]; !ok {
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
	if err := ac.validateFilter(meta.Id, meta.Filter); err != nil {
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
