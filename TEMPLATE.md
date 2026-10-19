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

`.skelly.yml` uses three skelly features to produce both variants from one tree:

- `choices` on `ProjectType`, so a typo is rejected instead of producing a half-empty project.
- `exclude` entries, each removing a path when `ProjectType` is not the one that needs it. `main.go`, the tests and the `Makefile` are small enough to switch with `{{if}}` inside the file.
- `rename`, which moves `github/` to `.github` when the project is generated. The workflows and `dependabot.yml` for generated projects live in `github/`, so GitHub does not run them, and Dependabot does not read them, in this repository. The real `.github/` here holds only the template's own CI.

`github` is listed under `raw` because workflows use `${{ }}`. `raw` and `exclude` use the template's own names, which is why they say `github/...`. `TEMPLATE.md` and the real `.github` are excluded unconditionally: they belong to this repository, not to generated projects.

When you edit a workflow for generated projects, edit it under `github/workflows/`. It is not run here, so check changes with the test below.

## Using it from a script or an AI agent

```bash
# What does the template ask for? JSON with defaults, choices and which inputs are required
skelly inspect go

# Create a project without any prompt
skelly init go -dst ./my-app -no-input -set ProjectName=my-app -set ProjectType=service
```

## Testing the template

`.github/workflows/template-test.yml` renders each variant with skelly and runs `make all` on the result. To do the same locally:

```bash
skelly init . -dst /tmp/example-app -no-input -set ProjectName=example-app -set ProjectType=cli
cd /tmp/example-app && make deps all
```
