# Serial Monitor Design

## Overview

集成 `github.com/tarm/serial` 实现串口收发，并设计事件驱动的监视器机制，通过关键字匹配触发回调或文件写入。

## Architecture

```
┌─────────────────────────────────────────────────────┐
│  Serial (Subject)                                   │
│  - 后台 goroutine 持续读取                          │
│  - 维护 Observer 列表                               │
│  - 数据到达 → 通知所有 Observer                     │
└─────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────┐
│  Event                                              │
│  - Data: []byte                                     │
│  - Timestamp: time.Time                             │
│  - Rule: *Rule (if any)                             │
└─────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────┐
│  Observer (interface)                               │
│  - OnData(event)                                    │
│  - OnError(err)                                     │
│  - OnClose()                                        │
└─────────────────────────────────────────────────────┘
          │                   │
          ▼                   ▼
┌─────────────────┐   ┌─────────────────┐
│  CallbackHandler│   │  FileHandler    │
│  - 函数回调     │   │  - 写入文件     │
└─────────────────┘   └─────────────────┘
```

## Components

### 1. Serial (Subject)

```go
type Serial struct {
    mu        sync.RWMutex
    observers []Observer
    eventCh  chan Event
    closed   chan struct{}
    cfg      *Config
    port     *serial.Port
}
```

**职责：**
- 打开/关闭串口
- 启动 `readLoop()` goroutine 持续读取数据
- 管理 Observer 列表
- 通过 `dispatch()` goroutine 将事件分发给所有 Observer

### 2. Event

```go
type Event struct {
    Data      []byte
    Timestamp time.Time
    Rule      *Rule  // nil if no rule matched
}
```

### 3. Observer Interface

```go
type Observer interface {
    OnData(event Event)
    OnError(err error)
    OnClose()
}
```

### 4. Monitor (implements Observer)

```go
type MatchType int

const (
    MatchOnce      MatchType = iota  // 一次性，触发后自动禁用
    MatchContinuous                   // 持续触发
    MatchRateLimited                  // 限频触发
)

type Rule struct {
    Keyword      string
    MatchType
    RateInterval time.Duration  // 限频间隔
    IsRegex      bool
    enabled      bool
    lastTrigger  time.Time
    mu           sync.Mutex
}

type EventHandler func(Event)

type Monitor struct {
    rules    []*Rule
    handlers map[string]EventHandler
    mu       sync.RWMutex
}
```

### 5. FileHandler

```go
type FileHandler struct {
    Path string
    mu   sync.Mutex
}
```

## Data Flow

```
Serial.Open()
    │
    ├── 启动 readLoop() goroutine
    │       │
    │       └── for { Read() → Event → eventCh }
    │
    └── AddObserver() → 启动 dispatch() goroutine
            │
            └── for { <-eventCh → Observer.OnData() }

Serial.Close()
    │
    ├── 关闭 readLoop (close port)
    ├── 关闭 dispatch
    └── 通知所有 Observer.OnClose()
```

## Matching Logic

```go
func (m *Monitor) OnData(e Event) {
    for _, rule := range m.rules {
        if !rule.shouldTrigger(e.Data) {
            continue
        }
        // 构建带 Rule 的 Event
        matchedEvent := Event{Data: e.Data, Timestamp: e.Timestamp, Rule: rule}
        if handler, ok := m.handlers[rule.Keyword]; ok {
            handler(matchedEvent)
        }
    }
}

func (r *Rule) shouldTrigger(data []byte) bool {
    r.mu.Lock()
    defer r.mu.Unlock()

    if !r.enabled && r.MatchType == MatchOnce {
        return false
    }

    if r.IsRegex {
        if !regexp.MatchString(r.Keyword, string(data)) {
            return false
        }
    } else {
        if !bytes.Contains(data, []byte(r.Keyword)) {
            return false
        }
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
```

## API

### Serial

```go
func NewSerial(cfg *Config) (*Serial, error)
func (s *Serial) AddObserver(o Observer)
func (s *Serial) RemoveObserver(o Observer)
func (s *Serial) Read(p []byte) (int, error)
func (s *Serial) Write(p []byte) (int, error)
func (s *Serial) WriteString(str string) (int, error)
func (s *Serial) Close() error
func (s *Serial) Config() *Config
```

### Monitor

```go
func NewMonitor() *Monitor
func (m *Monitor) AddRule(keyword string, matchType MatchType, handler EventHandler)
func (m *Monitor) AddRateLimitedRule(keyword string, interval time.Duration, handler EventHandler)
func (m *Monitor) EnableRule(keyword string) error
func (m *Monitor) DisableRule(keyword string) error
func (m *Monitor) OnData(event Event)
func (m *Monitor) OnError(err error)
func (m *Monitor) OnClose()
```

### FileHandler

```go
func (h FileHandler) Handle(event Event)
```

## Usage Example

```go
s, _ := serial.NewSerial(&serial.Config{Name: "/dev/ttyUSB0", Baud: 115200})

m := serial.NewMonitor()

// 一次性回调
m.AddRule("OK", serial.MatchOnce, func(e serial.Event) {
    fmt.Println("Device ready")
})

// 限频写入文件
m.AddRateLimitedRule("ERROR", 5*time.Second, serial.FileHandler{Path: "/tmp/errors.log"}.Handle)

// 持续回调
m.AddRule("DATA:", serial.MatchContinuous, func(e serial.Event) {
    fmt.Printf("Received: %s\n", e.Data)
})

s.AddObserver(m)

// 数据开始接收...
s.Close()
```

## Error Handling

- 串口读取错误：关闭事件 channel，通知所有 Observer.OnError()
- 规则匹配错误：记录日志，不影响其他规则
- 文件写入错误：记录到标准错误，不阻塞事件处理

## Thread Safety

- Serial 使用 sync.RWMutex 保护 observers 列表
- Monitor 使用 sync.RWMutex 保护 rules/handlers
- Rule 内部使用 sync.Mutex 保护 enabled 和 lastTrigger
- FileHandler 使用 sync.Mutex 保护文件操作

## Considerations

- 事件 channel 容量：100（阻塞时丢弃旧事件）
- readLoop 使用 select + done channel 实现优雅关闭
- 支持正则匹配时缓存已编译的 regexp 对象
