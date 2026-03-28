package console

import (
	"fmt"
	"sync"
	"time"

	"github.com/liukunup/go-uop/core"
	"github.com/liukunup/go-uop/pkg/android"
	"github.com/liukunup/go-uop/pkg/android/adb"
	"github.com/liukunup/go-uop/pkg/ios"
)

type DeviceManager struct {
	mu      sync.RWMutex
	devices map[string]core.Device
	info    map[string]*Device
}

func NewDeviceManager() *DeviceManager {
	return &DeviceManager{
		devices: make(map[string]core.Device),
		info:    make(map[string]*Device),
	}
}

// ListDevices lists all available devices
func (m *DeviceManager) ListDevices() ([]*Device, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := []*Device{}

	// List Android devices
	androidDevices, err := adb.Devices()
	if err == nil {
		for _, d := range androidDevices {
			status := "available"
			if _, connected := m.devices["android-"+d.Serial]; connected {
				status = "connected"
			}
			result = append(result, &Device{
				ID:       "android-" + d.Serial,
				Platform: "android",
				Name:     d.Model,
				Serial:   d.Serial,
				Status:   status,
				Model:    d.Model,
			})
		}
	}

	// iOS devices need user to manually provide address
	// Currently only return connected iOS devices
	for id, info := range m.info {
		if info.Platform == "ios" {
			status := "available"
			if _, connected := m.devices[id]; connected {
				status = "connected"
			}
			result = append(result, &Device{
				ID:       id,
				Platform: "ios",
				Name:     info.Name,
				Serial:   info.Serial,
				Status:   status,
				Address:  info.Address,
			})
		}
	}

	return result, nil
}

// ConnectDevice connects to a device
func (m *DeviceManager) ConnectDevice(d *Device) (core.Device, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var device core.Device
	var err error

	switch d.Platform {
	case "android":
		device, err = android.NewDevice(
			android.WithUDID(d.Serial),
			android.WithPackage(d.PkgName),
		)
	case "ios":
		opts := []ios.Option{ios.WithAddress(d.Address)}
		if d.SkipSession {
			opts = append(opts, ios.SkipSession())
		}
		device, err = ios.NewDevice(
			d.Serial, // bundleID for iOS
			opts...,
		)
	default:
		err = ErrUnsupportedPlatform
	}

	if err != nil {
		return nil, err
	}

	m.devices[d.ID] = device
	m.info[d.ID] = d
	d.Status = "connected"

	return device, nil
}

// DisconnectDevice disconnects a device
func (m *DeviceManager) DisconnectDevice(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	device, ok := m.devices[id]
	if !ok {
		return ErrDeviceNotFound
	}

	if err := device.Close(); err != nil {
		return err
	}

	delete(m.devices, id)
	if info, ok := m.info[id]; ok {
		info.Status = "available"
	}

	return nil
}

// GetConnected gets a connected device
func (m *DeviceManager) GetConnected(id string) (core.Device, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	device, ok := m.devices[id]
	if !ok {
		return nil, ErrDeviceNotConnected
	}

	return device, nil
}

// ExecuteCommand executes a command on a device
func (m *DeviceManager) ExecuteCommand(id string, cmd string, params map[string]interface{}) (*CommandRecord, error) {
	device, err := m.GetConnected(id)
	if err != nil {
		return nil, err
	}

	record := &CommandRecord{
		ID:        generateID(),
		Timestamp: time.Now().Format(time.RFC3339),
		Type:      cmd,
		Params:    params,
	}

	start := time.Now()
	defer func() {
		record.Duration = time.Since(start).String()
	}()

	switch cmd {
	case "tap":
		x, _ := toInt(params["x"])
		y, _ := toInt(params["y"])
		err = device.Tap(x, y)
	case "input":
		text, _ := toString(params["text"])
		err = device.SendKeys(text)
	case "launch":
		err = device.Launch()
	case "terminate":
		if t, ok := device.(interface{ Terminate() error }); ok {
			err = t.Terminate()
		} else {
			err = ErrUnknownCommand
		}
	case "screenshot":
		// Screenshot handled separately
	default:
		err = ErrUnknownCommand
	}

	record.Success = err == nil
	if err != nil {
		record.Output = err.Error()
	}

	return record, err
}

// Helper functions
func generateID() string {
	return fmt.Sprintf("cmd-%d", time.Now().UnixNano())
}

func toInt(v interface{}) (int, bool) {
	switch val := v.(type) {
	case int:
		return val, true
	case float64:
		return int(val), true
	case string:
		i, err := fmt.Sscanf(val, "%d", new(int))
		if err == nil && i > 0 {
			return i, true
		}
	}
	return 0, false
}

func toString(v interface{}) (string, bool) {
	if s, ok := v.(string); ok {
		return s, true
	}
	return "", false
}
