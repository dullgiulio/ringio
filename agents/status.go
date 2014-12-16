package agents

type AgentStatus int

const (
	AgentStatusNone AgentStatus = iota
	AgentStatusRunning
	AgentStatusKilled
	AgentStatusStopped
	AgentStatusFinished
)

func (s AgentStatus) IsRunning() bool {
	return s == AgentStatusRunning
}

func (s AgentStatus) String() string {
	switch s {
	case AgentStatusRunning:
		return "R"
	case AgentStatusStopped:
		return "S"
	case AgentStatusKilled:
		return "K"
	case AgentStatusFinished:
		return "F"
	}

	return "?"
}

func (s AgentStatus) Text() string {
	switch s {
	case AgentStatusRunning:
		return "running"
	case AgentStatusStopped:
		return "stopped"
	case AgentStatusKilled:
		return "killed"
	case AgentStatusFinished:
		return "finished"
	}

	return "none"
}
