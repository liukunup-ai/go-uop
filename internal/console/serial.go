package console

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/liukunup/go-uop/pkg/serial"
)

type SerialConfig struct {
	Name     string        `json:"name"`
	Baud     int           `json:"baud"`
	DataBits int           `json:"dataBits"`
	Parity   string        `json:"parity"`
	StopBits int           `json:"stopBits"`
	Timeout  time.Duration `json:"timeout"`
}

type SerialConnection struct {
	ID       string               `json:"id"`
	Config   *SerialConfig        `json:"config"`
	Port     *serial.Serial       `json:"-"`
	Status   string               `json:"status"`
	Commands *serial.CommandTable `json:"commands,omitempty"`
}

type SerialManager struct {
	mu          sync.RWMutex
	connections map[string]*SerialConnection
}

type SerialEvent struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Data      string `json:"data"`
	Direction string `json:"direction"`
}

func NewSerialManager() *SerialManager {
	return &SerialManager{
		connections: make(map[string]*SerialConnection),
	}
}

func (m *SerialManager) ListPorts() ([]*SerialPortInfo, error) {
	entries, err := os.ReadDir("/dev")
	if err != nil {
		return nil, err
	}
	ports := []*SerialPortInfo{}
	for _, e := range entries {
		if len(e.Name()) >= 5 && e.Name()[:5] == "tty." {
			ports = append(ports, &SerialPortInfo{Name: "/dev/" + e.Name()})
		}
	}
	return ports, nil
}

func (m *SerialManager) Connect(cfg *SerialConfig) (*SerialConnection, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cfg.Name == "" {
		return nil, fmt.Errorf("port name is required")
	}
	if cfg.Baud <= 0 {
		cfg.Baud = 115200
	}
	if cfg.DataBits <= 0 {
		cfg.DataBits = 8
	}
	if cfg.Parity == "" {
		cfg.Parity = "N"
	}
	if cfg.StopBits <= 0 {
		cfg.StopBits = 1
	}

	var parity serial.Parity
	switch cfg.Parity {
	case "N":
		parity = serial.ParityNone
	case "O":
		parity = serial.ParityOdd
	case "E":
		parity = serial.ParityEven
	case "M":
		parity = serial.ParityMark
	case "S":
		parity = serial.ParitySpace
	default:
		parity = serial.ParityNone
	}

	var stopBits serial.StopBits
	switch cfg.StopBits {
	case 1:
		stopBits = serial.Stop1
	case 2:
		stopBits = serial.Stop2
	case 3:
		stopBits = serial.Stop1Half
	default:
		stopBits = serial.Stop1
	}

	sCfg := &serial.Config{
		Name:        cfg.Name,
		Baud:        cfg.Baud,
		ReadTimeout: cfg.Timeout,
		Size:        byte(cfg.DataBits),
		Parity:      parity,
		StopBits:    stopBits,
	}

	s, err := serial.NewSerial(sCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to open serial port: %w", err)
	}

	conn := &SerialConnection{
		ID:     fmt.Sprintf("serial-%d", time.Now().UnixNano()),
		Config: cfg,
		Port:   s,
		Status: "open",
	}

	m.connections[conn.ID] = conn
	return conn, nil
}

func (m *SerialManager) Disconnect(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	conn, ok := m.connections[id]
	if !ok {
		return fmt.Errorf("connection not found")
	}

	if err := conn.Port.Close(); err != nil {
		return fmt.Errorf("failed to close port: %w", err)
	}

	conn.Status = "closed"
	delete(m.connections, id)
	return nil
}

func (m *SerialManager) GetConnection(id string) (*SerialConnection, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	conn, ok := m.connections[id]
	if !ok {
		return nil, fmt.Errorf("connection not found")
	}
	return conn, nil
}

func (m *SerialManager) SendRaw(id string, data string) (*SerialCommandResult, error) {
	m.mu.RLock()
	conn, ok := m.connections[id]
	m.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("connection not found")
	}

	if conn.Status != "open" {
		return nil, fmt.Errorf("connection is not open")
	}

	_, err := conn.Port.WriteString(data)
	if err != nil {
		return nil, fmt.Errorf("write failed: %w", err)
	}

	return &SerialCommandResult{
		ID:        id,
		Success:   true,
		Sent:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

func (m *SerialManager) SendByCommandID(id string, cmdID string) (*SerialCommandResult, error) {
	m.mu.RLock()
	conn, ok := m.connections[id]
	m.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("connection not found")
	}

	if conn.Commands == nil {
		return nil, fmt.Errorf("no command table loaded")
	}

	var result *serial.SendResult
	err := conn.Port.SendByID(cmdID, func(r *serial.SendResult) {
		result = r
	})

	if err != nil {
		return nil, fmt.Errorf("send failed: %w", err)
	}

	return &SerialCommandResult{
		ID:        id,
		Success:   result.Success,
		Sent:      string(result.Echo),
		Matched:   result.Matched,
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

func (m *SerialManager) LoadCommandTable(id string, yamlContent string) error {
	m.mu.RLock()
	conn, ok := m.connections[id]
	m.mu.RUnlock()

	if !ok {
		return fmt.Errorf("connection not found")
	}

	tmpFile, err := os.CreateTemp("", "commands-*.yaml")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(yamlContent); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}
	tmpFile.Close()

	ct := serial.NewCommandTable()
	if err := ct.LoadFromFile(tmpFile.Name()); err != nil {
		return fmt.Errorf("load command table: %w", err)
	}

	m.mu.Lock()
	conn.Commands = ct
	m.mu.Unlock()

	return nil
}

func (m *SerialManager) LoadCommandTableFromFile(id string, filePath string) error {
	m.mu.RLock()
	conn, ok := m.connections[id]
	m.mu.RUnlock()

	if !ok {
		return fmt.Errorf("connection not found")
	}

	ct := serial.NewCommandTable()
	if err := ct.LoadFromFile(filePath); err != nil {
		return fmt.Errorf("load command table: %w", err)
	}

	m.mu.Lock()
	conn.Commands = ct
	m.mu.Unlock()

	return nil
}

func (m *SerialManager) ListCommands(id string) ([]*serial.Command, error) {
	m.mu.RLock()
	conn, ok := m.connections[id]
	m.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("connection not found")
	}

	if conn.Commands == nil {
		return []*serial.Command{}, nil
	}

	return conn.Commands.List(), nil
}

type SerialPortInfo struct {
	Name string `json:"name"`
	Desc string `json:"description,omitempty"`
}

type SerialCommandResult struct {
	ID        string `json:"id"`
	Success   bool   `json:"success"`
	Sent      string `json:"sent"`
	Matched   bool   `json:"matched,omitempty"`
	Timestamp string `json:"timestamp"`
}
