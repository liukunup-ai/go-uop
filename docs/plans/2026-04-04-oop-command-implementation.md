# OOP Command Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Refactor command system to use Command Pattern + Strategy Pattern with proper OOP interfaces and implementations.

**Architecture:** Create `internal/command/` package with Command interface, CommandRegistry, and implementations for Device/Serial/System commands. Replace console handler switch-case with Handler interface.

**Tech Stack:** Go, standard library sync, context

---

## Phase 1: Create Command Package and Interfaces

### Task 1: Create internal/command/command.go

**Files:**
- Create: `internal/command/command.go`

**Step 1: Write the failing test**

```go
package command

import (
    "testing"
)

func TestCommandInterface(t *testing.T) {
    // Verify Command interface exists and has required methods
    var _ interface {
        Execute(ctx interface{}) error
        Validate() error
        Name() string
    } = (*BaseCommand)(nil)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/command/... -v`
Expected: FAIL - directory does not exist

**Step 3: Write minimal interface**

```go
package command

import "context"

// Command 命令接口
type Command interface {
    Execute(ctx context.Context) error
    Validate() error
    Name() string
}

// UndoableCommand 可撤销命令接口
type UndoableCommand interface {
    Command
    Undo(ctx context.Context) error
}

// BaseCommand 基础命令实现
type BaseCommand struct{}

func (c *BaseCommand) Validate() error {
    return nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/command/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/command/command.go
git commit -m "feat(command): add Command interface"
```

---

### Task 2: Create internal/command/errors.go

**Files:**
- Create: `internal/command/errors.go`

**Step 1: Write error definitions**

```go
package command

import "errors"

var (
    ErrUnknownCommand  = errors.New("unknown command")
    ErrInvalidCommand = errors.New("invalid command")
    ErrNoHandlerFound = errors.New("no handler found for command")
    ErrUnsupportedCommand = errors.New("unsupported command")
)
```

**Step 2: Commit**

```bash
git add internal/command/errors.go
git commit -m "feat(command): add error definitions"
```

---

## Phase 2: Implement CommandRegistry

### Task 3: Create internal/command/registry.go

**Files:**
- Create: `internal/command/registry.go`

**Step 1: Write the failing test**

```go
func TestCommandRegistry(t *testing.T) {
    reg := NewCommandRegistry()
    
    // Test registration
    cmd := &BaseCommand{}
    err := reg.RegisterCommand(cmd)
    if err != nil {
        t.Fatalf("RegisterCommand failed: %v", err)
    }
    
    // Test Get
    got := reg.Get("BaseCommand")
    if got == nil {
        t.Fatal("Get returned nil")
    }
    
    // Test unknown command
    unknown := reg.Get("NonExistent")
    if unknown != nil {
        t.Fatal("Get should return nil for unknown command")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/command/... -v`
Expected: FAIL - NewCommandRegistry not defined

**Step 3: Write CommandRegistry**

```go
package command

import (
    "context"
    "sync"
)

// Handler 命令处理器接口
type Handler interface {
    Handle(ctx context.Context, cmd Command) error
    CanHandle(cmd Command) bool
}

// CommandRegistry 命令注册表
type CommandRegistry struct {
    mu       sync.RWMutex
    commands map[string]Command
    handlers []Handler
}

// NewCommandRegistry 创建命令注册表
func NewCommandRegistry() *CommandRegistry {
    return &CommandRegistry{
        commands: make(map[string]Command),
    }
}

// RegisterCommand 注册命令
func (r *CommandRegistry) RegisterCommand(cmd Command) error {
    if cmd == nil {
        return ErrInvalidCommand
    }
    if err := cmd.Validate(); err != nil {
        return err
    }
    
    r.mu.Lock()
    defer r.mu.Unlock()
    r.commands[cmd.Name()] = cmd
    return nil
}

// RegisterHandler 注册处理器
func (r *CommandRegistry) RegisterHandler(h Handler) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.handlers = append(r.handlers, h)
}

// Get 获取命令
func (r *CommandRegistry) Get(name string) Command {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return r.commands[name]
}

// Execute 执行命令
func (r *CommandRegistry) Execute(ctx context.Context, name string) error {
    cmd := r.Get(name)
    if cmd == nil {
        return ErrUnknownCommand
    }
    return cmd.Execute(ctx)
}

// Dispatch 分发命令到合适的处理器
func (r *CommandRegistry) Dispatch(ctx context.Context, cmd Command) error {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    for _, h := range r.handlers {
        if h.CanHandle(cmd) {
            return h.Handle(ctx, cmd)
        }
    }
    return ErrNoHandlerFound
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/command/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/command/registry.go
git commit -m "feat(command): add CommandRegistry"
```

