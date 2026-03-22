package uop

import (
	"testing"
)

func TestNewDevice_InvalidPlatform(t *testing.T) {
	_, err := NewDevice("unknown")
	if err == nil {
		t.Error("expected error for unknown platform")
	}
}

func TestDeviceOption_WithSerial(t *testing.T) {
	opt := WithSerial("test-123")
	if opt == nil {
		t.Error("expected non-nil option")
	}
}

func TestDeviceOption_WithAddress(t *testing.T) {
	opt := WithAddress("localhost:8080")
	if opt == nil {
		t.Error("expected non-nil option")
	}
}
