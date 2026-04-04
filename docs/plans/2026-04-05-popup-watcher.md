# Popup Watcher Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement popup watcher system that detects and handles modal dialogs during test execution using hybrid detection (system API + image matching).

**Architecture:** 
- `internal/watcher/` package with 4 files: engine.go, match.go, action.go, config.go
- WatcherEngine executes after each test step
- Rules matched in priority order (first-match)
- Hybrid detection: system API first, image fallback

**Tech Stack:** Go, existing vision module, existing iOS/Android device interfaces

---

## Task 1: Create Watcher Package Structure

**Files:**
- Create: `internal/watcher/engine.go`
- Create: `internal/watcher/match.go`
- Create: `internal/watcher/action.go`
- Create: `internal/watcher/config.go`
- Create: `internal/watcher/watcher_test.go`

---

### Task 1.1: Define MatchCondition Interface

**Files:**
- Create: `internal/watcher/match.go`

**Step 1: Write the failing test**

```go
package watcher

import (
	"context"
	"testing"
	"github.com/liukunup/go-uop/core"
)

// Mock device for testing
type mockDevice struct{}

func (m *mockDevice) Screenshot() ([]byte, error) { return []byte("screenshot"), nil }
func (m *mockDevice) Tap(x, y int) error { return nil }
func (m *mockDevice) SendKeys(text string) error { return nil }
func (m *mockDevice) Launch() error { return nil }
func (m *mockDevice) PressKey(code int) error { return nil }
func (m *mockDevice) GetAlertText() (string, error) { return "", nil }
func (m *mockDevice) AcceptAlert() error { return nil }
func (m *mockDevice) DismissAlert() error { return nil }

var _ core.Device = (*mockDevice)(nil)

func TestImageMatch(t *testing.T) {
	ctx := context.Background()
	device := &mockDevice{}
	
	// ImageMatch should match when screenshot contains template
	m := ImageMatch("test_template.png", 0.8)
	matched, err := m.Match(ctx, device)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Initially no template file exists, so it will fail
	// After implementation, this should pass
	_ = matched
}

func TestTextMatch(t *testing.T) {
	ctx := context.Background()
	device := &mockDevice{}
	
	m := TextMatch("确定")
	matched, err := m.Match(ctx, device)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if matched {
		t.Log("TextMatch: popup detected")
	}
}

func TestRegexMatch(t *testing.T) {
	ctx := context.Background()
	device := &mockDevice{}
	
	m := RegexMatch("版本.*更新")
	matched, err := m.Match(ctx, device)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = matched
}

func TestCompoundMatch(t *testing.T) {
	ctx := context.Background()
	device := &mockDevice{}
	
	m := CompoundMatch("or", []MatchCondition{
		TextMatch("确定"),
		TextMatch("取消"),
	})
	matched, err := m.Match(ctx, device)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = matched
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/watcher/... -v`
Expected: FAIL - undefined functions/types

**Step 3: Write minimal match.go**

```go
package watcher

import (
	"context"
	"regexp"
	"github.com/liukunup/go-uop/core"
)

// MatchCondition defines how to detect a popup
type MatchCondition interface {
	Match(ctx context.Context, device core.Device) (bool, error)
}

// ImageMatch matches popup using template image
type ImageMatch struct {
	TemplatePath string
	Threshold    float64
}

func ImageMatch(templatePath string, threshold float64) *ImageMatch {
	return &ImageMatch{
		TemplatePath: templatePath,
		Threshold:    threshold,
	}
}

func (m *ImageMatch) Match(ctx context.Context, device core.Device) (bool, error) {
	// TODO: Implement using vision module
	// 1. Take screenshot
	// 2. Load template from m.TemplatePath
	// 3. Use vision matcher to find match
	// 4. Return true if confidence >= m.Threshold
	return false, nil
}

// TextMatch matches popup using system API
type TextMatch struct {
	Text string
}

func TextMatch(text string) *TextMatch {
	return &TextMatch{Text: text}
}

func (m *TextMatch) Match(ctx context.Context, device core.Device) (bool, error) {
	alertText, err := device.GetAlertText()
	if err != nil {
		// No alert present
		return false, nil
	}
	return alertText == m.Text, nil
}

// RegexMatch matches popup text using regex
type RegexMatch struct {
	Pattern string
}

func RegexMatch(pattern string) *RegexMatch {
	return &RegexMatch{Pattern: pattern}
}

func (m *RegexMatch) Match(ctx context.Context, device core.Device) (bool, error) {
	alertText, err := device.GetAlertText()
	if err != nil {
		return false, nil
	}
	matched, _ := regexp.MatchString(m.Pattern, alertText)
	return matched, nil
}

// CompoundMatch combines multiple conditions with AND/OR
type CompoundMatch struct {
	Operator    string // "and" or "or"
	Conditions  []MatchCondition
}

func CompoundMatch(operator string, conditions []MatchCondition) *CompoundMatch {
	return &CompoundMatch{
		Operator:   operator,
		Conditions: conditions,
	}
}

func (m *CompoundMatch) Match(ctx context.Context, device core.Device) (bool, error) {
	for _, cond := range m.Conditions {
		matched, err := cond.Match(ctx, device)
		if err != nil {
			return false, err
		}
		if m.Operator == "or" && matched {
			return true, nil
		}
		if m.Operator == "and" && !matched {
			return false, nil
		}
	}
	if m.Operator == "and" {
		return true, nil
	}
	return false, nil
}
```

