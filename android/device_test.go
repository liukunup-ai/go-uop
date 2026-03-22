package android

import (
	"testing"

	"github.com/liukunup/go-uop/core"
)

func TestNewDevice_WithOptions(t *testing.T) {
	_, err := NewDevice(WithUDID("test-serial"), WithPackage("com.example.app"))
	if err != nil {
		t.Skip("requires adb:", err)
	}
}

func TestOption_WithUDID(t *testing.T) {
	opt := WithUDID("test-udid")
	cfg := &config{}
	opt(cfg)
	if cfg.udid != "test-udid" {
		t.Errorf("expected udid 'test-udid', got '%s'", cfg.udid)
	}
}

func TestOption_WithPackage(t *testing.T) {
	opt := WithPackage("com.example.app")
	cfg := &config{}
	opt(cfg)
	if cfg.packageName != "com.example.app" {
		t.Errorf("expected package 'com.example.app', got '%s'", cfg.packageName)
	}
}

func TestOption_Multiple(t *testing.T) {
	cfg := &config{}
	WithUDID("serial-123")(cfg)
	WithPackage("com.example.app")(cfg)

	if cfg.udid != "serial-123" {
		t.Errorf("expected udid 'serial-123', got '%s'", cfg.udid)
	}
	if cfg.packageName != "com.example.app" {
		t.Errorf("expected package 'com.example.app', got '%s'", cfg.packageName)
	}
}

func TestDevice_Platform(t *testing.T) {
	device := &Device{}
	if device.Platform() != core.Android {
		t.Errorf("expected Android platform, got %v", device.Platform())
	}
}

func TestDevice_ImplementsCoreDevice(t *testing.T) {
	var _ core.Device = (*Device)(nil)
}

func TestDevice_Launch_WithoutPackage(t *testing.T) {
	device := &Device{}
	err := device.Launch()
	if err == nil {
		t.Error("expected error when launching without package")
	}
}

func TestDevice_Terminate_WithoutPackage(t *testing.T) {
	device := &Device{}
	err := device.Terminate()
	if err == nil {
		t.Error("expected error when terminating without package")
	}
}

func TestDevice_Close(t *testing.T) {
	device := &Device{}
	err := device.Close()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
