# OOP Command Design Plan

Date: 2026-04-04
Status: Draft

## Overview

Refactor the command system to use proper OOP patterns: Command Pattern + Strategy Pattern.

## Problem

Current architecture issues:

| Component | Current State | Issue |
|-----------|--------------|-------|
| `pkg/serial/command.go` | Plain struct | No validation, no methods on Command |
| `internal/action/action.go` | Interface + structs | Implementations have no `Do()` methods |
| `internal/console/handler.go` | 30+ handle* methods | Massive switch-case, no polymorphism |
| `internal/runner/executor.go` | `map[string]CommandExecutor` func map | Hard to extend, no state |
| `internal/console/device.go` | `ExecuteCommand()` switch-case | Should use polymorphic Commands |

## Solution: Command Pattern + Strategy Pattern

### Layered Architecture

```
┌─────────────────────────────────────────────────────┐
│  Handler Layer (HTTP/Console)                       │
│  handleDevices, handleSerialConnect, etc.           │
└─────────────────────┬───────────────────────────────┘
                      │ Command
                      ▼
┌─────────────────────────────────────────────────────┐
│  Command Registry                                   │
│  - commands map[string]Command                      │
│  - handlers []Handler                               │
└─────────────────────┬───────────────────────────────┘
                      │ Execute()
                      ▼
┌─────────────────────────────────────────────────────┐
│  Command Implementations                            │
│  ├─ DeviceCommand (Tap, Launch, SendKeys...)       │
│  ├─ SerialCommand (SendByID, SendRaw...)           │
│  └─ SystemCommand (Wait, Screenshot...)            │
└─────────────────────┬───────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────┐
│  Device Abstraction (core.Device)                   │
│  └─ iOS, Android, Serial implementations          │
└─────────────────────────────────────────────────────┘
```

### Core Interfaces

#### Command Interface

```go
type Command interface {
    Execute(ctx context.Context) error  // Execute command
    Validate() error                      // Parameter validation
    Name() string                         // Command name
    Undo(ctx context.Context) error       // Undo (optional)
}
```

#### Handler Interface

```go
type Handler interface {
    Handle(ctx context.Context, cmd Command) error  // Handle command
    CanHandle(cmd Command) bool                       // Can this handler process?
    Priority() int                                   // Handler priority
}
```

### CommandRegistry

```go
type CommandRegistry struct {
    mu       sync.RWMutex
    commands map[string]Command
    handlers []Handler
}

func (r *CommandRegistry) RegisterCommand(cmd Command) error
func (r *CommandRegistry) RegisterHandler(h Handler)
func (r *CommandRegistry) Get(name string) Command
func (r *CommandRegistry) Execute(ctx context.Context, name string, args map[string]any) error
func (r *CommandRegistry) Dispatch(ctx context.Context, cmd Command) error
```

### Command Implementations

#### DeviceCommand

```go
type BaseDeviceCommand struct {
    device core.Device
    args   map[string]any
}

type TapCommand struct {
    BaseDeviceCommand
    X int
    Y int
}

func (c *TapCommand) Name() string { return "tapOn" }
func (c *TapCommand) Execute(ctx context.Context) error { return c.device.Tap(c.X, c.Y) }
func (c *TapCommand) Validate() error {
    if c.X < 0 || c.Y < 0 {
        return errors.New("coordinates must be non-negative")
    }
    return nil
}

type LaunchCommand struct {
    BaseDeviceCommand
    AppID    string
    Args     []string
    WaitIdle bool
}

func (c *LaunchCommand) Name() string { return "launch" }
func (c *LaunchCommand) Execute(ctx context.Context) error { return c.device.Launch() }
```

#### SerialCommand

```go
type SendByIDCommand struct {
    Serial    *serial.Serial
    CommandID string
    Timeout   time.Duration
}

func (c *SendByIDCommand) Name() string { return "sendByID" }
func (c *SendByIDCommand) Execute(ctx context.Context) error {
    return c.Serial.SendByID(c.CommandID, nil)
}
```

#### SystemCommand

```go
type WaitCommand struct {
    Duration time.Duration
}

func (c *WaitCommand) Name() string { return "wait" }
func (c *WaitCommand) Execute(ctx context.Context) error {
    time.Sleep(c.Duration)
    return nil
}
```

### Handler Implementations

```go
type DeviceOpsHandler struct {
    deviceMgr *DeviceManager
}

func (h *DeviceOpsHandler) CanHandle(cmd Command) bool {
    _, ok := cmd.(DeviceCommand)
    return ok
}

func (h *DeviceOpsHandler) Handle(ctx context.Context, cmd Command) error {
    switch c := cmd.(type) {
    case *TapCommand:
        return c.Execute(ctx)
    case *LaunchCommand:
        return c.Execute(ctx)
    default:
        return ErrUnsupportedCommand
    }
}
```

## File Structure Changes

### New Files

```
internal/command/
├── command.go          # Command interface definition
├── registry.go         # CommandRegistry implementation
├── device_command.go   # Device command implementations
├── serial_command.go   # Serial command implementations
├── system_command.go   # System command implementations
└── handler.go          # Handler interface and implementations

internal/console/
├── server.go           # Unchanged
├── handler.go           # Refactored: use CommandRouter
├── device_manager.go   # Renamed: device.go → device_manager.go
├── types.go            # Unchanged
├── ios_manager.go      # Unchanged
├── serial_manager.go   # Unchanged
└── history.go          # Unchanged
```

### Delete/Archive

- `internal/action/action.go` (replaced by command package)
- `internal/runner/executor.go` (replaced by CommandRegistry)

## Implementation Phases

| Phase | Task | Priority |
|-------|------|----------|
| 1 | Create `internal/command/` package, define interfaces | High |
| 2 | Implement `CommandRegistry` | High |
| 3 | Implement DeviceCommand, SerialCommand, SystemCommand | High |
| 4 | Implement Handler interface, refactor console handler | Medium |
| 5 | Delete deprecated `action.go` and `executor.go` | Low |
| 6 | Write tests | High |