---

## Phase 3: Implement DeviceCommand

### Task 4: Create internal/command/device_command.go

**Files:**
- Create: `internal/command/device_command.go`

**Step 1: Write the failing test**

```go
func TestTapCommand(t *testing.T) {
    cmd := NewTapCommand(100, 200)
    
    if cmd.Name() != "tapOn" {
        t.Errorf("Name() = %s, want tapOn", cmd.Name())
    }
    
    if err := cmd.Validate(); err != nil {
        t.Errorf("Validate() error = %v", err)
    }
    
    // Test invalid coordinates
    invalidCmd := NewTapCommand(-1, -1)
    if err := invalidCmd.Validate(); err == nil {
        t.Error("Validate() should return error for negative coordinates")
    }
}

func TestLaunchCommand(t *testing.T) {
    cmd := NewLaunchCommand("com.example.app")
    
    if cmd.Name() != "launch" {
        t.Errorf("Name() = %s, want launch", cmd.Name())
    }
    
    if cmd.AppID != "com.example.app" {
        t.Errorf("AppID = %s, want com.example.app", cmd.AppID)
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/command/... -v -run TestTap`
Expected: FAIL - functions not defined

**Step 3: Write device commands**

```go
package command

import (
    "context"
    "errors"

    "github.com/liukunup/go-uop/core"
)

// DeviceCommand 设备命令接口
type DeviceCommand interface {
    Command
    SetDevice(device core.Device)
}

// BaseDeviceCommand 基础设备命令
type BaseDeviceCommand struct {
    device core.Device
}

func (c *BaseDeviceCommand) SetDevice(device core.Device) {
    c.device = device
}

// TapCommand 点击命令
type TapCommand struct {
    BaseDeviceCommand
    X int
    Y int
}

// NewTapCommand 创建点击命令
func NewTapCommand(x, y int) *TapCommand {
    return &TapCommand{X: x, Y: y}
}

func (c *TapCommand) Name() string { return "tapOn" }

func (c *TapCommand) Validate() error {
    if c.X < 0 || c.Y < 0 {
        return errors.New("coordinates must be non-negative")
    }
    return nil
}

func (c *TapCommand) Execute(ctx context.Context) error {
    if c.device == nil {
        return errors.New("device not set")
    }
    return c.device.Tap(c.X, c.Y)
}

// LaunchCommand 启动命令
type LaunchCommand struct {
    BaseDeviceCommand
    AppID    string
    Args     []string
    WaitIdle bool
}

// NewLaunchCommand 创建启动命令
func NewLaunchCommand(appID string) *LaunchCommand {
    return &LaunchCommand{AppID: appID}
}

func (c *LaunchCommand) Name() string { return "launch" }

func (c *LaunchCommand) Execute(ctx context.Context) error {
    if c.device == nil {
        return errors.New("device not set")
    }
    return c.device.Launch()
}

// SendKeysCommand 输入文本命令
type SendKeysCommand struct {
    BaseDeviceCommand
    Text     string
    Secure   bool
    Enter    bool
}

// NewSendKeysCommand 创建输入文本命令
func NewSendKeysCommand(text string) *SendKeysCommand {
    return &SendKeysCommand{Text: text}
}

func (c *SendKeysCommand) Name() string { return "inputText" }

func (c *SendKeysCommand) Execute(ctx context.Context) error {
    if c.device == nil {
        return errors.New("device not set")
    }
    return c.device.SendKeys(c.Text)
}

// PressKeyCommand 按键命令
type PressKeyCommand struct {
    BaseDeviceCommand
    KeyCode int
}

// NewPressKeyCommand 创建按键命令
func NewPressKeyCommand(keyCode int) *PressKeyCommand {
    return &PressKeyCommand{KeyCode: keyCode}
}

func (c *PressKeyCommand) Name() string { return "pressKey" }

func (c *PressKeyCommand) Execute(ctx context.Context) error {
    // TODO: implement when device supports pressKey
    return nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/command/... -v -run TestTap`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/command/device_command.go
