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
