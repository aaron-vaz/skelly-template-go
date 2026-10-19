package cli

import (
	"bytes"
	"context"
	"flag"
	"io"
	"strings"
	"testing"

	"github.com/{{.Inputs.Owner.Value}}/{{.Inputs.ProjectName.Value}}/internal/view"
)

func TestFlagCommandInvoker_Execute(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantOut    string
		wantErrMsg string
	}{
		{name: "no arguments shows help", wantOut: "USAGE:"},
		{name: "runs a command", args: []string{"greet"}, wantOut: "Hello, World!"},
		{name: "flag after the command", args: []string{"greet", "-name", "Gopher"}, wantOut: "Hello, Gopher!"},
		{name: "command help flag", args: []string{"greet", "-h"}, wantOut: "Who to greet"},
		{name: "unknown command", args: []string{"nope"}, wantErrMsg: "unknown command: nope"},
		{name: "unknown flag", args: []string{"greet", "-nope"}, wantErrMsg: "failed to parse flags for command 'greet'"},
		{name: "unexpected argument", args: []string{"greet", "extra"}, wantErrMsg: "greet takes no arguments"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			var stdout, stderr bytes.Buffer
			invoker := NewFlagCommandInvoker(view.NewStdUI(&stdout, &stderr))

			// When
			err := invoker.Execute(context.Background(), tt.args)

			// Then
			if tt.wantErrMsg != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Fatalf("Execute() error = %v, want error containing %q", err, tt.wantErrMsg)
				}
				return
			}
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			// And
			if !strings.Contains(stdout.String(), tt.wantOut) {
				t.Errorf("stdout = %q, want it to contain %q", stdout.String(), tt.wantOut)
			}
		})
	}
}

func TestParseInterspersed(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{name: "flags before arguments", args: []string{"-name", "x", "a", "b"}, want: []string{"a", "b"}},
		{name: "flags after arguments", args: []string{"a", "-name", "x", "b"}, want: []string{"a", "b"}},
		{name: "double dash ends flag parsing", args: []string{"--", "-name", "x"}, want: []string{"-name", "x"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			invoker := NewFlagCommandInvoker(view.NewStdUI(&bytes.Buffer{}, &bytes.Buffer{})).(*FlagCommandInvoker)
			cmd := invoker.commands["greet"]
			flags := flag.NewFlagSet("greet", flag.ContinueOnError)
			flags.SetOutput(io.Discard)
			cmd.Init(flags)

			// When
			got, err := parseInterspersed(flags, tt.args)

			// Then
			if err != nil {
				t.Fatalf("parseInterspersed() error = %v", err)
			}

			// And
			if strings.Join(got, " ") != strings.Join(tt.want, " ") {
				t.Errorf("parseInterspersed() = %v, want %v", got, tt.want)
			}
		})
	}
}
