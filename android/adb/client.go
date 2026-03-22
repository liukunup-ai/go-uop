package adb

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
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

func (c *Client) Tap(x, y int) error {
	_, err := c.exec("shell", "input", "tap", strconv.Itoa(x), strconv.Itoa(y))
	return err
}

func (c *Client) Swipe(x1, y1, x2, y2 int, durationMs int) error {
	args := []string{"shell", "input", "swipe",
		strconv.Itoa(x1), strconv.Itoa(y1),
		strconv.Itoa(x2), strconv.Itoa(y2)}
	if durationMs > 0 {
		args = append(args, strconv.Itoa(durationMs))
	}
	_, err := c.exec(args...)
	return err
}

func (c *Client) SendText(text string) error {
	escaped := strings.ReplaceAll(text, " ", "%s")
	_, err := c.exec("shell", "input", "text", escaped)
	return err
}

func (c *Client) PressKey(keyCode int) error {
	_, err := c.exec("shell", "input", "keyevent", strconv.Itoa(keyCode))
	return err
}

const (
	KeyHome       = 3
	KeyBack       = 4
	KeyEnter      = 66
	KeyVolumeUp   = 24
	KeyVolumeDown = 25
	KeyPower      = 26
)

func (c *Client) Screenshot() ([]byte, error) {
	out := "/sdcard/screenshot.png"
	if _, err := c.exec("shell", "screencap", "-p", out); err != nil {
		return nil, err
	}

	cmd := exec.Command("adb", append(strings.Fields(c.serialArg()), "exec-out", "screencap", "-p")...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("screencap: %w", err)
	}

	return buf.Bytes(), nil
}

func (c *Client) Shell(command string) (string, error) {
	return c.exec("shell", command)
}

func (c *Client) StartActivity(component string) error {
	_, err := c.exec("shell", "am", "start", "-n", component)
	return err
}

func (c *Client) StopPackage(packageName string) error {
	_, err := c.exec("shell", "am", "force-stop", packageName)
	return err
}

func (c *Client) Install(apkPath string, grantPerms bool) error {
	args := []string{"install"}
	if grantPerms {
		args = append(args, "-g")
	}
	args = append(args, apkPath)
	_, err := c.exec(args...)
	return err
}

func (c *Client) Uninstall(packageName string) error {
	_, err := c.exec("uninstall", packageName)
	return err
}

func (c *Client) CurrentPackage() (string, error) {
	output, err := c.Shell("dumpsys activity activities | grep mResumedActivity")
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`([\w.]+)/`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		return "", fmt.Errorf("could not find foreground package")
	}

	return matches[1], nil
}

func (c *Client) WaitForIdle(timeout time.Duration) error {
	end := time.Now().Add(timeout)
	for time.Now().Before(end) {
		_, err := c.Shell("sleep 0.1 && dumpsys activity activities")
		if err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for idle")
}
