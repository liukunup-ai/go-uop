package maestro

import (
	"fmt"
	"time"

	"github.com/liukunup/go-uop/core"
	"github.com/liukunup/go-uop/internal/action"
	"github.com/liukunup/go-uop/internal/selector"
)

type Translator struct {
}

func NewTranslator() *Translator {
	return &Translator{}
}

func (t *Translator) TranslateCommand(cmd *MaestroCommand, device core.Device) (action.Action, error) {
	if cmd.Launch != "" {
		return &LaunchWrapper{appID: cmd.Launch, waitIdle: true, device: device}, nil
	}

	if cmd.TapOn != nil {
		return t.translateTapOn(cmd.TapOn, device)
	}

	if cmd.Tap != nil {
		return &TapWrapper{x: cmd.Tap.X, y: cmd.Tap.Y, device: device}, nil
	}

	if cmd.InputText != nil {
		return t.translateInputText(cmd.InputText, device)
	}

	if cmd.Swipe != nil {
		return t.translateSwipe(cmd.Swipe, device)
	}

	if cmd.Terminate != "" {
		return &LaunchWrapper{appID: cmd.Terminate, waitIdle: false, device: device}, nil
	}

	if cmd.Wait > 0 {
		return &WaitWrapper{duration: time.Duration(cmd.Wait) * time.Millisecond}, nil
	}

	return nil, fmt.Errorf("%w: no matching command found", ErrUnsupportedCommand)
}

func (t *Translator) translateTapOn(cmd *TapOnCommand, device core.Device) (action.Action, error) {
	var elem *selector.Selector

	if cmd.ID != "" {
		elem = selector.ByID(cmd.ID)
	} else if cmd.Text != "" {
		elem = selector.ByText(cmd.Text)
	}

	if cmd.Index > 0 && elem != nil {
		elem.SetIndex(cmd.Index)
	}

	return &TapOnWrapper{element: elem, device: device}, nil
}

func (t *Translator) translateInputText(cmd *InputTextCommand, device core.Device) (action.Action, error) {
	var elem *selector.Selector

	if cmd.ID != "" {
		elem = selector.ByID(cmd.ID)
	}

	return &SendKeysWrapper{text: cmd.Text, element: elem, device: device}, nil
}

func (t *Translator) translateSwipe(cmd *SwipeCommand, device core.Device) (action.Action, error) {
	return &SwipeWrapper{
		startX:   cmd.StartX,
		startY:   cmd.StartY,
		endX:     cmd.EndX,
		endY:     cmd.EndY,
		duration: time.Duration(cmd.Duration) * time.Millisecond,
		device:   device,
	}, nil
}

func (t *Translator) TranslateFlow(flow *MaestroFlow, device core.Device) ([]action.Action, error) {
	if flow == nil {
		return nil, nil
	}

	actions := make([]action.Action, 0, len(flow.Steps))

	for i, cmd := range flow.Steps {
		act, err := t.TranslateCommand(&cmd, device)
		if err != nil {
			return nil, fmt.Errorf("step %d: %w", i, err)
		}
		actions = append(actions, act)
	}

	return actions, nil
}

type LaunchWrapper struct {
	appID    string
	waitIdle bool
	device   core.Device
}

func (w *LaunchWrapper) Do() error {
	if w.device == nil {
		return fmt.Errorf("device not set")
	}
	return w.device.Launch()
}

type TapWrapper struct {
	x, y   int
	device core.Device
}

func (w *TapWrapper) Do() error {
	if w.device == nil {
		return fmt.Errorf("device not set")
	}
	return w.device.Tap(w.x, w.y)
}

type TapOnWrapper struct {
	element *selector.Selector
	device  core.Device
}

func (w *TapOnWrapper) Do() error {
	if w.device == nil {
		return fmt.Errorf("device not set")
	}
	// TODO: Find element by selector and tap
	// For now, placeholder - need element finding logic
	return fmt.Errorf("tapOn: element finding not implemented")
}

type SendKeysWrapper struct {
	text    string
	element *selector.Selector
	device  core.Device
}

func (w *SendKeysWrapper) Do() error {
	if w.device == nil {
		return fmt.Errorf("device not set")
	}
	return w.device.SendKeys(w.text)
}

type SwipeWrapper struct {
	startX, startY int
	endX, endY     int
	duration       time.Duration
	device         core.Device
}

func (w *SwipeWrapper) Do() error {
	if w.device == nil {
		return fmt.Errorf("device not set")
	}
	// TODO: Use device swipe if available, otherwise tap sequence
	return fmt.Errorf("swipe: not implemented")
}

type WaitWrapper struct {
	duration time.Duration
}

func (w *WaitWrapper) Do() error {
	time.Sleep(w.duration)
	return nil
}
