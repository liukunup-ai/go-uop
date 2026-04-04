package system

import (
	"context"
	"time"
)

type SwipeCommand struct {
	device interface {
		Swipe(x1, y1, x2, y2 int, duration time.Duration) error
	}
	StartX, StartY int
	EndX, EndY     int
	Duration       time.Duration
}

func NewSwipeCommand(startX, startY, endX, endY int) *SwipeCommand {
	return &SwipeCommand{
		StartX:   startX,
		StartY:   startY,
		EndX:     endX,
		EndY:     endY,
		Duration: 300 * time.Millisecond,
	}
}

func (c *SwipeCommand) Name() string        { return "swipe" }
func (c *SwipeCommand) Description() string { return "Swipe from one point to another" }

func (c *SwipeCommand) Validate() error {
	return nil
}

func (c *SwipeCommand) Execute(ctx context.Context) error {
	if c.device == nil {
		return nil
	}
	return c.device.Swipe(c.StartX, c.StartY, c.EndX, c.EndY, c.Duration)
}
