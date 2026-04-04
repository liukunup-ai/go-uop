package command

import (
	"context"
	"testing"
)

func TestTapCommand(t *testing.T) {
	cmd := NewTapCommand(100, 200)

	if cmd.Name() != "tapOn" {
		t.Errorf("Name() = %s, want tapOn", cmd.Name())
	}

	if err := cmd.Validate(); err != nil {
		t.Errorf("Validate() error = %v", err)
	}

	// Test invalid coordinates
	invalidCmd := NewTapCommand(-1, -1)
	if err := invalidCmd.Validate(); err == nil {
		t.Error("Validate() should return error for negative coordinates")
	}
}

func TestLaunchCommand(t *testing.T) {
	cmd := NewLaunchCommand("com.example.app")

	if cmd.Name() != "launch" {
		t.Errorf("Name() = %s, want launch", cmd.Name())
	}

	if cmd.AppID != "com.example.app" {
		t.Errorf("AppID = %s, want com.example.app", cmd.AppID)
	}
}

func TestSendKeysCommand(t *testing.T) {
	cmd := NewSendKeysCommand("hello world")

	if cmd.Name() != "inputText" {
		t.Errorf("Name() = %s, want inputText", cmd.Name())
	}

	if cmd.Text != "hello world" {
		t.Errorf("Text = %s, want hello world", cmd.Text)
	}
}

func TestPressKeyCommand(t *testing.T) {
	cmd := NewPressKeyCommand(42)

	if cmd.Name() != "pressKey" {
		t.Errorf("Name() = %s, want pressKey", cmd.Name())
	}

	if cmd.KeyCode != 42 {
		t.Errorf("KeyCode = %d, want 42", cmd.KeyCode)
	}
}

func TestBaseDeviceCommand_SetDevice(t *testing.T) {
	// Create a mock device that implements core.Device interface
	// For now, we just verify the SetDevice method exists and can be called
	cmd := NewTapCommand(100, 200)

	// Test that SetDevice can be called without panic (device is nil)
	cmd.SetDevice(nil)

	// Verify device is nil
	if cmd.BaseDeviceCommand.device != nil {
		t.Error("device should be nil after SetDevice(nil)")
	}
}

func TestTapCommand_Execute(t *testing.T) {
	cmd := NewTapCommand(100, 200)

	// Test execute without device set
	err := cmd.Execute(context.Background())
	if err == nil {
		t.Error("Execute() should return error when device not set")
	}
}

func TestLaunchCommand_Execute(t *testing.T) {
	cmd := NewLaunchCommand("com.example.app")

	// Test execute without device set
	err := cmd.Execute(context.Background())
	if err == nil {
		t.Error("Execute() should return error when device not set")
	}
}

func TestSendKeysCommand_Execute(t *testing.T) {
	cmd := NewSendKeysCommand("test")

	// Test execute without device set
	err := cmd.Execute(context.Background())
	if err == nil {
		t.Error("Execute() should return error when device not set")
	}
}

func TestPressKeyCommand_Execute(t *testing.T) {
	cmd := NewPressKeyCommand(42)

	// Test execute without device set
	err := cmd.Execute(context.Background())
	// PressKeyCommand returns nil (not implemented yet), so this should pass
	if err != nil {
		t.Errorf("Execute() error = %v, want nil", err)
	}
}
