package commands

import (
	"testing"
	"time"

	"github.com/liukunup/go-uop/maestro"
)

func TestWaitForAnimationToEndWithTimeout(t *testing.T) {
	translator := NewNavigationTranslator()

	cmd := &maestro.WaitForAnimationEnd{
		Timeout: 15000,
	}

	act := translator.TranslateWaitForAnimationToEnd(cmd)

	if act == nil {
		t.Fatal("expected WaitAction, got nil")
	}

	if act.Duration != 15000*time.Millisecond {
		t.Errorf("expected Duration 15000ms, got %v", act.Duration)
	}
}

func TestWaitForAnimationToEndWithCustomTimeout(t *testing.T) {
	translator := NewNavigationTranslator()

	cmd := &maestro.WaitForAnimationEnd{
		Timeout: 30000,
	}

	act := translator.TranslateWaitForAnimationToEnd(cmd)

	if act == nil {
		t.Fatal("expected WaitAction, got nil")
	}

	if act.Duration != 30000*time.Millisecond {
		t.Errorf("expected Duration 30000ms, got %v", act.Duration)
	}
}

func TestWaitForAnimationToEndWithoutTimeout(t *testing.T) {
	translator := NewNavigationTranslator()

	cmd := &maestro.WaitForAnimationEnd{}

	act := translator.TranslateWaitForAnimationToEnd(cmd)

	if act == nil {
		t.Fatal("expected WaitAction, got nil")
	}

	if act.Duration != DefaultTimeout {
		t.Errorf("expected Duration DefaultTimeout (%v), got %v", DefaultTimeout, act.Duration)
	}
}

func TestWaitForAnimationToEndNilCommand(t *testing.T) {
	translator := NewNavigationTranslator()

	act := translator.TranslateWaitForAnimationToEnd(nil)

	if act == nil {
		t.Fatal("expected WaitAction, got nil")
	}

	if act.Duration != DefaultTimeout {
		t.Errorf("expected Duration DefaultTimeout (%v), got %v", DefaultTimeout, act.Duration)
	}
}

func TestPressKeyHome(t *testing.T) {
	translator := NewNavigationTranslator()

	cmd := &maestro.PressKeyCommand{
		Key: "home",
	}

	act := translator.TranslatePressKey(cmd)

	if act == nil {
		t.Fatal("expected PressKeyAction, got nil")
	}

	if act.KeyCode != 3 {
		t.Errorf("expected KeyCode 3 (KEYCODE_HOME), got %d", act.KeyCode)
	}
}

func TestPressKeyBack(t *testing.T) {
	translator := NewNavigationTranslator()

	cmd := &maestro.PressKeyCommand{
		Key: "back",
	}

	act := translator.TranslatePressKey(cmd)

	if act == nil {
		t.Fatal("expected PressKeyAction, got nil")
	}

	if act.KeyCode != 4 {
		t.Errorf("expected KeyCode 4 (KEYCODE_BACK), got %d", act.KeyCode)
	}
}

func TestPressKeyEnter(t *testing.T) {
	translator := NewNavigationTranslator()

	cmd := &maestro.PressKeyCommand{
		Key: "enter",
	}

	act := translator.TranslatePressKey(cmd)

	if act == nil {
		t.Fatal("expected PressKeyAction, got nil")
	}

	if act.KeyCode != 66 {
		t.Errorf("expected KeyCode 66 (KEYCODE_ENTER), got %d", act.KeyCode)
	}
}

func TestPressKeyRecentApps(t *testing.T) {
	translator := NewNavigationTranslator()

	cmd := &maestro.PressKeyCommand{
		Key: "recent_apps",
	}

	act := translator.TranslatePressKey(cmd)

	if act == nil {
		t.Fatal("expected PressKeyAction, got nil")
	}

	if act.KeyCode != 187 {
		t.Errorf("expected KeyCode 187 (KEYCODE_APP_SWITCH), got %d", act.KeyCode)
	}
}

func TestPressKeyNilCommand(t *testing.T) {
	translator := NewNavigationTranslator()

	act := translator.TranslatePressKey(nil)

	if act != nil {
		t.Errorf("expected nil for nil command, got %v", act)
	}
}

func TestPressKeyUnknownKey(t *testing.T) {
	translator := NewNavigationTranslator()

	cmd := &maestro.PressKeyCommand{
		Key: "unknown_key",
	}

	act := translator.TranslatePressKey(cmd)

	if act == nil {
		t.Fatal("expected PressKeyAction, got nil")
	}

	if act.KeyCode != 0 {
		t.Errorf("expected KeyCode 0 for unknown key, got %d", act.KeyCode)
	}
}

func TestTranslateBack(t *testing.T) {
	translator := NewNavigationTranslator()

	act := translator.TranslateBack()

	if act == nil {
		t.Fatal("expected PressKeyAction, got nil")
	}

	if act.KeyCode != 4 {
		t.Errorf("expected KeyCode 4 (KEYCODE_BACK), got %d", act.KeyCode)
	}
}

func TestTranslatePressHome(t *testing.T) {
	translator := NewNavigationTranslator()

	act := translator.TranslatePressHome()

	if act == nil {
		t.Fatal("expected PressKeyAction, got nil")
	}

	if act.KeyCode != 3 {
		t.Errorf("expected KeyCode 3 (KEYCODE_HOME), got %d", act.KeyCode)
	}
}

func TestTranslatePressRecentApps(t *testing.T) {
	translator := NewNavigationTranslator()

	act := translator.TranslatePressRecentApps()

	if act == nil {
		t.Fatal("expected PressKeyAction, got nil")
	}

	if act.KeyCode != 187 {
		t.Errorf("expected KeyCode 187 (KEYCODE_APP_SWITCH), got %d", act.KeyCode)
	}
}

func TestMapKeyToCode(t *testing.T) {
	tests := []struct {
		key      string
		expected int
	}{
		{"home", 3},
		{"back", 4},
		{"enter", 66},
		{"recent_apps", 187},
		{"recent", 187},
		{"volume_up", 24},
		{"volume_down", 25},
		{"power", 26},
		{"delete", 67},
		{"tab", 61},
		{"escape", 111},
		{"esc", 111},
		{"space", 62},
		{"menu", 82},
		{"unknown", 0},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := MapKeyToCode(tt.key)
			if result != tt.expected {
				t.Errorf("MapKeyToCode(%q) = %d, want %d", tt.key, result, tt.expected)
			}
		})
	}
}
