// Package serial provides serial port communication with event-driven monitoring.
//
// Example usage:
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
package serial

import (
	"fmt"
	"io"
	"sync"
	"time"

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
		eventCh: make(chan Event, 100),
		done:    make(chan struct{}),
	}

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
