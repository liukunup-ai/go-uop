package uiautomator2

import "time"

// Config holds connection configuration
type Config struct {
	Serial  string
	Address string // IP:port for WiFi connection
	Timeout int    // HTTP timeout in seconds
	Package string // Target app package
}

// Option configures device creation
type Option func(*Config)

// WithSerial sets device serial number
func WithSerial(serial string) Option {
	return func(c *Config) {
		c.Serial = serial
	}
}

// WithAddress sets device address (IP:port for WiFi)
func WithAddress(addr string) Option {
	return func(c *Config) {
		c.Address = addr
	}
}

// WithTimeout sets HTTP timeout in seconds
func WithTimeout(seconds int) Option {
	return func(c *Config) {
		c.Timeout = seconds
	}
}

// WithPackage sets target app package
func WithPackage(pkg string) Option {
	return func(c *Config) {
		c.Package = pkg
	}
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Timeout: 60,
	}
}

// HTTPTimeout returns timeout as duration
func (c *Config) HTTPTimeout() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}
