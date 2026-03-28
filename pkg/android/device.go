package android

import (
	"fmt"
	"time"

	"github.com/liukunup/go-uop/android/adb"
	"github.com/liukunup/go-uop/core"
)

type Device struct {
	client *adb.Client
	pkg    string
}

func NewDevice(opts ...Option) (*Device, error) {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}

	client, err := adb.NewClient(cfg.udid)
	if err != nil {
		return nil, fmt.Errorf("create ADB client: %w", err)
	}

	return &Device{
		client: client,
		pkg:    cfg.packageName,
	}, nil
}

func (d *Device) Platform() core.Platform {
	return core.Android
}

func (d *Device) Info() (map[string]interface{}, error) {
	devices, err := adb.Devices()
	if err != nil {
		return nil, err
	}

	for _, dev := range devices {
		if dev.Serial == d.client.Serial() || d.client.Serial() == "" {
			return map[string]interface{}{
				"platform": "android",
				"serial":   dev.Serial,
				"model":    dev.Model,
				"product":  dev.Product,
			}, nil
		}
	}

	return map[string]interface{}{
		"platform": "android",
		"serial":   d.client.Serial(),
	}, nil
}

func (d *Device) Screenshot() ([]byte, error) {
	return d.client.Screenshot()
}

func (d *Device) Close() error {
	return nil
}

func (d *Device) Tap(x, y int) error {
	return d.client.Tap(x, y)
}

func (d *Device) SendKeys(text string) error {
	return d.client.SendText(text)
}

func (d *Device) Launch() error {
	if d.pkg == "" {
		return fmt.Errorf("package name not set")
	}
	return d.client.StartActivity(d.pkg + "/.MainActivity")
}

func (d *Device) Terminate() error {
	if d.pkg == "" {
		return fmt.Errorf("package name not set")
	}
	return d.client.StopPackage(d.pkg)
}

func (d *Device) GetSource() (string, error) {
	return d.client.Shell("uiautomator dump /sdcard/ui.xml && cat /sdcard/ui.xml")
}

func (d *Device) PressKey(keyCode int) error {
	return d.client.PressKey(keyCode)
}

func (d *Device) Swipe(x1, y1, x2, y2 int, duration time.Duration) error {
	return d.client.Swipe(x1, y1, x2, y2, int(duration.Milliseconds()))
}

var _ core.Device = (*Device)(nil)
