package device

import (
	"context"
	"errors"
)

type SendKeysCommand struct {
	baseDeviceCommand
	Text   string
	Secure bool
	Enter  bool
}

func NewSendKeysCommand(text string) *SendKeysCommand {
	return &SendKeysCommand{Text: text}
}

func (c *SendKeysCommand) Name() string        { return "inputText" }
func (c *SendKeysCommand) Description() string { return "Input text" }

func (c *SendKeysCommand) Validate() error {
	return nil
}

func (c *SendKeysCommand) Execute(ctx context.Context) error {
	if c.Device() == nil {
		return errors.New("device not set")
	}
	return c.Device().SendKeys(c.Text)
}
