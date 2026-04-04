package runner

import (
	"fmt"

	"github.com/liukunup/go-uop/core"
	"github.com/liukunup/go-uop/pkg/android"
	"github.com/liukunup/go-uop/pkg/ios"
)

// deviceFactory creates device instances
type deviceFactory interface {
	Create(deviceType, serial string) (core.Device, error)
}

// defaultDeviceFactory creates real devices
type defaultDeviceFactory struct{}

func (f *defaultDeviceFactory) Create(deviceType, serial string) (core.Device, error) {
	switch deviceType {
	case "ios":
		return ios.NewDevice("")
	case "android":
		return android.NewDevice(android.WithUDID(serial))
	default:
		return nil, fmt.Errorf("unsupported device type: %s", deviceType)
	}
}

// PoolDevice represents a device in the pool
type PoolDevice struct {
	ID     string
	Type   string
	Serial string
	device core.Device
}

// DevicePool manages multiple devices
type DevicePool struct {
	devices       map[string]*PoolDevice
	currentDevice string
	defaultDevice string
	factory       deviceFactory
}

// NewDevicePool creates a new device pool
func NewDevicePool() *DevicePool {
	return &DevicePool{
		devices: make(map[string]*PoolDevice),
		factory: &defaultDeviceFactory{},
	}
}

// newDevicePoolWithFactory creates a device pool with a custom factory (for testing)
func newDevicePoolWithFactory(factory deviceFactory) *DevicePool {
	return &DevicePool{
		devices: make(map[string]*PoolDevice),
		factory: factory,
	}
}

// AddDevice adds a device to the pool
func (p *DevicePool) AddDevice(id, deviceType, serial string) error {
	if _, exists := p.devices[id]; exists {
		return fmt.Errorf("device %s already exists", id)
	}

	dev, err := p.factory.Create(deviceType, serial)
	if err != nil {
		return fmt.Errorf("failed to create device: %w", err)
	}

	poolDev := &PoolDevice{
		ID:     id,
		Type:   deviceType,
		Serial: serial,
		device: dev,
	}
	p.devices[id] = poolDev

	if p.defaultDevice == "" {
		p.defaultDevice = id
		p.currentDevice = id
	}

	return nil
}

// GetDevice returns a device by ID
func (p *DevicePool) GetDevice(id string) (*PoolDevice, error) {
	dev, exists := p.devices[id]
	if !exists {
		return nil, fmt.Errorf("device %s not found", id)
	}
	return dev, nil
}

// CurrentDevice returns the currently selected device
func (p *DevicePool) CurrentDevice() *PoolDevice {
	if p.currentDevice == "" {
		return nil
	}
	return p.devices[p.currentDevice]
}

// DefaultDevice returns the default device
func (p *DevicePool) DefaultDevice() *PoolDevice {
	if p.defaultDevice == "" {
		return nil
	}
	return p.devices[p.defaultDevice]
}

// SwitchDevice switches to a different device
func (p *DevicePool) SwitchDevice(id string) error {
	dev, exists := p.devices[id]
	if !exists {
		return fmt.Errorf("device %s not found", id)
	}
	p.currentDevice = dev.ID
	return nil
}

// Close closes all devices in the pool
func (p *DevicePool) Close() error {
	for _, dev := range p.devices {
		if dev.device != nil {
			dev.device.Close()
		}
	}
	p.devices = make(map[string]*PoolDevice)
	p.currentDevice = ""
	p.defaultDevice = ""
	return nil
}

func (p *DevicePool) ListDevices() map[string]*PoolDevice {
	return p.devices
}
