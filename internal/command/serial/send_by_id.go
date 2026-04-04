package serial

import (
	"context"
	"errors"
	"time"
)

type SendByIDCommand struct {
	baseSerialCommand
	CommandID string
	Timeout   time.Duration
}

func NewSendByIDCommand(commandID string, timeout time.Duration) *SendByIDCommand {
	return &SendByIDCommand{
		CommandID: commandID,
		Timeout:   timeout,
	}
}

func (c *SendByIDCommand) Name() string        { return "sendByID" }
func (c *SendByIDCommand) Description() string { return "Send command by ID" }

func (c *SendByIDCommand) Validate() error {
	if c.CommandID == "" {
		return errors.New("command ID is required")
	}
	return nil
}

func (c *SendByIDCommand) Execute(ctx context.Context) error {
	if c.Serial() == nil {
		return errors.New("serial connection not set")
	}
	return c.Serial().SendByID(c.CommandID, nil)
}
