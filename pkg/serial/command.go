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
	ID          string        `yaml:"id"`
	Name        string        `yaml:"name"`
	DefaultName string        `yaml:"default_name"`
	Command     string        `yaml:"command"`
	Log         string        `yaml:"log"`
	Timeout     time.Duration `yaml:"timeout"`
}

// CommandTable 命令表
type CommandTable struct {
	mu            sync.RWMutex
	commands      map[string]*Command
	byName        map[string]*Command
	byDefaultName map[string]*Command
}

// commandTableFile YAML 文件格式
type commandTableFile struct {
	Commands []*Command `yaml:"commands"`
}

// NewCommandTable 创建空命令表
func NewCommandTable() *CommandTable {
	return &CommandTable{
		commands:      make(map[string]*Command),
		byName:        make(map[string]*Command),
		byDefaultName: make(map[string]*Command),
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
		if cmd.DefaultName != "" {
			ct.byDefaultName[cmd.DefaultName] = cmd
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

func (ct *CommandTable) GetByDefaultName(name string) (*Command, bool) {
	ct.mu.RLock()
	defer ct.mu.RUnlock()
	cmd, ok := ct.byDefaultName[name]
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
