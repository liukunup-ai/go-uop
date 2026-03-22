package maestro

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/liukunup/go-uop/core"
	"github.com/liukunup/go-uop/internal/action"
)

// Executor runs translated actions on a device
type Executor struct {
	device    core.Device
	outputDir string
}

// NewExecutor creates a new executor
func NewExecutor(device core.Device, outputDir string) *Executor {
	return &Executor{
		device:    device,
		outputDir: outputDir,
	}
}

// Execute runs all actions and reports progress
func (e *Executor) Execute(actions []action.Action, flowName string) error {
	if len(actions) == 0 {
		return nil
	}

	for i, act := range actions {
		stepNum := i + 1
		totalSteps := len(actions)
		stepName := getActionName(act)

		fmt.Printf("[STEP %d/%d] %s\n", stepNum, totalSteps, stepName)

		if err := act.Do(); err != nil {
			return fmt.Errorf("step %d failed: %w", stepNum, err)
		}
	}

	return nil
}

// ExecuteWithScreenshots runs actions and captures screenshots on failure
func (e *Executor) ExecuteWithScreenshots(actions []action.Action, flowName string) error {
	if len(actions) == 0 {
		return nil
	}

	for i, act := range actions {
		stepNum := i + 1
		totalSteps := len(actions)
		stepName := getActionName(act)

		fmt.Printf("[STEP %d/%d] %s\n", stepNum, totalSteps, stepName)

		if err := act.Do(); err != nil {
			e.captureScreenshot(flowName, stepNum)
			return fmt.Errorf("step %d failed: %w", stepNum, err)
		}
	}

	return nil
}

func (e *Executor) captureScreenshot(flowName string, stepNum int) {
	if e.outputDir == "" || e.device == nil {
		return
	}

	screenshot, err := e.device.Screenshot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to capture screenshot: %v\n", err)
		return
	}

	filename := fmt.Sprintf("%s_step_%d.png", flowName, stepNum)
	path := filepath.Join(e.outputDir, filename)

	if err := os.WriteFile(path, screenshot, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save screenshot: %v\n", err)
		return
	}

	fmt.Printf("Screenshot saved: %s\n", path)
}

func getActionName(act action.Action) string {
	switch a := act.(type) {
	case *LaunchWrapper:
		if a.waitIdle {
			return fmt.Sprintf("launch: %s", a.appID)
		}
		return fmt.Sprintf("terminate: %s", a.appID)
	case *TapWrapper:
		return fmt.Sprintf("tap: (%d, %d)", a.x, a.y)
	case *TapOnWrapper:
		if a.element != nil {
			return fmt.Sprintf("tapOn: %q", a.element.Value)
		}
		return "tapOn"
	case *SendKeysWrapper:
		return fmt.Sprintf("inputText: %q", a.text)
	case *SwipeWrapper:
		return fmt.Sprintf("swipe: (%d,%d)->(%d,%d)", a.startX, a.startY, a.endX, a.endY)
	case *WaitWrapper:
		return fmt.Sprintf("wait: %v", a.duration)
	default:
		return "unknown action"
	}
}
