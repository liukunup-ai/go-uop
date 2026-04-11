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
	state         DebugState
	mu            sync.RWMutex
	currentTC     int
	currentStep   int
	testCases     []TestCase
	breakpoints   map[string]bool
	onBreak       func(tc, step int, stepType string)
}

func NewDebugger(testCases []TestCase) *Debugger {
	return &Debugger{
		state:       DebugRunning,
		testCases:   testCases,
		breakpoints: make(map[string]bool),
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

func (d *Debugger) CurrentPosition() (tc, step int) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.currentTC, d.currentStep
}

func (d *Debugger) SetBreakpoint(tc, step int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.breakpoints[fmt.Sprintf("%d:%d", tc, step)] = true
}

func (d *Debugger) ClearBreakpoint(tc, step int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.breakpoints, fmt.Sprintf("%d:%d", tc, step))
}

func (d *Debugger) HasBreakpoint(tc, step int) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.breakpoints[fmt.Sprintf("%d:%d", tc, step)]
}

func (d *Debugger) WaitForResume() error {
	d.mu.Lock()
	for d.state == DebugPaused {
		d.mu.Unlock()
		return fmt.Errorf("debugger paused at testcase %d, step %d", d.currentTC, d.currentStep)
	}
	d.mu.Unlock()
	return nil
}

func (d *Debugger) ExecuteWithDebug(suite *TestSuite, exec func(tc, step int, s Step) error) error {
	for i := range suite.TestCases {
		d.mu.Lock()
		d.currentTC = i

		for d.state == DebugPaused {
			d.mu.Unlock()
			return fmt.Errorf("debugger paused at testcase %d", i)
		}

		if d.state == DebugStopped {
			d.mu.Unlock()
			return fmt.Errorf("debugger stopped at testcase %d", i)
		}
		d.mu.Unlock()

		tc := suite.TestCases[i]
		for j := range tc.Steps {
			if d.HasBreakpoint(i, j) && d.onBreak != nil {
				d.onBreak(i, j, tc.Steps[j].Type)
			}

			if err := exec(i, j, tc.Steps[j]); err != nil {
				return fmt.Errorf("testcase %d, step %d failed: %w", i, j, err)
			}
		}
	}
	return nil
}
