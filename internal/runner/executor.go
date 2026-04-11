package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/liukunup/go-uop/core"
	"github.com/liukunup/go-uop/internal/report"
	"github.com/liukunup/go-uop/internal/watcher"
)

type CommandExecutor func(device core.Device, args map[string]any) error

type Executor struct {
	pool      *DevicePool
	reportGen *report.Generator
	executors map[string]CommandExecutor
	watcher   *watcher.WatcherEngine
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

func (e *Executor) WithWatcher(w *watcher.WatcherEngine) *Executor {
	e.watcher = w
	watcher.CommandExecutor = func(name string, args map[string]any, device core.Device) error {
		executor, exists := e.executors[name]
		if !exists {
			return fmt.Errorf("unknown command: %s", name)
		}
		return executor(device, args)
	}
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
		if content, ok := args["content"]; ok {
			textStr := fmt.Sprintf("%v", content)
			return device.SendKeys(textStr)
		}
		return fmt.Errorf("inputText: missing 'text' or 'content' argument")
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

func (e *Executor) ExecuteSuite(suite *TestSuite) error {
	for i, tc := range suite.TestCases {
		if err := e.ExecuteTestCase(i, tc); err != nil {
			return fmt.Errorf("testcase %d failed: %w", i, err)
		}
	}
	return nil
}

func (e *Executor) ExecuteTestCase(index int, tc TestCase) error {
	e.reportGen.StartTest(tc.Name)

	for i, step := range tc.Steps {
		stepName := fmt.Sprintf("%s-%d", step.Type, i)
		if err := e.ExecuteStep(stepName, step); err != nil {
			e.reportGen.EndTest("failed", err)
			return fmt.Errorf("step %d failed: %w", i, err)
		}
	}

	e.reportGen.EndTest("passed", nil)
	return nil
}

func (e *Executor) ExecuteStep(stepName string, step Step) error {
	executor, exists := e.executors[step.Type]
	if !exists {
		return fmt.Errorf("unknown step: %s", step.Type)
	}

	start := time.Now()

	var device core.Device
	if dev := e.pool.CurrentDevice(); dev != nil {
		device = dev.device
	}

	var err error
	if device != nil {
		err = executor(device, step.Params)
	}

	duration := time.Since(start)
	status := "passed"
	if err != nil {
		status = "failed"
	}
	e.reportGen.AddStep(stepName, duration, status, err)

	if err != nil {
		return err
	}

	if e.watcher != nil && e.watcher.Enabled() {
		if dev := e.pool.CurrentDevice(); dev != nil && dev.device != nil {
			if watcherErr := e.watcher.Check(context.Background(), dev.device); watcherErr != nil {
				fmt.Printf("watcher warning: %v\n", watcherErr)
			}
		}
	}

	return nil
}
