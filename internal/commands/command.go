// Package commands holds the implementation of each CLI command.
package commands

import (
	"context"
	"flag"
)

// AppName is the name of the binary, used in help and version output.
const AppName = "{{.Inputs.ProjectName.Value}}"

// Command represents an executable command in the application
type Command interface {
	// Name returns the name of the command.
	Name() string
	// Description returns a short description of the command.
	Description() string
	// Init registers the command's flags on the given flag set.
	Init(flags *flag.FlagSet)
	// Run executes the command with the positional arguments left after flag parsing.
	Run(ctx context.Context, args []string) error
}

// Invoker executes commands with provided arguments
type Invoker interface {
	// Execute runs a command with the given arguments and returns an error if the execution fails
	Execute(ctx context.Context, args []string) error
}
