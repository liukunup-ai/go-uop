package command

import "context"

type Command interface {
	Execute(ctx context.Context) error
	Validate() error
	Name() string
}

type UndoableCommand interface {
	Command
	Undo(ctx context.Context) error
}

type BaseCommand struct{}

func (c *BaseCommand) Execute(ctx context.Context) error {
	return nil
}

func (c *BaseCommand) Validate() error {
	return nil
}

func (c *BaseCommand) Name() string {
	return "BaseCommand"
}

func (c *BaseCommand) Undo(ctx context.Context) error {
	return nil
}
