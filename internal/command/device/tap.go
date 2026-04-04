package device

import (
	"context"
	"errors"
)

type TapCommand struct {
	baseDeviceCommand
	X int
	Y int
}

func NewTapCommand(x, y int) *TapCommand {
	return &TapCommand{X: x, Y: y}
}

func (c *TapCommand) Name() string        { return "tapOn" }
func (c *TapCommand) Description() string { return "Tap at coordinates" }

func (c *TapCommand) Validate() error {
	if c.X < 0 || c.Y < 0 {
		return errors.New("coordinates must be non-negative")
	}
	return nil
}

func (c *TapCommand) Execute(ctx context.Context) error {
	if c.Device() == nil {
		return errors.New("device not set")
	}
	return c.Device().Tap(c.X, c.Y)
}
