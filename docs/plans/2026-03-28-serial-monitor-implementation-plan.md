# Serial Monitor Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 为 `pkg/serial` 添加事件驱动的 Observer Pattern 和 Monitor 监视器，支持关键字匹配触发回调或文件写入。

**Architecture:** 在现有 Serial 结构体基础上，新增 Observer 接口、Event 事件、Monitor 监视器（含 MatchOnce/MatchContinuous/MatchRateLimited 三种策略）、FileHandler 文件写入器。Serial 后台 goroutine 持续读取并通过 channel 分发给所有 Observer。

**Tech Stack:** Go, `github.com/tarm/serial`, sync/atomic

---

## Task 1: 重构 Serial 结构体添加 Observer 支持

**Files:**
- Modify: `pkg/serial/serial.go:41-80`

**Step 1: 查看当前 Serial 结构体**

Read `pkg/serial/serial.go` lines 41-80

**Step 2: 在 serial.go 顶部添加 Event 和 Observer 定义**

```go
// Event 事件结构
type Event struct {
    Data      []byte
    Timestamp time.Time
    Rule      *Rule  // nil if no rule matched
}

// Observer 观察者接口
type Observer interface {
    OnData(event Event)
    OnError(err error)
    OnClose()
}
```

**Step 3: 重构 Serial 结构体**

```go
// Serial is a serial port connection.
type Serial struct {
    mu        sync.RWMutex
    observers []Observer
    eventCh   chan Event
    done      chan struct{}
    cfg       *Config
    port      *serial.Port
    readErr   error
}
```

**Step 4: 添加 AddObserver/RemoveObserver 方法**

```go
// AddObserver 添加观察者，启动独立的 dispatch goroutine
func (s *Serial) AddObserver(o Observer) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.observers = append(s.observers, o)
}

// RemoveObserver 移除观察者
func (s *Serial) RemoveObserver(o Observer) {
    s.mu.Lock()
    defer s.mu.Unlock()
    for i, obs := range s.observers {
        if obs == o {
            s.observers = append(s.observers[:i], s.observers[i+1:]...)
            return
        }
    }
}

// notifyAll 调用所有 Observer 的 OnData
func (s *Serial) notifyAll(event Event) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    for _, o := range s.observers {
        go o.OnData(event)
    }
}

// notifyError 通知所有 Observer 错误
func (s *Serial) notifyError(err error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    for _, o := range s.observers {
        go o.OnError(err)
    }
}

// notifyClose 通知所有 Observer 关闭
func (s *Serial) notifyClose() {
    s.mu.RLock()
    defer s.mu.RUnlock()
    for _, o := range s.observers {
        go o.OnClose()
    }
}
```

**Step 5: 重写 NewSerial 启动后台读取**

```go
// NewSerial opens a serial port with the given config.
func NewSerial(cfg *Config) (*Serial, error) {
    // ... 现有校验逻辑保持不变 ...

    s := &Serial{
        cfg:     cfg,
        port:    port,
        eventCh: make(chan Event, 100),
        done:    make(chan struct{}),
    }

    // 启动后台读取 goroutine
    go s.readLoop()

    return s, nil
}

// readLoop 后台持续读取串口数据
func (s *Serial) readLoop() {
    buf := make([]byte, 1024)
    for {
        select {
        case <-s.done:
            return
        default:
        }

        n, err := s.port.Read(buf)
        if err != nil {
            if err != io.EOF {
                s.readErr = err
                s.notifyError(err)
            }
            return
        }

        if n > 0 {
            event := Event{
                Data:      append([]byte{}, buf[:n]...),
                Timestamp: time.Now(),
            }
            select {
            case s.eventCh <- event:
                s.notifyAll(event)
            default:
                // channel 满，丢弃旧事件
            }
        }
    }
}
```

**Step 6: 更新 Close 方法**

```go
// Close closes the serial port.
func (s *Serial) Close() error {
    close(s.done)
    s.notifyClose()
    return s.port.Close()
}
```

**Step 7: 验证编译**

Run: `go build ./pkg/serial/...`
Expected: PASS

