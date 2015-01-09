package agents

import (
	"testing"
)

func TestPipeIsAgent(t *testing.T) {
	meta := NewAgentMetadata()
	meta.Role = AgentRoleSource
	p := NewAgentPipe("test", meta)
	if Agent(p) == nil {
		t.Error("Cast failed")
	}
}
