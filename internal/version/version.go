// Package version holds the build version of {{.Inputs.ProjectName.Value}}.
package version

// Version can be injected at build time using
// -ldflags "-X github.com/{{.Inputs.Owner.Value}}/{{.Inputs.ProjectName.Value}}/internal/version.Version=x.y.z"
var Version = "dev"
