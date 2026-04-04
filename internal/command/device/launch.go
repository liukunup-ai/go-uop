package device

import (
	"context"
	"errors"
)

type LaunchCommand struct {
	baseDeviceCommand
	AppID    string
	Args     []string
	WaitIdle bool
}

func NewLaunchCommand(appID string) *LaunchCommand {
	return &LaunchCommand{AppID: appID}
}

func (c *LaunchCommand) Name() string        { return "launch" }
func (c *LaunchCommand) Description() string { return "Launch an application" }

func (c *LaunchCommand) Validate() error {
	return nil
}

func (c *LaunchCommand) Execute(ctx context.Context) error {
	if c.Device() == nil {
		return errors.New("device not set")
	}
	return c.Device().Launch()
}
