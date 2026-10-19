package commands

import (
	"context"
	"flag"
	"strings"
	"testing"
)

type mockUI struct {
	infoMsgs []string
}

func (m *mockUI) RenderInfo(message string) error {
	m.infoMsgs = append(m.infoMsgs, message)
	return nil
}
func (m *mockUI) RenderError(message string) error { return nil }

func TestVersionCommand_Run(t *testing.T) {
	// Given
	ui := &mockUI{}
	cmd := NewVersionCommand(ui)

	// When
	err := cmd.Run(context.Background(), nil)

	// Then
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	// And
	if len(ui.infoMsgs) != 1 || !strings.Contains(ui.infoMsgs[0], AppName+" version") {
		t.Errorf("Expected version output, got %v", ui.infoMsgs)
	}
}

func TestGreetCommand_Run(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		flags   []string
		want    string
		wantErr bool
	}{
		{name: "default name", want: "Hello, World!"},
		{name: "custom name", flags: []string{"-name", "Gopher"}, want: "Hello, Gopher!"},
		{name: "unexpected argument", args: []string{"extra"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			ui := &mockUI{}
			cmd := NewGreetCommand(ui)
			flags := flag.NewFlagSet("greet", flag.ContinueOnError)
			cmd.Init(flags)
			if err := flags.Parse(tt.flags); err != nil {
				t.Fatalf("failed to parse flags: %v", err)
			}

			// When
			err := cmd.Run(context.Background(), tt.args)

			// Then
			if (err != nil) != tt.wantErr {
				t.Fatalf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			// And
			if !tt.wantErr && (len(ui.infoMsgs) != 1 || ui.infoMsgs[0] != tt.want) {
				t.Errorf("output = %v, want %q", ui.infoMsgs, tt.want)
			}
		})
	}
}

func TestHelpCommand_Run(t *testing.T) {
	// Given
	ui := &mockUI{}
	cmds := map[string]Command{}
	for _, cmd := range []Command{NewGreetCommand(ui), NewVersionCommand(ui)} {
		cmds[cmd.Name()] = cmd
	}
	help := NewHelpCommand(cmds, ui)

	t.Run("lists every command", func(t *testing.T) {
		// When
		err := help.Run(context.Background(), nil)

		// Then
		if err != nil {
			t.Fatalf("Run() error = %v", err)
		}

		// And
		out := ui.infoMsgs[len(ui.infoMsgs)-1]
		for _, want := range []string{"USAGE:", "greet", "version"} {
			if !strings.Contains(out, want) {
				t.Errorf("Expected %q in help output, got %q", want, out)
			}
		}
	})

	t.Run("shows the flags of one command", func(t *testing.T) {
		// When
		err := help.Run(context.Background(), []string{"greet"})

		// Then
		if err != nil {
			t.Fatalf("Run() error = %v", err)
		}

		// And
		out := ui.infoMsgs[len(ui.infoMsgs)-1]
		if !strings.Contains(out, "-name") {
			t.Errorf("Expected the -name flag in help output, got %q", out)
		}
	})

	t.Run("unknown command is an error", func(t *testing.T) {
		// When
		err := help.Run(context.Background(), []string{"nope"})

		// Then
		if err == nil || !strings.Contains(err.Error(), "unknown command: nope") {
			t.Errorf("Run() error = %v, want unknown command", err)
		}
	})
}
