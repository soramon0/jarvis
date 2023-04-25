package commands

import (
	"flag"
	"fmt"

	"github.com/soramon0/jarvis/pkg/project"
)

var (
	createCmd = flag.NewFlagSet("create", flag.ExitOnError)
	deleteCmd = flag.NewFlagSet("delete", flag.ExitOnError)
	commands  = map[string]*flag.FlagSet{
		createCmd.Name(): createCmd,
		deleteCmd.Name(): deleteCmd,
	}
)

type Cmd struct {
	*flag.FlagSet
}

func Run(args []string) error {
	cmd, err := parseCommand(args)
	if err != nil {
		return err
	}

	return cmd.Execute()
}

func (c *Cmd) Execute() error {
	if c == nil {
		return fmt.Errorf("uknown command")
	}

	switch c.Name() {
	case "create":
		args, err := c.ParseCreateProjectArgs()
		if err != nil {
			return err
		}

		p, err := project.Create(args)
		if err != nil {
			return err
		}
		fmt.Printf("Your %q project is ready at\n%s\n", p.FriendlyName(), p.AbsPath())
	case "delete":
		args, err := c.ParseDeleteProjectArgs()
		if err != nil {
			return err
		}

		p, err := project.Delete(args)
		if err != nil {
			return err
		}
		if p != nil {
			fmt.Printf("Your %q project has been deleted\n", p.FriendlyName())
		} else {
			fmt.Printf("Your %q project has been deleted\n", args.ProjectName)
		}
	default:
		return fmt.Errorf("only create command supported")
	}

	return nil
}

func parseCommand(args []string) (*Cmd, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("a command is required")
	}

	cmd := commands[args[1]]
	if cmd == nil {
		return nil, fmt.Errorf("uknown %q command", args[1])
	}

	return &Cmd{cmd}, nil
}