git commit -m "feat(command): add DeviceCommand implementations"
```

---

## Phase 4: Implement SerialCommand and SystemCommand

### Task 5: Create internal/command/serial_command.go

**Files:**
- Create: `internal/command/serial_command.go`

**Step 1: Write serial command implementation**

```go
package command

import (
    "context"
    "errors"
    "time"

    "github.com/liukunup/go-uop/pkg/serial"
)

// SerialCommand 串口命令接口
type SerialCommand interface {
    Command
    SetSerial(s *serial.Serial)
}

// BaseSerialCommand 基础串口命令
type BaseSerialCommand struct {
    serial *serial.Serial
}

// SetSerial 设置串口连接
func (c *BaseSerialCommand) SetSerial(s *serial.Serial) {
    c.serial = s
}

// SendByIDCommand 通过ID发送命令
type SendByIDCommand struct {
    BaseSerialCommand
    CommandID string
    Timeout   time.Duration
}

// NewSendByIDCommand 创建按ID发送命令
func NewSendByIDCommand(commandID string, timeout time.Duration) *SendByIDCommand {
    return &SendByIDCommand{
        CommandID: commandID,
        Timeout:   timeout,
    }
}

func (c *SendByIDCommand) Name() string { return "sendByID" }

func (c *SendByIDCommand) Validate() error {
    if c.CommandID == "" {
        return errors.New("command ID is required")
    }
    return nil
}

func (c *SendByIDCommand) Execute(ctx context.Context) error {
    if c.serial == nil {
        return errors.New("serial connection not set")
    }
    return c.serial.SendByID(c.CommandID, nil)
}

// SendRawCommand 发送原始数据命令
type SendRawCommand struct {
    BaseSerialCommand
    Data string
}

// NewSendRawCommand 创建发送原始数据命令
func NewSendRawCommand(data string) *SendRawCommand {
    return &SendRawCommand{Data: data}
}

func (c *SendRawCommand) Name() string { return "sendRaw" }

func (c *SendRawCommand) Execute(ctx context.Context) error {
    if c.serial == nil {
        return errors.New("serial connection not set")
    }
    _, err := c.serial.WriteString(c.Data)
    return err
}
```

**Step 2: Commit**

```bash
git add internal/command/serial_command.go
git commit -m "feat(command): add SerialCommand implementations"
```

---

### Task 6: Create internal/command/system_command.go

**Files:**
- Create: `internal/command/system_command.go`

**Step 1: Write system command implementation**

```go
package command

import (
    "context"
    "errors"
    "time"
)

// WaitCommand 等待命令
type WaitCommand struct {
    Duration time.Duration
}

// NewWaitCommand 创建等待命令
func NewWaitCommand(ms int) *WaitCommand {
    return &WaitCommand{Duration: time.Duration(ms) * time.Millisecond}
}

func (c *WaitCommand) Name() string { return "wait" }

func (c *WaitCommand) Validate() error {
    if c.Duration < 0 {
        return errors.New("duration must be non-negative")
    }
    return nil
}

func (c *WaitCommand) Execute(ctx context.Context) error {
    time.Sleep(c.Duration)
    return nil
}

// ScreenshotCommand 截图命令
type ScreenshotCommand struct {
    BaseDeviceCommand
    Name string
    Path string
}

// NewScreenshotCommand 创建截图命令
func NewScreenshotCommand(name, path string) *ScreenshotCommand {
    return &ScreenshotCommand{Name: name, Path: path}
}

func (c *ScreenshotCommand) Name() string { return "screenshot" }

func (c *ScreenshotCommand) Execute(ctx context.Context) error {
    if c.device == nil {
        return errors.New("device not set")
    }
    _, err := c.device.Screenshot()
    return err
}

