package core

// Device represents a connected mobile device
type Device interface {
	// Platform returns the device platform
	Platform() Platform

	// Info returns device information
	Info() (map[string]interface{}, error)

	// Screenshot captures current screen
	Screenshot() ([]byte, error)

	// Tap performs tap at coordinates
	Tap(x, y int) error

	// SendKeys inputs text
	SendKeys(text string) error

	// Launch launches the app
	Launch() error

	// Close releases device resources
	Close() error
}

// Platform represents target platform
type Platform string

const (
	IOS     Platform = "ios"
	Android Platform = "android"
)

// NewDevice creates a new device connection
func NewDevice(platform Platform, opts ...DeviceOption) (Device, error) {
	// TODO: implement
	return nil, ErrNotImplemented
}
