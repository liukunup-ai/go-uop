package maestro

import "errors"

var (
	ErrUnsupportedCommand = errors.New("unsupported maestro command")
	ErrElementNotFound    = errors.New("element not found")
	ErrParsingFailed      = errors.New("maestro flow parsing failed")
)

func IsUnsupportedCommand(err error) bool {
	return errors.Is(err, ErrUnsupportedCommand)
}