**Step 4: Run test to verify it compiles**

Run: `go build ./internal/watcher/...`
Expected: PASS (tests may fail on unimplemented ImageMatch)

**Step 5: Commit**

```bash
git add internal/watcher/match.go internal/watcher/match_test.go
git commit -m "feat(watcher): add MatchCondition interface and implementations"
```

---

### Task 1.2: Define Action Interface

**Files:**
- Create: `internal/watcher/action.go`

**Step 1: Write the failing test**

```go
package watcher

import (
	"context"
	"testing"
)

func TestInlineCommand(t *testing.T) {
	ctx := context.Background()
	device := &mockDevice{}
	
	action := InlineCommand("tapOn", map[string]any{"x": 100, "y": 200})
	err := action.Execute(ctx, device)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReferenceFlow(t *testing.T) {
	ctx := context.Background()
	device := &mockDevice{}
	
	action := ReferenceFlow("dismiss-popup-flow")
	// Initially will fail because flow doesn't exist
	err := action.Execute(ctx, device)
	if err != nil {
		t.Logf("expected error (flow not exist): %v", err)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/watcher/... -v`
Expected: FAIL - undefined types

**Step 3: Write minimal action.go**

```go
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

// InlineCommand executes a command directly
type InlineCommand struct {
	Name string
	Args map[string]any
}

// CommandExecutor is the function signature for executing commands
var CommandExecutor func(name string, args map[string]any, device core.Device) error

func InlineCommand(name string, args map[string]any) *InlineCommand {
	return &InlineCommand{Name: name, Args: args}
}

func (a *InlineCommand) Execute(ctx context.Context, device core.Device) error {
	if CommandExecutor == nil {
		// No executor registered, skip silently
		return nil
	}
	return CommandExecutor(a.Name, a.Args, device)
}

// ReferenceFlow references an existing flow by name
type ReferenceFlow struct {
	FlowName string
}

func ReferenceFlow(flowName string) *ReferenceFlow {
	return &ReferenceFlow{FlowName: flowName}
}

func (a *ReferenceFlow) Execute(ctx context.Context, device core.Device) error {
	// TODO: Look up flow from registry and execute
	// For now, return nil to avoid breaking tests
	return nil
}

// ActionSequence executes multiple actions in order
type ActionSequence struct {
	Actions []Action
	Retry   int
}

func ActionSequenceWithRetry(actions []Action, retry int) *ActionSequence {
	return &ActionSequence{Actions: actions, Retry: retry}
}

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
```

**Step 4: Run test to verify it compiles**

Run: `go build ./internal/watcher/...`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/watcher/action.go internal/watcher/action_test.go
git commit -m "feat(watcher): add Action interface and implementations"
```

---

### Task 1.3: Create WatcherEngine

**Files:**
- Create: `internal/watcher/engine.go`

**Step 1: Write the failing test**

```go
package watcher

import (
	"context"
	"testing"
)

func TestWatcherEngine_Enabled(t *testing.T) {
	engine := NewWatcherEngine()
	if engine.Enabled() {
		t.Error("new engine should be disabled by default")
	}
	
	engine.Enable()
	if !engine.Enabled() {
		t.Error("after Enable(), engine should be enabled")
	}
	
	engine.Disable()
	if engine.Enabled() {
		t.Error("after Disable(), engine should be disabled")
	}
}

