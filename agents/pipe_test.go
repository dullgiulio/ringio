package agents

import (
	"testing"
)

func TestPipeIsAgent(t *testing.T) {
	p := NewAgentPipe("test", AgentRoleSource, nil)
	if Agent(p) == nil {
		t.Error("Cast failed")
	}
}
