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
