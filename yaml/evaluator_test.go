package yaml

import (
	"testing"
)

func TestContext_SetGetVariable(t *testing.T) {
	ctx := NewContext()

	ctx.SetVariable("name", "test")
	if ctx.GetVariable("name") != "test" {
		t.Error("expected 'test'")
	}

	ctx.SetVariable("count", 42)
	if ctx.GetVariable("count") != 42 {
		t.Error("expected 42")
	}
}

func TestContext_Evaluate_SimpleVariable(t *testing.T) {
	ctx := NewContext()
	ctx.SetVariable("username", "alice")

	result, err := ctx.Evaluate("hello ${username}")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result != "hello alice" {
		t.Errorf("expected 'hello alice', got '%s'", result)
	}
}

func TestContext_Evaluate_MissingVariable(t *testing.T) {
	ctx := NewContext()

	result, err := ctx.Evaluate("hello ${nonexistent}")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result != "hello ${nonexistent}" {
		t.Errorf("expected 'hello ${nonexistent}', got '%s'", result)
	}
}

func TestContext_Evaluate_MultipleVariables(t *testing.T) {
	ctx := NewContext()
	ctx.SetVariable("first", "hello")
	ctx.SetVariable("second", "world")

	result, err := ctx.Evaluate("${first} ${second}")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", result)
	}
}

func TestContext_Evaluate_NoVariables(t *testing.T) {
	ctx := NewContext()

	result, err := ctx.Evaluate("plain text with no variables")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result != "plain text with no variables" {
		t.Errorf("expected 'plain text with no variables', got '%s'", result)
	}
}

func TestContext_Evaluate_EmptyInput(t *testing.T) {
	ctx := NewContext()

	result, err := ctx.Evaluate("")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}
}
