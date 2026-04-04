package serial

import (
	"context"
	"errors"
)

type SendRawCommand struct {
	baseSerialCommand
	Data string
}

func NewSendRawCommand(data string) *SendRawCommand {
	return &SendRawCommand{Data: data}
}

func (c *SendRawCommand) Name() string        { return "sendRaw" }
func (c *SendRawCommand) Description() string { return "Send raw data" }

func (c *SendRawCommand) Execute(ctx context.Context) error {
	if c.Serial() == nil {
		return errors.New("serial connection not set")
	}
	_, err := c.Serial().WriteString(c.Data)
	return err
}
