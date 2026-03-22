package action

import (
	"time"

	"github.com/liukunup/go-uop/internal/locator"
)

type Action interface {
	Do() error
}

type TapAction struct {
	X, Y    int
	Element *locator.Locator
}

type SwipeAction struct {
	StartX, StartY int
	EndX, EndY     int
	Duration       time.Duration
}

type SendKeysAction struct {
	Text    string
	Element *locator.Locator
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
	Element  *locator.Locator
	Optional bool
}
