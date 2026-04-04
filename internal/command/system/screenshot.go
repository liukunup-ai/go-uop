package system

import (
	"context"
	"errors"

	"github.com/liukunup/go-uop/core"
)

type ScreenshotCommand struct {
	device core.Device
	Name_  string
	Path   string
}

func NewScreenshotCommand(name, path string) *ScreenshotCommand {
	return &ScreenshotCommand{Name_: name, Path: path}
}

func (c *ScreenshotCommand) Name() string        { return "screenshot" }
func (c *ScreenshotCommand) Description() string { return "Take a screenshot" }

func (c *ScreenshotCommand) Validate() error {
	return nil
}

func (c *ScreenshotCommand) Execute(ctx context.Context) error {
	if c.device == nil {
		return errors.New("device not set")
	}
	_, err := c.device.Screenshot()
	return err
}

func (c *ScreenshotCommand) SetDevice(device core.Device) {
	c.device = device
}
