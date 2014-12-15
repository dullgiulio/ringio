package client

import (
	"errors"
	"flag"
	"fmt"
	"net/rpc"
	"strconv"

	"github.com/dullgiulio/ringio/msg"
	"github.com/dullgiulio/ringio/utils"
)

type Command int

const (
	None Command = iota
	Output
	Input
	IO
	Error
	Log
	Run
	Open
	Close
	List
	Stop
	Kill
)

type Client interface {
	Help() string
	Init(fs *flag.FlagSet) error
	Run(cli *Cli) error
}

type Cli struct {
	NArgs      []string
	Args       []string
	Session    string
	Command    Command
	CommandStr string
	Filter     *msg.Filter
	argsLen    int
	flagset    *flag.FlagSet
	client     Client
}

func NewCli() *Cli {
	return &Cli{}
}

func (cli *Cli) ParseArgs(args []string) error {
	var cmdCommand map[string]Command = map[string]Command{
		"input":  Input,
		"in":     Input,
		"output": Output,
		"out":    Output,
		"errors": Error,
		"error":  Error,
		"err":    Error,
		"log":    Log,
		"run":    Run,
		"open":   Open,
		"close":  Close,
		"list":   List,
		"ls":     List,
		"stop":   Stop,
		"kill":   Kill,
	}

	cli.argsLen = len(args) - 3
	if cli.argsLen < 0 {
		return nil
	}

	cli.Session = args[1]
	cli.CommandStr = args[2]
	cli.Command = cmdCommand[cli.CommandStr]
	cli.NArgs = args[3:]
	cli.flagset = flag.NewFlagSet("ringio-"+cli.CommandStr, flag.ExitOnError)
	cli.client = cli.getClient()

	if cli.client == nil {
		return errors.New(fmt.Sprintf("Unsupported command %s", cli.CommandStr))
	}

	if err := cli.client.Init(cli.flagset); err != nil {
		return err
	}

	cli.Filter = cli.parseFilter(cli.NArgs)

	if err := cli.flagset.Parse(cli.NArgs); err != nil {
		return err
	}

	n := cli.flagset.NArg()

	// For us, what is parsed and not as the opposite meaning.
	cli.Args = cli.NArgs[len(cli.NArgs)-n:]
	cli.NArgs = cli.NArgs[0 : len(cli.NArgs)-n]

	return nil
}

func (cli *Cli) parseFilter(args []string) (f *msg.Filter) {
	skipped := 0

	for i := range args {
		arg := args[i]

		if args[i][0] == '%' {
			arg = args[i][1:]
		}

		if d, err := strconv.Atoi(arg); err == nil {
			if skipped == 0 {
				f = msg.NewFilter()
			}

			if d < 0 {
				f.Out(-d)
			} else {
				f.In(d)
			}

			skipped++
		} else {
			break
		}
	}

	if skipped > 0 {
		cli.NArgs = cli.NArgs[skipped:]
	}

	return
}

func (cli *Cli) getClient() Client {
	if cli.Command == None {
		return nil
	}

	switch cli.Command {
	case Output:
		return NewCommandOutput()
	case Input:
		return NewCommandInput()
	case IO:
		return NewCommandIO()
	case Error:
		return NewCommandError()
	case Log:
		return NewCommandLog()
	case Run:
		return NewCommandRun()
	case Open:
		return NewCommandOpen()
	case Close:
		return NewCommandClose()
	case List:
		return NewCommandList()
	case Stop:
		return NewCommandStop()
	case Kill:
		return NewCommandKill()
	}

	return nil
}

func (cli *Cli) Run() {
	if err := cli.client.Run(cli); err != nil {
		utils.Fatal(err)
	}
}

func (cli *Cli) GetClient() *rpc.Client {
	client, err := rpc.Dial("unix", utils.FileInDotpath(cli.Session))
	if err != nil {
		utils.Fatal(err)
	}

	return client
}
