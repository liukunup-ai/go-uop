package adb

import (
	"testing"
)

func TestNewClient_Empty(t *testing.T) {
	t.Skip("requires adb")
}

func TestDevices_Format(t *testing.T) {
	t.Skip("requires adb")

	devices, err := Devices()
	if err != nil {
		t.Skipf("adb not available: %v", err)
	}

	for _, d := range devices {
		if d.Serial == "" {
			t.Error("device serial should not be empty")
		}
		if d.Status == "" {
			t.Error("device status should not be empty")
		}
	}
}

func TestClient_ExecValidation(t *testing.T) {
	c := &Client{serial: "test-serial"}
	arg := c.serialArg()
	if arg != "-s test-serial" {
		t.Errorf("unexpected serial arg: %s", arg)
	}
}

func TestClient_NoSerial(t *testing.T) {
	c := &Client{serial: ""}
	arg := c.serialArg()
	if arg != "" {
		t.Errorf("unexpected serial arg: %s", arg)
	}
}
