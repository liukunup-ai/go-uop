package watcher

import (
	"context"
	"testing"

	"github.com/liukunup/go-uop/core"
)

func TestWatcherEngine_Enabled(t *testing.T) {
	engine := NewWatcherEngine()
	if engine.Enabled() {
		t.Error("new engine should be disabled by default")
	}

	engine.Enable()
	if !engine.Enabled() {
		t.Error("after Enable(), engine should be enabled")
	}

	engine.Disable()
	if engine.Enabled() {
		t.Error("after Disable(), engine should be disabled")
	}
}

func TestWatcherEngine_AddRule(t *testing.T) {
	engine := NewWatcherEngine()

	engine.AddRule(Rule{
		Name:     "test rule",
		Priority: 10,
		Match:    NewTextMatch("test"),
		Actions:  []Action{NewInlineCommand("tapOn", nil)},
	})

	if len(engine.rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(engine.rules))
	}

	engine.AddRule(Rule{
		Name:     "high priority rule",
		Priority: 1,
		Match:    NewTextMatch("test"),
		Actions:  []Action{},
	})

	if engine.rules[0].Name != "high priority rule" {
		t.Error("rules should be sorted by priority")
	}
}

func TestWatcherEngine_Check(t *testing.T) {
	ctx := context.Background()
	device := &mockDevice{}

	engine := NewWatcherEngine()
	engine.Enable()

	RegisterActionExecutor("tapOn", func(args map[string]any, d core.Device) error {
		return nil
	})

	engine.AddRule(Rule{
		Name:     "tap test",
		Priority: 10,
		Match:    NewTextMatch("nonexistent"),
		Actions:  []Action{NewInlineCommand("tapOn", map[string]any{"x": 100, "y": 200})},
	})

	err := engine.Check(ctx, device)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWatcherEngine_CheckWithMatch(t *testing.T) {
	ctx := context.Background()
	device := &mockDevice{}

	CommandExecutor = func(name string, args map[string]any, d core.Device) error {
		return nil
	}

	engine := NewWatcherEngine()
	engine.Enable()

	engine.AddRule(Rule{
		Name:     "match test",
		Priority: 10,
		Match:    NewTextMatch(""),
		Actions:  []Action{NewInlineCommand("tapOn", map[string]any{"x": 100, "y": 200})},
	})

	err := engine.Check(ctx, device)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
