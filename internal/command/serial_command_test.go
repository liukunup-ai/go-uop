package command

import (
	"context"
	"testing"
	"time"
)

func TestSendByIDCommand_Validate(t *testing.T) {
	tests := []struct {
		name      string
		commandID string
		wantErr   bool
	}{
		{"valid command ID", "reset", false},
		{"empty command ID", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewSendByIDCommand(tt.commandID, time.Second)
			err := cmd.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendByIDCommand_Name(t *testing.T) {
	cmd := NewSendByIDCommand("reset", time.Second)
	if cmd.Name() != "sendByID" {
		t.Errorf("Name() = %s, want sendByID", cmd.Name())
	}
}

func TestSendByIDCommand_Execute_NoSerial(t *testing.T) {
	cmd := NewSendByIDCommand("reset", time.Second)
	err := cmd.Execute(context.Background())
	if err == nil {
		t.Error("Execute() should return error when serial not set")
	}
}

func TestSendRawCommand_Name(t *testing.T) {
	cmd := NewSendRawCommand("test data")
	if cmd.Name() != "sendRaw" {
		t.Errorf("Name() = %s, want sendRaw", cmd.Name())
	}
}

func TestSendRawCommand_Execute_NoSerial(t *testing.T) {
	cmd := NewSendRawCommand("test data")
	err := cmd.Execute(context.Background())
	if err == nil {
		t.Error("Execute() should return error when serial not set")
	}
}

func TestBaseSerialCommand_SetSerial(t *testing.T) {
	cmd := NewSendByIDCommand("reset", time.Second)
	cmd.SetSerial(nil)
	if cmd.BaseSerialCommand.serial != nil {
		t.Error("serial should be nil after SetSerial(nil)")
	}
}
