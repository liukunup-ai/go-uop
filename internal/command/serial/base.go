package serial

import (
	"github.com/liukunup/go-uop/pkg/serial"
)

type Command interface {
	Name() string
	Description() string
	Validate() error
}

type SerialCommand interface {
	Command
	SetSerial(s *serial.Serial)
	Serial() *serial.Serial
}

type baseSerialCommand struct {
	serial *serial.Serial
}

func (c *baseSerialCommand) SetSerial(s *serial.Serial) {
	c.serial = s
}

func (c *baseSerialCommand) Serial() *serial.Serial {
	return c.serial
}
