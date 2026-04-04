package command

import (
	"context"
	"errors"

	"github.com/liukunup/go-uop/core"
	"github.com/liukunup/go-uop/pkg/serial"
)

// DeviceOpsHandler handles device operations
type DeviceOpsHandler struct {
	GetDevice func() (core.Device, error)
}

func (h *DeviceOpsHandler) CanHandle(cmd Command) bool {
	_, ok := cmd.(DeviceCommand)
	return ok
}

func (h *DeviceOpsHandler) Handle(ctx context.Context, cmd Command) error {
	devCmd, ok := cmd.(DeviceCommand)
	if !ok {
		return ErrUnsupportedCommand
	}

	device, err := h.GetDevice()
	if err != nil {
		return err
	}
	devCmd.SetDevice(device)

	return cmd.Execute(ctx)
}

// SerialOpsHandler handles serial operations
type SerialOpsHandler struct {
	GetSerial func() (*serial.Serial, error)
}

func (h *SerialOpsHandler) CanHandle(cmd Command) bool {
	_, ok := cmd.(SerialCommand)
	return ok
}

func (h *SerialOpsHandler) Handle(ctx context.Context, cmd Command) error {
	serCmd, ok := cmd.(SerialCommand)
	if !ok {
		return ErrUnsupportedCommand
	}

	serial, err := h.GetSerial()
	if err != nil {
		return err
	}
	serCmd.SetSerial(serial)

	return cmd.Execute(ctx)
}

type SystemOpsHandler struct{}

func (h *SystemOpsHandler) CanHandle(cmd Command) bool {
	switch cmd.(type) {
	case *WaitCommand:
		return true
	default:
		return false
	}
}

func (h *SystemOpsHandler) Handle(ctx context.Context, cmd Command) error {
	return cmd.Execute(ctx)
}

type UnknownCommandHandler struct{}

func (h *UnknownCommandHandler) CanHandle(cmd Command) bool {
	return true
}

func (h *UnknownCommandHandler) Handle(ctx context.Context, cmd Command) error {
	return errors.New("no handler available for command")
}
