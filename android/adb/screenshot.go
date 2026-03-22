package adb

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
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