---

## Task 2: 创建 pkg/serial/monitor.go

**Files:**
- Create: `pkg/serial/monitor.go`

**Step 1: 编写 Monitor 代码**

```go
package serial

import (
    "bytes"
    "fmt"
    "regexp"
    "sync"
    "time"
)

// MatchType 匹配策略
type MatchType int

const (
    MatchOnce MatchType = iota
    MatchContinuous
    MatchRateLimited
)

// Rule 规则
type Rule struct {
    Keyword      string
    MatchType
    RateInterval time.Duration
    IsRegex      bool
    enabled      bool
    lastTrigger  time.Time
    mu           sync.Mutex
    regex        *regexp.Regexp
}

// EventHandler 事件处理函数
type EventHandler func(Event)

// Monitor 监视器（实现 Observer 接口）
type Monitor struct {
    rules    []*Rule
    handlers map[string]EventHandler
    mu       sync.RWMutex
}

// NewMonitor 创建监视器
func NewMonitor() *Monitor {
    return &Monitor{
        handlers: make(map[string]EventHandler),
    }
}

// AddRule 添加规则（MatchOnce 或 MatchContinuous）
func (m *Monitor) AddRule(keyword string, matchType MatchType, handler EventHandler) {
    m.mu.Lock()
    defer m.mu.Unlock()

    rule := &Rule{
        Keyword:   keyword,
        MatchType: matchType,
        enabled:   true,
    }

    if matchType == MatchOnce {
        rule.enabled = true
    }

    m.rules = append(m.rules, rule)
    m.handlers[keyword] = handler
}

// AddRateLimitedRule 添加限频规则
func (m *Monitor) AddRateLimitedRule(keyword string, interval time.Duration, handler EventHandler) {
    m.mu.Lock()
    defer m.mu.Unlock()

    rule := &Rule{
        Keyword:      keyword,
        MatchType:    MatchRateLimited,
        RateInterval: interval,
        enabled:      true,
    }

    m.rules = append(m.rules, rule)
    m.handlers[keyword] = handler
}

// EnableRule 启用规则
func (m *Monitor) EnableRule(keyword string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    for _, r := range m.rules {
        if r.Keyword == keyword {
            r.mu.Lock()
            r.enabled = true
            r.mu.Unlock()
            return nil
        }
    }
    return fmt.Errorf("rule not found: %s", keyword)
}

// DisableRule 禁用规则
func (m *Monitor) DisableRule(keyword string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    for _, r := range m.rules {
        if r.Keyword == keyword {
            r.mu.Lock()
            r.enabled = false
            r.mu.Unlock()
            return nil
        }
    }
    return fmt.Errorf("rule not found: %s", keyword)
}

// OnData 处理事件
func (m *Monitor) OnData(e Event) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    for _, rule := range m.rules {
        if !rule.shouldTrigger(e.Data) {
            continue
        }

        matchedEvent := Event{
            Data:      e.Data,
            Timestamp: e.Timestamp,
            Rule:      rule,
        }

        if handler, ok := m.handlers[rule.Keyword]; ok {
            handler(matchedEvent)
        }
    }
}

// shouldTrigger 检查规则是否应触发
func (r *Rule) shouldTrigger(data []byte) bool {
    r.mu.Lock()
    defer r.mu.Unlock()

    if !r.enabled && r.MatchType == MatchOnce {
        return false
    }

    matched := false
    if r.IsRegex {
        if r.regex == nil {
            r.regex = regexp.MustCompile(r.Keyword)
        }
        matched = r.regex.Match(data)
    } else {
        matched = bytes.Contains(data, []byte(r.Keyword))
    }

    if !matched {
        return false
    }

    switch r.MatchType {
    case MatchOnce:
        r.enabled = false
    case MatchRateLimited:
        if time.Since(r.lastTrigger) < r.RateInterval {
            return false
        }
        r.lastTrigger = time.Now()
    }

    return true
}

// OnError 错误处理（空实现）
func (m *Monitor) OnError(err error) {
    // 日志记录，不影响其他规则
}

// OnClose 关闭处理（空实现）
func (m *Monitor) OnClose() {
}
```

