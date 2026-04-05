package watcher

import (
	"context"
	"sort"
	"sync"

	"github.com/liukunup/go-uop/core"
)

// Rule defines a popup detection rule
type Rule struct {
	Name     string
	Priority int
	Match    MatchCondition
	Actions  []Action
	Retry    int
}

// WatcherEngine manages popup detection rules
type WatcherEngine struct {
	mu      sync.RWMutex
	rules   []Rule
	enabled bool
}

// NewWatcherEngine creates a new watcher engine
func NewWatcherEngine() *WatcherEngine {
	return &WatcherEngine{
		rules:   make([]Rule, 0),
		enabled: false,
	}
}

// Enable enables the watcher
func (e *WatcherEngine) Enable() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.enabled = true
}

// Disable disables the watcher
func (e *WatcherEngine) Disable() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.enabled = false
}

// Enabled returns whether the watcher is enabled
func (e *WatcherEngine) Enabled() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.enabled
}

// AddRule adds a rule to the engine
func (e *WatcherEngine) AddRule(rule Rule) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules = append(e.rules, rule)
	sort.Slice(e.rules, func(i, j int) bool {
		return e.rules[i].Priority < e.rules[j].Priority
	})
}

// Rules returns a copy of current rules
func (e *WatcherEngine) Rules() []Rule {
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make([]Rule, len(e.rules))
	copy(result, e.rules)
	return result
}

// Check checks for popups and executes matching rules
func (e *WatcherEngine) Check(ctx context.Context, device core.Device) error {
	e.mu.RLock()
	if !e.enabled {
		e.mu.RUnlock()
		return nil
	}
	rules := e.Rules()
	e.mu.RUnlock()

	for _, rule := range rules {
		matched, err := rule.Match.Match(ctx, device)
		if err != nil {
			continue
		}

		if matched {
			actionSeq := ActionSequenceWithRetry(rule.Actions, rule.Retry)
			if err := actionSeq.Execute(ctx, device); err != nil {
				continue
			}
			return nil
		}
	}

	return nil
}
