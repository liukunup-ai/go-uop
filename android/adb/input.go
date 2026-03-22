package adb

import (
	"strconv"
	"strings"
)

const (
	KeyHome       = 3
	KeyBack       = 4
	KeyEnter      = 66
	KeyVolumeUp   = 24
	KeyVolumeDown = 25
	KeyPower      = 26
)

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
