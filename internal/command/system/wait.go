package system

import (
	"context"
	"errors"
	"time"
)

type Command interface {
	Name() string
	Description() string
	Validate() error
}

type WaitCommand struct {
	Duration time.Duration
}

func NewWaitCommand(ms int) *WaitCommand {
	return &WaitCommand{Duration: time.Duration(ms) * time.Millisecond}
}

func (c *WaitCommand) Name() string        { return "wait" }
func (c *WaitCommand) Description() string { return "Wait for a duration" }

func (c *WaitCommand) Validate() error {
	if c.Duration < 0 {
		return errors.New("duration must be non-negative")
	}
	return nil
}

func (c *WaitCommand) Execute(ctx context.Context) error {
	time.Sleep(c.Duration)
	return nil
}
