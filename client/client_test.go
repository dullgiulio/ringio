package client

import (
	"testing"
)

func TestArgsParsing(t *testing.T) {
	cli := NewCli()
	err := cli.ParseArgs([]string{"ringio", "session-name", "input"})

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if cli.Session != "session-name" {
		t.Error("Expected 'session-name' as session name")
	}

	if cli.CommandStr != "input" {
		t.Error("Expected 'input' as command name")
	}

	err = cli.ParseArgs([]string{"ringio", "session-name", "output", "tail", "-f", "/some/file"})

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if cli.CommandStr != "output" {
		t.Error("Expected 'output' as command name")
	}

	if cli.Session != "session-name" {
		t.Error("Expected 'session-name' as session name")
	}

	if len(cli.NArgs) != 0 {
		t.Error("Expected all arguments to be parsed")
	}

	if cli.Args[0] != "tail" ||
		cli.Args[1] != "-f" ||
		cli.Args[2] != "/some/file" {
		t.Error("Expected all args after command name")
	}
}

func TestCorrectClient(t *testing.T) {
	cli := NewCli()
	c := cli.getClient()

	if c != nil {
		t.Error("Invalid client instance")
	}

	cli.Command = Output
	c = cli.getClient()
	if _, ok := c.(*CommandOutput); !ok {
		t.Error("Invalid client instance")
	}

	cli.Command = Input
	c = cli.getClient()
	if _, ok := c.(*CommandInput); !ok {
		t.Error("Invalid client instance")
	}

	cli.Command = IO
	c = cli.getClient()
	if _, ok := c.(*CommandIO); !ok {
		t.Error("Invalid client instance")
	}

	cli.Command = Log
	c = cli.getClient()
	if _, ok := c.(*CommandLog); !ok {
		t.Error("Invalid client instance")
	}

	cli.Command = Run
	c = cli.getClient()
	if _, ok := c.(*CommandRun); !ok {
		t.Error("Invalid client instance")
	}

	cli.Command = Open
	c = cli.getClient()
	if _, ok := c.(*CommandOpen); !ok {
		t.Error("Invalid client instance")
	}

	cli.Command = Close
	c = cli.getClient()
	if _, ok := c.(*CommandClose); !ok {
		t.Error("Invalid client instance")
	}

	cli.Command = List
	c = cli.getClient()
	if _, ok := c.(*CommandList); !ok {
		t.Error("Invalid client instance")
	}

	cli.Command = Stop
	c = cli.getClient()
	if _, ok := c.(*CommandStop); !ok {
		t.Error("Invalid client instance")
	}

	cli.Command = Set
	c = cli.getClient()
	if _, ok := c.(*CommandSet); !ok {
		t.Error("Invalid client instance")
	}
}
