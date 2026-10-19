// Package version holds the build version of {{.Inputs.project_name.Value}}.
package version

// Version can be injected at build time using
// -ldflags "-X github.com/{{.Inputs.owner.Value}}/{{.Inputs.project_name.Value}}/internal/version.Version=x.y.z"
var Version = "dev"
