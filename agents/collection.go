package agents

import (
	"fmt"
	"sort"
	"time"

	"github.com/dullgiulio/ringbuf"
	"github.com/dullgiulio/ringio/config"
	"github.com/dullgiulio/ringio/log"
)

type messageType int

const (
	agentMessageAdd messageType = iota
	agentMessageStart
	agentMessageStop
	agentMessageKill
	agentMessageFinished
	agentMessageList
	agentMessageStartAll
	agentMessageCancel
	agentMessageWriteReady
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

func (c *Collection) Start(a Agent, response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageStart, response, a)
}

func (c *Collection) StartAll(response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageStartAll, response, nil)
}

func (c *Collection) Finished(a Agent, response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageFinished, response, a)
}

func (c *Collection) Kill(a Agent, response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageKill, response, a)
}

func (c *Collection) Stop(a Agent, response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageStop, response, a)
}

func (c *Collection) List(response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageList, response, nil)
}

func (c *Collection) WriteReady(a Agent, response AgentMessageResponse) {
	c.requestCh <- newAgentMessage(agentMessageWriteReady, response, a)
}

func (c *Collection) add(newAgent Agent) int {
	maxID := 0

	for _, a := range c.agents {
		meta := a.Meta()

		if id := meta.ID; id > maxID {
			maxID = id
		}
	}

	maxID++

	meta := newAgent.Meta()
	meta.ID = maxID

	c.agents = append(c.agents, newAgent)
	return maxID
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
		metaI.ID < metaJ.ID {
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
		if a == agent || a.Meta().ID == meta.ID {
			realAgent = a
			break
		}
	}

	return realAgent
}

func (c *Collection) waitFinish(agent Agent) {
	if err := agent.WaitFinish(); err != nil {
		// Something went wrong, this agent is not finished.
		log.Error(log.FacilityAgent, fmt.Sprintf("Error waiting for %d: %s", agent.Meta().ID, err))
	}

	// We signal back to the collection that this agents is finished.
	resp := NewAgentMessageResponseBool()
	c.Finished(agent, &resp)

	if _, err := resp.Get(); err != nil {
		log.Error(log.FacilityAgent, fmt.Sprintf("Error cleaning up %d: %s", agent.Meta().ID, err))
	}
}

func (c *Collection) stopOrKillAgent(msg agentMessage, kill bool) {
	realAgent := c.getRealAgent(msg.agent)

	if realAgent == nil {
		msg.response.Err(fmt.Errorf("Agent %d not found", msg.agent.Meta().ID))
		return
	}

	meta := realAgent.Meta()
	isRunning := meta.Status.IsRunning()

	if !kill {
		// Cannot stop an Agent that is not in running state.
		if !isRunning {
			msg.response.Err(fmt.Errorf("Agent %d is not marked as running", meta.ID))
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

func (c *Collection) startAgent(agent Agent) error {
	meta := agent.Meta()

	switch meta.Status {
	case AgentStatusKilled, AgentStatusStopped, AgentStatusRunning:
		return fmt.Errorf("Agent %d is already running", meta.ID)
	}

	agent.Init()

	if err := c.runAgent(agent); err != nil {
		return err
	}

	log.Info(log.FacilityAgent, fmt.Sprintf("Started agent %d", meta.ID))
	return nil
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

			msg.response.Data(id)
			msg.response.Ok()
		case agentMessageKill:
			c.stopOrKillAgent(msg, true)
		case agentMessageStop:
			c.stopOrKillAgent(msg, false)
		case agentMessageStart:
			realAgent := c.getRealAgent(msg.agent)

			if realAgent == nil {
				err := fmt.Errorf("Agent %d not found", msg.agent.Meta().ID)
				log.Error(log.FacilityAgent, err)
				msg.response.Err(err)
				continue
			}

			if err := c.startAgent(realAgent); err != nil {
				log.Error(log.FacilityAgent, err)
				msg.response.Err(err)
			} else {
				msg.response.Ok()
			}
		case agentMessageCancel:
			// Kill all outstanding processes.
			for _, a := range c.agents {
				meta := a.Meta()

				if !meta.Status.IsRunning() {
					continue
				}

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

			c.Close()
			msg.response.Ok()
			return
		case agentMessageStartAll:
			// Start all agents that haven't been started yet.
			for _, a := range c.agents {
				meta := a.Meta()

				if meta.Status != AgentStatusNone {
					continue
				}

				if err := c.runAgent(a); err != nil {
					log.Error(log.FacilityAgent, err.Error())
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
		case agentMessageWriteReady:
			realAgent := c.getRealAgent(msg.agent)
			realAgent.StartWrite()
			msg.response.Ok()
		}
	}
}

func (c *Collection) Close() {
	c.output.Cancel()
	c.errors.Cancel()
	c.stdout.Cancel()
}
