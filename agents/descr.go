package agents

import (
	"bytes"
	"fmt"
	"strings"
)

type AgentDescr struct {
	Args []string
	Meta AgentMetadata
	Type AgentType
}

func (a *AgentDescr) String() string {
	var args string

	if a.Type == AgentTypeCmd {
		args = strings.Join(a.Args, " ")
	} else if a.Type == AgentTypePipe {
		args = "[pipe]"
	}

	filter := ""

	if a.Meta.Role == AgentRoleSink &&
		a.Meta.Filter != nil {
		filter = fmt.Sprintf(" [%s]", a.Meta.Filter.String())
	}

	name := ""

	if a.Meta.Name != "" {
		name = " # " + a.Meta.Name
	}

	return fmt.Sprintf("%d %s %s%s %s%s",
		a.Meta.Id, a.Meta.Status.String(), a.Meta.Role.String(), filter, args, name)
}

func (a *AgentDescr) Text() string {
	var b bytes.Buffer

	dateFormat := "2006-01-02 15:04:05 -0700 MST"
	descr := "[pipe]"

	if a.Type == AgentTypeCmd {
		descr = strings.Join(a.Args, " ")
	}

	if a.Meta.Name != "" {
		descr += " # " + a.Meta.Name
	}

	fmt.Fprintf(&b, "%%%d %s agent\n", a.Meta.Id, a.Meta.Role.Text())
	fmt.Fprintf(&b, "status: %s\n", a.Meta.Status.Text())
	fmt.Fprintf(&b, "descr: %s\n", descr)

	if a.Meta.Status != AgentStatusNone {
		fmt.Fprintf(&b, "started: %s\n", a.Meta.Started.Format(dateFormat))
	} else {
		fmt.Fprintf(&b, "started: not started\n")
	}

	if a.Meta.Status == AgentStatusFinished {
		fmt.Fprintf(&b, "finished: %s\n", a.Meta.Finished.Format(dateFormat))
	} else {
		fmt.Fprintf(&b, "finished: not finished\n")
	}

	if a.Meta.User != "" {
		fmt.Fprintf(&b, "user: %s\n", a.Meta.User)
	}

	return b.String()
}
