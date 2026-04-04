package runner

import (
	"errors"
	"testing"

	"github.com/liukunup/go-uop/core"
)

// mockDevice implements core.Device for testing
type mockDevice struct {
	id         string
	deviceType core.Platform
	serial     string
}

func (m *mockDevice) Platform() core.Platform {
	return m.deviceType
}

func (m *mockDevice) Info() (map[string]interface{}, error) {
	return map[string]interface{}{
		"id":     m.id,
		"type":   m.deviceType,
		"serial": m.serial,
	}, nil
}

func (m *mockDevice) Screenshot() ([]byte, error) {
	return nil, nil
}

func (m *mockDevice) Tap(x, y int) error {
	return nil
}

func (m *mockDevice) SendKeys(text string) error {
	return nil
}

func (m *mockDevice) Launch() error {
	return nil
}

func (m *mockDevice) Close() error {
	return nil
}

// mockDeviceFactory creates mock devices for testing
type mockDeviceFactory struct {
	devices map[string]*mockDevice
}

func newMockDeviceFactory() *mockDeviceFactory {
	return &mockDeviceFactory{
		devices: make(map[string]*mockDevice),
	}
}

func (f *mockDeviceFactory) Create(deviceType, serial string) (core.Device, error) {
	var platform core.Platform
	switch deviceType {
	case "ios":
		platform = core.IOS
	case "android":
		platform = core.Android
	default:
		return nil, errors.New("unsupported device type: " + deviceType)
	}
	dev := &mockDevice{
		id:         deviceType + "-" + serial,
		deviceType: platform,
		serial:     serial,
	}
	f.devices[deviceType+serial] = dev
	return dev, nil
}

// ErrMockDeviceFailure is used to simulate device creation failure
var ErrMockDeviceFailure = errors.New("mock device creation failure")

func TestNewDevicePool(t *testing.T) {
	pool := NewDevicePool()
	if pool == nil {
		t.Fatal("expected non-nil DevicePool")
	}
	if len(pool.devices) != 0 {
		t.Errorf("expected empty pool, got %d devices", len(pool.devices))
	}
}

func TestDevicePool_AddDevice(t *testing.T) {
	factory := newMockDeviceFactory()
	pool := newDevicePoolWithFactory(factory)

	err := pool.AddDevice("iphone", "ios", "00001234-00123456789")
	if err != nil {
		t.Fatalf("AddDevice failed: %v", err)
	}

	dev, ok := pool.devices["iphone"]
	if !ok {
		t.Fatal("expected iphone in pool")
	}
	if dev.ID != "iphone" {
		t.Errorf("expected ID 'iphone', got '%s'", dev.ID)
	}
	if dev.Type != "ios" {
		t.Errorf("expected Type 'ios', got '%s'", dev.Type)
	}
	if dev.Serial != "00001234-00123456789" {
		t.Errorf("expected Serial '00001234-00123456789', got '%s'", dev.Serial)
	}

	if pool.defaultDevice != "iphone" {
		t.Errorf("expected default device 'iphone', got '%s'", pool.defaultDevice)
	}

	if pool.currentDevice != "iphone" {
		t.Errorf("expected current device 'iphone', got '%s'", pool.currentDevice)
	}

	err = pool.AddDevice("android-tablet", "android", "emulator-5554")
	if err != nil {
		t.Fatalf("AddDevice failed: %v", err)
	}

	dev, ok = pool.devices["android-tablet"]
	if !ok {
		t.Fatal("expected android-tablet in pool")
	}
	if dev.ID != "android-tablet" {
		t.Errorf("expected ID 'android-tablet', got '%s'", dev.ID)
	}

	if pool.defaultDevice != "iphone" {
		t.Errorf("expected default device still 'iphone', got '%s'", pool.defaultDevice)
	}

	if pool.currentDevice != "iphone" {
		t.Errorf("expected current device still 'iphone', got '%s'", pool.currentDevice)
	}
}

func TestDevicePool_GetDevice(t *testing.T) {
	factory := newMockDeviceFactory()
	pool := newDevicePoolWithFactory(factory)

	err := pool.AddDevice("test-device", "ios", "12345")
	if err != nil {
		t.Fatalf("AddDevice failed: %v", err)
	}

	dev, err := pool.GetDevice("test-device")
	if err != nil {
		t.Fatalf("GetDevice failed: %v", err)
	}
	if dev == nil {
		t.Fatal("expected device, got nil")
	}
	if dev.ID != "test-device" {
		t.Errorf("expected ID 'test-device', got '%s'", dev.ID)
	}

	_, err = pool.GetDevice("non-existent")
	if err == nil {
		t.Error("expected error for non-existent device")
	}
}

func TestDevicePool_CurrentDevice(t *testing.T) {
	factory := newMockDeviceFactory()
	pool := newDevicePoolWithFactory(factory)

	current := pool.CurrentDevice()
	if current != nil {
		t.Error("expected nil current device for empty pool")
	}

	err := pool.AddDevice("iphone", "ios", "12345")
	if err != nil {
		t.Fatalf("AddDevice failed: %v", err)
	}

	current = pool.CurrentDevice()
	if current == nil {
		t.Fatal("expected current device, got nil")
	}
	if current.ID != "iphone" {
		t.Errorf("expected current device 'iphone', got '%s'", current.ID)
	}

	err = pool.AddDevice("android", "android", "emulator-5554")
	if err != nil {
		t.Fatalf("AddDevice failed: %v", err)
	}

	current = pool.CurrentDevice()
	if current.ID != "iphone" {
		t.Errorf("expected current device still 'iphone', got '%s'", current.ID)
	}
}

