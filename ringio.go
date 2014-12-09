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
	} else {
		cli.Run()
	}

	onexit.Exit(0)
}

func help() int {
	fmt.Printf(
		`Usage: ringio <session-name> open &
       ringio <session-name> input [COMMAND...]
       ringio <session-name> output [COMMAND...]
       ringio <session-name> run
       ringio <session-name> close
`)
	return 1
}
