package commands

import (
	"context"
	"flag"
	"fmt"

	"github.com/{{.Inputs.owner.Value}}/{{.Inputs.project_name.Value}}/internal/version"
	"github.com/{{.Inputs.owner.Value}}/{{.Inputs.project_name.Value}}/internal/view"
)

// VersionCommand outputs the current build version.
type VersionCommand struct {
	ui view.UI
}

// Name returns the CLI identifier for this command.
func (c *VersionCommand) Name() string {
	return "version"
}

// Description returns a short summary shown in the help output.
func (c *VersionCommand) Description() string {
	return "Show " + AppName + " version"
}

// Init is a no-op as the version command has no flags.
func (c *VersionCommand) Init(flags *flag.FlagSet) {}

// Run displays the injected application version string to the user.
func (c *VersionCommand) Run(ctx context.Context, args []string) error {
	return c.ui.RenderInfo(fmt.Sprintf("%s version %s", AppName, version.Version))
}

// NewVersionCommand initializes and returns the version command.
func NewVersionCommand(ui view.UI) Command {
	return &VersionCommand{ui: ui}
}