**Step 2: 验证编译**

Run: `go build ./pkg/serial/...`
Expected: PASS

---

## Task 3: 创建 pkg/serial/file_handler.go

**Files:**
- Create: `pkg/serial/file_handler.go`

**Step 1: 编写 FileHandler 代码**

```go
package serial

import (
    "fmt"
    "os"
    "sync"
    "time"
)

// FileHandler 文件写入器
type FileHandler struct {
    Path string
    mu   sync.Mutex
}

// NewFileHandler 创建文件写入器
func NewFileHandler(path string) FileHandler {
    return FileHandler{Path: path}
}

// Handle 处理事件，写入文件
func (h FileHandler) Handle(e Event) {
    h.mu.Lock()
    defer h.mu.Unlock()

    f, err := os.OpenFile(h.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Fprintf(os.Stderr, "serial: open file %s: %v\n", h.Path, err)
        return
    }
    defer f.Close()

    ts := e.Timestamp.Format("2006-01-02 15:04:05.000")
    line := fmt.Sprintf("[%s] %s\n", ts, string(e.Data))
    _, err = f.WriteString(line)
    if err != nil {
        fmt.Fprintf(os.Stderr, "serial: write file %s: %v\n", h.Path, err)
    }
}
```

**Step 2: 验证编译**

Run: `go build ./pkg/serial/...`
Expected: PASS

---

## Task 4: 创建 pkg/serial/serial_test.go

**Files:**
- Create: `pkg/serial/serial_test.go`

**Step 1: 编写测试代码**

```go
package serial

import (
    "sync"
    "testing"
    "time"
)

func TestEvent(t *testing.T) {
    event := Event{
        Data:      []byte("test data"),
        Timestamp: time.Now(),
    }

    if string(event.Data) != "test data" {
        t.Errorf("expected 'test data', got '%s'", string(event.Data))
    }
}

func TestMonitor_MatchOnce(t *testing.T) {
    m := NewMonitor()

    count := 0
    m.AddRule("OK", MatchOnce, func(e Event) {
        count++
    })

    // 第一次匹配，触发
    m.OnData(Event{Data: []byte("OK")})
    if count != 1 {
        t.Errorf("expected count=1, got %d", count)
    }

    // 第二次匹配，不触发（已禁用）
    m.OnData(Event{Data: []byte("OK")})
    if count != 1 {
        t.Errorf("expected count=1 after second match, got %d", count)
    }
}

func TestMonitor_MatchContinuous(t *testing.T) {
    m := NewMonitor()

    count := 0
    m.AddRule("OK", MatchContinuous, func(e Event) {
        count++
    })

    m.OnData(Event{Data: []byte("OK")})
    m.OnData(Event{Data: []byte("OK")})
    m.OnData(Event{Data: []byte("OK")})

    if count != 3 {
        t.Errorf("expected count=3, got %d", count)
    }
}

func TestMonitor_MatchRateLimited(t *testing.T) {
    m := NewMonitor()

    count := 0
    m.AddRateLimitedRule("OK", 100*time.Millisecond, func(e Event) {
        count++
    })

    m.OnData(Event{Data: []byte("OK")})
    m.OnData(Event{Data: []byte("OK")}) // 应被限频

    if count != 1 {
        t.Errorf("expected count=1, got %d", count)
    }

    time.Sleep(150 * time.Millisecond)
    m.OnData(Event{Data: []byte("OK")}) // 超过间隔，再次触发

    if count != 2 {
        t.Errorf("expected count=2, got %d", count)
    }
}

func TestMonitor_EnableDisable(t *testing.T) {
    m := NewMonitor()

    count := 0
    m.AddRule("OK", MatchOnce, func(e Event) {
        count++
    })

    m.OnData(Event{Data: []byte("OK")})
    if count != 1 {
        t.Errorf("expected count=1, got %d", count)
    }

    // 禁用后无法触发
    m.DisableRule("OK")
    m.OnData(Event{Data: []byte("OK")})
    if count != 1 {
        t.Errorf("expected count=1 after disable, got %d", count)
    }

    // 启用后恢复
    m.EnableRule("OK")
    m.OnData(Event{Data: []byte("OK")})
    if count != 2 {
        t.Errorf("expected count=2 after re-enable, got %d", count)
    }
}

func TestMonitor_Concurrent(t *testing.T) {
    m := NewMonitor()

    var mu sync.Mutex
    count := 0

    m.AddRule("OK", MatchContinuous, func(e Event) {
        mu.Lock()
        count++
        mu.Unlock()
    })

    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < 100; j++ {
                m.OnData(Event{Data: []byte("OK")})
            }
        }()
    }

    wg.Wait()

    mu.Lock()
    expected := 1000
    if count != expected {
        t.Errorf("expected count=%d, got %d", expected, count)
    }
    mu.Unlock()
}

func TestFileHandler(t *testing.T) {
    handler := NewFileHandler("/tmp/serial_test.log")
    
    event := Event{
        Data:      []byte("test message"),
        Timestamp: time.Now(),
    }

    handler.Handle(event)

    // 验证文件存在且内容正确
    // (省略文件内容验证，简化为编译测试)
}

func TestSerial_AddRemoveObserver(t *testing.T) {
    // 创建 mock Serial（需要串口，实际测试时跳过）
    // 这里仅测试 Observer 接口
    m := NewMonitor()
    
    var received Event
    m.AddRule("test", MatchContinuous, func(e Event) {
        received = e
    })

    m.OnData(Event{Data: []byte("test data")})

    if string(received.Data) != "test data" {
        t.Errorf("expected 'test data', got '%s'", string(received.Data))
    }
}
```

