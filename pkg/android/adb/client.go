package adb

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Client struct {
	serial string
}

func NewClient(serial ...string) (*Client, error) {
	s := ""
	if len(serial) > 0 {
		s = serial[0]
	}

	if s != "" {
		devices, err := Devices()
		if err != nil {
			return nil, err
		}
		found := false
		for _, d := range devices {
			if d.Serial == s {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("device not found: %s", s)
		}
	}

	return &Client{serial: s}, nil
}

func (c *Client) serialArg() string {
	if c.serial != "" {
		return "-s " + c.serial
	}
	return ""
}

func (c *Client) Serial() string {
	return c.serial
}

func (c *Client) exec(args ...string) (string, error) {
	cmd := exec.Command("adb", append(strings.Fields(c.serialArg()), args...)...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("adb %v: %w\n%s", args, err, buf.String())
	}

	return strings.TrimSpace(buf.String()), nil
}

type DeviceInfo struct {
	Serial    string
	Status    string
	Product   string
	Model     string
	Device    string
	Transport string
}

func Devices() ([]DeviceInfo, error) {
	cmd := exec.Command("adb", "devices", "-l")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("adb devices: %w", err)
	}

	var devices []DeviceInfo
	lines := strings.Split(buf.String(), "\n")
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "List") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		d := DeviceInfo{Serial: parts[0], Status: parts[1]}

		for _, p := range parts[2:] {
			kv := strings.SplitN(p, ":", 2)
			if len(kv) != 2 {
				continue
			}
			switch kv[0] {
			case "product":
				d.Product = kv[1]
			case "model":
				d.Model = strings.ReplaceAll(kv[1], "_", " ")
			case "device":
				d.Device = kv[1]
			case "transport_id":
				d.Transport = kv[1]
			}
		}

		devices = append(devices, d)
	}

	return devices, nil
}
