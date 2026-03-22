package action

import (
	"time"

	"github.com/liukunup/go-uop/internal/selector"
)

type Action interface {
	Do() error
}

type TapAction struct {
	X, Y    int
	Element *selector.Selector
}

type SwipeAction struct {
	StartX, StartY int
	EndX, EndY     int
	Duration       time.Duration
}

type SendKeysAction struct {
	Text    string
	Element *selector.Selector
	Secure  bool
}

type LaunchAction struct {
	AppID     string
	Arguments []string
	WaitIdle  bool
}

type PressKeyAction struct {
	KeyCode int
}

type WaitAction struct {
	Duration time.Duration
	Element  *selector.Selector
	Optional bool
}
