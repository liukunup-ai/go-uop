package command

import (
	"context"
	"errors"
	"testing"

	"github.com/liukunup/go-uop/core"
	"github.com/liukunup/go-uop/pkg/serial"
)

type mockDevice struct{}

func (m *mockDevice) Platform() core.Platform               { return core.IOS }
func (m *mockDevice) Info() (map[string]interface{}, error) { return nil, nil }
func (m *mockDevice) Tap(x, y int) error                    { return nil }
func (m *mockDevice) Screenshot() ([]byte, error)           { return nil, nil }
func (m *mockDevice) SendKeys(text string) error            { return nil }
func (m *mockDevice) Launch() error                         { return nil }
func (m *mockDevice) Close() error                          { return nil }

func TestDeviceOpsHandler_CanHandle(t *testing.T) {
	h := &DeviceOpsHandler{}
	tapCmd := NewTapCommand(100, 200)

	if !h.CanHandle(tapCmd) {
		t.Error("CanHandle should return true for DeviceCommand")
	}
}

func TestDeviceOpsHandler_CanHandle_NotDeviceCommand(t *testing.T) {
	h := &DeviceOpsHandler{}
	cmd := &BaseCommand{}

	if h.CanHandle(cmd) {
		t.Error("CanHandle should return false for non-DeviceCommand")
	}
}

func TestDeviceOpsHandler_Handle(t *testing.T) {
	h := &DeviceOpsHandler{
		GetDevice: func() (core.Device, error) {
			return &mockDevice{}, nil
		},
	}
	tapCmd := NewTapCommand(100, 200)

	err := h.Handle(context.Background(), tapCmd)
	if err != nil {
		t.Errorf("Handle() error = %v", err)
	}
}

func TestDeviceOpsHandler_Handle_GetDeviceError(t *testing.T) {
	h := &DeviceOpsHandler{
		GetDevice: func() (core.Device, error) {
			return nil, errors.New("device error")
		},
	}
	tapCmd := NewTapCommand(100, 200)

	err := h.Handle(context.Background(), tapCmd)
	if err == nil {
		t.Error("Handle() should return error when GetDevice fails")
	}
}

func TestDeviceOpsHandler_Handle_NotDeviceCommand(t *testing.T) {
	h := &DeviceOpsHandler{}
	cmd := &BaseCommand{}

	err := h.Handle(context.Background(), cmd)
	if err != ErrUnsupportedCommand {
		t.Errorf("Handle() = %v, want ErrUnsupportedCommand", err)
	}
}

func TestSerialOpsHandler_CanHandle(t *testing.T) {
	h := &SerialOpsHandler{}
	cmd := NewSendByIDCommand("reset", 0)

	if !h.CanHandle(cmd) {
		t.Error("CanHandle should return true for SerialCommand")
	}
}

func TestSerialOpsHandler_CanHandle_NotSerialCommand(t *testing.T) {
	h := &SerialOpsHandler{}
	cmd := &BaseCommand{}

	if h.CanHandle(cmd) {
		t.Error("CanHandle should return false for non-SerialCommand")
	}
}

func TestSerialOpsHandler_Handle(t *testing.T) {
	h := &SerialOpsHandler{
		GetSerial: func() (*serial.Serial, error) {
			return nil, nil
		},
	}
	cmd := NewSendByIDCommand("reset", 0)

	err := h.Handle(context.Background(), cmd)
	if err == nil {
		t.Error("Handle() should return error when serial returns nil")
	}
}

func TestSerialOpsHandler_Handle_NotSerialCommand(t *testing.T) {
	h := &SerialOpsHandler{}
	cmd := &BaseCommand{}

	err := h.Handle(context.Background(), cmd)
	if err != ErrUnsupportedCommand {
		t.Errorf("Handle() = %v, want ErrUnsupportedCommand", err)
	}
}

func TestSystemOpsHandler_CanHandle_WaitCommand(t *testing.T) {
	h := &SystemOpsHandler{}
	cmd := NewWaitCommand(100)

	if !h.CanHandle(cmd) {
		t.Error("CanHandle should return true for WaitCommand")
	}
}

func TestSystemOpsHandler_CanHandle_NotWaitCommand(t *testing.T) {
	h := &SystemOpsHandler{}
	cmd := &BaseCommand{}

	if h.CanHandle(cmd) {
		t.Error("CanHandle should return false for non-WaitCommand")
	}
}

func TestSystemOpsHandler_Handle(t *testing.T) {
	h := &SystemOpsHandler{}
	cmd := NewWaitCommand(10)

	err := h.Handle(context.Background(), cmd)
	if err != nil {
		t.Errorf("Handle() error = %v", err)
	}
}

func TestUnknownCommandHandler_CanHandle(t *testing.T) {
	h := &UnknownCommandHandler{}
	cmd := &BaseCommand{}

	if !h.CanHandle(cmd) {
		t.Error("CanHandle should return true for any command")
	}
}

func TestUnknownCommandHandler_Handle(t *testing.T) {
	h := &UnknownCommandHandler{}
	cmd := &BaseCommand{}

	err := h.Handle(context.Background(), cmd)
	if err == nil {
		t.Error("Handle() should return error")
	}
}
