package ios

import (
	"fmt"

	"github.com/danielpaulus/go-ios/ios"
	"github.com/danielpaulus/go-ios/ios/instruments"
)

func (d *Device) screenshotGoIOS() ([]byte, error) {
	device, err := d.findDevice()
	if err != nil {
		return nil, fmt.Errorf("find device: %w", err)
	}

	screenshotSvc, err := instruments.NewScreenshotService(device)
	if err != nil {
		return nil, fmt.Errorf("create screenshot service: %w", err)
	}
	defer screenshotSvc.Close()

	return screenshotSvc.TakeScreenshot()
}

func (d *Device) findDevice() (ios.DeviceEntry, error) {
	deviceList, err := ios.ListDevices()
	if err != nil {
		return ios.DeviceEntry{}, fmt.Errorf("list devices: %w", err)
	}

	if d.config.udid != "" {
		for _, dev := range deviceList.DeviceList {
			if dev.Properties.SerialNumber == d.config.udid {
				return dev, nil
			}
		}
		return ios.DeviceEntry{}, fmt.Errorf("device not found: %s", d.config.udid)
	}

	if len(deviceList.DeviceList) > 0 {
		return deviceList.DeviceList[0], nil
	}

	return ios.DeviceEntry{}, fmt.Errorf("no devices found")
}