// SwipeCommand 滑动命令
type SwipeCommand struct {
    BaseDeviceCommand
    StartX, StartY int
    EndX, EndY     int
    Duration       time.Duration
}

// NewSwipeCommand 创建滑动命令
func NewSwipeCommand(startX, startY, endX, endY int) *SwipeCommand {
    return &SwipeCommand{
        StartX: startX,
        StartY: startY,
        EndX:   endX,
        EndY:   endY,
    }
}

func (c *SwipeCommand) Name() string { return "swipe" }

func (c *SwipeCommand) Execute(ctx context.Context) error {
    // TODO: implement when device supports swipe
    return nil
}
```

**Step 2: Commit**

```bash
git add internal/command/system_command.go
git commit -m "feat(command): add SystemCommand implementations"
```

---

## Phase 5: Implement Handler and Refactor Console

### Task 7: Create internal/command/handler.go

**Files:**
- Create: `internal/command/handler.go`

**Step 1: Write handler implementations**

```go
package command

import (
    "context"
    "errors"
)

// DeviceOpsHandler 处理设备操作
type DeviceOpsHandler struct {
    GetDevice func() (DeviceCommand, error)
}

func (h *DeviceOpsHandler) CanHandle(cmd Command) bool {
    _, ok := cmd.(DeviceCommand)
    return ok
}

func (h *DeviceOpsHandler) Handle(ctx context.Context, cmd Command) error {
    devCmd, ok := cmd.(DeviceCommand)
    if !ok {
        return ErrUnsupportedCommand
    }
    
    device, err := h.GetDevice()
    if err != nil {
        return err
    }
    devCmd.SetDevice(device)
    
    return cmd.Execute(ctx)
}

// SerialOpsHandler 处理串口操作
type SerialOpsHandler struct {
    GetSerial func() (SerialCommand, error)
}

func (h *SerialOpsHandler) CanHandle(cmd Command) bool {
    _, ok := cmd.(SerialCommand)
    return ok
}

func (h *SerialOpsHandler) Handle(ctx context.Context, cmd Command) error {
    serCmd, ok := cmd.(SerialCommand)
    if !ok {
        return ErrUnsupportedCommand
    }
    
    serial, err := h.GetSerial()
    if err != nil {
        return err
    }
    serCmd.SetSerial(serial)
    
    return cmd.Execute(ctx)
}

// SystemOpsHandler 处理系统操作
type SystemOpsHandler struct{}

func (h *SystemOpsHandler) CanHandle(cmd Command) bool {
    switch cmd.(type) {
    case *WaitCommand:
        return true
    default:
        return false
    }
}

func (h *SystemOpsHandler) Handle(ctx context.Context, cmd Command) error {
    return cmd.Execute(ctx)
}

// UnknownCommandHandler 处理未知命令
type UnknownCommandHandler struct{}

func (h *UnknownCommandHandler) CanHandle(cmd Command) bool {
    return true
}

func (h *UnknownCommandHandler) Handle(ctx context.Context, cmd Command) error {
    return errors.New("no handler available for command")
}
```

**Step 2: Commit**

```bash
git add internal/command/handler.go
git commit -m "feat(command): add Handler implementations"
```

---

## Phase 6: Refactor Console Handlers

### Task 8: Create internal/command/router.go (CommandRouter)

**Files:**
- Create: `internal/command/router.go`

**Step 1: Write CommandRouter**

```go
package command

import (
    "context"
)

// CommandRouter 命令路由器
type CommandRouter struct {
    registry *CommandRegistry
}

// NewCommandRouter 创建命令路由器
func NewCommandRouter(reg *CommandRegistry) *CommandRouter {
    return &CommandRouter{registry: reg}
}

// Route 根据名称路由命令
func (r *CommandRouter) Route(name string) (Command, error) {
    cmd := r.registry.Get(name)
    if cmd == nil {
        return nil, ErrUnknownCommand
    }
    return cmd, nil
}

