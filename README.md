# {{.Inputs.project_name.Value}}

{{if eq .Inputs.project_type.Value "cli" -}}
A Go command line application.
{{- else -}}
A Go web service.
{{- end}}

## Development

Requires Go {{.Inputs.go_version.Value}} or later.

```bash
make deps      # install tooling (golangci-lint)
make all       # lint, test and build
make test      # tests with the race detector and coverage -> build/coverage.html
make lint      # golangci-lint
make fmt       # format with golangci-lint
make build     # build -> build/{{.Inputs.project_name.Value}}
{{- if eq .Inputs.project_type.Value "cli"}}
make release   # build for linux and darwin on amd64 and arm64 -> build/bin
{{- else}}
make run       # run the service locally
make docker    # build the container image
{{- end}}
```
{{if eq .Inputs.project_type.Value "cli"}}
## Usage

```bash
{{.Inputs.project_name.Value}} help
{{.Inputs.project_name.Value}} greet -name Gopher
{{.Inputs.project_name.Value}} version
```

To add a command, implement `commands.Command` in `internal/commands` and register it in `internal/cli/cli_flag.go`.
{{- else}}
## Running

```bash
make run
curl localhost:8080/healthz
curl localhost:8080/hello/Gopher
```

The service is configured through environment variables:

| Variable           | Default | Description                                         |
|--------------------|---------|-----------------------------------------------------|
| `ADDR`             | `:8080` | Address to listen on                                |
| `LOG_LEVEL`        | `info`  | `debug`, `info`, `warn` or `error`                  |
| `SHUTDOWN_TIMEOUT` | `10s`   | How long in-flight requests get to finish on SIGTERM |

Endpoints live in `internal/server/handlers.go`.
{{- end}}

## Releases

Commits follow [Conventional Commits](https://www.conventionalcommits.org/). [release-please](https://github.com/googleapis/release-please) keeps a release pull request up to date on `main`; merging it tags a release, which {{if eq .Inputs.project_type.Value "cli" -}}
builds the binaries and attaches them to the GitHub release.
{{- else -}}
builds the container image and pushes it to `ghcr.io/{{.Inputs.owner.Value}}/{{.Inputs.project_name.Value}}`.
{{- end}}

The release-please workflow needs a `RELEASE_PLEASE_PAT` repository secret (a personal access token) so that the tag it creates triggers the release workflow.