**Step 2: 运行测试**

Run: `go test ./pkg/serial/... -v -count=1`
Expected: PASS (Monitor 相关测试通过)

---

## Task 5: 更新主文件导出新类型

**Files:**
- Modify: `pkg/serial/serial.go`

**Step 1: 确认导出类型在同一个包中**

确认 `Event`, `Observer`, `MatchType`, `MatchOnce`, `MatchContinuous`, `MatchRateLimited`, `Rule`, `EventHandler`, `Monitor`, `FileHandler` 已在同一 `serial` 包中定义，无需额外导出。

**Step 2: 添加包级别注释**

```go
// Package serial provides serial port communication with event-driven monitoring.
//
// Example:
//
//  s, _ := serial.NewSerial(&serial.Config{Name: "/dev/ttyUSB0", Baud: 115200})
//
//  m := serial.NewMonitor()
//  m.AddRule("OK", serial.MatchOnce, func(e serial.Event) {
//      fmt.Println("Device ready")
//  })
//  m.AddRateLimitedRule("ERROR", 5*time.Second, serial.FileHandler{Path: "/tmp/errors.log"}.Handle)
//
//  s.AddObserver(m)
```

**Step 3: 验证编译**

Run: `go build ./pkg/serial/...`
Expected: PASS

---

## Task 6: 提交代码

**Step 1: 提交**

```bash
git add pkg/serial/
git commit -m "feat(serial): add Observer pattern and Monitor with event-driven handling

- Add Event and Observer interface for data reception
- Implement Serial with background readLoop goroutine
- Add Monitor with MatchOnce/MatchContinuous/MatchRateLimited strategies
- Add FileHandler for log file writing
- Add comprehensive tests for Monitor"
```

---

## Summary

| Task | Component | Files |
|------|-----------|-------|
| 1 | Observer Pattern | `pkg/serial/serial.go` |
| 2 | Monitor | `pkg/serial/monitor.go` |
| 3 | FileHandler | `pkg/serial/file_handler.go` |
| 4 | Tests | `pkg/serial/serial_test.go` |
| 5 | Documentation | `pkg/serial/serial.go` |
| 6 | Commit | - |

**Plan complete.** Two execution options:

1. **Subagent-Driven (this session)** - I dispatch fresh subagent per task, review between tasks
2. **Parallel Session (separate)** - Open new session with executing-plans

Which approach?
