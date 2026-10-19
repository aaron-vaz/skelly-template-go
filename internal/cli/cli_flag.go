// Package cli parses command line arguments and dispatches them to commands.
package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/{{.Inputs.Owner.Value}}/{{.Inputs.ProjectName.Value}}/internal/commands"
	"github.com/{{.Inputs.Owner.Value}}/{{.Inputs.ProjectName.Value}}/internal/view"
)

// FlagCommandInvoker is an implementation of commands.Invoker that maps
// CLI commands to their respective handlers and parses flag inputs.
type FlagCommandInvoker struct {
	commands map[string]commands.Command
}

// Execute parses the provided arguments, determines the target command,
// parses its flags, and runs it.
func (i *FlagCommandInvoker) Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return i.commands["help"].Run(ctx, nil)
	}

	cmdName := args[0]
	cmd, ok := i.commands[cmdName]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmdName)
	}

	flags := flag.NewFlagSet(cmd.Name(), flag.ContinueOnError)
	// Help and errors are rendered by us, so silence the flag package's own output
	flags.SetOutput(io.Discard)
	cmd.Init(flags)

	positional, err := parseInterspersed(flags, args[1:])
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return i.commands["help"].Run(ctx, []string{cmdName})
		}
		return fmt.Errorf("failed to parse flags for command '%s': %w", cmdName, err)
	}

	return cmd.Run(ctx, positional)
}

// parseInterspersed parses flags that appear before or after positional
// arguments, e.g. "cmd <arg> -flag value". The flag package alone stops at
// the first positional argument. A "--" still ends flag parsing.
func parseInterspersed(flags *flag.FlagSet, args []string) ([]string, error) {
	var positional []string
	for {
		if err := flags.Parse(args); err != nil {
			return nil, err
		}

		rest := flags.Args()
		consumed := len(args) - len(rest)
		if len(rest) == 0 || (consumed > 0 && args[consumed-1] == "--") {
			return append(positional, rest...), nil
		}

		positional = append(positional, rest[0])
		args = rest[1:]
	}
}

// NewFlagCommandInvoker creates a FlagCommandInvoker with all supported CLI
// commands registered. Add new commands to the list below.
func NewFlagCommandInvoker(ui view.UI) commands.Invoker {
	invoker := &FlagCommandInvoker{
		commands: make(map[string]commands.Command),
	}

	cmds := []commands.Command{
		commands.NewGreetCommand(ui),
		commands.NewHelpCommand(invoker.commands, ui),
		commands.NewVersionCommand(ui),
	}

	for _, cmd := range cmds {
		invoker.commands[cmd.Name()] = cmd
	}

	return invoker
}
