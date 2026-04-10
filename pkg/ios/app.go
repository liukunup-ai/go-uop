package ios

import (
	"fmt"

	"github.com/danielpaulus/go-ios/ios/instruments"
)

func (d *Device) launchAppGoIOS(bundleID string) error {
	device, err := d.findDevice()
	if err != nil {
		return fmt.Errorf("find device: %w", err)
	}

	pctrl, err := instruments.NewProcessControl(device)
	if err != nil {
		return fmt.Errorf("create process control: %w", err)
	}
	defer pctrl.Close()

	_, err = pctrl.LaunchApp(bundleID, nil)
	return err
}
