package commands

import (
	"testing"

	"github.com/liukunup/go-uop/internal/selector"
	"github.com/liukunup/go-uop/maestro"
)

func TestInputTextShorthand(t *testing.T) {
	translator := NewInputTextTranslator()

	act := translator.TranslateShorthand("user@example.com")

	if act == nil {
		t.Fatal("expected SendKeysAction, got nil")
	}

	if act.Text != "user@example.com" {
		t.Errorf("expected Text 'user@example.com', got '%s'", act.Text)
	}

	if act.Element != nil {
		t.Error("expected Element to be nil for shorthand input")
	}
}

func TestInputTextExtended(t *testing.T) {
	translator := NewInputTextTranslator()

	cmd := &maestro.InputTextCommand{
		Text: "hello world",
	}

	act := translator.TranslateExtended(cmd)

	if act == nil {
		t.Fatal("expected SendKeysAction, got nil")
	}

	if act.Text != "hello world" {
		t.Errorf("expected Text 'hello world', got '%s'", act.Text)
	}
}

func TestInputTextExtendedWithID(t *testing.T) {
	translator := NewInputTextTranslator()

	cmd := &maestro.InputTextCommand{
		Text: "pass",
		ID:   "password_input",
	}

	act := translator.TranslateExtended(cmd)

	if act == nil {
		t.Fatal("expected SendKeysAction, got nil")
	}

	if act.Element == nil {
		t.Fatal("expected Element to be set")
	}

	if act.Element.Type != selector.SelectorTypeID {
		t.Errorf("expected SelectorTypeID, got %v", act.Element.Type)
	}

	if act.Element.Value != "password_input" {
		t.Errorf("expected ID 'password_input', got '%s'", act.Element.Value)
	}
}

func TestInputTextExtendedWithEnter(t *testing.T) {
	translator := NewInputTextTranslator()

	cmd := &maestro.InputTextCommand{
		Text:  "password123",
		Enter: true,
	}

	act := translator.TranslateExtended(cmd)

	if act == nil {
		t.Fatal("expected SendKeysAction, got nil")
	}

	if !act.Enter {
		t.Error("expected Enter to be true")
	}

	if act.Text != "password123" {
		t.Errorf("expected Text 'password123', got '%s'", act.Text)
	}
}

func TestInputTextExtendedWithClearExistingText(t *testing.T) {
	translator := NewInputTextTranslator()

	cmd := &maestro.InputTextCommand{
		Text:              "new value",
		ClearExistingText: true,
	}

	act := translator.TranslateExtended(cmd)

	if act == nil {
		t.Fatal("expected SendKeysAction, got nil")
	}

	if !act.ClearExistingText {
		t.Error("expected ClearExistingText to be true")
	}

	if act.Text != "new value" {
		t.Errorf("expected Text 'new value', got '%s'", act.Text)
	}
}

func TestInputTextExtendedWithIDAndEnter(t *testing.T) {
	translator := NewInputTextTranslator()

	cmd := &maestro.InputTextCommand{
		Text:  "submit",
		ID:    "submit_button",
		Enter: true,
	}

	act := translator.TranslateExtended(cmd)

	if act == nil {
		t.Fatal("expected SendKeysAction, got nil")
	}

	if act.Element == nil {
		t.Fatal("expected Element to be set")
	}

	if act.Element.Type != selector.SelectorTypeID {
		t.Errorf("expected SelectorTypeID, got %v", act.Element.Type)
	}

	if !act.Enter {
		t.Error("expected Enter to be true")
	}
}
