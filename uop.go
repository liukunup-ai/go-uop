package uop

import "errors"

// Platform represents target platform
type Platform string

const (
	IOS     Platform = "ios"
	Android Platform = "android"
)

// ErrNotImplemented indicates feature not implemented
var ErrNotImplemented = errors.New("not implemented")

// ErrDeviceNotFound indicates device not found
var ErrDeviceNotFound = errors.New("device not found")
