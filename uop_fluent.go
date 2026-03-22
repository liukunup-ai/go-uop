package uop

import (
	"github.com/liukunup/go-uop/internal/locator"
)

type ActionBuilder struct {
	device Device
}

func NewActionBuilder(device Device) *ActionBuilder {
	return &ActionBuilder{device: device}
}

func (ab *ActionBuilder) Tap(x, y int) *ActionBuilder {
	return ab
}

func (ab *ActionBuilder) TapElement(loc *locator.Locator) *ActionBuilder {
	return ab
}

func (ab *ActionBuilder) Swipe(x1, y1, x2, y2 int) *ActionBuilder {
	return ab
}

func (ab *ActionBuilder) SwipeUp() *ActionBuilder {
	return ab
}

func (ab *ActionBuilder) SendKeys(text string) *ActionBuilder {
	return ab
}

func (ab *ActionBuilder) Launch(appID string) *ActionBuilder {
	return ab
}

func (ab *ActionBuilder) Wait(duration string) *ActionBuilder {
	return ab
}

func (ab *ActionBuilder) Do() error {
	return nil
}
