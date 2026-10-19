package view

import (
	"fmt"
	"io"
)

// StdUI is a standard I/O implementation of the UI interface, writing to the
// provided streams.
type StdUI struct {
	stdout io.Writer
	stderr io.Writer
}

// RenderInfo writes an informational message to stdout.
func (s *StdUI) RenderInfo(message string) error {
	_, err := fmt.Fprintf(s.stdout, "%s\n", message)
	return err
}

// RenderError writes an error message to stderr.
func (s *StdUI) RenderError(message string) error {
	_, err := fmt.Fprintf(s.stderr, "%s\n", message)
	return err
}

// NewStdUI returns a ready-to-use StdUI backed by the provided streams.
func NewStdUI(stdout io.Writer, stderr io.Writer) UI {
	return &StdUI{
		stdout: stdout,
		stderr: stderr,
	}
}
