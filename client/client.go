package client

import (
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
	Ping
	listSessions
	help
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
	haveArgs   bool
}

var cmdCommand = map[string]Command{
	"input":  Input,
	"in":     Input,
	"output": Output,
	"out":    Output,
	"io":     IO,
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
	"ping":   Ping,
}

func NewCli() *Cli {
	return &Cli{}
}

func (cli *Cli) ParseArgs(args []string) error {
	if len(args) <= 1 {
		cli.Help()
	} else {
		cli.Session = args[1]
	}

	switch cli.Session {
	case "help":
		if err := cli.initCommandArgs(args); err != nil {
			return err
		}
		cli.commandHelp()
	case "-ls":
		cli.listSessions()
	}

	if err := cli.initCommandArgs(args); err != nil {
		return err
	}

	cli.NArgs, cli.Args = cli.parseOptions(cli.NArgs)
	if err := cli.flagset.Parse(cli.NArgs); err != nil {
		return err
	}

	return nil
}

func (cli *Cli) listSessions() {
	sessions := getAllSessions()

	if len(sessions) == 0 {
		fmt.Printf("No sessions open. Please run 'ringio <session-name> open &' to start a session\n")
		onexit.Exit(0)
	}

	fmt.Printf("Available open sessions:\n\n")
	printList(sessions)

	fmt.Printf("\nPlease use 'ringio <session-name> ping' to verify a session is active.\n")
	onexit.Exit(1)
}

func (cli *Cli) initCommandArgs(args []string) error {
	errNoCommand := fmt.Errorf("No command specified")

	cli.argsLen = len(args) - 3
	if cli.argsLen < 0 {
		return errNoCommand
	}

	if args[2] == "" {
		return errNoCommand
	}

	cli.CommandStr = args[2]

	if command, ok := cmdCommand[cli.CommandStr]; !ok {
		return fmt.Errorf("Unsupported command %s", cli.CommandStr)
	} else {
		cli.Command = command
	}

	cli.CommandStr = args[2]
	cli.NArgs = args[3:]

	cli.flagset = flag.NewFlagSet("`ringio "+cli.CommandStr+"'", flag.ExitOnError)

	if cli.client = cli.getCommand(); cli.client == nil {
		return fmt.Errorf("Unsupported command %s", cli.CommandStr)
	}

	cli.haveArgs = cli.client.Init(cli.flagset)
	return nil
}

func (cli *Cli) Help() {
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
	onexit.Exit(1)
}

func (cli *Cli) commandHelp() {
	fmt.Fprintf(os.Stderr, "ringio: %s: ", cli.CommandStr)
	fmt.Fprintf(os.Stderr, cli.client.Help())
	fmt.Fprint(os.Stderr, "\n")

	if cli.haveArgs {
		fmt.Fprintf(os.Stderr, "\nSupported options:\n\n")
		cli.flagset.PrintDefaults()
	}

	fmt.Fprintf(os.Stderr, "\nRun `ringio' without argument to see basic usage information.\n")
	onexit.Exit(1)
}

// XXX: "-option arg" is not accepted, but only "-option=arg".
func (cli *Cli) parseOptions(args []string) ([]string, []string) {
	var nargs []string

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
	case Ping:
		return NewCommandPing()
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
