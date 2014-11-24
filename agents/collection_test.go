package agents

import (
	"testing"

	"bitbucket.org/dullgiulio/ringio/config"
	"bitbucket.org/dullgiulio/ringio/log"
)

func TestAgentAdding(t *testing.T) {
	config.C.AutoExit = false

	go func() {
		nw := new(log.NullWriter)
		log.AddWriter(nw)

		if !log.Run(log.LevelError) {
			t.Log("Did not expect logger to return an error")
			t.Fail()
		}
	}()

	ac := NewCollection()
	go ac.Run(false)

	a := NewAgentCmd([]string{"/bin/ls", "/home"}, AgentRoleSource)
	resp := NewAgentMessageResponseBool()
	ac.Add(a, &resp)

	resp.Get()

	ac.Cancel(&resp)
	resp.Get()
}

// TODO: More tests using AgentNull
