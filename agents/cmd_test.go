package agents

import (
	"testing"
)

func TestCmdIsAgent(t *testing.T) {
	meta := NewAgentMetadata()
	meta.Role = AgentRoleSource
	p := NewAgentCmd([]string{"test", "test"}, meta)
	if Agent(p) == nil {
		t.Error("Cast failed")
	}
}
