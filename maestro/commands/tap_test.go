package commands

import (
	"testing"

	"github.com/liukunup/go-uop/internal/selector"
	"github.com/liukunup/go-uop/maestro"
)

func TestTapOnShorthand(t *testing.T) {
	translator := NewTapOnTranslator()

	act := translator.TranslateShorthand("Login")

	if act == nil {
		t.Fatal("expected TapAction, got nil")
	}

	if act.Element == nil {
		t.Fatal("expected Element to be set")
	}

	if act.Element.Type != selector.SelectorTypeText {
		t.Errorf("expected SelectorTypeText, got %v", act.Element.Type)
	}

	if act.Element.Value != "Login" {
		t.Errorf("expected 'Login', got %v", act.Element.Value)
	}
}

func TestTapOnExtendedByID(t *testing.T) {
	translator := NewTapOnTranslator()

	cmd := &maestro.TapOnCommand{
		ID: "submit",
	}

	act := translator.TranslateExtended(cmd)

	if act == nil {
		t.Fatal("expected TapAction, got nil")
	}

	if act.Element == nil {
		t.Fatal("expected Element to be set")
	}

	if act.Element.Type != selector.SelectorTypeID {
		t.Errorf("expected SelectorTypeID, got %v", act.Element.Type)
	}

	if act.Element.Value != "submit" {
		t.Errorf("expected 'submit', got %v", act.Element.Value)
	}
}

func TestTapOnExtendedByIDWithIndex(t *testing.T) {
	translator := NewTapOnTranslator()

	cmd := &maestro.TapOnCommand{
		ID:    "submit",
		Index: 2,
	}

	act := translator.TranslateExtended(cmd)

	if act == nil {
		t.Fatal("expected TapAction, got nil")
	}

	if act.Element == nil {
		t.Fatal("expected Element to be set")
	}

	if act.Element.Index != 2 {
		t.Errorf("expected Index 2, got %d", act.Element.Index)
	}
}

func TestTapOnExtendedByText(t *testing.T) {
	translator := NewTapOnTranslator()

	cmd := &maestro.TapOnCommand{
		Text: "Submit Button",
	}

	act := translator.TranslateExtended(cmd)

	if act == nil {
		t.Fatal("expected TapAction, got nil")
	}

	if act.Element == nil {
		t.Fatal("expected Element to be set")
	}

	if act.Element.Type != selector.SelectorTypeText {
		t.Errorf("expected SelectorTypeText, got %v", act.Element.Type)
	}

	if act.Element.Value != "Submit Button" {
		t.Errorf("expected 'Submit Button', got %v", act.Element.Value)
	}
}

func TestTapOnPoint(t *testing.T) {
	translator := NewTapOnTranslator()

	act := translator.TranslatePoint(100, 200)

	if act == nil {
		t.Fatal("expected TapAction, got nil")
	}

	if act.X != 100 {
		t.Errorf("expected X 100, got %d", act.X)
	}

	if act.Y != 200 {
		t.Errorf("expected Y 200, got %d", act.Y)
	}

	if act.Element != nil {
		t.Error("expected Element to be nil for point-based tap")
	}
}
