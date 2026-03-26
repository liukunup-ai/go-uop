package console

import "errors"

var (
	ErrDeviceNotFound      = errors.New("device not found")
	ErrDeviceNotConnected  = errors.New("device not connected")
	ErrUnsupportedPlatform = errors.New("unsupported platform")
	ErrUnknownCommand      = errors.New("unknown command")
)
