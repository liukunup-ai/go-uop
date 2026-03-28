package uop

import "github.com/liukunup/go-uop/core"

type Device = core.Device
type Platform = core.Platform
type DeviceOption = core.DeviceOption

const (
	IOS     = core.IOS
	Android = core.Android
	Serial  = core.Serial
)

var (
	ErrNotImplemented = core.ErrNotImplemented
	ErrDeviceNotFound = core.ErrDeviceNotFound
)

var NewDevice = core.NewDevice
var WithSerial = core.WithSerial
var WithAddress = core.WithAddress
var WithTimeout = core.WithTimeout
