# UOP Framework 重构实现计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 将 `maestro/` 和 `yaml/` 包合并为统一的 `internal/runner/`，新增 `internal/reporter/` 支持多格式报告，CLI 支持 run/debug/test 子命令

**Architecture:** 
- 新 `runner/` 包：统一 YAML 解析器 + 引擎执行器 + 调试器 + 设备池管理
- 新 `reporter/` 包：基于现有 `internal/report/generator.go` 扩展 HTML/JUnit XML 支持
- CLI：`cmd/uop/main.go` (新建) 支持 run/debug/test 三个子命令

**Tech Stack:** Go 1.21+, gopkg.in/yaml.v3, go-junit-report, github.com/grokify/htmlopts

---

## 阶段 1: 创建 `internal/runner/` 包

### Task 1.1: 创建 `internal/runner/parser.go` - YAML 解析器

**Files:**
- Create: `internal/runner/parser.go`
- Test: `internal/runner/parser_test.go`

**Step 1: 写测试**

```go
package runner

import (
	"testing"
	"strings"
)

func TestParseFlow(t *testing.T) {
	yaml := `
name: 登录测试
devices:
  - id: iphone
    type: ios
    serial: 00001234-00123456789
  - id: android-tablet
    type: android
    serial: emulator-5554
defaultDevice: iphone
steps:
  - launch: com.example.app
  - tapOn: { text: "用户名" }
  - device: android-tablet
    tapOn: { text: "跳过" }
`
	r := strings.NewReader(yaml)
	flow, err := ParseFlow(r)
	if err != nil {
		t.Fatalf("ParseFlow failed: %v", err)
	}
	if flow.Name != "登录测试" {
		t.Errorf("expected name 登录测试, got %s", flow.Name)
	}
	if len(flow.Devices) != 2 {
		t.Errorf("expected 2 devices, got %d", len(flow.Devices))
	}
	if flow.DefaultDevice != "iphone" {
		t.Errorf("expected defaultDevice iphone, got %s", flow.DefaultDevice)
	}
	if len(flow.Steps) != 3 {
		t.Errorf("expected 3 steps, got %d", len(flow.Steps))
	}
}

func TestParseSuite(t *testing.T) {
	yaml := `
name: 回归测试套件
devices:
  - id: ios-device
    type: ios
    serial: 00001234-00123456789
defaultDevice: ios-device
flows:
  - name: 登录流程
    path: ./flows/login.yaml
`
	r := strings.NewReader(yaml)
	suite, err := ParseSuite(r)
	if err != nil {
		t.Fatalf("ParseSuite failed: %v", err)
	}
	if suite.Name != "回归测试套件" {
		t.Errorf("expected name 回归测试套件, got %s", suite.Name)
	}
	if len(suite.Flows) != 1 {
		t.Errorf("expected 1 flow, got %d", len(suite.Flows))
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/runner/... -run TestParse -v`
Expected: FAIL - "undefined: ParseFlow"

**Step 3: Write minimal implementation**

```go
package runner

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// Flow represents a test flow
type Flow struct {
	Name         string   `yaml:"name"`
	Devices      []Device `yaml:"devices"`
	DefaultDevice string   `yaml:"defaultDevice"`
	Steps        []Step   `yaml:"steps"`
}

// Suite represents a test suite
type Suite struct {
	Name         string      `yaml:"name"`
	Devices      []Device    `yaml:"devices"`
	DefaultDevice string     `yaml:"defaultDevice"`
	Flows        []SuiteFlow `yaml:"flows"`
}

// SuiteFlow references an external flow file
type SuiteFlow struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

// Device represents a target device
type Device struct {
	ID     string `yaml:"id"`
	Type   string `yaml:"type"`
	Serial string `yaml:"serial"`
}

// Step represents a single command step
type Step map[string]interface{}

// ParseFlow parses a YAML flow definition
func ParseFlow(r io.Reader) (*Flow, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read flow: %w", err)
	}

	var flow Flow
	if err := yaml.Unmarshal(data, &flow); err != nil {
		return nil, fmt.Errorf("parse flow: %w", err)
	}

	return &flow, nil
}

// ParseSuite parses a YAML suite definition
func ParseSuite(r io.Reader) (*Suite, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read suite: %w", err)
	}

	var suite Suite
	if err := yaml.Unmarshal(data, &suite); err != nil {
		return nil, fmt.Errorf("parse suite: %w", err)
	}

	return &suite, nil
}

// ParseFlowFile parses a flow from a file path
func ParseFlowFile(path string) (*Flow, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open flow file: %w", err)
	}
	defer f.Close()
	return ParseFlow(f)
}

// ParseSuiteFile parses a suite from a file path
func ParseSuiteFile(path string) (*Suite, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open suite file: %w", err)
	}
	defer f.Close()
	return ParseSuite(f)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/runner/... -run TestParse -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/runner/parser.go internal/runner/parser_test.go
git commit -m "feat(runner): add YAML parser for flows and suites"
```

---

### Task 1.2: 创建 `internal/runner/device.go` - 设备池管理

**Files:**
- Create: `internal/runner/device.go`
- Test: `internal/runner/device_test.go`

**Step 1: 写测试**

