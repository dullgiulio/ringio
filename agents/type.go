package agents

type AgentType int

const (
	AgentTypeNull AgentType = iota
	AgentTypeCmd
	AgentTypePipe
)

func (t AgentType) String() string {
	switch t {
	case AgentTypeCmd:
		return "$"
	case AgentTypePipe:
		return "|"
	}

	return "?"
}
