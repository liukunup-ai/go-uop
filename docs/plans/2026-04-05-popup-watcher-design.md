# Popup Watcher Design

## Overview

A popup/dialog watcher system that detects and handles modal dialogs during test execution. Uses hybrid detection (system API + image matching) and supports configurable rules with action sequences.

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│  Test Flow Executor (internal/runner/executor.go)       │
│  - Calls WatcherEngine.Check(device) after each step   │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│  Watcher Engine (internal/watcher/engine.go)           │
│  - Manages rule list                                   │
│  - Matches in priority order (first-match)             │
└─────────────────────────────────────────────────────────┘
           │                           │
           ▼                           ▼
┌──────────────────────┐   ┌────────────────────────────┐
│ Match Conditions     │   │ Actions                    │
│ - ImageMatch         │   │ - InlineCommand            │
│ - TextMatch          │   │ - ReferencedFlow           │
│ - RegexMatch         │   │ - ActionSequence (retry)    │
│ - CompoundMatch      │   └────────────────────────────┘
└──────────────────────┘
```

## Package Structure

```
internal/watcher/
├── engine.go        # WatcherEngine, Rule, main logic
├── match.go         # MatchCondition interface + implementations
├── action.go        # Action interface + implementations
├── config.go        # YAML config parsing
└── watcher_test.go
```

## Data Structures

### Rule
```go
type Rule struct {
    Name     string           // Rule name for logging
    Priority int              // Lower number = higher priority
    Match    MatchCondition   // Match condition
    Actions  []Action         // Actions to execute
    Retry    int              // Retry count on failure
}
```

### MatchCondition Interface
```go
type MatchCondition interface {
    Match(ctx context.Context, device core.Device) (bool, error)
}
```

Implementations:
- `ImageMatch(templatePath string, threshold float64)` — Image template matching
- `TextMatch(text string)` — Exact text match via system API
- `RegexMatch(pattern string)` — Regex pattern match on alert text
- `CompoundMatch(operator string, conditions []MatchCondition)` — AND/OR combination

### Action Interface
```go
type Action interface {
    Execute(ctx context.Context, device core.Device) error
}
```

Implementations:
- `InlineCommand(name string, args map[string]any)` — Execute registered command inline
- `ReferenceFlow(flowName string)` — Reference an existing flow by name

## Configuration

### YAML Format
```yaml
watcher:
  enabled: true
  rules:
    - name: "permission popup"
      priority: 10
      match:
        type: image
        template: "popup_permission.png"
        threshold: 0.8
      actions:
        - tapOn: {x: 500, y: 800}
        - wait: {ms: 300}
    
    - name: "upgrade popup"
      priority: 20
      match:
        type: text
        text: "发现新版本"
      actions:
        - ref: "dismiss-upgrade-flow"
      retry: 3
    
    - name: "ad popup"
      priority: 5
      match:
        type: compound
        operator: or
        conditions:
          - type: image
            template: "ad_popup_1.png"
          - type: image
            template: "ad_popup_2.png"
      actions:
        - tapOn: {x: 800, y: 100}
        - wait: {ms: 500}
        - tapOn: {text: "跳过"}
```

### Go API Format
```go
watcher := NewWatcherEngine().
    AddRule(Rule{
        Name: "permission popup",
        Match: ImageMatch("popup_permission.png", 0.8),
        Actions: []Action{
            InlineCommand("tapOn", map[string]any{"x": 500, "y": 800}),
            InlineCommand("wait", map[string]any{"ms": 300}),
        },
    })
```

## Detection Flow

1. **Executor** executes each step, then calls `WatcherEngine.Check(device)`
2. **Engine** iterates rules in priority order
3. For each rule:
   - Try **system API first** (iOS: `GetAlertText()`, Android: uiautomator)
   - If API returns empty/no dialog, fallback to **image matching**
4. If match found:
   - Execute action sequence
   - If failed and Retry > 0, retry entire sequence
5. After all rules checked, return to Executor

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Package location | `internal/watcher/` | Alongside runner, clear ownership |
| Rule matching | First-match (priority order) | Simple, intuitive |
| Priority | Lower number = higher priority | Precise control |
| System API fallback | API first, image兜底 | Most dialogs detectable via API |
| Retry scope | Entire action sequence | Most popups need multi-step dismiss |
| Compound operator | `and` / `or` | Flexible condition composition |

## Integration Points

### Executor Integration (internal/runner/executor.go)
```go
type Executor struct {
    pool      *DevicePool
    reportGen *report.Generator
    executors map[string]CommandExecutor
    watcher   *WatcherEngine  // NEW
}

func (e *Executor) executeStep(...) {
    // ... existing logic ...
    
    // After step executes, check for popups
    if e.watcher != nil && e.watcher.Enabled() {
        if err := e.watcher.Check(device); err != nil {
            // Log warning, don't fail step
        }
    }
}
```

### Flow YAML Integration
```yaml
name: my flow
watcher:
  enabled: true
  rules:
    - name: "close ad"
      match:
        type: image
        template: "close_btn.png"
      actions:
        - tapOn: {x: 800, y: 100}
steps:
  - launch: com.example.app
  - tapOn: {text: "开始"}
```

## File Responsibilities

| File | Responsibility |
|------|----------------|
| `engine.go` | WatcherEngine struct, rule management, Check() logic |
| `match.go` | MatchCondition interface, 4 implementations |
| `action.go` | Action interface, InlineCommand, ReferenceFlow |
| `config.go` | YAML parsing, rule loading |

## Acceptance Criteria

- [ ] WatcherEngine can be enabled/disabled
- [ ] Rules support priority ordering
- [ ] ImageMatch works with template matching
- [ ] TextMatch works via system API
- [ ] RegexMatch works on alert text
- [ ] CompoundMatch supports AND/OR
- [ ] InlineCommand executes registered commands
- [ ] ReferenceFlow references existing flows
- [ ] Retry mechanism retries entire action sequence
- [ ] YAML configuration parses correctly
- [ ] Executor calls Watcher.Check() after each step
- [ ] Watcher does not fail the step if handling fails
