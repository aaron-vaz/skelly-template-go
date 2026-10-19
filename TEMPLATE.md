# skelly-template-go

A [skelly](https://github.com/aaron-vaz/skelly) template for Go projects. One template, two kinds of project, picked with the `ProjectType` input:

| `ProjectType` | You get |
|---------------|---------|
| `cli`         | A command line app: command interface, flag parsing that allows flags after arguments, `help`, `version` and an example `greet` command, plus a release workflow that attaches binaries to the GitHub release |
| `service`     | A web service on the standard library: `net/http` routing, request ids and access logs (`log/slog`), panic recovery, health/readiness/version endpoints, env-based config, graceful shutdown, a Dockerfile and a workflow that publishes the image to GHCR |

Both come with a `Makefile`, golangci-lint v2 config, CI, release-please, Dependabot, an `AGENTS.md` describing the conventions, and tests (including an end-to-end test of the compiled binary).

## Usage

```bash
skelly sources add go github.com/aaron-vaz/skelly-template-go
skelly init go -dst ./my-app
```

Inputs are asked in alphabetical order: `GoVersion` (default `1.26`), `Owner` (default `aaron-vaz`), `ProjectName` (required) and `ProjectType` (`cli` or `service`, default `cli`).

## How it works

`.skelly.yml` uses two skelly features to produce both variants from one tree:

- `choices` on `ProjectType`, so a typo is rejected instead of producing a half-empty project.
- `exclude` entries, each removing a path when `ProjectType` is not the one that needs it. `main.go`, the tests and the `Makefile` are small enough to switch with `{{if}}` inside the file.

`.github` is listed under `raw` because workflows use `${{ }}`. `TEMPLATE.md` and `.github/workflows/template-test.yml` are excluded unconditionally: they belong to this repository, not to generated projects.

## Testing the template

`.github/workflows/template-test.yml` renders each variant with skelly and runs `make all` on the result. To do the same locally:

```bash
printf '\n\nexample-app\ncli\n' | skelly init . -dst /tmp/example-app
cd /tmp/example-app && make deps all
```
