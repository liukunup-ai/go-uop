package vision

import (
	"testing"
)

func TestNewMatcher_UnknownAlgorithm(t *testing.T) {
	_, err := NewMatcher("unknown")
	if err == nil {
		t.Error("NewMatcher(unknown) should return error")
	}
	if err.Error() != "unknown algorithm: unknown" {
		t.Errorf("error = %q, want %q", err.Error(), "unknown algorithm: unknown")
	}
}

func TestNewMatcher_Template(t *testing.T) {
	m, err := NewMatcher("template")
	if err != nil {
		t.Fatalf("NewMatcher(template) failed: %v", err)
	}
	if m.Name() != "template" {
		t.Errorf("Name() = %q, want %q", m.Name(), "template")
	}
}
