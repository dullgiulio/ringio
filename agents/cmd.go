package agents

import (
	"os/exec"
	"strings"
	"sync"

	"github.com/dullgiulio/ringbuf"
	"github.com/dullgiulio/ringio/log"
)

type AgentCmd struct {
	cmd    *exec.Cmd
	meta   *AgentMetadata
	cancel chan bool
}

func NewAgentCmd(cmd []string, role AgentRole) *AgentCmd {
	return &AgentCmd{
		cmd:    exec.Command(cmd[0], cmd[1:]...),
		meta:   &AgentMetadata{Role: role},
		cancel: make(chan bool),
	}
}

func (a *AgentCmd) Init() {
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

func (a *AgentCmd) Cancel() error {
	a.cancel <- true
	a.cancel <- true

	if a.meta.Role != AgentRoleSource {
		a.cancel <- true
	}

	return a.cmd.Process.Kill()
}

func (a *AgentCmd) InputToRingbuf(errors, output *ringbuf.Ringbuf) {
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

	id := a.meta.Id

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go writeToRingbuf(id, stderr, errors, a.cancel, wg)
	go writeToRingbuf(id, stdout, output, a.cancel, wg)

	wg.Wait()

	// Close stdout to trigger a SIGPIPE on next write.
	stderr.Close()
	stdout.Close()

	close(a.cancel)
}

func (a *AgentCmd) Stop() {
	if err := a.cmd.Wait(); err != nil {
		log.Error(log.FacilityAgent, err)
	}
}

func (a *AgentCmd) OutputFromRingbuf(rStdout, rErrors, rOutput *ringbuf.Ringbuf) {
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

	id := a.meta.Id

	wg := new(sync.WaitGroup)
	wg.Add(3)

	go writeToRingbuf(id, stdout, rStdout, a.cancel, wg)
	go writeToRingbuf(id, stderr, rErrors, a.cancel, wg)
	go readFromRingbuf(stdin, rOutput, a.cancel, wg)

	wg.Wait()

	close(a.cancel)
}
