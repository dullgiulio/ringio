package main

import (
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

	cli.Run()
	onexit.Exit(0)
}
