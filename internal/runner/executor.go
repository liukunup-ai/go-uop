package runner

import (
	"fmt"
	"time"

	"github.com/liukunup/go-uop/core"
	"github.com/liukunup/go-uop/internal/report"
)

type CommandExecutor func(device core.Device, args map[string]any) error

type Executor struct {
	pool      *DevicePool
	reportGen *report.Generator
	executors map[string]CommandExecutor
}

func NewExecutor(pool *DevicePool, reportGen *report.Generator) *Executor {
	e := &Executor{
		pool:      pool,
		reportGen: reportGen,
		executors: make(map[string]CommandExecutor),
	}
	e.registerCommands()
	return e
}

func (e *Executor) registerCommands() {
	e.executors["launch"] = func(device core.Device, args map[string]any) error {
		return device.Launch()
	}

	e.executors["tapOn"] = func(device core.Device, args map[string]any) error {
		x := 0
		y := 0
		if xv, ok := args["x"]; ok {
			switch v := xv.(type) {
			case int:
				x = v
			case float64:
				x = int(v)
			}
		}
		if yv, ok := args["y"]; ok {
			switch v := yv.(type) {
			case int:
				y = v
			case float64:
				y = int(v)
			}
		}
		return device.Tap(x, y)
	}

	e.executors["inputText"] = func(device core.Device, args map[string]any) error {
		if text, ok := args["text"]; ok {
			textStr := fmt.Sprintf("%v", text)
			return device.SendKeys(textStr)
		}
		return fmt.Errorf("inputText: missing 'text' argument")
	}

	e.executors["swipe"] = func(device core.Device, args map[string]any) error {
		return nil
	}

	e.executors["pressKey"] = func(device core.Device, args map[string]any) error {
		return nil
	}

	e.executors["wait"] = func(device core.Device, args map[string]any) error {
		ms := 1000
		if msv, ok := args["ms"].(float64); ok {
			ms = int(msv)
		}
		time.Sleep(time.Duration(ms) * time.Millisecond)
		return nil
	}

	e.executors["screenshot"] = func(device core.Device, args map[string]any) error {
		_, err := device.Screenshot()
		return err
	}

	e.executors["device"] = func(device core.Device, args map[string]any) error {
		return nil
	}
}

func (e *Executor) ExecuteFlow(flow *Flow) error {
	for i, step := range flow.Steps {
		if err := e.executeStep(i, step); err != nil {
			return fmt.Errorf("step %d failed: %w", i, err)
		}
	}
	return nil
}

func (e *Executor) executeStep(index int, step Step) error {
	stepName := fmt.Sprintf("step-%d", index)

	if cmd, ok := step["device"]; ok {
		deviceID, ok := cmd.(string)
		if !ok {
			return fmt.Errorf("device argument must be a string")
		}
		if err := e.pool.SwitchDevice(deviceID); err != nil {
			return fmt.Errorf("failed to switch to device %s: %w", deviceID, err)
		}
		e.reportGen.AddStep(stepName, 0, "passed", nil)
		return nil
	}

	for cmdName, cmdArgs := range step {
		executor, exists := e.executors[cmdName]
		if !exists {
			return fmt.Errorf("unknown command: %s", cmdName)
		}

		var args map[string]any
		switch v := cmdArgs.(type) {
		case Step:
			args = v
		case map[string]any:
			args = v
		default:
			args = make(map[string]any)
		}

		start := time.Now()
		device := e.pool.CurrentDevice()
		if device == nil {
			return fmt.Errorf("no device available")
		}

		var err error
		if device.device != nil {
			err = executor(device.device, args)
		}

		duration := time.Since(start)
		status := "passed"
		if err != nil {
			status = "failed"
		}
		e.reportGen.AddStep(fmt.Sprintf("%s-%s", stepName, cmdName), duration, status, err)

		if err != nil {
			return err
		}
	}

	return nil
}