```go
package runner

import (
	"testing"
)

func TestDevicePool(t *testing.T) {
	pool := NewDevicePool()

	// Add devices
	err := pool.AddDevice("iphone", "ios", "00001234-00123456789")
	if err != nil {
		t.Fatalf("AddDevice failed: %v", err)
	}
	err = pool.AddDevice("android-tablet", "android", "emulator-5554")
	if err != nil {
		t.Fatalf("AddDevice failed: %v", err)
	}

	// Check default device (first added becomes default)
	defaultDev := pool.DefaultDevice()
	if defaultDev == nil {
		t.Fatal("expected default device, got nil")
	}
	if defaultDev.ID != "iphone" {
		t.Errorf("expected default iphone, got %s", defaultDev.ID)
	}

	// Switch device
	err = pool.SwitchDevice("android-tablet")
	if err != nil {
		t.Fatalf("SwitchDevice failed: %v", err)
	}
	current := pool.CurrentDevice()
	if current.ID != "android-tablet" {
		t.Errorf("expected current android-tablet, got %s", current.ID)
	}

	// Get non-existent device
	_, err = pool.GetDevice("non-existent")
	if err == nil {
		t.Error("expected error for non-existent device")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/runner/... -run TestDevicePool -v`
Expected: FAIL - "undefined: NewDevicePool"

**Step 3: Write minimal implementation**

```go
package runner

import (
	"fmt"

	"github.com/liukunup/go-uop/core"
	"github.com/liukunup/go-uop/pkg/android"
	"github.com/liukunup/go-uop/pkg/ios"
)

// PoolDevice represents a device in the pool
type PoolDevice struct {
	ID     string
	Type   string
	Serial string
	device core.Device
}

// DevicePool manages multiple devices
type DevicePool struct {
	devices       map[string]*PoolDevice
	currentDevice string
	defaultDevice string
}

// NewDevicePool creates a new device pool
func NewDevicePool() *DevicePool {
	return &DevicePool{
		devices: make(map[string]*PoolDevice),
	}
}

// AddDevice adds a device to the pool
func (p *DevicePool) AddDevice(id, deviceType, serial string) error {
	if _, exists := p.devices[id]; exists {
		return fmt.Errorf("device %s already exists", id)
	}

	var device core.Device
	var err error

	switch deviceType {
	case "ios":
		device, err = ios.NewDevice("")
	case "android":
		device, err = android.NewDevice(android.WithSerial(serial))
	default:
		return fmt.Errorf("unsupported device type: %s", deviceType)
	}
	if err != nil {
		return fmt.Errorf("create device: %w", err)
	}

	p.devices[id] = &PoolDevice{
		ID:     id,
		Type:   deviceType,
		Serial: serial,
		device: device,
	}

	// First device becomes default
	if p.defaultDevice == "" {
		p.defaultDevice = id
		p.currentDevice = id
	}

	return nil
}

// GetDevice returns a device by ID
func (p *DevicePool) GetDevice(id string) (*PoolDevice, error) {
	d, ok := p.devices[id]
	if !ok {
		return nil, fmt.Errorf("device %s not found", id)
	}
	return d, nil
}

// CurrentDevice returns the currently selected device
func (p *DevicePool) CurrentDevice() *PoolDevice {
	if p.currentDevice == "" {
		return nil
	}
	return p.devices[p.currentDevice]
}

// DefaultDevice returns the default device
func (p *DevicePool) DefaultDevice() *PoolDevice {
	if p.defaultDevice == "" {
		return nil
	}
	return p.devices[p.defaultDevice]
}

// SwitchDevice switches to a different device
func (p *DevicePool) SwitchDevice(id string) error {
	if _, ok := p.devices[id]; !ok {
		return fmt.Errorf("device %s not found", id)
	}
	p.currentDevice = id
	return nil
}

// Close closes all devices in the pool
func (p *DevicePool) Close() error {
	for _, d := range p.devices {
		if d.device != nil {
			d.device.Close()
		}
	}
	return nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/runner/... -run TestDevicePool -v`
Expected: PASS (may skip if iOS/Android init fails in test environment)

**Step 5: Commit**

```bash
git add internal/runner/device.go internal/runner/device_test.go
git commit -m "feat(runner): add device pool management"
```

---

### Task 1.3: 创建 `internal/runner/executor.go` - 步骤执行器

**Files:**
- Create: `internal/runner/executor.go`
- Test: `internal/runner/executor_test.go`

**Step 1: 写测试**

```go
package runner

import (
	"testing"
	"strings"
)

func TestExecutor(t *testing.T) {
	yaml := `
name: 测试流程
devices:
  - id: test
    type: ios
    serial: test-serial
defaultDevice: test
steps:
  - launch: com.example.app
`
	flow, err := ParseFlow(strings.NewReader(yaml))
	if err != nil {
		t.Fatalf("ParseFlow failed: %v", err)
	}

	// Note: Full executor test requires device mock
	// This test verifies step parsing works correctly
	if len(flow.Steps) != 1 {
		t.Errorf("expected 1 step, got %d", len(flow.Steps))
	}
}
```

**Step 2: Run test to verify it passes**

Run: `go test ./internal/runner/... -run TestExecutor -v`
Expected: PASS

**Step 3: Write executor implementation (基础版)**