// Execute 执行命令
func (r *CommandRouter) Execute(ctx context.Context, name string) error {
    return r.registry.Execute(ctx, name)
}
```

**Step 2: Commit**

```bash
git add internal/command/router.go
git commit -m "feat(command): add CommandRouter"
```

---

### Task 9: Refactor internal/console/handler.go

**Files:**
- Modify: `internal/console/handler.go`

**Changes:**
- Replace large switch-case in `ExecuteCommand` with CommandRouter
- Keep existing handler methods but use Command pattern internally

**Step 1: Review current implementation**

Read `internal/console/handler.go` lines 148-195 (ExecuteCommand method)

**Step 2: Refactor to use Command pattern**

```go
// ExecuteCommand executes a command on a device using Command pattern
func (m *DeviceManager) ExecuteCommand(id string, cmd string, params map[string]interface{}) (*CommandRecord, error) {
    device, err := m.GetConnected(id)
    if err != nil {
        return nil, err
    }

    record := &CommandRecord{
        ID:        generateID(),
        Timestamp: time.Now().Format(time.RFC3339),
        Type:      cmd,
        Params:    params,
    }

    start := time.Now()
    defer func() {
        record.Duration = time.Since(start).String()
    }()

    // Create command based on type
    var command command.DeviceCommand
    switch cmd {
    case "tap":
        x, _ := toInt(params["x"])
        y, _ := toInt(params["y"])
        command = command.NewTapCommand(x, y)
    case "input":
        text, _ := toString(params["text"])
        command = command.NewSendKeysCommand(text)
    case "launch":
        command = command.NewLaunchCommand("")
    case "screenshot":
        command = command.NewScreenshotCommand("", "")
    default:
        err = ErrUnknownCommand
        record.Success = false
        record.Output = err.Error()
        return record, err
    }

    command.SetDevice(device)
    
    err = command.Execute(context.Background())
    record.Success = err == nil
    if err != nil {
        record.Output = err.Error()
    }

    return record, err
}
```

**Step 3: Commit**

```bash
git add internal/console/handler.go
git commit -m "refactor(console): use Command pattern in ExecuteCommand"
```

---

## Phase 7: Cleanup Deprecated Files

### Task 10: Delete deprecated files

**Files:**
- Delete: `internal/action/action.go`
- Delete: `internal/runner/executor.go`

**Step 1: Verify no other files depend on them**

Run: `grep -r "internal/action" --include="*.go"`
Run: `grep -r "internal/runner/executor" --include="*.go"`

**Step 2: Delete files**

```bash
rm internal/action/action.go
rm internal/runner/executor.go
git add -A
git commit -m "refactor: remove deprecated action and executor files"
```

---

## Phase 8: Write Tests

### Task 11: Write comprehensive tests

**Files:**
- Create: `internal/command/command_test.go`
- Create: `internal/command/registry_test.go`
- Create: `internal/command/device_command_test.go`
- Create: `internal/command/serial_command_test.go`
- Create: `internal/command/system_command_test.go`

**Step 1: Run all tests**

Run: `go test ./internal/command/... -v -cover`

**Step 2: Verify coverage > 70%**

If coverage < 70%, add more test cases.

**Step 3: Commit**

```bash
git add internal/command/*_test.go
git commit -m "test(command): add comprehensive tests"
```

---

## Summary

| Task | Description | Status |
|------|-------------|--------|
| 1 | Create command interface | Pending |
| 2 | Create errors | Pending |
| 3 | Create CommandRegistry | Pending |
| 4 | Create DeviceCommand | Pending |
| 5 | Create SerialCommand | Pending |
| 6 | Create SystemCommand | Pending |
| 7 | Create Handlers | Pending |
| 8 | Create CommandRouter | Pending |
| 9 | Refactor console handlers | Pending |
| 10 | Delete deprecated files | Pending |
| 11 | Write tests | Pending |

---

**Plan complete and saved to `docs/plans/2026-04-04-oop-command-implementation.md`.**

Two execution options:

**1. Subagent-Driven (this session)** - I dispatch fresh subagent per task, review between tasks, fast iteration

**2. Parallel Session (separate)** - Open new session with executing-plans, batch execution with checkpoints

Which approach?
