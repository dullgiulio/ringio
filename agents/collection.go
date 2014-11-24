package agents

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"bitbucket.org/dullgiulio/ringbuf"
	"bitbucket.org/dullgiulio/ringio/config"
	"bitbucket.org/dullgiulio/ringio/log"
)

type CollectionRingType int

const (
	CollectionRingTypeErrors = iota
	CollectionRingTypeStdout
	CollectionRingTypeOutput
)

type agentMessageStatus int

const (
	agentMessageStatusAdd = iota
	agentMessageStatusStop
	agentMessageStatusKill
	agentMessageStatusFinished
	agentMessageStatusRunning
	agentMessageStatusList
	agentMessageStatusDone
	agentMessageStatusAutorun
	agentMessageStatusCancel
)

type agentMessage struct {
	status   agentMessageStatus
	response AgentMessageResponse
	agent    Agent
}

func newAgentMessage(status agentMessageStatus, response AgentMessageResponse, agent Agent) agentMessage {
	return agentMessage{status: status, response: response, agent: agent}
}

type Collection struct {
	agents    []Agent
	requestCh chan agentMessage
	output    *ringbuf.Ringbuf
	errors    *ringbuf.Ringbuf
	stdout    *ringbuf.Ringbuf
}

func NewCollection() *Collection {
	return &Collection{
		agents:    make([]Agent, 0),
		requestCh: make(chan agentMessage),
		// Output from source agents.
		output: ringbuf.NewRingbuf(config.C.RingbufSize),
		// All errors.
		errors: ringbuf.NewRingbuf(config.C.RingbufSize),
		// Output from sink agents.
		stdout: ringbuf.NewRingbuf(config.C.RingbufSize),
	}
}

func (c *Collection) Cancel(response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageStatusCancel, response, nil)
}

func (c *Collection) Done() {
	c.requestCh <- newAgentMessage(agentMessageStatusDone, nil, nil)
}

func (c *Collection) Add(a Agent, response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageStatusAdd, response, a)
}

func (c *Collection) StartAll(response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageStatusAutorun, response, nil)
}

func (c *Collection) SetAgentStatusFinished(a Agent, response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageStatusFinished, response, a)
}

func (c *Collection) SetAgentStatusKill(a Agent, response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageStatusKill, response, a)
}

func (c *Collection) SetAgentStatusRunning(a Agent, response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageStatusRunning, response, a)
}

func (c *Collection) List(response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageStatusList, response, nil)
}

func (c *Collection) Reader(which CollectionRingType) <-chan interface{} {
	var r *ringbuf.Ringbuf

	switch which {
	case CollectionRingTypeErrors:
		r = c.errors
	case CollectionRingTypeStdout:
		r = c.stdout
	case CollectionRingTypeOutput:
		r = c.output
	}

	reader := ringbuf.NewRingbufReader(r)
	return reader.ReadCh()
}

func (c *Collection) add(newAgent Agent) int {
	maxId := 0

	for _, a := range c.agents {
		meta := a.Meta()

		if id := meta.Id; id > maxId {
			maxId = id
		}
	}

	maxId++

	meta := newAgent.Meta()
	meta.Id = maxId

	c.agents = append(c.agents, newAgent)
	return maxId
}

// Implement a sort.Interface that sorts by (role,started date) and prints.
func (c *Collection) Len() int {
	l := len(c.agents)
	return l
}

func (c *Collection) Less(i, j int) bool {
	metaI := c.agents[i].Meta()
	metaJ := c.agents[j].Meta()

	if metaI.Status < metaJ.Status {
		return true
	}

	if metaI.Status == metaJ.Status &&
		metaI.Id < metaJ.Id {
		return true
	}

	return false
}

func (c *Collection) Swap(i, j int) {
	c.agents[i], c.agents[j] = c.agents[j], c.agents[i]
}

