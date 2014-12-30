package agents

import (
	"testing"
)

func TestCmdIsAgent(t *testing.T) {
	p := NewAgentCmd([]string{"test", "test"}, AgentRoleSource, nil, &AgentOptions{})
	if Agent(p) == nil {
		t.Error("Cast failed")
	}
}
