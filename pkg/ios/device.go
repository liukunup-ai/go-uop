package ios

import (
	"fmt"

	"github.com/liukunup/go-uop/core"
	"github.com/liukunup/go-uop/pkg/wda"
)

type Device struct {
	client   *wda.Client
	bundleID string
}

func NewDevice(bundleID string, opts ...Option) (*Device, error) {
	cfg := &config{
		address: "http://localhost:8100",
	}
	for _, opt := range opts {
		opt(cfg)
	}

	client, err := wda.NewClient(cfg.address)
	if err != nil {
		return nil, fmt.Errorf("create WDA client: %w", err)
	}

	if !cfg.skipSession {
		if err := client.StartSession(bundleID); err != nil {
			return nil, fmt.Errorf("start session: %w", err)
		}
	}

	return &Device{
		client:   client,
		bundleID: bundleID,
	}, nil
}

func (d *Device) Platform() core.Platform {
	return core.IOS
}

func (d *Device) Info() (map[string]interface{}, error) {
	return map[string]interface{}{
		"platform": "ios",
		"bundleId": d.bundleID,
	}, nil
}

func (d *Device) Screenshot() ([]byte, error) {
	return d.client.Screenshot()
}

func (d *Device) Close() error {
	return d.client.StopSession()
}

func (d *Device) Tap(x, y int) error {
	return d.client.Tap(x, y)
}

func (d *Device) SendKeys(text string) error {
	return d.client.SendKeys(text)
}

func (d *Device) Launch() error {
	return d.client.LaunchApp(d.bundleID)
}

func (d *Device) Terminate() error {
	return d.client.TerminateApp(d.bundleID)
}

func (d *Device) GetSource() (string, error) {
	return d.client.GetSource()
}

func (d *Device) GetAlertText() (string, error) {
	return d.client.GetAlertText()
}

func (d *Device) AcceptAlert() error {
	return d.client.AcceptAlert()
}

func (d *Device) DismissAlert() error {
	return d.client.DismissAlert()
}

var _ core.Device = (*Device)(nil)
