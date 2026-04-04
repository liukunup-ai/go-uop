package command

import (
	"context"
	"testing"
	"time"
)

// WaitCommand tests

func TestWaitCommand_Validate(t *testing.T) {
	tests := []struct {
		name     string
		duration int
		wantErr  bool
	}{
		{"valid positive duration", 100, false},
		{"zero duration", 0, false},
		{"negative duration", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewWaitCommand(tt.duration)
			err := cmd.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWaitCommand_Name(t *testing.T) {
	cmd := NewWaitCommand(100)
	if cmd.Name() != "wait" {
		t.Errorf("Name() = %s, want wait", cmd.Name())
	}
}

func TestWaitCommand_Execute(t *testing.T) {
	cmd := NewWaitCommand(10) // 10ms
	start := time.Now()
	err := cmd.Execute(context.Background())
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}
	// Should have slept for at least 10ms
	if elapsed < 10*time.Millisecond {
		t.Errorf("Execute() did not sleep long enough: %v", elapsed)
	}
}

// ScreenshotCommand tests

func TestScreenshotCommand_Name(t *testing.T) {
	cmd := NewScreenshotCommand("test", "/tmp/screenshot.png")
	if cmd.Name() != "screenshot" {
		t.Errorf("Name() = %s, want screenshot", cmd.Name())
	}
}

func TestScreenshotCommand_Validate(t *testing.T) {
	cmd := NewScreenshotCommand("test", "/tmp/screenshot.png")
	if err := cmd.Validate(); err != nil {
		t.Errorf("Validate() error = %v", err)
	}
}

func TestScreenshotCommand_Execute_NoDevice(t *testing.T) {
	cmd := NewScreenshotCommand("test", "/tmp/screenshot.png")
	err := cmd.Execute(context.Background())
	if err == nil {
		t.Error("Execute() should return error when device not set")
	}
}

func TestScreenshotCommand_SetDevice(t *testing.T) {
	cmd := NewScreenshotCommand("test", "/tmp/screenshot.png")
	cmd.SetDevice(nil)
	if cmd.BaseDeviceCommand.device != nil {
		t.Error("device should be nil after SetDevice(nil)")
	}
}

// SwipeCommand tests

func TestSwipeCommand_Name(t *testing.T) {
	cmd := NewSwipeCommand(0, 0, 100, 100)
	if cmd.Name() != "swipe" {
		t.Errorf("Name() = %s, want swipe", cmd.Name())
	}
}

func TestSwipeCommand_Execute(t *testing.T) {
	cmd := NewSwipeCommand(0, 0, 100, 100)
	err := cmd.Execute(context.Background())
	// Currently returns nil (not implemented)
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}
}

func TestSwipeCommand_SetDevice(t *testing.T) {
	cmd := NewSwipeCommand(0, 0, 100, 100)
	cmd.SetDevice(nil)
	if cmd.BaseDeviceCommand.device != nil {
		t.Error("device should be nil after SetDevice(nil)")
	}
}
