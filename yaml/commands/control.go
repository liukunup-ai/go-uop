package commands

import (
	"fmt"
	"strings"

	"github.com/liukunup/go-uop/core"
	"github.com/liukunup/go-uop/yaml"
)

type Executor struct {
	device    core.Device
	context   *yaml.Context
	variables map[string]interface{}
}

func NewExecutor(device core.Device, vars map[string]interface{}) *Executor {
	return &Executor{
		device:    device,
		context:   yaml.NewContext(),
		variables: vars,
	}
}

func (e *Executor) ExecuteCommand(cmd yaml.Command) error {
	if cmd.If != nil {
		return e.executeIf(*cmd.If)
	}
	if cmd.Foreach != nil {
		return e.executeForeach(*cmd.Foreach)
	}
	if cmd.While != nil {
		return e.executeWhile(*cmd.While)
	}
	return e.executeBasicCommand(cmd)
}

func (e *Executor) ExecuteCommands(cmds []yaml.Command) error {
	for i := range cmds {
		if err := e.ExecuteCommand(cmds[i]); err != nil {
			return fmt.Errorf("step %d: %w", i, err)
		}
	}
	return nil
}

func (e *Executor) executeIf(cmd yaml.IfCommand) error {
	cond, err := e.context.Evaluate(cmd.Condition)
	if err != nil {
		return fmt.Errorf("evaluate condition: %w", err)
	}

	if e.isTrue(cond) {
		return e.ExecuteCommands(cmd.Then)
	}
	if len(cmd.Else) > 0 {
		return e.ExecuteCommands(cmd.Else)
	}
	return nil
}

func (e *Executor) executeForeach(cmd yaml.ForeachCommand) error {
	listStr, err := e.context.Evaluate(cmd.In)
	if err != nil {
		return fmt.Errorf("evaluate in: %w", err)
	}

	items := strings.Split(listStr, ",")
	for _, item := range items {
		item = strings.TrimSpace(item)
		e.variables[cmd.Variable] = item
		if err := e.ExecuteCommands(cmd.Do); err != nil {
			return fmt.Errorf("foreach item %s: %w", item, err)
		}
	}
	return nil
}

func (e *Executor) executeWhile(cmd yaml.WhileCommand) error {
	maxIter := cmd.MaxIter
	if maxIter <= 0 {
		maxIter = 100
	}

	for i := 0; i < maxIter; i++ {
		cond, err := e.context.Evaluate(cmd.Condition)
		if err != nil {
			return fmt.Errorf("evaluate condition: %w", err)
		}

		if !e.isTrue(cond) {
			break
		}

		if err := e.ExecuteCommands(cmd.Do); err != nil {
			return fmt.Errorf("while iteration %d: %w", i, err)
		}
	}
	return nil
}

func (e *Executor) executeBasicCommand(cmd yaml.Command) error {
	if cmd.Launch != "" {
		return e.device.Launch()
	}
	if cmd.Wait > 0 {
		return nil
	}
	return nil
}

func (e *Executor) isTrue(cond string) bool {
	cond = strings.ToLower(strings.TrimSpace(cond))
	return cond == "true" || cond == "1" || cond == "yes" || cond == "y"
}
