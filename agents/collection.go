package agents

import (
	"fmt"
	"sort"
	"time"

	"github.com/dullgiulio/ringbuf"
	"github.com/dullgiulio/ringio/config"
	"github.com/dullgiulio/ringio/log"
)

type CollectionRingType int

const (
	CollectionRingTypeErrors CollectionRingType = iota
	CollectionRingTypeStdout
	CollectionRingTypeOutput
)

type messageType int

const (
	agentMessageAdd messageType = iota
	agentMessageStop
	agentMessageKill
	agentMessageFinished
	agentMessageRunning
	agentMessageList
	agentMessageAutorun
	agentMessageCancel
)

type agentMessage struct {
	status   messageType
	response AgentMessageResponse
	// This might not be a real agent but just a container for the metadata.
	agent Agent
}

func newAgentMessage(status messageType, response AgentMessageResponse, agent Agent) agentMessage {
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
	c.requestCh <- newAgentMessage(agentMessageCancel, response, nil)
}

func (c *Collection) Add(a Agent, response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageAdd, response, a)
}

func (c *Collection) StartAll(response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageAutorun, response, nil)
}

func (c *Collection) SetAgentStatusFinished(a Agent, response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageFinished, response, a)
}

func (c *Collection) SetAgentStatusKill(a Agent, response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageKill, response, a)
}

func (c *Collection) SetAgentStatusStop(a Agent, response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageStop, response, a)
}

func (c *Collection) SetAgentStatusRunning(a Agent, response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageRunning, response, a)
}

func (c *Collection) List(response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageList, response, nil)
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

	reader := ringbuf.NewReader(r)
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

func (c *Collection) getRealAgent(agent Agent) Agent {
	var realAgent Agent

	meta := agent.Meta()

	for _, a := range c.agents {
		if a.Meta().Id == meta.Id {
			realAgent = a
			break
		}
	}

	return realAgent
}

func (c *Collection) waitFinish(agent Agent) {
	if err := agent.WaitFinish(); err != nil {
		// Something went wrong, this agent is not finished.
		log.Error(log.FacilityAgent, fmt.Sprintf("Error waiting for %d: %s", agent.Meta().Id, err))
	}

	// We signal back to the collection that this agents is finished.
	resp := NewAgentMessageResponseBool()
	c.SetAgentStatusFinished(agent, &resp)

	if _, err := resp.Get(); err != nil {
		log.Error(log.FacilityAgent, fmt.Sprintf("Error cleaning up %d: %s", agent.Meta().Id, err))
	}
}

func (c *Collection) stopOrKillAgent(msg agentMessage, kill bool) {
	realAgent := c.getRealAgent(msg.agent)

	if realAgent == nil {
		msg.response.Err(fmt.Errorf("Agent %d not found", msg.agent.Meta().Id))
		return
	}

	meta := realAgent.Meta()
	isRunning := meta.Status.IsRunning()

	if !kill {
		// Cannot stop an Agent that is not in running state.
		if !isRunning {
			msg.response.Err(fmt.Errorf("Agent %d is not marked as running", meta.Id))
			return
		}

		if err := realAgent.Stop(); err != nil {
			log.Error(log.FacilityAgent, err)
			msg.response.Err(err)
			return
		}

		meta.Status = AgentStatusStopped
	} else {
		if err := realAgent.Kill(); err != nil {
			log.Error(log.FacilityAgent, err)
			msg.response.Err(err)
			return
		}

		meta.Status = AgentStatusKilled
	}

	msg.response.Ok()

	// Only if it was running before it was stopped, we need
	// to wait for the Agent to clean up. If it wasn't running,
	// we must have a routine already waiting for it.
	if isRunning {
		go c.waitFinish(realAgent)
	}
}

func (c *Collection) Run(autorun bool) {
	go c.errors.Run()
	go c.output.Run()
	go c.stdout.Run()

	for msg := range c.requestCh {
		switch msg.status {
		case agentMessageAdd:
			id := c.add(msg.agent)

			log.Info(log.FacilityAgent, "Added new agent", msg.agent)

			if autorun {
				msg.agent.Init()

				if err := c.runAgent(msg.agent); err != nil {
					msg.response.Err(err)
					continue
				}

				log.Debug(log.FacilityAgent, "New agent started automatically")
			}

			msg.response.Data(id)
			msg.response.Ok()
		case agentMessageKill:
			c.stopOrKillAgent(msg, true)
		case agentMessageStop:
			c.stopOrKillAgent(msg, false)
		case agentMessageCancel:
			// Kill all outstanding processes.
			for _, a := range c.agents {
				meta := a.Meta()

				if meta.Status.IsRunning() {
					// Kill running agents.
					if err := a.Kill(); err != nil {
						log.Error(log.FacilityAgent, err)
					} else {
						meta.Status = AgentStatusKilled

						if err := a.WaitFinish(); err != nil {
							log.Error(log.FacilityAgent, err)
						}

						meta.Finished = time.Now()
						meta.Status = AgentStatusFinished
					}
				}
			}

			c.Close()
			msg.response.Ok()
			return
		case agentMessageAutorun:
			if autorun {
				// If we are already in autorun, ignore this message.
				msg.response.Ok()
				continue
			}

			autorun = true

			// Start all agents that haven't been started yet.
			for _, a := range c.agents {
				meta := a.Meta()

				if meta.Status == AgentStatusNone {
					if err := c.runAgent(a); err != nil {
						log.Error(log.FacilityAgent, err.Error())
					}
				}
			}

			log.Info(log.FacilityAgent, "Starting all agents")

			msg.response.Ok()
		case agentMessageFinished:
			meta := msg.agent.Meta()
			meta.Finished = time.Now()
			meta.Status = AgentStatusFinished

			log.Info(log.FacilityAgent, "Agent", msg.agent, "finished")

			msg.response.Ok()
		case agentMessageList:
			var agents []AgentDescr

			sort.Sort(c)

			for _, a := range c.agents {
				agents = append(agents, a.Descr())
			}

			msg.response.Data(&agents)
			msg.response.Ok()
		}
	}
}

func (c *Collection) Close() {
	c.output.Cancel()
	c.errors.Cancel()
	c.stdout.Cancel()
}
