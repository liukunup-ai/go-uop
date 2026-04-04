package runner

import (
	"fmt"
	"sync"
)

type DebugState int

const (
	DebugRunning DebugState = iota
	DebugPaused
	DebugStepping
	DebugStopped
)

type Debugger struct {
	state       DebugState
	mu          sync.RWMutex
	currentStep int
	steps       []Step
	breakpoints map[int]bool
	onBreak     func(step int, cmd string)
}

func NewDebugger(steps []Step) *Debugger {
	return &Debugger{
		state:       DebugRunning,
		steps:       steps,
		breakpoints: make(map[int]bool),
	}
}

func (d *Debugger) State() DebugState {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.state
}

func (d *Debugger) Pause() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.state == DebugRunning {
		d.state = DebugPaused
	}
}

func (d *Debugger) Resume() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.state == DebugPaused || d.state == DebugStepping {
		d.state = DebugRunning
	}
}

func (d *Debugger) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.state = DebugStopped
}

func (d *Debugger) Step() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.state = DebugStepping
	d.currentStep++
}

func (d *Debugger) Skip() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.currentStep++
}

func (d *Debugger) Retry() {
	d.mu.Lock()
	defer d.mu.Unlock()
}

func (d *Debugger) CurrentStep() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.currentStep
}

func (d *Debugger) SetBreakpoint(step int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.breakpoints[step] = true
}

func (d *Debugger) ClearBreakpoint(step int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.breakpoints, step)
}

func (d *Debugger) HasBreakpoint(step int) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.breakpoints[step]
}

func (d *Debugger) WaitForResume() error {
	d.mu.Lock()
	for d.state == DebugPaused {
		d.mu.Unlock()
		return fmt.Errorf("debugger paused at step %d", d.currentStep)
	}
	d.mu.Unlock()
	return nil
}

func (d *Debugger) ExecuteWithDebug(flow *Flow, exec func(int, Step) error) error {
	for i := range flow.Steps {
		d.mu.Lock()
		d.currentStep = i

		for d.state == DebugPaused {
			d.mu.Unlock()
			return fmt.Errorf("debugger paused at step %d", i)
		}

		if d.state == DebugStopped {
			d.mu.Unlock()
			return fmt.Errorf("debugger stopped at step %d", i)
		}
		d.mu.Unlock()

		if d.HasBreakpoint(i) && d.onBreak != nil {
			d.onBreak(i, fmt.Sprintf("%v", flow.Steps[i]))
		}

		if err := exec(i, flow.Steps[i]); err != nil {
			return fmt.Errorf("step %d failed: %w", i, err)
		}
	}
	return nil
}
