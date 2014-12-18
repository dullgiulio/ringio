package main

import (
	"fmt"
	"os"

	"github.com/dullgiulio/ringio/client"
	"github.com/dullgiulio/ringio/config"
	"github.com/dullgiulio/ringio/onexit"
	"github.com/dullgiulio/ringio/utils"
)

func main() {
	// Handle interrupt signals.
	onexit.HandleInterrupt()

	config.Init()

	cli := client.NewCli()
	if err := cli.ParseArgs(os.Args); err != nil {
		utils.Fatal(err)
	}

	if cli.Session == "" {
		help()
		onexit.Exit(1)
	}

	cli.Run()
	onexit.Exit(0)
}

func help() {
	fmt.Print(
		`Usage: ringio <session-name> open &
       ringio <session-name> in|input [%job...] [-%job...] [COMMAND...]
       ringio <session-name> out|output [%job...] [-%job...] [COMMAND...]
       ringio <session-name> io [%job...] [-%job...]
       ringio <session-name> run
       ringio <session-name> list
       ringio <session-name> start %job...
       ringio <session-name> stop %job...
       ringio <session-name> kill %job...
       ringio <session-name> log
       ringio <session-name> close

Type 'ringio help <command>' for help on any command.
`)
}
