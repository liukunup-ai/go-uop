package command

import (
	"context"
	"testing"
)

func TestCommandInterface(t *testing.T) {
	var _ interface {
		Execute(ctx context.Context) error
		Validate() error
		Name() string
	} = (*BaseCommand)(nil)
}

func TestUndoableCommandInterface(t *testing.T) {
	var _ interface {
		Execute(ctx context.Context) error
		Validate() error
		Name() string
		Undo(ctx context.Context) error
	} = (*BaseCommand)(nil)
}

func TestBaseCommandValidate(t *testing.T) {
	cmd := &BaseCommand{}
	err := cmd.Validate()
	if err != nil {
		t.Errorf("Validate() returned error: %v", err)
	}
}
