package device

import (
	"context"

	"github.com/liukunup/go-uop/core"
)

type Command interface {
	Name() string
	Description() string
	Validate() error
}

type DeviceCommand interface {
	Command
	Execute(ctx context.Context) error
	SetDevice(device core.Device)
	Device() core.Device
}

type baseDeviceCommand struct {
	device core.Device
}

func (c *baseDeviceCommand) SetDevice(device core.Device) {
	c.device = device
}

func (c *baseDeviceCommand) Device() core.Device {
	return c.device
}
