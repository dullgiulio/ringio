package agents

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/dullgiulio/ringbuf"
	"github.com/dullgiulio/ringio/log"
	"github.com/dullgiulio/ringio/msg"
)

type AgentCmd struct {
	cmd         *exec.Cmd
	meta        *AgentMetadata
	cancelInCh  chan bool
	cancelOutCh chan bool
}

func NewAgentCmd(cmd []string, meta *AgentMetadata) *AgentCmd {
	a := &AgentCmd{
		cmd:         exec.Command(cmd[0], cmd[1:]...),
		meta:        meta,
		cancelInCh:  make(chan bool),
		cancelOutCh: make(chan bool),
	}
	a.cmd.Env = a.meta.Env
	return a
}

func (a *AgentCmd) Init() {
	if a.meta.Options == nil {
		a.meta.Options = &AgentOptions{}
	}
}

func (a *AgentCmd) Meta() *AgentMetadata {
	return a.meta
}

func (a *AgentCmd) Descr() AgentDescr {
	return AgentDescr{
		Args: a.cmd.Args,
		Meta: *a.meta,
		Type: AgentTypeCmd,
	}
}

func (a *AgentCmd) String() string {
	return strings.Join(a.cmd.Args, " ")
}

func (a *AgentCmd) cancelReading() {
	a.cancelInCh <- true // Stop reading from stderr
	a.cancelInCh <- true // Stop reading from stdout
}

func (a *AgentCmd) Stop() error {
	// Stop data pipes to this process.
	if a.meta.Status.IsRunning() {
		if a.meta.Role == AgentRoleSink {
			a.cancelOutCh <- true // Stop writing to this agent

			// We don't cancel reading because a sink might still print
			// useful information, for example aggregated results.
		} else {
			// In all other cases, stop reading immediately.
			a.cancelReading()
		}
	}

	// Subprocess should exit normally now, a successive Wait() will finish it.

	log.Info(log.FacilityAgent, fmt.Sprintf("CmdAgent %d has been stopped", a.meta.ID))
	return nil
}

func (a *AgentCmd) Kill() error {
	// If the agent is running, we need to stop all goroutines that
	// are reading from or writing to the agent.
	if a.meta.Status.IsRunning() {
		a.cancelReading()

		if a.meta.Role == AgentRoleSink {
			a.cancelOutCh <- true
		}
	}

	// An agent that was stopped only has a goroutine running.
	if a.meta.Status == AgentStatusStopped {
		if a.meta.Role == AgentRoleSink {
			a.cancelOutCh <- true
		}
	}

	if err := a.cmd.Process.Kill(); err != nil {
		return err
	}

	log.Info(log.FacilityAgent, fmt.Sprintf("CmdAgent %d has been killed", a.meta.ID))
	return nil
}

func (a *AgentCmd) WaitFinish() error {
	return a.cmd.Wait()
}

func (a *AgentCmd) InputToRingbuf(rErrors, rOutput *ringbuf.Ringbuf) {
	stdout, err := a.cmd.StdoutPipe()

	if err != nil {
		log.Error(log.FacilityAgent, err)
		return
	}

	stderr, err := a.cmd.StderrPipe()

	if err != nil {
		log.Error(log.FacilityAgent, err)
		return
	}

	if err := a.cmd.Start(); err != nil {
		log.Error(log.FacilityAgent, err)
		return
	}

	id := a.meta.ID

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go writeToRingbuf(id, stderr, rErrors, a.cancelInCh, wg)
	go writeToRingbuf(id, stdout, rOutput, a.cancelInCh, wg)

	wg.Wait()

	// Close stdout to trigger a SIGPIPE on next write.
	stderr.Close()
	stdout.Close()

	close(a.cancelInCh)
	close(a.cancelOutCh)
}

func (a *AgentCmd) OutputFromRingbuf(rStdout, rErrors, rOutput *ringbuf.Ringbuf, filter *msg.Filter) {
	stdout, err := a.cmd.StdoutPipe()

	if err != nil {
		log.Error(log.FacilityAgent, err)
		return
	}

	stderr, err := a.cmd.StderrPipe()

	if err != nil {
		log.Error(log.FacilityAgent, err)
		return
	}

	stdin, err := a.cmd.StdinPipe()

	if err != nil {
		log.Error(log.FacilityAgent, err)
		return
	}

	if err := a.cmd.Start(); err != nil {
		log.Error(log.FacilityAgent, err)
		return
	}

	id := a.meta.ID

	wg := new(sync.WaitGroup)
	wg.Add(3)

	go writeToRingbuf(id, stdout, rStdout, a.cancelInCh, wg)
	go writeToRingbuf(id, stderr, rErrors, a.cancelInCh, wg)
	go readFromRingbuf(stdin, filter, a.meta.Options.getMask(), rOutput, makeReaderOptions(a.meta.Options), a.cancelOutCh, wg)

	wg.Wait()

	close(a.cancelInCh)
	close(a.cancelOutCh)
}
