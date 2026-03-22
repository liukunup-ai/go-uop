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
	Text              string
	Element           *selector.Selector
	Secure            bool
	Enter             bool
	ClearExistingText bool
}

type LaunchAction struct {
	AppID      string
	Arguments  []string
	WaitIdle   bool
	ClearState bool
}

type PressKeyAction struct {
	KeyCode int
}

type WaitAction struct {
	Duration time.Duration
	Element  *selector.Selector
	Optional bool
}

type KillAction struct {
	AppID string
}

type StopAction struct {
	AppID    string
	Graceful bool
}

type ClearStateAction struct {
	AppID string
}

type AssertAction struct {
	Element   *selector.Selector
	MustExist bool
	Timeout   time.Duration
}

type ScreenshotAction struct {
	Name string
	Path string
}

type RunFlowAction struct {
	SubflowPath string
	EnvVars     map[string]string
	Depth       int
}

func (a *RunFlowAction) Do() error {
	return nil
}