```go
package runner

import (
	"fmt"
	"reflect"
	"time"

	"github.com/liukunup/go-uop/core"
	"github.com/liukunup/go-uop/internal/report"
)

// CommandExecutor executes a single command
type CommandExecutor func(device core.Device, args map[string]interface{}) error

// Executor executes flow steps
type Executor struct {
	pool          *DevicePool
	reportGen     *report.Generator
	executors     map[string]CommandExecutor
}

// NewExecutor creates a new executor
func NewExecutor(pool *DevicePool, reportGen *report.Generator) *Executor {
	e := &Executor{
		pool:      pool,
		reportGen: reportGen,
		executors: make(map[string]CommandExecutor),
	}
	e.registerCommands()
	return e
}

// registerCommands registers all supported commands
func (e *Executor) registerCommands() {
	e.executors["launch"] = e.cmdLaunch
	e.executors["tapOn"] = e.cmdTapOn
	e.executors["inputText"] = e.cmdInputText
	e.executors["swipe"] = e.cmdSwipe
	e.executors["pressKey"] = e.cmdPressKey
	e.executors["wait"] = e.cmdWait
	e.executors["screenshot"] = e.cmdScreenshot
	e.executors["device"] = e.cmdDevice
}

// ExecuteFlow executes a flow
func (e *Executor) ExecuteFlow(flow *Flow) error {
	e.reportGen.StartTest(flow.Name)
	defer func() {
		e.reportGen.EndTest("completed", nil)
	}()

	for i, step := range flow.Steps {
		if err := e.executeStep(i, step); err != nil {
			e.reportGen.AddStep(fmt.Sprintf("step-%d", i), 0, "failed", err)
			return fmt.Errorf("step %d failed: %w", i, err)
		}
	}
	return nil
}

// executeStep executes a single step
func (e *Executor) executeStep(index int, step Step) error {
	start := time.Now()
	defer func() {
		e.reportGen.AddStep(fmt.Sprintf("step-%d", index), time.Since(start), "passed", nil)
	}()

	// Check for device switch
	if devID, ok := step["device"].(string); ok {
		if err := e.pool.SwitchDevice(devID); err != nil {
			return err
		}
		return nil
	}

	// Find and execute command
	for cmd, args := range step {
		executor, ok := e.executors[cmd]
		if !ok {
			return fmt.Errorf("unknown command: %s", cmd)
		}

		device := e.pool.CurrentDevice().device
		var argsMap map[string]interface{}
		if args != nil {
			argsMap = args.(map[string]interface{})
		}
		return executor(device, argsMap)
	}
	return nil
}

// Command implementations
func (e *Executor) cmdLaunch(device core.Device, args map[string]interface{}) error {
	if device == nil {
		return fmt.Errorf("no device selected")
	}
	return device.Launch()
}

func (e *Executor) cmdTapOn(device core.Device, args map[string]interface{}) error {
	if device == nil {
		return fmt.Errorf("no device selected")
	}
	// Simplified - would use selector to find coordinates
	x := int(args["x"].(float64))
	y := int(args["y"].(float64))
	return device.Tap(x, y)
}

func (e *Executor) cmdInputText(device core.Device, args map[string]interface{}) error {
	if device == nil {
		return fmt.Errorf("no device selected")
	}
	text := args["text"].(string)
	return device.SendKeys(text)
}

func (e *Executor) cmdSwipe(device core.Device, args map[string]interface{}) error {
	// Implementation for swipe
	return nil
}

func (e *Executor) cmdPressKey(device core.Device, args map[string]interface{}) error {
	// Implementation for pressKey
	return nil
}

func (e *Executor) cmdWait(device core.Device, args map[string]interface{}) error {
	ms := int(args["ms"].(float64))
	time.Sleep(time.Duration(ms) * time.Millisecond)
	return nil
}

func (e *Executor) cmdScreenshot(device core.Device, args map[string]interface{}) error {
	_, err := device.Screenshot()
	return err
}

func (e *Executor) cmdDevice(device core.Device, args map[string]interface{}) error {
	// Handled in executeStep
	return nil
}
```

**Step 4: Run test to verify it compiles**

Run: `go build ./internal/runner/...`
Expected: SUCCESS

**Step 5: Commit**

```bash
git add internal/runner/executor.go internal/runner/executor_test.go
git commit -m "feat(runner): add flow executor with command registry"
```

---

### Task 1.4: 创建 `internal/runner/debugger.go` - 调试器

**Files:**
- Create: `internal/runner/debugger.go`
- Test: `internal/runner/debugger_test.go`

**Step 1: 写测试**

```go
package runner

import (
	"testing"
	"strings"
)

func TestDebugger(t *testing.T) {
	yaml := `
name: 调试流程
devices:
  - id: test
    type: ios
    serial: test-serial
defaultDevice: test
steps:
  - launch: com.example.app
  - tapOn: { x: 100, y: 200 }
