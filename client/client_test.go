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

	err = cli.ParseArgs([]string{"ringio", "session-name", "output", "-no-wait", "tail", "-f", "/some/file"})

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

	if len(cli.NArgs) != 1 {
		t.Error("Expected all arguments to be parsed")
	}

	if cli.NArgs[0] != "-no-wait" {
		t.Error("Expected all arguments to be parsed")
	}

	if cli.Args[0] != "tail" ||
		cli.Args[1] != "-f" ||
		cli.Args[2] != "/some/file" {
		t.Error("Expected all args after command name")
	}
}

func TestFilterArgs(t *testing.T) {
	cli := NewCli()
	nargs, args := cli.parseOptions([]string{
		"-my-arg", "%1", "%-2", "--other-arg", "4", "-3", "-last-arg", "command", "-cmdarg"})

	if cli.Filter.String() != "1,4,-2,-3" {
		t.Error("Expected '1,4,-2,-3' as filter")
	}

	if nargs[0] != "-my-arg" {
		t.Error("NArgs don't start properly")
	}

	if nargs[1] != "-other-arg" {
		t.Error("NArgs not parsed in the middle")
	}

	if nargs[2] != "-last-arg" {
		t.Error("NArgs not parsed in last position")
	}

	if args[0] != "command" {
		t.Error("Args don't start properly")
	}

	if args[1] != "-cmdarg" {
		t.Error("Args don't finish properly")
	}
}

func TestArgsAndFilterAfter(t *testing.T) {
	cli := NewCli()
	nargs, args := cli.parseOptions([]string{"-my-arg", "%1"})

	if cli.Filter.String() != "1" {
		t.Error("Expected '1' as filter")
	}

	if nargs[0] != "-my-arg" {
		t.Error("NArgs don't start properly")
	}

	if len(args) > 0 {
		t.Error("Args are not empty as expected")
	}
}

func TestArgsAndFilterBefore(t *testing.T) {
	cli := NewCli()
	nargs, args := cli.parseOptions([]string{"%1", "-my-arg"})

	if cli.Filter.String() != "1" {
		t.Error("Expected '1' as filter")
	}

	if nargs[0] != "-my-arg" {
		t.Error("NArgs don't start properly")
	}

	if len(args) > 0 {
		t.Error("Args are not empty as expected")
	}
}

func TestCorrectClient(t *testing.T) {
	cli := NewCli()
	c := cli.getCommand()

	if c != nil {
		t.Error("Invalid client instance")
	}

	cli.Command = Output
	c = cli.getCommand()
	if _, ok := c.(*CommandOutput); !ok {
		t.Error("Invalid client instance")
	}

	cli.Command = Input
	c = cli.getCommand()
	if _, ok := c.(*CommandInput); !ok {
		t.Error("Invalid client instance")
	}

	cli.Command = IO
	c = cli.getCommand()
	if _, ok := c.(*CommandIO); !ok {
		t.Error("Invalid client instance")
	}

	cli.Command = Log
	c = cli.getCommand()
	if _, ok := c.(*CommandLog); !ok {
		t.Error("Invalid client instance")
	}

	cli.Command = Run
	c = cli.getCommand()
	if _, ok := c.(*CommandRun); !ok {
		t.Error("Invalid client instance")
	}

	cli.Command = Open
	c = cli.getCommand()
	if _, ok := c.(*CommandOpen); !ok {
		t.Error("Invalid client instance")
	}

	cli.Command = Close
	c = cli.getCommand()
	if _, ok := c.(*CommandClose); !ok {
		t.Error("Invalid client instance")
	}

	cli.Command = List
	c = cli.getCommand()
	if _, ok := c.(*CommandList); !ok {
		t.Error("Invalid client instance")
	}

	cli.Command = Stop
	c = cli.getCommand()
	if _, ok := c.(*CommandStop); !ok {
		t.Error("Invalid client instance")
	}
}
