package watcher

import (
	"context"
	"fmt"
	"sync"

	"github.com/liukunup/go-uop/core"
)

// Action defines what to do when a popup is matched
type Action interface {
	Execute(ctx context.Context, device core.Device) error
}

// InlineCmd executes a command directly
type InlineCmd struct {
	Name string
	Args map[string]any
}

// CommandExecutor is the function signature for executing commands
var CommandExecutor func(name string, args map[string]any, device core.Device) error

// NewInlineCommand creates a new InlineCmd action
func NewInlineCommand(name string, args map[string]any) *InlineCmd {
	return &InlineCmd{Name: name, Args: args}
}

// Execute runs the inline command
func (a *InlineCmd) Execute(ctx context.Context, device core.Device) error {
	if CommandExecutor == nil {
		return nil
	}
	return CommandExecutor(a.Name, a.Args, device)
}

// RefFlow references an existing flow by name
type RefFlow struct {
	FlowName string
}

// NewReferenceFlow creates a new RefFlow action
func NewReferenceFlow(flowName string) *RefFlow {
	return &RefFlow{FlowName: flowName}
}

// Execute runs the referenced flow
func (a *RefFlow) Execute(ctx context.Context, device core.Device) error {
	// TODO: Look up flow from registry and execute
	return nil
}

// ActionSequence executes multiple actions in order
type ActionSequence struct {
	Actions []Action
	Retry   int
}

// ActionSequenceWithRetry creates a new ActionSequence with retry
func ActionSequenceWithRetry(actions []Action, retry int) *ActionSequence {
	return &ActionSequence{Actions: actions, Retry: retry}
}

// Execute runs all actions in sequence with optional retry
func (s *ActionSequence) Execute(ctx context.Context, device core.Device) error {
	var lastErr error
	for attempt := 0; attempt <= s.Retry; attempt++ {
		for _, action := range s.Actions {
			if err := action.Execute(ctx, device); err != nil {
				lastErr = err
				break
			}
		}
		if lastErr == nil {
			return nil
		}
		if attempt < s.Retry {
			lastErr = fmt.Errorf("retry %d: %w", attempt+1, lastErr)
		}
	}
	return lastErr
}

// actionRegistry holds registered action executors (for testing)
var actionRegistry sync.Map

// RegisterActionExecutor registers a function to execute named actions
func RegisterActionExecutor(name string, fn func(args map[string]any, device core.Device) error) {
	actionRegistry.Store(name, fn)
}

// GetActionExecutor retrieves a registered action executor
func GetActionExecutor(name string) (func(args map[string]any, device core.Device) error, bool) {
	val, ok := actionRegistry.Load(name)
	if !ok {
		return nil, false
	}
	return val.(func(args map[string]any, device core.Device) error), true
}
