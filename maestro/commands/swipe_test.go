package commands

import (
	"testing"
	"time"

	"github.com/liukunup/go-uop/maestro"
)

func TestSwipeDirectionUp(t *testing.T) {
	translator := NewSwipeTranslator()

	act := translator.TranslateDirection(maestro.SwipeUp)

	if act == nil {
		t.Fatal("expected SwipeAction, got nil")
	}

	if act.StartX != 50 {
		t.Errorf("expected StartX 50, got %d", act.StartX)
	}

	if act.StartY != 80 {
		t.Errorf("expected StartY 80, got %d", act.StartY)
	}

	if act.EndX != 50 {
		t.Errorf("expected EndX 50, got %d", act.EndX)
	}

	if act.EndY != 20 {
		t.Errorf("expected EndY 20, got %d", act.EndY)
	}

	if act.Duration != 300*time.Millisecond {
		t.Errorf("expected Duration 300ms, got %v", act.Duration)
	}
}

func TestSwipeDirectionDown(t *testing.T) {
	translator := NewSwipeTranslator()

	act := translator.TranslateDirection(maestro.SwipeDown)

	if act == nil {
		t.Fatal("expected SwipeAction, got nil")
	}

	if act.StartY != 20 {
		t.Errorf("expected StartY 20, got %d", act.StartY)
	}

	if act.EndY != 80 {
		t.Errorf("expected EndY 80, got %d", act.EndY)
	}
}

func TestSwipeDirectionLeft(t *testing.T) {
	translator := NewSwipeTranslator()

	act := translator.TranslateDirection(maestro.SwipeLeft)

	if act == nil {
		t.Fatal("expected SwipeAction, got nil")
	}

	if act.StartX != 70 {
		t.Errorf("expected StartX 70, got %d", act.StartX)
	}

	if act.EndX != 30 {
		t.Errorf("expected EndX 30, got %d", act.EndX)
	}
}

func TestSwipeDirectionRight(t *testing.T) {
	translator := NewSwipeTranslator()

	act := translator.TranslateDirection(maestro.SwipeRight)

	if act == nil {
		t.Fatal("expected SwipeAction, got nil")
	}

	if act.StartX != 30 {
		t.Errorf("expected StartX 30, got %d", act.StartX)
	}

	if act.EndX != 70 {
		t.Errorf("expected EndX 70, got %d", act.EndX)
	}
}

func TestSwipeDirectionWithDuration(t *testing.T) {
	translator := NewSwipeTranslator()

	act := translator.TranslateDirectionWithDuration(maestro.SwipeLeft, 500)

	if act == nil {
		t.Fatal("expected SwipeAction, got nil")
	}

	if act.Duration != 500*time.Millisecond {
		t.Errorf("expected Duration 500ms, got %v", act.Duration)
	}
}

func TestSwipeExtendedWithDirection(t *testing.T) {
	translator := NewSwipeTranslator()

	cmd := &maestro.SwipeCommand{
		Direction: maestro.SwipeUp,
	}

	act := translator.TranslateExtended(cmd)

	if act == nil {
		t.Fatal("expected SwipeAction, got nil")
	}

	if act.StartX != 50 || act.StartY != 80 {
		t.Errorf("unexpected start coordinates: (%d, %d)", act.StartX, act.StartY)
	}
}

func TestSwipeExtendedWithCoordinates(t *testing.T) {
	translator := NewSwipeTranslator()

	cmd := &maestro.SwipeCommand{
		StartX:   10,
		StartY:   50,
		EndX:     90,
		EndY:     50,
		Duration: 200,
	}

	act := translator.TranslateExtended(cmd)

	if act == nil {
		t.Fatal("expected SwipeAction, got nil")
	}

	if act.StartX != 10 {
		t.Errorf("expected StartX 10, got %d", act.StartX)
	}

	if act.StartY != 50 {
		t.Errorf("expected StartY 50, got %d", act.StartY)
	}

	if act.EndX != 90 {
		t.Errorf("expected EndX 90, got %d", act.EndX)
	}

	if act.EndY != 50 {
		t.Errorf("expected EndY 50, got %d", act.EndY)
	}

	if act.Duration != 200*time.Millisecond {
		t.Errorf("expected Duration 200ms, got %v", act.Duration)
	}
}

func TestSwipeUntilVisible(t *testing.T) {
	translator := NewSwipeTranslator()

	cmd := &maestro.ElementSelector{
		Text: "Submit",
	}

	act := translator.TranslateSwipeUntilVisible(maestro.SwipeUp, cmd)

	if act == nil {
		t.Fatal("expected SwipeAction, got nil")
	}

	if act.StartX != 50 || act.StartY != 80 {
		t.Errorf("unexpected start coordinates: (%d, %d)", act.StartX, act.StartY)
	}

	if act.EndX != 50 || act.EndY != 20 {
		t.Errorf("unexpected end coordinates: (%d, %d)", act.EndX, act.EndY)
	}
}

func TestSwipeDefaultDuration(t *testing.T) {
	translator := NewSwipeTranslator()

	cmd := &maestro.SwipeCommand{
		StartX: 0,
		StartY: 0,
		EndX:   100,
		EndY:   100,
	}

	act := translator.TranslateExtended(cmd)

	if act == nil {
		t.Fatal("expected SwipeAction, got nil")
	}

	if act.Duration != 300*time.Millisecond {
		t.Errorf("expected default Duration 300ms, got %v", act.Duration)
	}
}
