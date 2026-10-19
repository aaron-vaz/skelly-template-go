package view

import (
	"bytes"
	"testing"
)

func TestStdUI_RenderInfo(t *testing.T) {
	// Given
	var stdout, stderr bytes.Buffer
	ui := NewStdUI(&stdout, &stderr)

	// When
	err := ui.RenderInfo("hello")

	// Then
	if err != nil {
		t.Fatalf("RenderInfo() error = %v", err)
	}

	// And
	if got := stdout.String(); got != "hello\n" {
		t.Errorf("stdout = %q, want %q", got, "hello\n")
	}
	if stderr.Len() != 0 {
		t.Errorf("stderr = %q, want empty", stderr.String())
	}
}

func TestStdUI_RenderError(t *testing.T) {
	// Given
	var stdout, stderr bytes.Buffer
	ui := NewStdUI(&stdout, &stderr)

	// When
	err := ui.RenderError("boom")

	// Then
	if err != nil {
		t.Fatalf("RenderError() error = %v", err)
	}

	// And
	if got := stderr.String(); got != "boom\n" {
		t.Errorf("stderr = %q, want %q", got, "boom\n")
	}
	if stdout.Len() != 0 {
		t.Errorf("stdout = %q, want empty", stdout.String())
	}
}
