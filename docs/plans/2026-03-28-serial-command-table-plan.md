# Serial Command Table Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 实现命令表管理，支持从 YAML 文件加载命令，通过 ID/Name 发送命令并通过 Monitor 正则校验回显。

**Architecture:** 新增 Command/CommandTable 到 command.go，Monitor 新增 AddTemporaryRule，Serial 新增 SendByID/SendByName。

**Tech Stack:** Go, gopkg.in/yaml.v3, sync

---

## Task 1: 创建 pkg/serial/command.go

**Files:**
- Create: `pkg/serial/command.go`

**Step 1: 编写 command.go**

```go
package serial

import (
    "fmt"
    "os"
    "sync"

    "gopkg.in/yaml.v3"
)

// Command 串口命令
type Command struct {
    ID      string        `yaml:"id"`      // 唯一标识符
    Name    string        `yaml:"name"`    // 易读名称
    Command string        `yaml:"command"`  // 发送字节序列（字符串格式）
    Log     string        `yaml:"log"`     // 回显校验正则（可选）
    Timeout time.Duration `yaml:"timeout"`  // 超时时间
}

// CommandTable 命令表
type CommandTable struct {
    mu       sync.RWMutex
    commands map[string]*Command  // by ID
    byName   map[string]*Command  // by Name
}

// commandTableFile YAML 文件格式
type commandTableFile struct {
    Commands []*Command `yaml:"commands"`
}

// NewCommandTable 创建空命令表
func NewCommandTable() *CommandTable {
    return &CommandTable{
        commands: make(map[string]*Command),
        byName:   make(map[string]*Command),
    }
}

// LoadFromFile 从 YAML 文件加载命令表
func (ct *CommandTable) LoadFromFile(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("read command table file: %w", err)
    }

    var f commandTableFile
    if err := yaml.Unmarshal(data, &f); err != nil {
        return fmt.Errorf("parse command table file: %w", err)
    }

    ct.mu.Lock()
    defer ct.mu.Unlock()

    for _, cmd := range f.Commands {
        if cmd.ID == "" {
            return fmt.Errorf("command missing id")
        }
        ct.commands[cmd.ID] = cmd
        if cmd.Name != "" {
            ct.byName[cmd.Name] = cmd
        }
    }

    return nil
}

// GetByID 通过 ID 获取命令
func (ct *CommandTable) GetByID(id string) (*Command, bool) {
    ct.mu.RLock()
    defer ct.mu.RUnlock()
    cmd, ok := ct.commands[id]
    return cmd, ok
}

// GetByName 通过名称获取命令
func (ct *CommandTable) GetByName(name string) (*Command, bool) {
    ct.mu.RLock()
    defer ct.mu.RUnlock()
    cmd, ok := ct.byName[name]
    return cmd, ok
}

// List 返回所有命令
func (ct *CommandTable) List() []*Command {
    ct.mu.RLock()
    defer ct.mu.RUnlock()
    cmds := make([]*Command, 0, len(ct.commands))
    for _, cmd := range ct.commands {
        cmds = append(cmds, cmd)
    }
    return cmds
}
```

**Step 2: 添加 time import 并验证**

需要添加 `time` 到 import。

Run: `go build ./pkg/serial/...`
Expected: PASS

---

## Task 2: 改造 pkg/serial/monitor.go

**Files:**
- Modify: `pkg/serial/monitor.go`

**Step 1: 在文件末尾添加 AddTemporaryRule**

```go
// cleanupFunc 清理函数，用于移除临时规则
type cleanupFunc func()

// AddTemporaryRule 添加临时规则（命令执行完自动移除）
// 返回 cleanupFunc，调用后立即移除规则
func (m *Monitor) AddTemporaryRule(keyword string, matchType MatchType, handler EventHandler) cleanupFunc {
    m.mu.Lock()
    defer m.mu.Unlock()

    rule := &Rule{
        Keyword:   keyword,
        MatchType: matchType,
        enabled:   true,
    }

    m.rules = append(m.rules, rule)
    m.handlers[keyword] = handler

    return func() {
        m.mu.Lock()
        defer m.mu.Unlock()

        // 找到规则并移除
        for i, r := range m.rules {
            if r == rule {
                m.rules = append(m.rules[:i], m.rules[i+1:]...)
                delete(m.handlers, keyword)
                return
            }
        }
    }
}
```

