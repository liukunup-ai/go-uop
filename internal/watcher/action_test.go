package watcher

import (
	"context"
	"fmt"
	"testing"

	"github.com/liukunup/go-uop/core"
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

func TestActionSequence(t *testing.T) {
	ctx := context.Background()
	device := &mockDevice{}
	executed := 0

	CommandExecutor = func(name string, args map[string]any, d core.Device) error {
		executed++
		return nil
	}

	seq := ActionSequenceWithRetry([]Action{
		NewInlineCommand("count", nil),
		NewInlineCommand("count", nil),
	}, 0)

	err := seq.Execute(ctx, device)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if executed != 2 {
		t.Errorf("expected 2 executions, got %d", executed)
	}
}

func TestActionSequenceRetry(t *testing.T) {
	ctx := context.Background()
	device := &mockDevice{}
	attempts := 0

	CommandExecutor = func(name string, args map[string]any, d core.Device) error {
		attempts++
		if attempts < 2 {
			return fmt.Errorf("intentional failure on attempt %d", attempts)
		}
		return nil
	}

	seq := ActionSequenceWithRetry([]Action{
		NewInlineCommand("test", nil),
	}, 2)

	err := seq.Execute(ctx, device)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attempts != 2 {
		t.Errorf("expected 2 attempts (fail then succeed), got %d", attempts)
	}
}
