package command

import (
	"github.com/liukunup/go-uop/internal/command/device"
	"github.com/liukunup/go-uop/internal/command/serial"
	"github.com/liukunup/go-uop/internal/command/system"
)

type DeviceCommand = device.DeviceCommand
type SerialCommand = serial.SerialCommand

type TapCommand = device.TapCommand
type LaunchCommand = device.LaunchCommand
type SendKeysCommand = device.SendKeysCommand
type PressKeyCommand = device.PressKeyCommand

type SendByIDCommand = serial.SendByIDCommand
type SendRawCommand = serial.SendRawCommand

type WaitCommand = system.WaitCommand
type ScreenshotCommand = device.ScreenshotCommand
type SwipeCommand = system.SwipeCommand

var NewTapCommand = device.NewTapCommand
var NewLaunchCommand = device.NewLaunchCommand
var NewSendKeysCommand = device.NewSendKeysCommand
var NewPressKeyCommand = device.NewPressKeyCommand

var NewSendByIDCommand = serial.NewSendByIDCommand
var NewSendRawCommand = serial.NewSendRawCommand

var NewWaitCommand = system.NewWaitCommand
var NewScreenshotCommand = device.NewScreenshotCommand
var NewSwipeCommand = system.NewSwipeCommand