**Step 2: 验证编译**

Run: `go build ./pkg/serial/...`
Expected: PASS

---

## Task 3: 改造 pkg/serial/serial.go

**Files:**
- Modify: `pkg/serial/serial.go`

**Step 1: 添加 SendResult 和 SendCallback 类型（在 Event 类型附近）**

```go
// SendResult 发送结果
type SendResult struct {
    Command *Command
    Success bool
    Echo    []byte
    Matched bool
    Error   error
}

// SendCallback 发送完成回调
type SendCallback func(*SendResult)
```

**Step 2: 在 Serial 结构体中添加 monitor 字段**

```go
// Serial is a serial port connection.
type Serial struct {
    mu        sync.RWMutex
    observers []Observer
    monitor   *Monitor  // 内置监视器用于命令回显校验
    eventCh   chan Event
    done      chan struct{}
    cfg       *Config
    port      *serial.Port
    readErr   error
}
```

**Step 3: 修改 NewSerial 初始化内置 Monitor**

```go
func NewSerial(cfg *Config) (*Serial, error) {
    // ... 现有代码直到 port 打开 ...

    s := &Serial{
        cfg:     cfg,
        port:    port,
        monitor: NewMonitor(),  // 初始化内置 Monitor
        eventCh: make(chan Event, 100),
        done:    make(chan struct{}),
    }

    // 将内置 Monitor 添加为 Observer
    s.observers = append(s.observers, s.monitor)

    // 启动后台读取 goroutine
    go s.readLoop()

    return s, nil
}
```

**Step 4: 添加 SendByID 和 SendByName 方法**

```go
// SendByID 通过 ID 发送命令
func (s *Serial) SendByID(id string, callback SendCallback) error {
    ct := NewCommandTable() // 需要外部传入或从 Serial 持有
    // 改为：Serial 需要持有 CommandTable 或通过参数传入
    // 设计调整：CommandTable 作为参数传入
    return nil
}
```

**实际设计：CommandTable 作为 Serial 的可选字段，或作为方法参数**

```go
// 方案：CommandTable 作为 Serial 的可选配置

// 修改 Config 添加 CommandTable
type Config struct {
    // ... 现有字段 ...
    Commands *CommandTable  // 可选，命令表
}

// 修改 SendByID
func (s *Serial) SendByID(id string, callback SendCallback) error {
    s.mu.RLock()
    ct := s.cfg.Commands
    s.mu.RUnlock()

    if ct == nil {
        return fmt.Errorf("command table not configured")
    }

    cmd, ok := ct.GetByID(id)
    if !ok {
        return fmt.Errorf("command not found: %s", id)
    }

    return s.sendCommand(cmd, callback)
}

// SendByName 通过名称发送命令
func (s *Serial) SendByName(name string, callback SendCallback) error {
    s.mu.RLock()
    ct := s.cfg.Commands
    s.mu.RUnlock()

    if ct == nil {
        return fmt.Errorf("command table not configured")
    }

    cmd, ok := ct.GetByName(name)
    if !ok {
        return fmt.Errorf("command not found: %s", name)
    }

    return s.sendCommand(cmd, callback)
}

// sendCommand 内部发送方法
func (s *Serial) sendCommand(cmd *Command, callback SendCallback) error {
    result := &SendResult{
        Command: cmd,
    }

    var cleanup cleanupFunc

    // 如果有 Log，设置 Monitor 规则进行回显校验
    if cmd.Log != "" {
        cleanup = s.monitor.AddTemporaryRule(cmd.Log, MatchOnce, func(e Event) {
            result.Success = true
            result.Matched = true
            result.Echo = e.Data
            if callback != nil {
                callback(result)
            }
        })

        // 启动超时处理
        if cmd.Timeout > 0 {
            go func() {
                time.Sleep(cmd.Timeout)
                if cleanup != nil {
                    cleanup()
                    cleanup = nil
                }
                if !result.Success && callback != nil {
                    result.Success = false
                    result.Error = fmt.Errorf("command timeout: %s", cmd.ID)
                    callback(result)
                }
            }()
        }
    }

    // 发送命令
    _, err := s.port.Write([]byte(cmd.Command))
    if err != nil {
        if cleanup != nil {
            cleanup()
        }
        return fmt.Errorf("send command: %w", err)
    }

    // 如果没有 Log，立即返回成功
    if cmd.Log == "" {
        result.Success = true
        if callback != nil {
            callback(result)
        }
    }

    return nil
}
```

