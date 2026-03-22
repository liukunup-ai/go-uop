package commands

import (
	"testing"

	"github.com/liukunup/go-uop/maestro"
)

func TestLaunchAppShorthand(t *testing.T) {
	translator := NewAppCommandTranslator()

	act := translator.TranslateLaunchShorthand("com.example")

	if act == nil {
		t.Fatal("expected LaunchAction, got nil")
	}

	if act.AppID != "com.example" {
		t.Errorf("expected AppID 'com.example', got %v", act.AppID)
	}
}

func TestLaunchAppExtended(t *testing.T) {
	translator := NewAppCommandTranslator()

	cmd := &maestro.LaunchAppCommand{
		AppID: "com.example",
	}

	act := translator.TranslateLaunchExtended(cmd)

	if act == nil {
		t.Fatal("expected LaunchAction, got nil")
	}

	if act.AppID != "com.example" {
		t.Errorf("expected AppID 'com.example', got %v", act.AppID)
	}

	if act.ClearState {
		t.Error("expected ClearState false by default")
	}
}

func TestLaunchAppExtendedWithClearState(t *testing.T) {
	translator := NewAppCommandTranslator()

	cmd := &maestro.LaunchAppCommand{
		AppID:      "com.example",
		ClearState: true,
	}

	act := translator.TranslateLaunchExtended(cmd)

	if act == nil {
		t.Fatal("expected LaunchAction, got nil")
	}

	if act.AppID != "com.example" {
		t.Errorf("expected AppID 'com.example', got %v", act.AppID)
	}

	if !act.ClearState {
		t.Error("expected ClearState true")
	}
}

func TestKillApp(t *testing.T) {
	translator := NewAppCommandTranslator()

	act := translator.TranslateKill("com.example")

	if act == nil {
		t.Fatal("expected KillAction, got nil")
	}

	if act.AppID != "com.example" {
		t.Errorf("expected AppID 'com.example', got %v", act.AppID)
	}
}

func TestStopApp(t *testing.T) {
	translator := NewAppCommandTranslator()

	act := translator.TranslateStop("com.example")

	if act == nil {
		t.Fatal("expected StopAction, got nil")
	}

	if act.AppID != "com.example" {
		t.Errorf("expected AppID 'com.example', got %v", act.AppID)
	}

	if !act.Graceful {
		t.Error("expected Graceful true for StopAction")
	}
}

func TestClearState(t *testing.T) {
	translator := NewAppCommandTranslator()

	act := translator.TranslateClearState("com.example")

	if act == nil {
		t.Fatal("expected ClearStateAction, got nil")
	}

	if act.AppID != "com.example" {
		t.Errorf("expected AppID 'com.example', got %v", act.AppID)
	}
}
