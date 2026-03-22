package uop

import (
	"fmt"
	"time"

	"github.com/liukunup/go-uop/internal/locator"
)

type ActionBuilder struct {
	device Device
	err    error
}

func NewActionBuilder(device Device) *ActionBuilder {
	return &ActionBuilder{device: device}
}

func (ab *ActionBuilder) Tap(x, y int) *ActionBuilder {
	if ab.err != nil {
		return ab
	}
	if err := ab.device.Tap(x, y); err != nil {
		ab.err = fmt.Errorf("tap: %w", err)
	}
	return ab
}

func (ab *ActionBuilder) TapElement(loc *locator.Locator) *ActionBuilder {
	if ab.err != nil {
		return ab
	}
	return ab
}

func (ab *ActionBuilder) Swipe(x1, y1, x2, y2 int) *ActionBuilder {
	if ab.err != nil {
		return ab
	}
	if dev, ok := ab.device.(interface {
		Swipe(x1, y1, x2, y2 int, duration time.Duration) error
	}); ok {
		if err := dev.Swipe(x1, y1, x2, y2, 300*time.Millisecond); err != nil {
			ab.err = fmt.Errorf("swipe: %w", err)
		}
	}
	return ab
}

func (ab *ActionBuilder) SwipeUp() *ActionBuilder {
	if ab.err != nil {
		return ab
	}
	return ab
}

func (ab *ActionBuilder) SendKeys(text string) *ActionBuilder {
	if ab.err != nil {
		return ab
	}
	if err := ab.device.SendKeys(text); err != nil {
		ab.err = fmt.Errorf("sendKeys: %w", err)
	}
	return ab
}

func (ab *ActionBuilder) Launch(appID string) *ActionBuilder {
	if ab.err != nil {
		return ab
	}
	if err := ab.device.Launch(); err != nil {
		ab.err = fmt.Errorf("launch: %w", err)
	}
	return ab
}

func (ab *ActionBuilder) Wait(duration string) *ActionBuilder {
	if ab.err != nil {
		return ab
	}
	return ab
}

func (ab *ActionBuilder) Do() error {
	return ab.err
}
