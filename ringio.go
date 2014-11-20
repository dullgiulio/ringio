package main

import (
	"fmt"
	"os"

	"bitbucket.org/dullgiulio/ringio/client"
	"bitbucket.org/dullgiulio/ringio/onexit"
	"bitbucket.org/dullgiulio/ringio/utils"
)

func main() {
	// Handle interrupt signals.
	onexit.HandleInterrupt()

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
       ringio <session-name> set [verbose|quiet|locked]
       ringio <session-name> close
`)
	return 1
}
