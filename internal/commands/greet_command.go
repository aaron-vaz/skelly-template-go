package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/{{.Inputs.Owner.Value}}/{{.Inputs.ProjectName.Value}}/internal/view"
)

// GreetCommand is an example command showing how flags and arguments are
// handled; replace it with your own.
type GreetCommand struct {
	ui   view.UI
	name string
}

// Name returns the CLI identifier for this command.
func (c *GreetCommand) Name() string {
	return "greet"
}

// Description returns a short summary shown in the help output.
func (c *GreetCommand) Description() string {
	return "Print a greeting"
}

// Init binds the command's flags.
func (c *GreetCommand) Init(flags *flag.FlagSet) {
	flags.StringVar(&c.name, "name", "World", "Who to greet")
}

// Run prints the greeting.
func (c *GreetCommand) Run(ctx context.Context, args []string) error {
	if len(args) > 0 {
		return errors.New("greet takes no arguments")
	}
	return c.ui.RenderInfo(fmt.Sprintf("Hello, %s!", c.name))
}

// NewGreetCommand initializes and returns the greet command.
func NewGreetCommand(ui view.UI) Command {
	return &GreetCommand{ui: ui}
}