`
	flow, err := ParseFlow(strings.NewReader(yaml))
	if err != nil {
		t.Fatalf("ParseFlow failed: %v", err)
	}

	dbg := NewDebugger(flow)

	// Initial state
	if dbg.CurrentStep() != 0 {
		t.Errorf("expected initial step 0, got %d", dbg.CurrentStep())
	}
	if dbg.Status() != DebugStatusPending {
		t.Errorf("expected pending status, got %s", dbg.Status())
	}

	// Step through
	step, err := dbg.NextStep()
	if err != nil {
		t.Fatalf("NextStep failed: %v", err)
	}
	if step == nil {
		t.Fatal("expected step, got nil")
	}

	// Pause
	dbg.Pause()
	if dbg.Status() != DebugStatusPaused {
		t.Errorf("expected paused status, got %s", dbg.Status())
	}

	// Resume
	dbg.Resume()
	if dbg.Status() != DebugStatusRunning {
		t.Errorf("expected running status, got %s", dbg.Status())
	}

	// Skip
	dbg.Skip()
	if dbg.CurrentStep() != 2 {
		t.Errorf("expected step 2 after skip, got %d", dbg.CurrentStep())
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/runner/... -run TestDebugger -v`
Expected: FAIL - "undefined: NewDebugger"

**Step 3: Write implementation**

```go
package runner

import (
	"fmt"
	"sync"
)

// DebugStatus represents debugger state
type DebugStatus string

const (
	DebugStatusPending   DebugStatus = "pending"
	DebugStatusRunning   DebugStatus = "running"
	DebugStatusPaused    DebugStatus = "paused"
	DebugStatusCompleted DebugStatus = "completed"
	DebugStatusFailed    DebugStatus = "failed"
)

// Debugger provides step-by-step debugging
type Debugger struct {
	flow        *Flow
	currentStep int
	status      DebugStatus
	mu          sync.Mutex
	breakpoints map[int]bool
}

// NewDebugger creates a new debugger
func NewDebugger(flow *Flow) *Debugger {
	return &Debugger{
		flow:        flow,
		currentStep: 0,
		status:      DebugStatusPending,
		breakpoints: make(map[int]bool),
	}
}

// CurrentStep returns the current step index
func (d *Debugger) CurrentStep() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.currentStep
}

// Status returns the current debug status
func (d *Debugger) Status() DebugStatus {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.status
}

// NextStep returns the next step and advances
func (d *Debugger) NextStep() (*Step, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.currentStep >= len(d.flow.Steps) {
		d.status = DebugStatusCompleted
		return nil, nil
	}

	step := &d.flow.Steps[d.currentStep]
	d.currentStep++
	d.status = DebugStatusRunning
	return step, nil
}

// Pause pauses execution
func (d *Debugger) Pause() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.status == DebugStatusRunning {
		d.status = DebugStatusPaused
	}
}

// Resume resumes execution
func (d *Debugger) Resume() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.status == DebugStatusPaused {
		d.status = DebugStatusRunning
	}
}

// Skip skips the current step
func (d *Debugger) Skip() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.currentStep < len(d.flow.Steps) {
		d.currentStep++
	}
}

// SetBreakpoint sets a breakpoint at a step
func (d *Debugger) SetBreakpoint(step int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.breakpoints[step] = true
}

// ClearBreakpoint clears a breakpoint
func (d *Debugger) ClearBreakpoint(step int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.breakpoints, step)
}

// IsBreakpoint checks if a step has a breakpoint
func (d *Debugger) IsBreakpoint(step int) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.breakpoints[step]
}

// Reset resets the debugger to initial state
func (d *Debugger) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.currentStep = 0
	d.status = DebugStatusPending
}

// Flow returns the flow being debugged
func (d *Debugger) Flow() *Flow {
	return d.flow
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/runner/... -run TestDebugger -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/runner/debugger.go internal/runner/debugger_test.go
git commit -m "feat(runner): add debugger with step-through control"
```

---

### Task 1.5: 创建 `internal/runner/suite.go` - 测试套件执行

**Files:**
- Create: `internal/runner/suite.go`
- Test: `internal/runner/suite_test.go`

**Step 1: 写测试**

```go
package runner

import (
	"testing"
	"strings"
)

func TestSuiteRunner(t *testing.T) {
	yaml := `
name: 测试套件
devices:
  - id: test
    type: ios
    serial: test-serial
defaultDevice: test
flows:
  - name: 流程A
    path: ./flow-a.yaml
`
	suite, err := ParseSuite(strings.NewReader(yaml))
	if err != nil {
		t.Fatalf("ParseSuite failed: %v", err)
	}

	if suite.Name != "测试套件" {
		t.Errorf("expected name 测试套件, got %s", suite.Name)
	}
	if len(suite.Flows) != 1 {
		t.Errorf("expected 1 flow, got %d", len(suite.Flows))
	}
}
```

**Step 2: Run test to verify it passes**

Run: `go test ./internal/runner/... -run TestSuiteRunner -v`
Expected: PASS

**Step 3: Write suite runner**

```go
package runner

import (
	"fmt"
	"os"
	"path/filepath"
)

// SuiteRunner runs a test suite
type SuiteRunner struct {
	suite   *Suite
	pool    *DevicePool
	reportGen *report.Generator
}

// NewSuiteRunner creates a new suite runner
func NewSuiteRunner(suite *Suite, pool *DevicePool, reportGen *report.Generator) *SuiteRunner {
	return &SuiteRunner{
		suite:     suite,
		pool:      pool,
		reportGen: reportGen,
	}
}

// Run executes all flows in the suite
func (s *SuiteRunner) Run() error {
	// Get base directory for relative paths
	baseDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get cwd: %w", err)
	}

	executor := NewExecutor(s.pool, s.reportGen)

	for _, flowRef := range s.suite.Flows {
		flowPath := flowRef.Path
		// Handle relative paths
		if !filepath.IsAbs(flowPath) {
			flowPath = filepath.Join(baseDir, flowPath)
		}

		flow, err := ParseFlowFile(flowPath)
		if err != nil {
			return fmt.Errorf("parse flow %s: %w", flowRef.Name, err)
		}

		if err := executor.ExecuteFlow(flow); err != nil {
			return fmt.Errorf("execute flow %s: %w", flowRef.Name, err)
		}
	}
	return nil
}
```