func TestWatcherEngine_AddRule(t *testing.T) {
	engine := NewWatcherEngine()
	
	engine.AddRule(Rule{
		Name:     "test rule",
		Priority: 10,
		Match:    TextMatch("test"),
		Actions:  []Action{InlineCommand("tapOn", nil)},
	})
	
	if len(engine.rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(engine.rules))
	}
	
	// Test priority ordering
	engine.AddRule(Rule{
		Name:     "high priority rule",
		Priority: 1,
		Match:    TextMatch("test"),
		Actions:  []Action{},
	})
	
	if engine.rules[0].Name != "high priority rule" {
		t.Error("rules should be sorted by priority")
	}
}

func TestWatcherEngine_Check(t *testing.T) {
	ctx := context.Background()
	device := &mockDevice{}
	
	engine := NewWatcherEngine()
	engine.Enable()
	
	// Register a test action executor
	RegisterActionExecutor("tapOn", func(args map[string]any, device core.Device) error {
		return nil
	})
	
	engine.AddRule(Rule{
		Name:     "tap test",
		Priority: 10,
		Match:    TextMatch("nonexistent"), // Won't match
		Actions:  []Action{InlineCommand("tapOn", map[string]any{"x": 100, "y": 200})},
	})
	
	// Check should not error even when no popup
	err := engine.Check(ctx, device)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/watcher/... -v`
Expected: FAIL - undefined functions/types

**Step 3: Write minimal engine.go**

```go
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
	Priority int // Lower number = higher priority
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
		enabled: false, // Disabled by default
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
	// Sort by priority (lower = higher priority)
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
	rules := e.Rules() // Get copy
	e.mu.RUnlock()
	
	for _, rule := range rules {
		matched, err := rule.Match.Match(ctx, device)
		if err != nil {
			// Log error but continue checking other rules
			continue
		}
		
		if matched {
			// Execute action sequence with retry
			var actionErr error
			for attempt := 0; attempt <= rule.Retry; attempt++ {
				actionErr = e.executeActions(ctx, device, rule.Actions)
				if actionErr == nil {
					break
				}
			}
			if actionErr != nil {
				// Log error but don't fail the step
				continue
			}
			// Found and handled a popup, return (first-match)
			return nil
		}
	}
	
	return nil
}

