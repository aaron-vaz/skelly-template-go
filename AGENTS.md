# AGENTS.md

Guidelines for agentic coding agents working in this repository.

## Project Overview

`{{.Inputs.project_name.Value}}` is a Go {{if eq .Inputs.project_type.Value "cli"}}command line application{{else}}web service{{end}}.

- **Language**: Go {{.Inputs.go_version.Value}}
- **Binary**: `{{.Inputs.project_name.Value}}`
- **Module**: `github.com/{{.Inputs.owner.Value}}/{{.Inputs.project_name.Value}}`

## Build/Lint/Test Commands

```bash
make all        # Runs lint, test, and build
make build      # Build binary for current platform -> build/{{.Inputs.project_name.Value}}
make test       # Run all tests with the race detector and coverage -> build/coverage.out
make lint       # Run golangci-lint (install it with make deps)
make fmt        # Format code
make deps       # Install tooling
go test -v -race -run TestName ./...      # Run a single test
```

## Code Style Guidelines

### Package Structure

```
cmd/{{.Inputs.project_name.Value}}/        # Main entry point, wiring only
internal/                                 # Private packages (not importable externally)
{{- if eq .Inputs.project_type.Value "cli"}}
  cli/                                    # CLI parsing and command invocation
  commands/                               # Command implementations
  view/                                   # User interface (I/O)
{{- else}}
  config/                                 # Configuration loaded from the environment
  server/                                 # Routes, handlers, middleware, graceful shutdown
{{- end}}
  version/                                # Build version, injected with -ldflags
```

### Imports

Group imports with blank lines between:
1. Standard library
2. External/third-party packages
3. Local packages (github.com/{{.Inputs.owner.Value}}/{{.Inputs.project_name.Value}}/...)

### Naming Conventions

- **Interfaces**: Single-word nouns (`Command`, `UI`)
- **Structs**: PascalCase; private fields camelCase
- **Constructors**: `NewXxx()`, returning the interface type when appropriate
- **Acronyms**: Keep consistent casing (`StdUI`, not `StdUi`)

### Error Handling

- Return errors, don't panic (except in main)
- Wrap errors with context using `fmt.Errorf` and `%w`
- Use `errors.Is()` / `errors.As()` for comparison
- `main` is the only place that prints errors and exits

### Interfaces

- Keep interfaces small and focused
- Define interfaces where they are used
- Document all exported types, functions, and interface methods

### Comments

- Add package-level comments
- Use minimal inline comments; prefer self-documenting code
- Comment "why" not "what"

### Testing Conventions

- Place test files in the same package (`package foo` not `package foo_test`)
- Name test functions `TestFunctionName` or `TestFunctionName_Scenario`
- Use table-driven tests for multiple cases
- **BDD Syntax**: Structure tests using BDD-style comments (`// Given`, `// When`, `// Then`, `// And`)
- **Integration Tests**: Place end-to-end tests alongside `main` (`cmd/{{.Inputs.project_name.Value}}/main_test.go`), compiling the binary in `TestMain` and running it as a sub-process
{{- if eq .Inputs.project_type.Value "service"}}
- Test handlers with `net/http/httptest`; no real network needed
{{- end}}

### Context Usage

Pass `context.Context` as the first parameter of anything that blocks or does I/O, and honour its cancellation.
{{- if eq .Inputs.project_type.Value "service"}}

### Service

- Everything is configured through environment variables, read once in `internal/config` and passed down; never call `os.Getenv` elsewhere
- Log with `log/slog`; use `loggerFrom(r.Context())` inside handlers so lines carry the request id
- Run with timeouts: the `http.Server` timeouts and the graceful shutdown timeout are deliberate, keep them when adding endpoints
{{- end}}

### Code Formatting

- Run `make fmt` and `make lint` before committing
- No unused imports or variables
- Prefer explicit error returns over ignoring them
