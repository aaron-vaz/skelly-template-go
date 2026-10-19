{{- if eq .Inputs.ProjectType.Value "cli" -}}
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/{{.Inputs.Owner.Value}}/{{.Inputs.ProjectName.Value}}/internal/cli"
	"github.com/{{.Inputs.Owner.Value}}/{{.Inputs.ProjectName.Value}}/internal/view"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ui := view.NewStdUI(os.Stdout, os.Stderr)
	invoker := cli.NewFlagCommandInvoker(ui)

	// Trapping the first Ctrl+C lets commands stop cleanly; once it fires, the
	// default handler is restored so a second Ctrl+C still kills the process.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	context.AfterFunc(ctx, stop)

	return invoker.Execute(ctx, os.Args[1:])
}
{{- else -}}
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/{{.Inputs.Owner.Value}}/{{.Inputs.ProjectName.Value}}/internal/config"
	"github.com/{{.Inputs.Owner.Value}}/{{.Inputs.ProjectName.Value}}/internal/server"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load(os.Getenv)
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))

	// The first SIGINT/SIGTERM starts a graceful shutdown; once it fires, the
	// default handler is restored so a second one still kills the process.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	context.AfterFunc(ctx, stop)

	ln, err := (&net.ListenConfig{}).Listen(ctx, "tcp", cfg.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", cfg.Addr, err)
	}

	return server.New(logger, cfg.ShutdownTimeout).Run(ctx, ln)
}
{{- end}}
