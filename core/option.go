package core

import "time"

// DeviceOption configures device creation
type DeviceOption func(*deviceConfig)

type deviceConfig struct {
	serial  string
	address string
	timeout time.Duration
}

// WithSerial sets device serial number
func WithSerial(serial string) DeviceOption {
	return func(c *deviceConfig) {
		c.serial = serial
	}
}

// WithAddress sets device address (IP:port for WiFi)
func WithAddress(addr string) DeviceOption {
	return func(c *deviceConfig) {
		c.address = addr
	}
}

// WithTimeout sets connection timeout
func WithTimeout(d time.Duration) DeviceOption {
	return func(c *deviceConfig) {
		c.timeout = d
	}
}
