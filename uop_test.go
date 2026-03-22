package uop

import (
	"testing"
	"time"
)

func TestNewDevice_InvalidPlatform(t *testing.T) {
	_, err := NewDevice("unknown")
	if err == nil {
		t.Error("expected error for unknown platform")
	}
}

func TestDeviceOption_WithSerial(t *testing.T) {
	opt := WithSerial("test-123")
	cfg := &deviceConfig{}
	opt(cfg)

	if cfg.serial != "test-123" {
		t.Errorf("expected serial 'test-123', got '%s'", cfg.serial)
	}
}

func TestDeviceOption_WithTimeout(t *testing.T) {
	opt := WithTimeout(30 * time.Second)
	cfg := &deviceConfig{}
	opt(cfg)

	if cfg.timeout != 30*time.Second {
		t.Errorf("expected 30s timeout, got %v", cfg.timeout)
	}
}
