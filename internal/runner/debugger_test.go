package runner

import (
	"testing"
)

func TestDebugger_Creation(t *testing.T) {
	steps := []Step{
		{"tapOn": Step{"x": 100, "y": 200}},
		{"launch": nil},
	}
	debugger := NewDebugger(steps)

	if debugger == nil {
		t.Fatal("NewDebugger returned nil")
	}

	if debugger.State() != DebugRunning {
		t.Errorf("Expected initial state DebugRunning, got %v", debugger.State())
	}

	if len(debugger.steps) != 2 {
		t.Errorf("Expected 2 steps, got %d", len(debugger.steps))
	}
}

func TestDebugger_PauseAndResume(t *testing.T) {
	steps := []Step{{"tapOn": Step{"x": 100, "y": 200}}}
	debugger := NewDebugger(steps)

	debugger.Pause()
	if debugger.State() != DebugPaused {
		t.Errorf("Expected DebugPaused, got %v", debugger.State())
	}

	debugger.Resume()
	if debugger.State() != DebugRunning {
		t.Errorf("Expected DebugRunning, got %v", debugger.State())
	}
}

func TestDebugger_Stop(t *testing.T) {
	steps := []Step{{"tapOn": Step{"x": 100, "y": 200}}}
	debugger := NewDebugger(steps)

	debugger.Stop()
	if debugger.State() != DebugStopped {
		t.Errorf("Expected DebugStopped, got %v", debugger.State())
	}
}

func TestDebugger_Step(t *testing.T) {
	steps := []Step{
		{"tapOn": Step{"x": 100, "y": 200}},
		{"inputText": Step{"text": "hello"}},
		{"launch": nil},
	}
	debugger := NewDebugger(steps)

	initialStep := debugger.CurrentStep()
	if initialStep != 0 {
		t.Errorf("Expected initial step 0, got %d", initialStep)
	}

	debugger.Step()
	if debugger.CurrentStep() != 1 {
		t.Errorf("Expected step 1 after Step(), got %d", debugger.CurrentStep())
	}
}

func TestDebugger_Skip(t *testing.T) {
	steps := []Step{
		{"tapOn": Step{"x": 100, "y": 200}},
		{"inputText": Step{"text": "hello"}},
	}
	debugger := NewDebugger(steps)

	debugger.Skip()
	if debugger.CurrentStep() != 1 {
		t.Errorf("Expected step 1 after Skip(), got %d", debugger.CurrentStep())
	}
}

func TestDebugger_Breakpoints(t *testing.T) {
	steps := []Step{
		{"tapOn": Step{"x": 100, "y": 200}},
		{"inputText": Step{"text": "hello"}},
		{"launch": nil},
	}
	debugger := NewDebugger(steps)

	debugger.SetBreakpoint(1)
	if !debugger.HasBreakpoint(1) {
		t.Error("Expected breakpoint at step 1")
	}

	if debugger.HasBreakpoint(0) {
		t.Error("Did not expect breakpoint at step 0")
	}

	debugger.ClearBreakpoint(1)
	if debugger.HasBreakpoint(1) {
		t.Error("Did not expect breakpoint at step 1 after clearing")
	}
}

func TestDebugger_ExecuteWithDebug(t *testing.T) {
	steps := []Step{
		{"tapOn": Step{"x": 100, "y": 200}},
		{"inputText": Step{"text": "hello"}},
	}
	debugger := NewDebugger(steps)

	executions := 0
	flow := &Flow{Name: "test-flow", Steps: steps}
	err := debugger.ExecuteWithDebug(flow, func(i int, step Step) error {
		executions++
		return nil
	})

	if err != nil {
		t.Fatalf("ExecuteWithDebug failed: %v", err)
	}

	if executions != 2 {
		t.Errorf("Expected 2 executions, got %d", executions)
	}
}

func TestDebugger_ExecuteWithDebugStop(t *testing.T) {
	steps := []Step{
		{"tapOn": Step{"x": 100, "y": 200}},
		{"inputText": Step{"text": "hello"}},
	}
	debugger := NewDebugger(steps)

	debugger.Stop()

	flow := &Flow{Name: "test-flow", Steps: steps}
	err := debugger.ExecuteWithDebug(flow, func(i int, step Step) error {
		return nil
	})

	if err == nil {
		t.Error("Expected error when debugger is stopped")
	}
}

func TestDebugger_StateTransitions(t *testing.T) {
	steps := []Step{{"tapOn": Step{"x": 100, "y": 200}}}
	debugger := NewDebugger(steps)

	if debugger.State() != DebugRunning {
		t.Errorf("Expected DebugRunning, got %v", debugger.State())
	}

	debugger.Pause()
	if debugger.State() != DebugPaused {
		t.Errorf("Expected DebugPaused, got %v", debugger.State())
	}

	debugger.Resume()
	if debugger.State() != DebugRunning {
		t.Errorf("Expected DebugRunning after Resume, got %v", debugger.State())
	}

	debugger.Pause()
	debugger.Step()
	if debugger.State() != DebugStepping {
		t.Errorf("Expected DebugStepping, got %v", debugger.State())
	}

	debugger.Resume()
	if debugger.State() != DebugRunning {
		t.Errorf("Expected DebugRunning, got %v", debugger.State())
	}

	debugger.Stop()
	if debugger.State() != DebugStopped {
		t.Errorf("Expected DebugStopped, got %v", debugger.State())
	}
}
