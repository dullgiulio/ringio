package client

import (
	"errors"
	"flag"
	"fmt"
	"net/rpc"
	"os"
	"strconv"

	"github.com/dullgiulio/ringio/msg"
	"github.com/dullgiulio/ringio/onexit"
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
	Start
	Stop
	Kill
)

type Client interface {
	Help() string
	Init(fs *flag.FlagSet) bool
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
		"start":  Start,
		"stop":   Stop,
		"kill":   Kill,
	}

	cli.argsLen = len(args) - 3
	if cli.argsLen < 0 {
		return nil
	}

	helpMode := false

	cli.Session = args[1]

	if cli.Session == "help" {
		cli.Session = ""
		helpMode = true
	}

	cli.CommandStr = args[2]
	cli.Command = cmdCommand[cli.CommandStr]
	cli.NArgs = args[3:]

	cli.flagset = flag.NewFlagSet("`ringio "+cli.CommandStr+"'", flag.ExitOnError)
	cli.client = cli.getCommand()

	if cli.client == nil {
		return errors.New(fmt.Sprintf("Unsupported command %s", cli.CommandStr))
	}

	haveArgs := cli.client.Init(cli.flagset)
	if helpMode {
		fmt.Fprintf(os.Stderr, "ringio: %s: ", cli.CommandStr)
		fmt.Fprintf(os.Stderr, cli.client.Help())
		fmt.Fprint(os.Stderr, "\n")

		if haveArgs {
			fmt.Fprintf(os.Stderr, "\nSupported options:\n\n")
			cli.flagset.PrintDefaults()
		}

		fmt.Fprintf(os.Stderr, "\nRun `ringio' without argument to see basic usage information.\n")
		onexit.Exit(1)
	}

	cli.NArgs, cli.Args = cli.parseOptions(cli.NArgs)

	if err := cli.flagset.Parse(cli.NArgs); err != nil {
		return err
	}

	return nil
}

// XXX: "-option arg" is not accepted, but only "-option=arg".
func (cli *Cli) parseOptions(args []string) ([]string, []string) {
	nargs := make([]string, 0)

	for i := range args {
		arg := args[i]

		if args[i][0] == '%' {
			arg = args[i][1:]
		}

		if d, err := strconv.Atoi(arg); err == nil {
			if cli.Filter == nil {
				cli.Filter = msg.NewFilter()
			}

			if d < 0 {
				cli.Filter.Out(-d)
			} else {
				cli.Filter.In(d)
			}
		} else if args[i][0] == '-' {
			// It's a command argument
			arg := args[i]

			// Make --gnu-style args into -g-args
			if args[i][1] == '-' {
				arg = args[i][1:]
			}

			nargs = append(nargs, arg)
		} else {
			return nargs, args[i:]
		}
	}

	return nargs, []string{}
}

func (cli *Cli) getCommand() Client {
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
	case Start:
		return NewCommandStart()
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
