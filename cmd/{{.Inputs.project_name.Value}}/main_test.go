{{- if eq .Inputs.project_type.Value "cli" -}}
package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

const binName = "{{.Inputs.project_name.Value}}_test_bin"

func TestMain(m *testing.M) {
	// Build the CLI binary
	build := exec.Command("go", "build", "-o", binName)
	if err := build.Run(); err != nil {
		_, _ = os.Stderr.WriteString("Could not build binary: " + err.Error() + "\n")
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Clean up
	_ = os.Remove(binName)

	os.Exit(code)
}

func TestIntegration(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantOut  string
		wantFail bool
	}{
		{name: "help", args: []string{"help"}, wantOut: "USAGE:"},
		{name: "no arguments shows help", wantOut: "USAGE:"},
		{name: "version", args: []string{"version"}, wantOut: "{{.Inputs.project_name.Value}} version dev"},
		{name: "greet", args: []string{"greet", "-name", "Gopher"}, wantOut: "Hello, Gopher!"},
		{name: "unknown command", args: []string{"nope"}, wantOut: "unknown command: nope", wantFail: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			cmd := exec.Command("./"+binName, tt.args...)

			// When
			out, err := cmd.CombinedOutput()

			// Then
			if (err != nil) != tt.wantFail {
				t.Fatalf("error = %v, wantFail %v: %s", err, tt.wantFail, out)
			}

			// And
			if !strings.Contains(string(out), tt.wantOut) {
				t.Errorf("output = %q, want it to contain %q", out, tt.wantOut)
			}
		})
	}
}
{{- else -}}
package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

const binName = "{{.Inputs.project_name.Value}}_test_bin"

func TestMain(m *testing.M) {
	// Build the service binary
	build := exec.Command("go", "build", "-o", binName)
	if err := build.Run(); err != nil {
		_, _ = os.Stderr.WriteString("Could not build binary: " + err.Error() + "\n")
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Clean up
	_ = os.Remove(binName)

	os.Exit(code)
}

func TestIntegration_ServesAndShutsDownGracefully(t *testing.T) {
	// Given a free port for the service to listen on
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	_ = ln.Close()

	cmd := exec.Command("./" + binName)
	cmd.Env = append(os.Environ(), "ADDR="+addr)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start service: %v", err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill() })

	// When it is up
	var resp *http.Response
	deadline := time.Now().Add(10 * time.Second)
	for {
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://"+addr+"/healthz", nil)
		resp, err = http.DefaultClient.Do(req)
		if err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("service did not start listening: %v", err)
		}
		time.Sleep(50 * time.Millisecond)
	}
	_ = resp.Body.Close()

	// Then it answers
	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET /healthz status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// And exits cleanly on SIGTERM
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("failed to signal service: %v", err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("service exited with %v, want a clean exit", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatalf("service did not exit after SIGTERM (%s)", addr)
	}
}

func TestIntegration_InvalidConfigFailsFast(t *testing.T) {
	// Given
	cmd := exec.Command("./" + binName)
	cmd.Env = append(os.Environ(), "LOG_LEVEL=loud")

	// When
	out, err := cmd.CombinedOutput()

	// Then
	if err == nil {
		t.Fatalf("Expected the service to fail, got output %q", out)
	}
}
{{- end}}
