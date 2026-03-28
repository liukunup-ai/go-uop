// Package serial provides serial port communication with event-driven monitoring.
//
// Basic usage:
//
//	s, _ := serial.NewSerial(&serial.Config{Name: "/dev/ttyUSB0", Baud: 115200})
//
//	m := serial.NewMonitor()
//	m.AddRule("OK", serial.MatchOnce, func(e serial.Event) {
//	    fmt.Println("Device ready")
//	})
//	m.AddRateLimitedRule("ERROR", 5*time.Second, serial.FileHandler{Path: "/tmp/errors.log"}.Handle)
//
//	s.AddObserver(m)
//
// Command table usage:
//
//	ct := serial.NewCommandTable()
//	ct.LoadFromFile("commands.yaml")
//
//	s, _ := serial.NewSerial(&serial.Config{
//	    Name: "/dev/ttyUSB0",
//	    Baud: 115200,
//	    Commands: ct,
//	})
//
//	s.SendByID("reset", func(result *serial.SendResult) {
//	    if result.Success {
//	        fmt.Println("Command executed")
//	    }
//	})
package serial

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/liukunup/go-uop/core"
	"github.com/tarm/serial"
)

// Config holds serial port configuration.
type Config struct {
	Name        string
	Baud        int
	ReadTimeout time.Duration
	Size        byte // data bits, default 8
	Parity      Parity
	StopBits    StopBits
	Commands    *CommandTable // 可选，命令表
}

// Parity parity mode.
type Parity byte

const (
	ParityNone  Parity = 'N' // no parity
	ParityOdd   Parity = 'O' // odd parity
	ParityEven  Parity = 'E' // even parity
	ParityMark  Parity = 'M' // mark parity (always 1)
	ParitySpace Parity = 'S' // space parity (always 0)
)

// StopBits stop bits.
type StopBits byte

const (
	Stop1     StopBits = 1  // 1 stop bit
	Stop1Half StopBits = 15 // 1.5 stop bits
	Stop2     StopBits = 2  // 2 stop bits
)

// Event 事件结构
type Event struct {
	Data      []byte
	Timestamp time.Time
	Rule      *Rule // nil if no rule matched
}

// Observer 观察者接口
type Observer interface {
	OnData(event Event)
	OnError(err error)
	OnClose()
}

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

// Serial is a serial port connection.
type Serial struct {
	mu        sync.RWMutex
	observers []Observer
	monitor   *Monitor // 内置监视器用于命令回显校验
	eventCh   chan Event
	done      chan struct{}
	cfg       *Config
	port      *serial.Port
	readErr   error
}

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

// NewSerial opens a serial port with the given config.
func NewSerial(cfg *Config) (*Serial, error) {
	if cfg.Name == "" {
		return nil, fmt.Errorf("serial port name is required")
	}
	if cfg.Baud <= 0 {
		cfg.Baud = 115200
	}
	if cfg.Size == 0 {
		cfg.Size = 8
	}
	if cfg.Parity == 0 {
		cfg.Parity = ParityNone
	}
	if cfg.StopBits == 0 {
		cfg.StopBits = Stop1
	}

	c := &serial.Config{
		Name:        cfg.Name,
		Baud:        cfg.Baud,
		ReadTimeout: cfg.ReadTimeout,
		Size:        cfg.Size,
		Parity:      serial.Parity(cfg.Parity),
		StopBits:    serial.StopBits(cfg.StopBits),
	}

	port, err := serial.OpenPort(c)
	if err != nil {
		return nil, fmt.Errorf("open serial port %s: %w", cfg.Name, err)
	}

	s := &Serial{
		cfg:     cfg,
		port:    port,
		monitor: NewMonitor(),
		eventCh: make(chan Event, 100),
		done:    make(chan struct{}),
	}

	s.observers = append(s.observers, s.monitor)

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
			}
		}
	}
}

// Read reads data from the serial port.
func (s *Serial) Read(p []byte) (int, error) {
	return s.port.Read(p)
}

// Write writes data to the serial port.
func (s *Serial) Write(p []byte) (int, error) {
	return s.port.Write(p)
}

// Close closes the serial port.
func (s *Serial) Close() error {
	close(s.done)
	s.notifyClose()
	return s.port.Close()
}

// Platform returns the device platform
func (s *Serial) Platform() core.Platform {
	return core.Serial
}

// Info returns device information
func (s *Serial) Info() (map[string]interface{}, error) {
	return map[string]interface{}{
		"platform": "serial",
		"port":     s.cfg.Name,
		"baud":     s.cfg.Baud,
	}, nil
}

// Screenshot captures current screen (not supported for serial)
func (s *Serial) Screenshot() ([]byte, error) {
	return nil, fmt.Errorf("screenshot not supported for serial device")
}

// Tap performs tap at coordinates (not supported for serial)
func (s *Serial) Tap(x, y int) error {
	return fmt.Errorf("tap not supported for serial device")
}

// SendKeys inputs text (not supported for serial)
func (s *Serial) SendKeys(text string) error {
	return fmt.Errorf("sendkeys not supported for serial device")
}

// Launch launches the app (not supported for serial)
func (s *Serial) Launch() error {
	return fmt.Errorf("launch not supported for serial device")
}

// SendCommand sends a command by DefaultName
func (s *Serial) SendCommand(name string, args ...interface{}) (interface{}, error) {
	ct := s.cfg.Commands
	if ct == nil {
		return nil, fmt.Errorf("command table not configured")
	}

	cmd, ok := ct.GetByDefaultName(name)
	if !ok {
		return nil, fmt.Errorf("command not found: %s", name)
	}

	// TODO: 参数化替换（后续实现）
	_ = args

	var result *SendResult
	err := s.sendCommand(cmd, func(r *SendResult) {
		result = r
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// SendByID 通过 ID 发送命令
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

	if cmd.Log != "" {
		cleanup = s.monitor.AddTemporaryRule(cmd.Log, MatchOnce, func(e Event) {
			result.Success = true
			result.Matched = true
			result.Echo = e.Data
			if callback != nil {
				callback(result)
			}
		})

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

	_, err := s.port.Write([]byte(cmd.Command))
	if err != nil {
		if cleanup != nil {
			cleanup()
		}
		return fmt.Errorf("send command: %w", err)
	}

	if cmd.Log == "" {
		result.Success = true
		if callback != nil {
			callback(result)
		}
	}

	return nil
}

// ReadByte reads a single byte.
func (s *Serial) ReadByte() (byte, error) {
	var b [1]byte
	_, err := s.port.Read(b[:])
	return b[0], err
}

// WriteString writes a string to the serial port.
func (s *Serial) WriteString(str string) (int, error) {
	return s.port.Write([]byte(str))
}

// Config returns the serial port configuration.
func (s *Serial) Config() *Config {
	return s.cfg
}

// ReadWithTimeout reads data with a specific timeout.
func (s *Serial) ReadWithTimeout(p []byte, timeout time.Duration) (int, error) {
	// Note: tarm/serial doesn't support per-read timeout,
	// this is a placeholder for future enhancement
	return s.port.Read(p)
}

// WriteAndRead writes data and then reads response.
func (s *Serial) WriteAndRead(w []byte, r []byte, timeout time.Duration) (int, error) {
	_, err := s.port.Write(w)
	if err != nil {
		return 0, err
	}
	return s.port.Read(r)
}

// Flush discards data written but not transmitted.
// This is a no-op for the standard tarm/serial port.
func (s *Serial) Flush() error {
	return nil
}

// ReaderFrom interface for parsing.
type ReaderFrom interface {
	ReadFrom(r io.Reader) (n int64, err error)
}
