package ios

import (
	"fmt"
	"sync"

	"github.com/liukunup/go-uop/core"
)

type Device struct {
	mu     sync.Mutex
	config *config
	wda    *wdaClient
}

func NewDevice(bundleID string, opts ...Option) (*Device, error) {
	cfg := &config{
		address: "http://localhost:8100",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	cfg.bundleID = bundleID

	wdaClient, err := newWDAClient(cfg.address, bundleID)
	if err != nil {
		return nil, fmt.Errorf("create WDA client: %w", err)
	}

	if !cfg.skipSession {
		if err := wdaClient.client.StartSession(bundleID); err != nil {
			return nil, fmt.Errorf("start WDA session: %w", err)
		}
	}

	return &Device{
		config: cfg,
		wda:    wdaClient,
	}, nil
}

func (d *Device) Platform() core.Platform {
	return core.IOS
}

func (d *Device) Info() (map[string]interface{}, error) {
	return map[string]interface{}{
		"platform": "ios",
		"bundleId": d.config.bundleID,
		"wda":      d.config.address,
	}, nil
}

func (d *Device) Screenshot() ([]byte, error) {
	if img, err := d.screenshotGoIOS(); err == nil {
		return img, nil
	}

	return d.wda.client.Screenshot()
}

func (d *Device) Close() error {
	if d.wda != nil {
		return d.wda.Close()
	}
	return nil
}

func (d *Device) Tap(x, y int) error {
	return d.wda.Tap(x, y)
}

func (d *Device) SendKeys(text string) error {
	return d.wda.SendKeys(text)
}

func (d *Device) Launch() error {
	return d.launchAppGoIOS(d.config.bundleID)
}

func (d *Device) GetSource() (string, error) {
	return d.wda.GetSource()
}

func (d *Device) GetAlertText() (string, error) {
	return d.wda.GetAlertText()
}

func (d *Device) AcceptAlert() error {
	return d.wda.AcceptAlert()
}

func (d *Device) DismissAlert() error {
	return d.wda.DismissAlert()
}

var _ core.Device = (*Device)(nil)