**Step 4: Run test to verify it compiles**

Run: `go build ./internal/runner/...`
Expected: SUCCESS

**Step 5: Commit**

```bash
git add internal/runner/suite.go internal/runner/suite_test.go
git commit -m "feat(runner): add suite runner for batch execution"
```

---

## 阶段 2: 创建 `internal/reporter/` 包

### Task 2.1: 创建 `internal/reporter/html.go` - HTML 报告

**Files:**
- Create: `internal/reporter/html.go`
- Test: `internal/reporter/html_test.go`

**Step 1: 写测试**

```go
package reporter

import (
	"os"
	"testing"
)

func TestHTMLReporter(t *testing.T) {
	rep := NewHTMLReporter("测试报告")

	// Add test data
	rep.StartTest("登录测试")
	rep.AddStep("启动APP", "passed", nil)
	rep.AddStep("输入用户名", "passed", nil)
	rep.EndTest("passed", nil)

	// Generate HTML
	html, err := rep.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(html) == 0 {
		t.Error("expected non-empty HTML")
	}

	// Write to file
	tmpfile, err := os.CreateTemp("", "report-*.html")
	if err != nil {
		t.Fatalf("CreateTemp failed: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	err = rep.WriteFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/reporter/... -run TestHTMLReporter -v`
Expected: FAIL - "undefined: NewHTMLReporter"

**Step 3: Write implementation**

```go
package reporter

import (
	"fmt"
	"os"
	"time"

	"github.com/liukunup/go-uop/internal/report"
)

// HTMLReporter generates HTML reports
type HTMLReporter struct {
	generator *report.Generator
	title     string
}

// NewHTMLReporter creates a new HTML reporter
func NewHTMLReporter(title string) *HTMLReporter {
	return &HTMLReporter{
		generator: report.NewGenerator(title),
		title:     title,
	}
}

// StartTest starts a new test
func (r *HTMLReporter) StartTest(name string) {
	r.generator.StartTest(name)
}

// AddStep adds a step result
func (r *HTMLReporter) AddStep(name, status string, err error) {
	r.generator.AddStep(name, 0, status, err)
}

// EndTest ends a test
func (r *HTMLReporter) EndTest(status string, err error) {
	r.generator.EndTest(status, err)
}

// AddScreenshot adds a screenshot
func (r *HTMLReporter) AddScreenshot(path, name string) {
	r.generator.AddScreenshot(path, name)
}

// Generate generates HTML report
func (r *HTMLReporter) Generate() ([]byte, error) {
	suite := r.generator.Generate()
	return []byte(htmlTemplate(suite)), nil
}

// WriteFile writes HTML to file
func (r *HTMLReporter) WriteFile(path string) error {
	data, err := r.Generate()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func htmlTemplate(s *report.SuiteResult) string {
	statusClass := "success"
	if s.FailedTests > 0 {
		statusClass = "danger"
	}

	rows := ""
	for _, r := range s.Results {
		testStatus := "success"
		if r.Status == "failed" {
			testStatus = "danger"
		}
		rows += fmt.Sprintf(`<tr class="%s"><td>%s</td><td>%s</td><td>%s</td></tr>`,
			testStatus, r.Name, r.Status, r.Duration)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<title>%s</title>
<style>
body { font-family: Arial, sans-serif; margin: 20px; }
.summary { background: #f5f5f5; padding: 15px; border-radius: 5px; margin-bottom: 20px; }
.summary h2 { margin-top: 0; }
.badge { padding: 5px 10px; border-radius: 3px; }
.badge.success { background: #28a745; color: white; }
.badge.danger { background: #dc3545; color: white; }
table { width: 100%%; border-collapse: collapse; }
th, td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
tr.danger { background: #f8d7da; }
tr.success { background: #d4edda; }
</style>
</head>
<body>
<h1>%s</h1>
<div class="summary">
<span class="badge %s">%d/%d Passed</span>
<span>Duration: %s</span>
<span>Start: %s</span>
</div>
<table>
<tr><th>Test</th><th>Status</th><th>Duration</th></tr>
%s
</table>
</body>
</html>`, r.title, r.title, statusClass, s.PassedTests, s.TotalTests, s.Duration, s.StartTime.Format(time.RFC3339), rows)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/reporter/... -run TestHTMLReporter -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/reporter/html.go internal/reporter/html_test.go
git commit -m "feat(reporter): add HTML report generation"
```

---

### Task 2.2: 创建 `internal/reporter/junit.go` - JUnit XML 报告

**Files:**
- Create: `internal/reporter/junit.go`
- Test: `internal/reporter/junit_test.go`

**Step 1: 写测试**

```go
package reporter

import (
	"os"
	"testing"
)

