package command

import "errors"

var (
	ErrUnknownCommand     = errors.New("unknown command")
	ErrInvalidCommand     = errors.New("invalid command")
	ErrNoHandlerFound     = errors.New("no handler found for command")
	ErrUnsupportedCommand = errors.New("unsupported command")
)
