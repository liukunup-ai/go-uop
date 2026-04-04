package command

import (
	"context"
	"testing"
)

func TestCommandRouter_New(t *testing.T) {
	reg := NewCommandRegistry()
	router := NewCommandRouter(reg)
	if router == nil {
		t.Fatal("NewCommandRouter returned nil")
	}
	if router.registry != reg {
		t.Error("router.registry should be equal to reg")
	}
}

func TestCommandRouter_Route(t *testing.T) {
	reg := NewCommandRegistry()
	cmd := &BaseCommand{}
	reg.RegisterCommand(cmd)
	router := NewCommandRouter(reg)

	got, err := router.Route("BaseCommand")
	if err != nil {
		t.Errorf("Route() error = %v", err)
	}
	if got == nil {
		t.Fatal("Route returned nil")
	}
}

func TestCommandRouter_Route_Unknown(t *testing.T) {
	reg := NewCommandRegistry()
	router := NewCommandRouter(reg)

	got, err := router.Route("NonExistent")
	if err != ErrUnknownCommand {
		t.Errorf("Route() = %v, want ErrUnknownCommand", err)
	}
	if got != nil {
		t.Error("Route should return nil for unknown command")
	}
}

func TestCommandRouter_Execute(t *testing.T) {
	reg := NewCommandRegistry()
	cmd := &BaseCommand{}
	reg.RegisterCommand(cmd)
	router := NewCommandRouter(reg)

	err := router.Execute(context.Background(), "BaseCommand")
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}
}