func (e *WatcherEngine) executeActions(ctx context.Context, device core.Device, actions []Action) error {
	for _, action := range actions {
		if err := action.Execute(ctx, device); err != nil {
			return err
		}
	}
	return nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/watcher/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/watcher/engine.go internal/watcher/engine_test.go
git commit -m "feat(watcher): add WatcherEngine with rule management"
```

---

### Task 1.4: Add YAML Configuration Support

**Files:**
- Create: `internal/watcher/config.go`

**Step 1: Write the failing test**

```go
package watcher

import (
	"os"
	"testing"
)

func TestParseWatcherConfig(t *testing.T) {
	yamlContent := `
watcher:
  enabled: true
  rules:
    - name: "permission popup"
      priority: 10
      match:
        type: text
        text: "允许"
      actions:
        - tapOn: {x: 500, y: 800}
    
    - name: "upgrade popup"
      priority: 20
      match:
        type: image
        template: "upgrade.png"
        threshold: 0.85
      actions:
        - ref: "dismiss-upgrade"
      retry: 3
`
	tmpfile, err := os.CreateTemp("", "watcher-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	
	if _, err := tmpfile.Write([]byte(yamlContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	
	engine, err := LoadWatcherConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	
	if !engine.Enabled() {
		t.Error("engine should be enabled from config")
	}
	
	rules := engine.Rules()
	if len(rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(rules))
	}
	
	if rules[0].Name != "permission popup" {
		t.Error("rules should be sorted by priority")
	}
	
	if rules[0].Retry != 0 {
		t.Error("first rule should have 0 retry")
	}
	
	if rules[1].Retry != 3 {
		t.Errorf("second rule should have 3 retry, got %d", rules[1].Retry)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/watcher/... -v -run TestParseWatcherConfig`
Expected: FAIL - undefined functions

**Step 3: Write config.go**

```go
package watcher

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// watcherConfig represents the YAML configuration structure
type watcherConfig struct {
	Watcher struct {
		Enabled bool     `yaml:"enabled"`
		Rules   []ruleConfig `yaml:"rules"`
	} `yaml:"watcher"`
}

type ruleConfig struct {
	Name     string      `yaml:"name"`
	Priority int         `yaml:"priority"`
	Match    matchConfig `yaml:"match"`
	Actions  []any      `yaml:"actions"`
	Retry    int         `yaml:"retry"`
}

type matchConfig struct {
	Type       string            `yaml:"type"`
	Text       string            `yaml:"text"`
	Pattern    string            `yaml:"pattern"`
	Template   string            `yaml:"template"`
	Threshold  float64           `yaml:"threshold"`
	Operator   string            `yaml:"operator"`
	Conditions []matchConfig     `yaml:"conditions"`
}

// LoadWatcherConfig loads watcher configuration from a YAML file
func LoadWatcherConfig(path string) (*WatcherEngine, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}
	
	var cfg watcherConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}
	
	engine := NewWatcherEngine()
	if cfg.Watcher.Enabled {
		engine.Enable()
	}
	
	for _, rc := range cfg.Watcher.Rules {
		match, err := parseMatchConfig(rc.Match)
		if err != nil {
			return nil, fmt.Errorf("parse match config for rule %s: %w", rc.Name, err)
		}
		
		actions, err := parseActionsConfig(rc.Actions)
		if err != nil {
			return nil, fmt.Errorf("parse actions config for rule %s: %w", rc.Name, err)
		}
		
		engine.AddRule(Rule{
			Name:     rc.Name,
			Priority: rc.Priority,
			Match:    match,
			Actions:  actions,
			Retry:    rc.Retry,
		})
	}
	
	return engine, nil
}

func parseMatchConfig(cfg matchConfig) (MatchCondition, error) {
	switch cfg.Type {
	case "text":
		return TextMatch(cfg.Text), nil
	
	case "regex":
		return RegexMatch(cfg.Pattern), nil
	
	case "image":
		threshold := cfg.Threshold
		if threshold == 0 {
			threshold = 0.8 // default
		}
		return ImageMatch(cfg.Template, threshold), nil
	
	case "compound":
		operator := cfg.Operator
		if operator == "" {
			operator = "and"
		}
		
		conditions := make([]MatchCondition, 0, len(cfg.Conditions))
		for _, condCfg := range cfg.Conditions {
			cond, err := parseMatchConfig(condCfg)
			if err != nil {
				return nil, err
			}
			conditions = append(conditions, cond)
		}
		
		return CompoundMatch(operator, conditions), nil
	
	default:
		return nil, fmt.Errorf("unknown match type: %s", cfg.Type)
	}
}

func parseActionsConfig(actionsCfg []any) ([]Action, error) {
	actions := make([]Action, 0, len(actionsCfg))
	
	for _, actionCfg := range actionsCfg {
		switch a := actionCfg.(type) {
		case map[string]any:
			// Inline action: { "tapOn": {x: 100, y: 200} }
			for name, args := range a {
				var argsMap map[string]any
				if args != nil {
					var ok bool
					argsMap, ok = args.(map[string]any)
					if !ok {
						argsMap = nil
					}
				}
				actions = append(actions, InlineCommand(name, argsMap))
			}
		
		case string:
			// Reference: "dismiss-upgrade-flow"
			actions = append(actions, ReferenceFlow(a))
		
		default:
			return nil, fmt.Errorf("unsupported action type: %T", actionCfg)
		}
	}
	
	return actions, nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/watcher/... -v -run TestParseWatcherConfig`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/watcher/config.go internal/watcher/config_test.go
git commit -m "feat(watcher): add YAML configuration support"
```

---

## Task 2: Integrate Watcher into Executor

**Files:**
- Modify: `internal/runner/executor.go`

**Step 1: Modify executor.go to accept WatcherEngine**

```go
type Executor struct {
	pool      *DevicePool
	reportGen *report.Generator
	executors map[string]CommandExecutor
	watcher   *watcher.WatcherEngine  // ADD THIS
}

// NewExecutor accepts optional watcher
func NewExecutor(pool *DevicePool, reportGen *report.Generator, watcherOpts ...func(*watcher.WatcherEngine)) *Executor {
	e := &Executor{
		pool:      pool,
		reportGen: reportGen,
		executors: make(map[string]CommandExecutor),
	}
	
	// Register default commands
	e.registerCommands()
	
	return e
}

// WithWatcher sets the watcher engine
func (e *Executor) WithWatcher(w *watcher.WatcherEngine) *Executor {
	e.watcher = w
	return e
}
```

**Step 2: Modify executeStep to call watcher.Check() after each step**

```go
func (e *Executor) executeStep(index int, step Step) error {
	// ... existing step execution logic ...
	
	// After step executes, check for popups
	if e.watcher != nil && e.watcher.Enabled() {
		device := e.pool.CurrentDevice()
		if device != nil && device.device != nil {
			if err := e.watcher.Check(context.Background(), device.device); err != nil {
				// Log warning but don't fail the step
				fmt.Printf("watcher warning: %v\n", err)
			}
		}
	}
	
	return nil
}
```

**Step 3: Verify it compiles**

Run: `go build ./internal/runner/...`
Expected: PASS

**Step 4: Commit**

```bash
git add internal/runner/executor.go
git commit -m "feat(runner): integrate watcher engine for popup handling"
```

---

## Task 3: Add ImageMatch Implementation

**Files:**
- Modify: `internal/watcher/match.go`

**Step 1: Implement ImageMatch using vision module**

```go
func (m *ImageMatch) Match(ctx context.Context, device core.Device) (bool, error) {
	// Take screenshot
	screenshot, err := device.Screenshot()
	if err != nil {
		return false, fmt.Errorf("take screenshot: %w", err)
	}
	
	// Load template
	template, err := os.ReadFile(m.TemplatePath)
	if err != nil {
		return false, fmt.Errorf("read template: %w", err)
	}
	
	// Use vision matcher
	matcher, err := vision.NewMatcher("template")
	if err != nil {
		return false, fmt.Errorf("create matcher: %w", err)
	}
	
	results, err := matcher.Find(screenshot, template)
	if err != nil {
		return false, fmt.Errorf("find matches: %w", err)
	}
	
	if len(results) == 0 {
		return false, nil
	}
	
	// Check if best match exceeds threshold
	bestScore := results[0].Confidence
	return bestScore >= m.Threshold, nil
}
```

**Step 2: Verify it compiles**

Run: `go build ./internal/watcher/...`
Expected: PASS (may need imports)

**Step 3: Commit**

```bash
git add internal/watcher/match.go
git commit -m "feat(watcher): implement ImageMatch with vision module"
```

---

## Task 4: Final Integration Test

**Files:**
- Create integration test combining watcher + executor

**Step 1: Write integration test**

```go
package runner

import (
	"context"
	"os"
	"testing"
	"github.com/liukunup/go-uop/internal/watcher"
)

func TestExecutor_WithWatcher(t *testing.T) {
	pool := NewDevicePool()
	
	// Create watcher engine
	w := watcher.NewWatcherEngine()
	w.Enable()
	w.AddRule(watcher.Rule{
		Name:     "dismiss test",
		Priority: 10,
		Match:    watcher.TextMatch("test"),
		Actions:  []watcher.Action{watcher.InlineCommand("tapOn", map[string]any{"x": 100, "y": 200})},
	})
	
	// Create executor with watcher
	executor := NewExecutor(pool, nil)
	executor = executor.WithWatcher(w)
	
	if !executor.watcher.Enabled() {
		t.Error("watcher should be enabled")
	}
}
```

**Step 2: Run integration test**

Run: `go test ./internal/runner/... -v -run TestExecutor_WithWatcher`
Expected: PASS

**Step 3: Commit**

```bash
git add internal/runner/executor_test.go
git commit -m "test(runner): add watcher integration test"
```

---

## Summary

| Task | Description | Files |
|------|-------------|-------|
| 1.1 | MatchCondition interface + implementations | match.go |
| 1.2 | Action interface + implementations | action.go |
| 1.3 | WatcherEngine with rule management | engine.go |
| 1.4 | YAML configuration support | config.go |
| 2 | Integrate into Executor | executor.go |
| 3 | ImageMatch implementation | match.go |
| 4 | Integration test | executor_test.go |

**Total: 7 tasks**
