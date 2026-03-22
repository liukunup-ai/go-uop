package commands

import (
	"testing"
	"time"

	"github.com/liukunup/go-uop/maestro"
)

func TestScrollShorthand(t *testing.T) {
	translator := NewMediaCommandTranslator()

	act := translator.TranslateScrollShorthand()

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

func TestScrollExtendedWithDownDirection(t *testing.T) {
	translator := NewMediaCommandTranslator()

	cmd := &maestro.ScrollCommand{
		Direction: maestro.ScrollDown,
	}

	act := translator.TranslateScrollExtended(cmd)

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

func TestScrollExtendedWithLeftDirection(t *testing.T) {
	translator := NewMediaCommandTranslator()

	cmd := &maestro.ScrollCommand{
		Direction: maestro.ScrollLeft,
	}

	act := translator.TranslateScrollExtended(cmd)

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

func TestScrollExtendedWithRightDirection(t *testing.T) {
	translator := NewMediaCommandTranslator()

	cmd := &maestro.ScrollCommand{
		Direction: maestro.ScrollRight,
	}

	act := translator.TranslateScrollExtended(cmd)

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

func TestScrollExtendedWithDuration(t *testing.T) {
	translator := NewMediaCommandTranslator()

	cmd := &maestro.ScrollCommand{
		Direction: maestro.ScrollUp,
		Duration:  500,
	}

	act := translator.TranslateScrollExtended(cmd)

	if act == nil {
		t.Fatal("expected SwipeAction, got nil")
	}

	if act.Duration != 500*time.Millisecond {
		t.Errorf("expected Duration 500ms, got %v", act.Duration)
	}
}

func TestScrollDefaultDirection(t *testing.T) {
	translator := NewMediaCommandTranslator()

	cmd := &maestro.ScrollCommand{}

	act := translator.TranslateScrollExtended(cmd)

	if act == nil {
		t.Fatal("expected SwipeAction, got nil")
	}

	if act.StartX != 50 || act.StartY != 80 {
		t.Errorf("unexpected start coordinates: (%d, %d)", act.StartX, act.StartY)
	}
}

func TestTakeScreenshotWithName(t *testing.T) {
	translator := NewMediaCommandTranslator()

	cmd := &maestro.TakeScreenshotCommand{
		Name: "login_screen",
	}

	act := translator.TranslateTakeScreenshot(cmd)

	if act == nil {
		t.Fatal("expected ScreenshotAction, got nil")
	}

	if act.Name != "login_screen" {
		t.Errorf("expected Name 'login_screen', got %q", act.Name)
	}
}

func TestTakeScreenshotWithoutName(t *testing.T) {
	translator := NewMediaCommandTranslator()

	act := translator.TranslateTakeScreenshot(nil)

	if act == nil {
		t.Fatal("expected ScreenshotAction, got nil")
	}

	if act.Name == "" {
		t.Error("expected non-empty auto-generated name")
	}
}

func TestTakeScreenshotShorthand(t *testing.T) {
	translator := NewMediaCommandTranslator()

	act := translator.TranslateTakeScreenshotShorthand()

	if act == nil {
		t.Fatal("expected ScreenshotAction, got nil")
	}

	if act.Name == "" {
		t.Error("expected non-empty auto-generated name")
	}
}

func TestTakeScreenshotDefaultPath(t *testing.T) {
	translator := NewMediaCommandTranslator()

	cmd := &maestro.TakeScreenshotCommand{
		Name: "test_screen",
	}

	act := translator.TranslateTakeScreenshot(cmd)

	if act.Path != DefaultScreenshotPath {
		t.Errorf("expected Path %q, got %q", DefaultScreenshotPath, act.Path)
	}
}

func TestTakeScreenshotCustomPath(t *testing.T) {
	translator := NewMediaCommandTranslatorWithPath("/custom/screenshots")

	cmd := &maestro.TakeScreenshotCommand{
		Name: "custom_screen",
	}

	act := translator.TranslateTakeScreenshot(cmd)

	if act.Path != "/custom/screenshots" {
		t.Errorf("expected Path '/custom/screenshots', got %q", act.Path)
	}
}

func TestMedia(t *testing.T) {
	t.Run("ScrollShorthand", TestScrollShorthand)
	t.Run("ScrollExtendedWithDownDirection", TestScrollExtendedWithDownDirection)
	t.Run("ScrollExtendedWithLeftDirection", TestScrollExtendedWithLeftDirection)
	t.Run("ScrollExtendedWithRightDirection", TestScrollExtendedWithRightDirection)
	t.Run("ScrollExtendedWithDuration", TestScrollExtendedWithDuration)
	t.Run("ScrollDefaultDirection", TestScrollDefaultDirection)
	t.Run("TakeScreenshotWithName", TestTakeScreenshotWithName)
	t.Run("TakeScreenshotWithoutName", TestTakeScreenshotWithoutName)
	t.Run("TakeScreenshotShorthand", TestTakeScreenshotShorthand)
	t.Run("TakeScreenshotDefaultPath", TestTakeScreenshotDefaultPath)
	t.Run("TakeScreenshotCustomPath", TestTakeScreenshotCustomPath)
}
