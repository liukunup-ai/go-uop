package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/liukunup/go-uop/core"
	"github.com/liukunup/go-uop/internal/action"
	"github.com/liukunup/go-uop/maestro"
)

const maxFlowDepth = 2

type FlowCommandTranslator struct{}

func NewFlowCommandTranslator() *FlowCommandTranslator {
	return &FlowCommandTranslator{}
}

func (t *FlowCommandTranslator) TranslateRunFlow(cmd *maestro.RunFlowCommand, parentVars map[string]string, currentDepth int) (*action.RunFlowAction, error) {
	if cmd == nil {
		return nil, fmt.Errorf("runFlow command is nil")
	}

	if cmd.Path == "" && cmd.Name == "" {
		return nil, fmt.Errorf("runFlow requires either path or name")
	}

	subflowPath := cmd.Path
	if subflowPath == "" {
		subflowPath = t.findSubflowByName(cmd.Name)
		if subflowPath == "" {
			return nil, fmt.Errorf("subflow not found: %s", cmd.Name)
		}
	}

	if currentDepth >= maxFlowDepth {
		return nil, fmt.Errorf("max flow depth (%d) exceeded: prevents infinite recursion", maxFlowDepth)
	}

	envVars := make(map[string]string)
	for k, v := range parentVars {
		envVars[k] = v
	}
	for k, v := range cmd.Params {
		envVars[k] = v
	}

	return &action.RunFlowAction{
		SubflowPath: subflowPath,
		EnvVars:     envVars,
		Depth:       currentDepth + 1,
	}, nil
}

func (t *FlowCommandTranslator) findSubflowByName(name string) string {
	extensions := []string{".yaml", ".yml", ".maestro.yaml"}
	for _, ext := range extensions {
		path := name + ext
		if _, err := os.Stat(path); err == nil {
			if absPath, err := filepath.Abs(path); err == nil {
				return absPath
			}
			return path
		}
	}
	return ""
}

func (t *FlowCommandTranslator) ExecuteRunFlow(flowAction *action.RunFlowAction, parentVars map[string]string, device core.Device) error {
	if flowAction == nil {
		return fmt.Errorf("flowAction is nil")
	}

	flowData, err := os.ReadFile(flowAction.SubflowPath)
	if err != nil {
		return fmt.Errorf("read subflow file: %w", err)
	}

	subflow, err := maestro.ParseFlowFromString(string(flowData))
	if err != nil {
		return fmt.Errorf("parse subflow: %w", err)
	}

	subflowVars := make(map[string]string)
	for k, v := range parentVars {
		subflowVars[k] = v
	}
	for k, v := range flowAction.EnvVars {
		subflowVars[k] = v
	}

	translator := &maestro.Translator{}
	actions, err := translator.TranslateFlow(subflow, device)
	if err != nil {
		return fmt.Errorf("translate subflow: %w", err)
	}

	for _, act := range actions {
		if runFlowAct, ok := act.(*action.RunFlowAction); ok {
			nestedAct, err := t.TranslateRunFlow(&maestro.RunFlowCommand{
				Path:   runFlowAct.SubflowPath,
				Params: runFlowAct.EnvVars,
			}, subflowVars, flowAction.Depth)
			if err != nil {
				return fmt.Errorf("nested runFlow: %w", err)
			}
			if err := t.ExecuteRunFlow(nestedAct, subflowVars, device); err != nil {
				return err
			}
			continue
		}
		if err := act.Do(); err != nil {
			return fmt.Errorf("execute action: %w", err)
		}
	}

	return nil
}

func (t *FlowCommandTranslator) ResolveSubflowPath(path string, baseDir string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(baseDir, path)
}
