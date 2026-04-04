package command

import (
	"context"
	"testing"
)

func TestCommandRegistry(t *testing.T) {
	reg := NewCommandRegistry()

	cmd := &BaseCommand{}
	err := reg.RegisterCommand(cmd)
	if err != nil {
		t.Fatalf("RegisterCommand failed: %v", err)
	}

	got := reg.Get("BaseCommand")
	if got == nil {
		t.Fatal("Get returned nil")
	}

	unknown := reg.Get("NonExistent")
	if unknown != nil {
		t.Fatal("Get should return nil for unknown command")
	}
}

func TestCommandRegistry_RegisterCommand_Nil(t *testing.T) {
	reg := NewCommandRegistry()
	err := reg.RegisterCommand(nil)
	if err != ErrInvalidCommand {
		t.Errorf("RegisterCommand(nil) = %v, want ErrInvalidCommand", err)
	}
}

func TestCommandRegistry_RegisterCommand_InvalidValidate(t *testing.T) {
	reg := NewCommandRegistry()
	cmd := NewTapCommand(-1, -1)
	err := reg.RegisterCommand(cmd)
	if err == nil {
		t.Error("RegisterCommand should return error for invalid command")
	}
}

func TestCommandRegistry_Execute(t *testing.T) {
	reg := NewCommandRegistry()
	cmd := &BaseCommand{}
	reg.RegisterCommand(cmd)

	err := reg.Execute(context.Background(), "BaseCommand")
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}
}

func TestCommandRegistry_Execute_UnknownCommand(t *testing.T) {
	reg := NewCommandRegistry()
	err := reg.Execute(context.Background(), "NonExistent")
	if err != ErrUnknownCommand {
		t.Errorf("Execute() = %v, want ErrUnknownCommand", err)
	}
}

func TestCommandRegistry_Dispatch(t *testing.T) {
	reg := NewCommandRegistry()
	cmd := &BaseCommand{}
	reg.RegisterCommand(cmd)

	handler := &testHandler{canHandle: false}
	reg.RegisterHandler(handler)

	err := reg.Dispatch(context.Background(), cmd)
	if err != ErrNoHandlerFound {
		t.Errorf("Dispatch() = %v, want ErrNoHandlerFound", err)
	}
}

func TestCommandRegistry_Dispatch_WithHandler(t *testing.T) {
	reg := NewCommandRegistry()
	cmd := &BaseCommand{}
	reg.RegisterCommand(cmd)

	handler := &testHandler{canHandle: true}
	reg.RegisterHandler(handler)

	err := reg.Dispatch(context.Background(), cmd)
	if err != nil {
		t.Errorf("Dispatch() error = %v", err)
	}
}

type testHandler struct {
	canHandle bool
}

func (h *testHandler) CanHandle(cmd Command) bool {
	return h.canHandle
}

func (h *testHandler) Handle(ctx context.Context, cmd Command) error {
	return nil
}
