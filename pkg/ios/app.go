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

func (d *Device) terminateAppGoIOS(bundleID string) error {
	device, err := d.findDevice()
	if err != nil {
		return fmt.Errorf("find device: %w", err)
	}

	infoSvc, err := instruments.NewDeviceInfoService(device)
	if err != nil {
		return fmt.Errorf("create device info service: %w", err)
	}
	defer infoSvc.Close()

	processes, err := infoSvc.ProcessList()
	if err != nil {
		return err
	}

	pctrl, err := instruments.NewProcessControl(device)
	if err != nil {
		return fmt.Errorf("create process control: %w", err)
	}
	defer pctrl.Close()

	for _, p := range processes {
		if p.IsApplication && p.Name == bundleID {
			return pctrl.KillProcess(p.Pid)
		}
	}

	return fmt.Errorf("app not found: %s", bundleID)
}

type iosProcess struct {
	Pid   int
	Name  string
	IsApp bool
}
