package commands

import (
	"testing"
	"time"

	"github.com/liukunup/go-uop/internal/selector"
	"github.com/liukunup/go-uop/maestro"
)

func TestAssertVisibleShorthand(t *testing.T) {
	translator := NewAssertTranslator()

	act := translator.TranslateAssertVisibleShorthand("Login")

	if act == nil {
		t.Fatal("expected AssertAction, got nil")
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

	if !act.MustExist {
		t.Error("expected MustExist to be true")
	}
}

func TestAssertNotVisibleShorthand(t *testing.T) {
	translator := NewAssertTranslator()

	act := translator.TranslateAssertNotVisibleShorthand("Loading")

	if act == nil {
		t.Fatal("expected AssertAction, got nil")
	}

	if act.Element == nil {
		t.Fatal("expected Element to be set")
	}

	if act.Element.Type != selector.SelectorTypeText {
		t.Errorf("expected SelectorTypeText, got %v", act.Element.Type)
	}

	if act.Element.Value != "Loading" {
		t.Errorf("expected 'Loading', got %v", act.Element.Value)
	}

	if act.MustExist {
		t.Error("expected MustExist to be false")
	}
}

func TestAssertVisibleByID(t *testing.T) {
	translator := NewAssertTranslator()

	cmd := &maestro.ElementSelector{
		ID: "submit",
	}

	act := translator.TranslateAssertVisible(cmd)

	if act == nil {
		t.Fatal("expected AssertAction, got nil")
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

	if !act.MustExist {
		t.Error("expected MustExist to be true")
	}
}

func TestAssertNotVisibleByID(t *testing.T) {
	translator := NewAssertTranslator()

	cmd := &maestro.ElementSelector{
		ID: "loading",
	}

	act := translator.TranslateAssertNotVisible(cmd)

	if act == nil {
		t.Fatal("expected AssertAction, got nil")
	}

	if act.Element == nil {
		t.Fatal("expected Element to be set")
	}

	if act.Element.Type != selector.SelectorTypeID {
		t.Errorf("expected SelectorTypeID, got %v", act.Element.Type)
	}

	if act.Element.Value != "loading" {
		t.Errorf("expected 'loading', got %v", act.Element.Value)
	}

	if act.MustExist {
		t.Error("expected MustExist to be false")
	}
}

func TestAssertVisibleByText(t *testing.T) {
	translator := NewAssertTranslator()

	cmd := &maestro.ElementSelector{
		Text: "Submit Button",
	}

	act := translator.TranslateAssertVisible(cmd)

	if act == nil {
		t.Fatal("expected AssertAction, got nil")
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

func TestAssertVisibleWithIndex(t *testing.T) {
	translator := NewAssertTranslator()

	cmd := &maestro.ElementSelector{
		ID:    "submit",
		Index: 2,
	}

	act := translator.TranslateAssertVisible(cmd)

	if act == nil {
		t.Fatal("expected AssertAction, got nil")
	}

	if act.Element == nil {
		t.Fatal("expected Element to be set")
	}

	if act.Element.Index != 2 {
		t.Errorf("expected Index 2, got %d", act.Element.Index)
	}
}

func TestAssertVisibleWithTimeout(t *testing.T) {
	translator := NewAssertTranslator()

	cmd := &maestro.TapOnCommand{
		ID:      "submit",
		Timeout: "5000",
	}

	act := translator.TranslateAssertVisibleWithTimeout(cmd)

	if act == nil {
		t.Fatal("expected AssertAction, got nil")
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

	if !act.MustExist {
		t.Error("expected MustExist to be true")
	}

	if act.Timeout != 5000*time.Millisecond {
		t.Errorf("expected Timeout 5s, got %v", act.Timeout)
	}
}

func TestAssertNotVisibleWithTimeout(t *testing.T) {
	translator := NewAssertTranslator()

	cmd := &maestro.TapOnCommand{
		Text:    "Loading",
		Timeout: "3000",
	}

	act := translator.TranslateAssertNotVisibleWithTimeout(cmd)

	if act == nil {
		t.Fatal("expected AssertAction, got nil")
	}

	if act.Element == nil {
		t.Fatal("expected Element to be set")
	}

	if act.Element.Type != selector.SelectorTypeText {
		t.Errorf("expected SelectorTypeText, got %v", act.Element.Type)
	}

	if act.Element.Value != "Loading" {
		t.Errorf("expected 'Loading', got %v", act.Element.Value)
	}

	if act.MustExist {
		t.Error("expected MustExist to be false")
	}

	if act.Timeout != 3000*time.Millisecond {
		t.Errorf("expected Timeout 3s, got %v", act.Timeout)
	}
}
