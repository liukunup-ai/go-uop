package serial

import (
	"fmt"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// Command 串口命令
type Command struct {
	ID      string        `yaml:"id"`      // 唯一标识符（也用于 SendCommand 查找）
	Name    string        `yaml:"name"`    // 易读名称
	Command string        `yaml:"command"` // 发送字节序列（字符串格式）
	Log     string        `yaml:"log"`     // 回显校验正则（可选）
	Timeout time.Duration `yaml:"timeout"` // 超时时间
}

// CommandTable 命令表
type CommandTable struct {
	mu       sync.RWMutex
	commands map[string]*Command
	byName   map[string]*Command
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