func TestDevicePool_DefaultDevice(t *testing.T) {
	factory := newMockDeviceFactory()
	pool := newDevicePoolWithFactory(factory)

	defaultDev := pool.DefaultDevice()
	if defaultDev != nil {
		t.Error("expected nil default device for empty pool")
	}

	err := pool.AddDevice("iphone", "ios", "12345")
	if err != nil {
		t.Fatalf("AddDevice failed: %v", err)
	}

	defaultDev = pool.DefaultDevice()
	if defaultDev == nil {
		t.Fatal("expected default device, got nil")
	}
	if defaultDev.ID != "iphone" {
		t.Errorf("expected default device 'iphone', got '%s'", defaultDev.ID)
	}

	err = pool.AddDevice("android", "android", "emulator-5554")
	if err != nil {
		t.Fatalf("AddDevice failed: %v", err)
	}

	defaultDev = pool.DefaultDevice()
	if defaultDev.ID != "iphone" {
		t.Errorf("expected default device still 'iphone', got '%s'", defaultDev.ID)
	}
}

func TestDevicePool_SwitchDevice(t *testing.T) {
	factory := newMockDeviceFactory()
	pool := newDevicePoolWithFactory(factory)

	err := pool.AddDevice("iphone", "ios", "00001234-00123456789")
	if err != nil {
		t.Fatalf("AddDevice failed: %v", err)
	}
	err = pool.AddDevice("android-tablet", "android", "emulator-5554")
	if err != nil {
		t.Fatalf("AddDevice failed: %v", err)
	}

	current := pool.CurrentDevice()
	if current.ID != "iphone" {
		t.Errorf("expected current iphone, got %s", current.ID)
	}

	err = pool.SwitchDevice("android-tablet")
	if err != nil {
		t.Fatalf("SwitchDevice failed: %v", err)
	}

	current = pool.CurrentDevice()
	if current.ID != "android-tablet" {
		t.Errorf("expected current android-tablet, got %s", current.ID)
	}

	err = pool.SwitchDevice("iphone")
	if err != nil {
		t.Fatalf("SwitchDevice failed: %v", err)
	}

	current = pool.CurrentDevice()
	if current.ID != "iphone" {
		t.Errorf("expected current iphone, got %s", current.ID)
	}

	err = pool.SwitchDevice("non-existent")
	if err == nil {
		t.Error("expected error for non-existent device")
	}

	current = pool.CurrentDevice()
	if current.ID != "iphone" {
		t.Errorf("expected current still iphone, got %s", current.ID)
	}
}

func TestDevicePool_Close(t *testing.T) {
	factory := newMockDeviceFactory()
	pool := newDevicePoolWithFactory(factory)

	err := pool.AddDevice("iphone", "ios", "12345")
	if err != nil {
		t.Fatalf("AddDevice failed: %v", err)
	}
	err = pool.AddDevice("android", "android", "emulator-5554")
	if err != nil {
		t.Fatalf("AddDevice failed: %v", err)
	}

	err = pool.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	if len(pool.devices) != 0 {
		t.Errorf("expected empty pool after close, got %d devices", len(pool.devices))
	}
}

func TestDevicePool_CloseEmpty(t *testing.T) {
	pool := NewDevicePool()

	err := pool.Close()
	if err != nil {
		t.Fatalf("Close failed on empty pool: %v", err)
	}
}

func TestDevicePool(t *testing.T) {
	factory := newMockDeviceFactory()
	pool := newDevicePoolWithFactory(factory)

	err := pool.AddDevice("iphone", "ios", "00001234-00123456789")
	if err != nil {
		t.Fatalf("AddDevice failed: %v", err)
	}
	err = pool.AddDevice("android-tablet", "android", "emulator-5554")
	if err != nil {
		t.Fatalf("AddDevice failed: %v", err)
	}

	defaultDev := pool.DefaultDevice()
	if defaultDev == nil {
		t.Fatal("expected default device, got nil")
	}
	if defaultDev.ID != "iphone" {
		t.Errorf("expected default iphone, got %s", defaultDev.ID)
	}

	err = pool.SwitchDevice("android-tablet")
	if err != nil {
		t.Fatalf("SwitchDevice failed: %v", err)
	}
	current := pool.CurrentDevice()
	if current.ID != "android-tablet" {
		t.Errorf("expected current android-tablet, got %s", current.ID)
	}

	_, err = pool.GetDevice("non-existent")
	if err == nil {
		t.Error("expected error for non-existent device")
	}
}

func TestDevicePool_AddDuplicate(t *testing.T) {
	factory := newMockDeviceFactory()
	pool := newDevicePoolWithFactory(factory)

	err := pool.AddDevice("iphone", "ios", "12345")
	if err != nil {
		t.Fatalf("AddDevice failed: %v", err)
	}

	err = pool.AddDevice("iphone", "android", "emulator-5554")
	if err == nil {
		t.Error("expected error for duplicate device ID")
	}
}

func TestDevicePool_UnsupportedDeviceType(t *testing.T) {
	factory := newMockDeviceFactory()
	pool := newDevicePoolWithFactory(factory)

	err := pool.AddDevice("device", "unsupported", "serial")
	if err == nil {
		t.Error("expected error for unsupported device type")
	}
}
