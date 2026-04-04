package command

import (
	"context"
	"sync"
)

type Handler interface {
	Handle(ctx context.Context, cmd Command) error
	CanHandle(cmd Command) bool
}

type CommandRegistry struct {
	mu       sync.RWMutex
	commands map[string]Command
	handlers []Handler
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]Command),
	}
}

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

func (r *CommandRegistry) RegisterHandler(h Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers = append(r.handlers, h)
}

func (r *CommandRegistry) Get(name string) Command {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.commands[name]
}

func (r *CommandRegistry) Execute(ctx context.Context, name string) error {
	cmd := r.Get(name)
	if cmd == nil {
		return ErrUnknownCommand
	}
	return cmd.Execute(ctx)
}

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
