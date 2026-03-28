package uop

import (
	"fmt"
	"strconv"
	"time"

	"github.com/liukunup/go-uop/internal/selector"
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

func (ab *ActionBuilder) TapElement(loc *selector.Selector) *ActionBuilder {
	if ab.err != nil {
		return ab
	}

	finder, ok := ab.device.(ElementFinder)
	if !ok {
		ab.err = fmt.Errorf("device does not support element finding")
		return ab
	}

	elem, err := finder.FindElement(loc)
	if err != nil {
		ab.err = fmt.Errorf("find element: %w", err)
		return ab
	}

	if err := ab.device.Tap(elem.X, elem.Y); err != nil {
		ab.err = fmt.Errorf("tap element: %w", err)
	}
	return ab
}

func (ab *ActionBuilder) Swipe(x1, y1, x2, y2 int) *ActionBuilder {
	return ab.SwipeWithDuration(x1, y1, x2, y2, 300*time.Millisecond)
}

func (ab *ActionBuilder) SwipeWithDuration(x1, y1, x2, y2 int, duration time.Duration) *ActionBuilder {
	if ab.err != nil {
		return ab
	}
	if dev, ok := ab.device.(interface {
		Swipe(x1, y1, x2, y2 int, duration time.Duration) error
	}); ok {
		if err := dev.Swipe(x1, y1, x2, y2, duration); err != nil {
			ab.err = fmt.Errorf("swipe: %w", err)
		}
	}
	return ab
}

func (ab *ActionBuilder) SwipeUp() *ActionBuilder {
	return ab.SwipeUpDistance(0.2)
}

func (ab *ActionBuilder) SwipeUpDistance(ratio float64) *ActionBuilder {
	if ab.err != nil {
		return ab
	}

	swiper, ok := ab.device.(Swiper)
	if !ok {
		ab.err = fmt.Errorf("device does not support swipe gestures")
		return ab
	}

	if err := swiper.SwipeUp(ratio); err != nil {
		ab.err = fmt.Errorf("swipe up: %w", err)
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

func (ab *ActionBuilder) Launch() *ActionBuilder {
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

	d, err := time.ParseDuration(duration)
	if err != nil {
		// Try parsing as plain milliseconds
		if ms, err := strconv.ParseInt(duration, 10, 64); err == nil {
			d = time.Duration(ms) * time.Millisecond
		} else {
			ab.err = fmt.Errorf("parse duration: %w", err)
			return ab
		}
	}

	time.Sleep(d)
	return ab
}

func (ab *ActionBuilder) Do() error {
	return ab.err
}

type Swiper interface {
	SwipeUp(ratio float64) error
}

type ElementFinder interface {
	FindElement(loc *selector.Selector) (*ElementInfo, error)
}

type ElementInfo struct {
	X, Y    int
	Width   int
	Height  int
	Text    string
	Visible bool
}