func TestJUnitReporter(t *testing.T) {
	rep := NewJUnitReporter("测试套件")

	rep.StartTest("登录测试")
	rep.AddStep("启动APP", "passed", nil)
	rep.EndTest("passed", nil)

	rep.StartTest("搜索测试")
	rep.AddStep("输入关键词", "failed", fmt.Errorf("element not found"))
	rep.EndTest("failed", fmt.Errorf("test failed"))

	xml, err := rep.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Verify XML contains expected elements
	if !strings.Contains(string(xml), "<testsuite") {
		t.Error("expected testsuite element in XML")
	}
	if !strings.Contains(string(xml), "failures=\"1\"") {
		t.Error("expected 1 failure in XML")
	}

	tmpfile, err := os.CreateTemp("", "report-*.xml")
	if err != nil {
		t.Fatalf("CreateTemp failed: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	err = rep.WriteFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/reporter/... -run TestJUnitReporter -v`
Expected: FAIL - "undefined: NewJUnitReporter"

**Step 3: Write implementation**

```go
package reporter

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/liukunup/go-uop/internal/report"
)

// JUnitReporter generates JUnit XML reports
type JUnitReporter struct {
	generator *report.Generator
	suiteName string
}

// JUnitTestSuite is the JUnit XML structure
type JUnitTestSuite struct {
	XMLName   xml.Name        `xml:"testsuite"`
	Name      string         `xml:"name,attr"`
	Tests     int            `xml:"tests,attr"`
	Failures  int            `xml:"failures,attr"`
	Skipped   int            `xml:"skipped,attr"`
	Time      string         `xml:"time,attr"`
	TestCases []JUnitTestCase `xml:"testcase"`
}

// JUnitTestCase is a single test case
type JUnitTestCase struct {
	Name      string      `xml:"name,attr"`
	ClassName string      `xml:"classname,attr"`
	Time      string      `xml:"time,attr"`
	Failure   *JUnitError `xml:"failure,omitempty"`
}

// JUnitError represents a test failure
type JUnitError struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}

// NewJUnitReporter creates a new JUnit reporter
func NewJUnitReporter(name string) *JUnitReporter {
	return &JUnitReporter{
		generator: report.NewGenerator(name),
		suiteName: name,
	}
}

// StartTest starts a new test
func (r *JUnitReporter) StartTest(name string) {
	r.generator.StartTest(name)
}

// AddStep adds a step result
func (r *JUnitReporter) AddStep(name, status string, err error) {
	r.generator.AddStep(name, 0, status, err)
}

// EndTest ends a test
func (r *JUnitReporter) EndTest(status string, err error) {
	r.generator.EndTest(status, err)
}

// AddScreenshot adds a screenshot (not used in JUnit)
func (r *JUnitReporter) AddScreenshot(path, name string) {
	r.generator.AddScreenshot(path, name)
}

// Generate generates JUnit XML
func (r *JUnitReporter) Generate() ([]byte, error) {
	suite := r.generator.Generate()

	failures := 0
	cases := make([]JUnitTestCase, 0, len(suite.Results))

	for _, r := range suite.Results {
		var failure *JUnitError
		if r.Status == "failed" {
			failures++
			failure = &JUnitError{
				Message: r.Error,
				Type:    "AssertionError",
				Content: r.Error,
			}
		}

		cases = append(cases, JUnitTestCase{
			Name:      r.Name,
			ClassName: r.Name,
			Time:      r.Duration.Seconds(),
			Failure:   failure,
		})
	}

	junitSuite := JUnitTestSuite{
		Name:      r.suiteName,
		Tests:     suite.TotalTests,
		Failures:  failures,
		Skipped:   suite.SkippedTests,
		Time:      suite.Duration.Seconds(),
		TestCases: cases,
	}

	return xml.MarshalIndent(junitSuite, "", "  ")
}

// WriteFile writes XML to file
func (r *JUnitReporter) WriteFile(path string) error {
	data, err := r.Generate()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/reporter/... -run TestJUnitReporter -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/reporter/junit.go internal/reporter/junit_test.go
git commit -m "feat(reporter): add JUnit XML report generation"
```

---

### Task 2.3: 创建 `internal/reporter/config.go` - 报告配置

**Files:**
- Create: `internal/reporter/config.go`

**Step 1: Write configuration struct**

```go
package reporter

// Config controls which reports are generated
type Config struct {
	// EnableHTML enables HTML report generation
	EnableHTML bool
	// EnableJSON enables JSON report generation
	EnableJSON bool
	// EnableJUnit enables JUnit XML report generation
	EnableJUnit bool
	// OutputDir is the directory for report files
	OutputDir string
	// HTMLPath is the path for HTML report (if enabled)
	HTMLPath string
	// JSONPath is the path for JSON report (if enabled)
	JSONPath string
	// JUnitPath is the path for JUnit XML report (if enabled)
	JUnitPath string
}

// DefaultConfig returns a config with all reports disabled
func DefaultConfig() *Config {
	return &Config{}
}

// WithHTML enables HTML report
func (c *Config) WithHTML(path string) *Config {
	c.EnableHTML = true
	c.HTMLPath = path
	return c
}

// WithJSON enables JSON report
func (c *Config) WithJSON(path string) *Config {
	c.EnableJSON = true
	c.JSONPath = path
	return c
}

// WithJUnit enables JUnit XML report
func (c *Config) WithJUnit(path string) *Config {
	c.EnableJUnit = true
	c.JUnitPath = path
	return c
}

// WithOutputDir sets the output directory
func (c *Config) WithOutputDir(dir string) *Config {
	c.OutputDir = dir
	return c
}
```

**Step 2: Run test to verify it compiles**

Run: `go build ./internal/reporter/...`
Expected: SUCCESS

**Step 3: Commit**

```bash
git add internal/reporter/config.go
git commit -m "feat(reporter): add report configuration options"
```

---

## 阶段 3: 创建 CLI (`cmd/uop/main.go`)

### Task 3.1: 创建 `cmd/uop/main.go` - 统一 CLI

**Files:**
- Create: `cmd/uop/main.go`
- Modify: `go.mod` (update main module if needed)

**Step 1: 写 CLI 实现**

```go
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/liukunup/go-uop/internal/report"
	"github.com/liukunup/go-uop/internal/runner"
)

const (
	version     = "0.2.0"
	exitSuccess = 0
	exitFailure = 1
)

var (
	// Flags
	deviceFlag  string
	outputFlag  string
	reportFlag  string
	debugFlag   bool

	// Subcommands
	cmdRun    string
	cmdDebug  string
	cmdTest   string
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printHelp()
		return exitSuccess
	}

	if args[0] == "--version" {
		fmt.Printf("uop version %s\n", version)
		return exitSuccess
	}

	cmd := args[0]
	remaining := args[1:]

	switch cmd {
	case "run":
		return runCmd(remaining)
	case "debug":
		return debugCmd(remaining)
	case "test":
		return testCmd(remaining)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printHelp()
		return exitFailure
	}
}

func runCmd(args []string) int {
	flagSet := flag.NewFlagSet("run", flag.ContinueOnError)
	flagSet.StringVar(&deviceFlag, "device", "", "Device type (ios|android)")
	flagSet.StringVar(&outputFlag, "output", "", "Output directory for reports")
	flagSet.StringVar(&reportFlag, "report", "", "Report formats: html,json,junit (comma-separated)")

	if err := flagSet.Parse(args); err != nil {
		return exitFailure
	}

	file := flagSet.Arg(0)
	if file == "" {
		fmt.Fprintln(os.Stderr, "Error: file path required")
		return exitFailure
	}

	return executeFlow(file, false)
}

func debugCmd(args []string) int {
	flagSet := flag.NewFlagSet("debug", flag.ContinueOnError)
	flagSet.StringVar(&outputFlag, "output", "", "Output directory for reports")

	if err := flagSet.Parse(args); err != nil {
		return exitFailure
	}

	file := flagSet.Arg(0)
	if file == "" {
		fmt.Fprintln(os.Stderr, "Error: file path required")
		return exitFailure
	}

	return executeFlow(file, true)
}

func testCmd(args []string) int {
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	flagSet.StringVar(&outputFlag, "output", "", "Output directory for reports")
	flagSet.StringVar(&reportFlag, "report", "json", "Report formats: html,json,junit (comma-separated)")

	if err := flagSet.Parse(args); err != nil {
		return exitFailure
	}

	file := flagSet.Arg(0)
	if file == "" {
		fmt.Fprintln(os.Stderr, "Error: file path required")
		return exitFailure
	}

	return executeFlow(file, false)
}

func executeFlow(file string, isDebug bool) int {
	// Parse flow
	flow, err := runner.ParseFlowFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to parse flow: %v\n", err)
		return exitFailure
	}

	// Initialize device pool
	pool := runner.NewDevicePool()
	defer pool.Close()

	for _, dev := range flow.Devices {
		if err := pool.AddDevice(dev.ID, dev.Type, dev.Serial); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to add device %s: %v\n", dev.ID, err)
			return exitFailure
		}
	}

	// Setup reporting
	repConfig := parseReportConfig(reportFlag, flow.Name)

	// Execute
	if isDebug {
		return executeDebug(flow, pool, repConfig)
	}
	return executeRun(flow, pool, repConfig)
}

func executeRun(flow *runner.Flow, pool *runner.DevicePool, cfg *reporter.Config) int {
	reportGen := report.NewGenerator(flow.Name)
	executor := runner.NewExecutor(pool, reportGen)

	if err := executor.ExecuteFlow(flow); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return exitFailure
	}

	// Generate reports
	if err := generateReports(reportGen, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: generating reports: %v\n", err)
		return exitFailure
	}

	fmt.Println("Flow completed successfully")
	return exitSuccess
}

func executeDebug(flow *runner.Flow, pool *runner.DevicePool, cfg *reporter.Config) int {
	fmt.Printf("Debug mode for flow: %s\n", flow.Name)
	fmt.Printf("Total steps: %d\n\n", len(flow.Steps))

	dbg := runner.NewDebugger(flow)
	reportGen := report.NewGenerator(flow.Name)
	executor := runner.NewExecutor(pool, reportGen)

	for {
		step, err := dbg.NextStep()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return exitFailure
		}

		if step == nil {
			fmt.Println("\nDebug completed")
			break
		}

		fmt.Printf("[Step %d] %v\n", dbg.CurrentStep(), *step)

		// In a real implementation, this would pause and wait for user input
		// For now, just execute each step
		if err := executor.ExecuteFlow(&runner.Flow{
			Name:         flow.Name,
			Devices:      flow.Devices,
			DefaultDevice: flow.DefaultDevice,
			Steps:        []runner.Step{*step},
		}); err != nil {
			fmt.Fprintf(os.Stderr, "Error at step %d: %v\n", dbg.CurrentStep(), err)
			return exitFailure
		}
	}

	return exitSuccess
}

func parseReportConfig(flag string, suiteName string) *reporter.Config {
	cfg := reporter.DefaultConfig()
	if flag == "" {
		return cfg
	}

	formats := strings.Split(flag, ",")
	outputDir := outputFlag
	if outputDir == "" {
		outputDir = "."
	}

	for _, f := range formats {
		switch strings.TrimSpace(f) {
		case "html":
			cfg.WithHTML(filepath.Join(outputDir, suiteName+".html"))
		case "json":
			cfg.WithJSON(filepath.Join(outputDir, suiteName+".json"))
		case "junit":
			cfg.WithJUnit(filepath.Join(outputDir, suiteName+"-junit.xml"))
		}
	}

	return cfg
}

func generateReports(gen *report.Generator, cfg *reporter.Config) error {
	if cfg.EnableJSON && cfg.JSONPath != "" {
		if err := gen.ToJSONFile(cfg.JSONPath); err != nil {
			return err
		}
		fmt.Printf("JSON report: %s\n", cfg.JSONPath)
	}
	// HTML and JUnit would use their respective reporters
	return nil
}

func printHelp() {
	fmt.Printf(`UOP - Unified Mobile Automation Framework %s

Usage:
  uop <command> [arguments] [flags]

Commands:
  run <file>      Execute a flow
  debug <file>    Debug a flow step-by-step
  test <file>     Run flow with test reporting

Flags:
  --device <type>   Target device: ios or android
  --output <dir>     Output directory for reports
  --report <formats> Report formats: html,json,junit (comma-separated)
  --help, -h        Show this help message
  --version         Show version information

Examples:
  uop run flow.yaml
  uop debug flow.yaml
  uop test flow.yaml --report html,json,junit --output ./reports

Exit codes:
  0  Success
  1  Failure`, version)
}
```

**Step 2: Run test to verify it compiles**

Run: `go build -o uop ./cmd/uop/`
Expected: SUCCESS

**Step 3: Commit**

```bash
git add cmd/uop/main.go
git commit -m "feat(cli): add unified uop CLI with run/debug/test commands"
```

---

## 阶段 4: 集成 Web Console

### Task 4.1: 集成 runner 到 console

**Files:**
- Modify: `internal/console/handler.go`

**Step 1: Read existing handler**

```bash
cat internal/console/handler.go
```

**Step 2: Add runner integration** (根据现有 handler 结构添加)

```go
// ExecuteFlow executes a flow via the console API
func (h *Handler) ExecuteFlow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ExecuteFlowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse and execute flow
	flow, err := runner.ParseFlow(strings.NewReader(req.YAML))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pool := runner.NewDevicePool()
	defer pool.Close()

	for _, dev := range flow.Devices {
		pool.AddDevice(dev.ID, dev.Type, dev.Serial)
	}

	reportGen := report.NewGenerator(flow.Name)
	executor := runner.NewExecutor(pool, reportGen)

	if err := executor.ExecuteFlow(flow); err != nil {
		jsonResponse(w, map[string]interface{}{"status": "failed", "error": err.Error()})
		return
	}

	jsonResponse(w, map[string]interface{}{"status": "success", "result": reportGen.Generate()})
}
```

**Step 3: Commit**

```bash
git add internal/console/handler.go
git commit -m "feat(console): integrate runner for web execution"
```

---

## 阶段 5: 清理旧代码

### Task 5.1: 删除 maestro/ 和 yaml/ 包

**Files:**
- Delete: `maestro/` (entire directory)
- Delete: `yaml/` (entire directory)

**Step 1: Verify no references remain**

```bash
grep -r "github.com/liukunup/go-uop/maestro" --include="*.go" .
grep -r "github.com/liukunup/go-uop/yaml" --include="*.go" .
```

**Step 2: Delete directories**

```bash
rm -rf maestro/ yaml/
```

**Step 3: Update imports and build**

```bash
go mod tidy
go build ./...
```

**Step 4: Commit**

```bash
git add -A
git commit -m "refactor!: remove legacy maestro and yaml packages"
```

---

## 阶段 6: 测试与验证

### Task 6.1: 运行所有测试

**Step 1: Run tests**

```bash
go test ./... -v
```

**Step 2: Build all binaries**

```bash
go build -o bin/uop ./cmd/uop/
go build -o bin/console ./cmd/console/
go build -o bin/server ./cmd/server/
```

**Step 3: Manual verification**

```bash
# Test CLI
./bin/uop --version
./bin/uop --help

# Validate flow syntax
# (Create a test flow.yaml and run)
```

---

## 执行选项

**Plan complete and saved to `docs/plans/2026-04-04-uop-framework-refactor.md`. Two execution options:**

**1. Subagent-Driven (this session)** - I dispatch fresh subagent per task, review between tasks, fast iteration

**2. Parallel Session (separate)** - Open new session with executing-plans, batch execution with checkpoints

**Which approach?**

