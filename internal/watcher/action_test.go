package watcher

import (
	"context"
	"testing"
)

func TestInlineCommand(t *testing.T) {
	ctx := context.Background()
	device := &mockDevice{}

	action := NewInlineCommand("tapOn", map[string]any{"x": 100, "y": 200})
	err := action.Execute(ctx, device)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReferenceFlow(t *testing.T) {
	ctx := context.Background()
	device := &mockDevice{}

	action := NewReferenceFlow("dismiss-popup-flow")
	err := action.Execute(ctx, device)
	if err != nil {
		t.Logf("expected error (flow not exist): %v", err)
	}
}
