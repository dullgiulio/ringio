package agents

type AgentRole int

const (
	AgentRoleSink AgentRole = iota
	AgentRoleSource
	AgentRoleErrors
	AgentRoleLog
)

func (r AgentRole) String() string {
	switch r {
	case AgentRoleSink:
		return "->"
	case AgentRoleSource:
		return "<-"
	case AgentRoleErrors:
		return "&&"
	case AgentRoleLog:
		return "||"
	}

	return ""
}

func (r AgentRole) Text() string {
	switch r {
	case AgentRoleSink:
		return "output"
	case AgentRoleSource:
		return "input"
	case AgentRoleErrors:
		return "errors"
	case AgentRoleLog:
		return "logs"
	}

	return ""
}
