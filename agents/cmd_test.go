package agents

import (
	"testing"
)

func TestCmdIsAgent(t *testing.T) {
	p := NewAgentCmd([]string{"test", "test"}, AgentRoleSource)
	if Agent(p) == nil {
		t.Error("Cast failed")
	}
}
