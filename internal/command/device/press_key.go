package device

import "context"

type PressKeyCommand struct {
	baseDeviceCommand
	KeyCode int
}

func NewPressKeyCommand(keyCode int) *PressKeyCommand {
	return &PressKeyCommand{KeyCode: keyCode}
}

func (c *PressKeyCommand) Name() string        { return "pressKey" }
func (c *PressKeyCommand) Description() string { return "Press a key" }

func (c *PressKeyCommand) Validate() error {
	return nil
}

func (c *PressKeyCommand) Execute(ctx context.Context) error {
	return nil
}
