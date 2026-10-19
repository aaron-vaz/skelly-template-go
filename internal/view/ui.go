// Package view is the application's user interface (I/O).
package view

// UI is how commands talk to the user.
type UI interface {
	// RenderInfo writes an informational message to the user.
	RenderInfo(message string) error
	// RenderError writes an error message to the user.
	RenderError(message string) error
}