func (c *Collection) stopAgent(a Agent) {
	if err := a.Cancel(); err != nil {
		fmt.Printf("%s\n", err)
		log.Error(log.FacilityAgent, err)
		<-c.requestCh
	} else {
		<-c.requestCh
		a.Stop()
	}
}

func (c *Collection) Run(autorun bool) {
	var addingLocked, sorted bool

	go c.errors.Run()
	go c.output.Run()
	go c.stdout.Run()

	waitedAgents := make(map[AgentRole]int)

	for msg := range c.requestCh {
		switch msg.status {
		case agentMessageStatusAdd:
			if !addingLocked {
				// Will need to sort elements again.
				sorted = false

				id := c.add(msg.agent)

				log.Debug(log.FacilityAgent, "Added new agent", msg.agent)

				if autorun {
					meta := msg.agent.Meta()
					msg.agent.Init()

					c.runAgent(msg.agent)

					waitedAgents[meta.Role]++

					log.Debug(log.FacilityAgent, "New agent started automatically")
				}

				msg.response.Data(id)
				msg.response.Ok()
			} else {
				log.Error(log.FacilityAgent, "Tried to add agent after locking with Done()")
				msg.response.Err(errors.New("Tried to add agent after locking process"))
			}
		case agentMessageStatusKill:
			var realAgent Agent

			meta := msg.agent.Meta()

			for _, a := range c.agents {
				if a.Meta().Id == meta.Id {
					realAgent = a
					break
				}
			}

			meta = realAgent.Meta()

			if realAgent == nil {
				msg.response.Err(errors.New("Agent not found"))
			} else {
				if meta.Status == AgentStatusRunning {
					c.stopAgent(realAgent)
				}

				msg.response.Ok()
			}
		case agentMessageStatusCancel:
			// Kill all outstanding processes.
			for _, a := range c.agents {
				if a.Meta().Status == AgentStatusRunning {
					c.stopAgent(a)
				}
			}

			c.Close()
			msg.response.Ok()
			return
		case agentMessageStatusDone:
			log.Debug(log.FacilityAgent, "Adding agents has been disabled")
			addingLocked = true
		case agentMessageStatusAutorun:
			if autorun {
				continue
			}

			if config.C.AutoLock {
				addingLocked = true
			}

			autorun = true

			// Start all agents that haven't been started yet.
			for _, a := range c.agents {
				meta := a.Meta()

				if meta.Status == AgentStatusNone {
					meta := a.Meta()

					c.runAgent(a)

					waitedAgents[meta.Role]++
				}
			}

			log.Debug(log.FacilityAgent, "Starting all agents")

			msg.response.Ok()
		case agentMessageStatusFinished:
			sorted = false

			meta := msg.agent.Meta()
			meta.Finished = time.Now()
			meta.Status = AgentStatusFinished

			waitedAgents[meta.Role]--

			log.Debug(log.FacilityAgent, "Agent", msg.agent, "finished")

			msg.response.Ok()
		case agentMessageStatusList:
			var agents []AgentDescr

			if !sorted {
				sort.Sort(c)
				sorted = true
			}

			for _, a := range c.agents {
				agents = append(agents, a.Descr())
			}

			msg.response.Data(&agents)
			msg.response.Ok()
		}

		if addingLocked && autorun &&
			waitedAgents[AgentRoleSource] == 0 {
			c.output.Eof()

			// Must not send EOF to the ringbuf twice.
			waitedAgents[AgentRoleSource] = -1
		}

		if config.C.AutoExit &&
			waitedAgents[AgentRoleSource] <= 0 &&
			waitedAgents[AgentRoleSink] <= 0 &&
			waitedAgents[AgentRoleLog] <= 0 {
			log.Debug(log.FacilityAgent, "No more agents to wait for, exiting main loop")
			c.Close()
			return
		}
	}
}

func (c *Collection) Close() {
	c.output.Cancel()
	c.errors.Cancel()
	c.stdout.Cancel()

	close(c.requestCh)
}