**Step 5: 添加 import "time"**

确保 `time` 包在 import 中。

**Step 6: 验证编译**

Run: `go build ./pkg/serial/...`
Expected: PASS

---

## Task 4: 创建单元测试

**Files:**
- Create: `pkg/serial/command_test.go`

**Step 1: 编写测试**

```go
package serial

import (
    "os"
    "testing"
)

func TestCommandTable_New(t *testing.T) {
    ct := NewCommandTable()
    if ct == nil {
        t.Fatal("expected non-nil CommandTable")
    }
}

func TestCommandTable_LoadFromFile(t *testing.T) {
    // 创建临时 YAML 文件
    content := `
commands:
  - id: reset
    name: 设备复位
    command: "AT+RST\\r\\n"
    log: "(?i)ok"
    timeout: 5s
  - id: version
    name: 查询版本
    command: "AT+VERSION?\\r\\n"
    timeout: 3s
`
    tmpfile, err := os.CreateTemp("", "commands_*.yaml")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmpfile.Name())

    if _, err := tmpfile.WriteString(content); err != nil {
        t.Fatal(err)
    }
    tmpfile.Close()

    ct := NewCommandTable()
    if err := ct.LoadFromFile(tmpfile.Name()); err != nil {
        t.Fatalf("LoadFromFile failed: %v", err)
    }

    // 测试 GetByID
    cmd, ok := ct.GetByID("reset")
    if !ok {
        t.Fatal("expected to find command 'reset'")
    }
    if cmd.Name != "设备复位" {
        t.Errorf("expected name '设备复位', got '%s'", cmd.Name)
    }
    if cmd.Command != "AT+RST\r\n" {
        t.Errorf("expected command 'AT+RST\\r\\n', got '%s'", cmd.Command)
    }

    // 测试 GetByName
    cmd, ok = ct.GetByName("设备复位")
    if !ok {
        t.Fatal("expected to find command by name '设备复位'")
    }
    if cmd.ID != "reset" {
        t.Errorf("expected id 'reset', got '%s'", cmd.ID)
    }

    // 测试不存在的命令
    _, ok = ct.GetByID("notexist")
    if ok {
        t.Error("expected not to find 'notexist'")
    }
}

func TestCommandTable_List(t *testing.T) {
    ct := NewCommandTable()
    // 直接添加命令进行测试
    ct.mu.Lock()
    ct.commands["cmd1"] = &Command{ID: "cmd1", Name: "Command 1"}
    ct.commands["cmd2"] = &Command{ID: "cmd2", Name: "Command 2"}
    ct.mu.Unlock()

    cmds := ct.List()
    if len(cmds) != 2 {
        t.Errorf("expected 2 commands, got %d", len(cmds))
    }
}
```

**Step 2: 运行测试**

Run: `go test ./pkg/serial/... -v -run "TestCommandTable" -count=1`
Expected: PASS

---

## Task 5: Git 提交

**Step 1: 提交**

```bash
git add pkg/serial/
git commit -m "feat(serial): add CommandTable for command management

- Add Command and CommandTable types
- Support loading commands from YAML file
- Add SendByID and SendByName for sending commands
- Integrate Monitor with regex validation for echo checking
- Add timeout handling for async command responses"
```

---

## Summary

| Task | Component | Files |
|------|-----------|-------|
| 1 | Command/CommandTable | `pkg/serial/command.go` |
| 2 | AddTemporaryRule | `pkg/serial/monitor.go` |
| 3 | SendByID/SendByName | `pkg/serial/serial.go` |
| 4 | Tests | `pkg/serial/command_test.go` |
| 5 | Commit | - |
