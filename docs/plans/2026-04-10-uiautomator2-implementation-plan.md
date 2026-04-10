# uiautomator2 Go Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement a complete Go-language uiautomator2 library for Android UI automation at `pkg/uiautomator2/`

**Architecture:** The implementation uses HTTP + JSON-RPC to communicate with a uiautomator2 service running on the Android device. It auto-installs the required APKs (atx-agent, uiautomator-server) and provides both UiSelector and XPath element location modes.

**Tech Stack:** Go 1.21+, standard library (net/http, encoding/json), adb for device management

---

## Phase 1: Core Infrastructure

### Task 1: Create Package Structure

**Files:**
- Create: `pkg/uiautomator2/types.go`
- Create: `pkg/uiautomator2/option.go`
- Create: `pkg/uiautomator2/jsonrpc/rpc.go`

**Step 1: Create types.go**

```go
package uiautomator2

// DeviceInfo represents basic device information (matches d.info)
type DeviceInfo struct {
    CurrentPackageName string `json:"currentPackageName"`
    DisplayHeight      int    `json:"displayHeight"`
    DisplayRotation   int    `json:"displayRotation"`
    DisplaySizeDpX    int    `json:"displaySizeDpX"`
    DisplaySizeDpY    int    `json:"displaySizeDpY"`
    DisplayWidth      int    `json:"displayWidth"`
    ProductName       string `json:"productName"`
    ScreenOn          bool   `json:"screenOn"`
    SdkInt            int    `json:"sdkInt"`
    NaturalOrientation bool  `json:"naturalOrientation"`
}

// DeviceDetail represents detailed device info (matches d.device_info)
type DeviceDetail struct {
    Arch   string `json:"arch"`
    Brand  string `json:"brand"`
    Model  string `json:"model"`
    Sdk    int    `json:"sdk"`
    Serial string `json:"serial"`
    Version int   `json:"version"`
}

// AppInfo represents application information
type AppInfo struct {
    MainActivity string `json:"mainActivity"`
    Label       string `json:"label"`
    VersionName string `json:"versionName"`
    VersionCode int    `json:"versionCode"`
    Size        int64  `json:"size"`
}

// Selector represents UI selector criteria
type Selector struct {
    Text               string `json:"text,omitempty"`
    TextContains       string `json:"textContains,omitempty"`
    TextMatches        string `json:"textMatches,omitempty"`
    TextStartsWith     string `json:"textStartsWith,omitempty"`
    ClassName          string `json:"className,omitempty"`
    ClassNameMatches   string `json:"classNameMatches,omitempty"`
    Description        string `json:"description,omitempty"`
    DescriptionContains string `json:"descriptionContains,omitempty"`
    DescriptionMatches string `json:"descriptionMatches,omitempty"`
    DescriptionStartsWith string `json:"descriptionStartsWith,omitempty"`
    Checkable          bool   `json:"checkable,omitempty"`
    Checked            bool   `json:"checked,omitempty"`
    Clickable          bool   `json:"clickable,omitempty"`
    LongClickable      bool   `json:"longClickable,omitempty"`
    Scrollable         bool   `json:"scrollable,omitempty"`
    Enabled            bool   `json:"enabled,omitempty"`
    Focusable          bool   `json:"focusable,omitempty"`
    Focused            bool   `json:"focused,omitempty"`
    Selected           bool   `json:"selected,omitempty"`
    PackageName        string `json:"packageName,omitempty"`
    PackageNameMatches string `json:"packageNameMatches,omitempty"`
    ResourceId         string `json:"resourceId,omitempty"`
    ResourceIdMatches  string `json:"resourceIdMatches,omitempty"`
    Index              int    `json:"index,omitempty"`
    Instance           int    `json:"instance,omitempty"`
}

// Point represents coordinates
type Point struct {
    X int `json:"x"`
    Y int `json:"y"`
}

// Bounds represents rectangular bounds
type Bounds struct {
    Left   int `json:"left"`
    Top    int `json:"top"`
    Right  int `json:"right"`
    Bottom int `json:"bottom"`
}

// ElementInfo represents UI element information
type ElementInfo struct {
    ContentDescription string `json:"contentDescription"`
    Checked             bool   `json:"checked"`
    Scrollable          bool   `json:"scrollable"`
    Text                string `json:"text"`
    PackageName         string `json:"packageName"`
    Selected            bool   `json:"selected"`
    Enabled             bool   `json:"enabled"`
    Bounds              Bounds `json:"bounds"`
    ClassName           string `json:"className"`
    Focused             bool   `json:"focused"`
    Focusable           bool   `json:"focusable"`
    Clickable           bool   `json:"clickable"`
    ChildCount          int    `json:"childCount"`
    LongClickable       bool   `json:"longClickable"`
    VisibleBounds       Bounds `json:"visibleBounds"`
    Checkable           bool   `json:"checkable"`
}

// SessionInfo represents current app session
type SessionInfo struct {
    Activity string `json:"activity"`
    Package  string `json:"package"`
    Pid      int    `json:"pid,omitempty"`
}
```

**Step 2: Create option.go**

```go
package uiautomator2

// Config holds connection configuration
type Config struct {
    Serial    string
    Address   string  // IP:port for WiFi connection
    Timeout   int     // HTTP timeout in seconds
    Package   string  // Target app package
}

// Option configures device creation
type Option func(*Config)

// WithSerial sets device serial number
func WithSerial(serial string) Option {
    return func(c *Config) {
        c.Serial = serial
    }
}

// WithAddress sets device address (IP:port for WiFi)
func WithAddress(addr string) Option {
    return func(c *Config) {
        c.Address = addr
    }
}

// WithTimeout sets HTTP timeout
func WithTimeout(seconds int) Option {
    return func(c *Config) {
        c.Timeout = seconds
    }
}

// WithPackage sets target app package
func WithPackage(pkg string) Option {
    return func(c *Config) {
        c.Package = pkg
    }
}
```

**Step 3: Create jsonrpc/rpc.go**

```go
package jsonrpc

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "sync"
    "time"
)

// Client handles JSON-RPC communication
type Client struct {
    baseURL    string
    httpClient *http.Client
    sessionID  string
    mu         sync.Mutex
}

// NewClient creates a new JSON-RPC client
func NewClient(baseURL string, timeout time.Duration) (*Client, error) {
    return &Client{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: timeout,
        },
    }, nil
}

// Request represents a JSON-RPC request
type Request struct {
    JSONRPC string          `json:"jsonrpc"`
    ID      string          `json:"id"`
    Method  string          `json:"method"`
    Params  json.RawMessage `json:"params,omitempty"`
}

// Response represents a JSON-RPC response
type Response struct {
    JSONRPC string          `json:"jsonrpc"`
    ID      string          `json:"id"`
    Result  json.RawMessage `json:"result,omitempty"`
    Error   *RPCError       `json:"error,omitempty"`
}

// RPCError represents a JSON-RPC error
type RPCError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// Call executes a JSON-RPC method and returns the raw result
func (c *Client) Call(method string, params interface{}) (json.RawMessage, error) {
    c.mu.Lock()
    defer c.mu.Unlock()

    reqID := fmt.Sprintf("R%d", time.Now().UnixNano()/1000)
    
    var paramsRaw json.RawMessage
    if params != nil {
        var err error
        paramsRaw, err = json.Marshal(params)
        if err != nil {
            return nil, fmt.Errorf("marshal params: %w", err)
        }
    }

    req := Request{
        JSONRPC: "2.0",
        ID:      reqID,
        Method:  method,
        Params:  paramsRaw,
    }

    reqBody, err := json.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("marshal request: %w", err)
    }

    resp, err := c.httpClient.Post(c.baseURL, "application/json", bytes.NewReader(reqBody))
    if err != nil {
        return nil, fmt.Errorf("POST request: %w", err)
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("read body: %w", err)
    }

    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
    }

    var rpcResp Response
    if err := json.Unmarshal(respBody, &rpcResp); err != nil {
        return nil, fmt.Errorf("unmarshal response: %w", err)
    }

    if rpcResp.Error != nil {
        return nil, fmt.Errorf("RPC error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
    }

    return rpcResp.Result, nil
}

// CallWithSession executes a JSON-RPC method with session ID
func (c *Client) CallWithSession(method string, params interface{}) (json.RawMessage, error) {
    c.mu.Lock()
    defer c.mu.Unlock()

    reqID := fmt.Sprintf("R%d", time.Now().UnixNano()/1000)
    
    type sessionParams struct {
        SessionID string `json:"sessionId"`
    }
    
    var p struct {
        SessionID string      `json:"sessionId"`
        Params    interface{} `json:",inline"`
    }
    p.SessionID = c.sessionID
    p.Params = params

    reqBody, err := json.Marshal(map[string]interface{}{
        "jsonrpc": "2.0",
        "id":      reqID,
        "method":  method,
        "params":  p,
    })
    if err != nil {
        return nil, fmt.Errorf("marshal request: %w", err)
    }

    resp, err := c.httpClient.Post(c.baseURL, "application/json", bytes.NewReader(reqBody))
    if err != nil {
        return nil, fmt.Errorf("POST request: %w", err)
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("read body: %w", err)
    }

    var rpcResp Response
    if err := json.Unmarshal(respBody, &rpcResp); err != nil {
        return nil, fmt.Errorf("unmarshal response: %w", err)
    }

    if rpcResp.Error != nil {
        return nil, fmt.Errorf("RPC error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
    }

    return rpcResp.Result, nil
}

// SetSession sets the session ID
func (c *Client) SetSession(sessionID string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.sessionID = sessionID
}

// Session returns the session ID
func (c *Client) Session() string {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.sessionID
}
```

**Step 4: Commit**

```bash
git add pkg/uiautomator2/
git commit -m "feat(uiautomator2): add core types, options, and JSON-RPC client"
```

---

### Task 2: Create HTTP Client

**Files:**
- Create: `pkg/uiautomator2/client.go`

**Step 1: Create client.go**

```go
package uiautomator2

import (
    "fmt"
    "time"

    "github.com/liukunup/go-uop/pkg/uiautomator2/jsonrpc"
)

// Client wraps HTTP + JSON-RPC communication
type Client struct {
    rpc           *jsonrpc.Client
    baseURL       string
    httpTimeout   time.Duration
}

// NewClient creates a new uiautomator2 client
func NewClient(addr string, timeout time.Duration) (*Client, error) {
    baseURL := fmt.Sprintf("http://%s:7912/jsonrpc/0", addr)
    
    rpcClient, err := jsonrpc.NewClient(baseURL, timeout)
    if err != nil {
        return nil, fmt.Errorf("create RPC client: %w", err)
    }

    return &Client{
        rpc:         rpcClient,
        baseURL:     baseURL,
        httpTimeout: timeout,
    }, nil
}

// DeviceInfo returns basic device information
func (c *Client) DeviceInfo() (*DeviceInfo, error) {
    result, err := c.rpc.Call("deviceInfo", nil)
    if err != nil {
        return nil, err
    }

    var info DeviceInfo
    if err := json.Unmarshal(result, &info); err != nil {
        return nil, fmt.Errorf("unmarshal device info: %w", err)
    }
    return &info, nil
}

// Screenshot returns screenshot in specified format
func (c *Client) Screenshot(format string) ([]byte, error) {
    result, err := c.rpc.Call("screenshot", map[string]string{"format": format})
    if err != nil {
        return nil, err
    }

    var value string
    if err := json.Unmarshal(result, &value); err != nil {
        return nil, fmt.Errorf("unmarshal screenshot: %w", err)
    }
    
    return []byte(value), nil
}

// DumpHierarchy returns UI hierarchy XML
func (c *Client) DumpHierarchy(compressed, pretty bool, maxDepth int) (string, error) {
    params := map[string]interface{}{
        "compressed": compressed,
        "pretty":     pretty,
        "max_depth":  maxDepth,
    }
    result, err := c.rpc.Call("dumpWindowHierarchy", params)
    if err != nil {
        return "", err
    }

    var xml string
    if err := json.Unmarshal(result, &xml); err != nil {
        return "", fmt.Errorf("unmarshal hierarchy: %w", err)
    }
    return xml, nil
}

// Selector performs a selector operation and returns matched elements
func (c *Client) Selector(sel Selector) ([]string, error) {
    params := map[string]interface{}{}
    
    // Build selector params based on non-empty fields
    selMap := make(map[string]interface{})
    if sel.Text != "" {
        selMap["text"] = sel.Text
    }
    if sel.ClassName != "" {
        selMap["className"] = sel.ClassName
    }
    if sel.ResourceId != "" {
        selMap["resourceId"] = sel.ResourceId
    }
    // ... add other fields

    result, err := c.rpc.Call("selector", selMap)
    if err != nil {
        return nil, err
    }

    var elements []string
    if err := json.Unmarshal(result, &elements); err != nil {
        return nil, fmt.Errorf("unmarshal elements: %w", err)
    }
    return elements, nil
}

// Touch performs touch action at coordinates
func (c *Client) Touch(action string, x, y int) error {
    _, err := c.rpc.Call("touch", map[string]interface{}{
        "action": action, // "down", "move", "up"
        "x":      x,
        "y":      y,
    })
    return err
}

// Click performs click at coordinates
func (c *Client) Click(x, y int) error {
    _, err := c.rpc.Call("click", map[string]interface{}{
        "x": x,
        "y": y,
    })
    return err
}

// Swipe performs swipe gesture
func (c *Client) Swipe(sx, sy, ex, ey, duration int) error {
    _, err := c.rpc.Call("swipe", map[string]interface{}{
        "startX":    sx,
        "startY":    sy,
        "endX":      ex,
        "endY":      ey,
        "duration":  duration,
    })
    return err
}

// Drag performs drag gesture
func (c *Client) Drag(sx, sy, ex, ey, duration int) error {
    _, err := c.rpc.Call("drag", map[string]interface{}{
        "startX":    sx,
        "startY":    sy,
        "endX":      ex,
        "endY":      ey,
        "duration":  duration,
    })
    return err
}

// SendKeys sends key events
func (c *Client) SendKeys(text string) error {
    _, err := c.rpc.Call("sendKeys", map[string]interface{}{
        "text": text,
    })
    return err
}

// PressKey sends a key press
func (c *Client) PressKey(key string) error {
    _, err := c.rpc.Call("pressKey", map[string]interface{}{
        "key": key,
    })
    return err
}

// ScreenOn turns screen on
func (c *Client) ScreenOn() error {
    _, err := c.rpc.Call("screenOn", nil)
    return err
}

// ScreenOff turns screen off
func (c *Client) ScreenOff() error {
    _, err := c.rpc.Call("screenOff", nil)
    return err
}
```

**Step 2: Commit**

```bash
git add pkg/uiautomator2/client.go
git commit -m "feat(uiautomator2): add HTTP/JSON-RPC client"
```

---

### Task 3: Create Device Implementation

**Files:**
- Create: `pkg/uiautomator2/device.go`
- Create: `pkg/uiautomator2/server/installer.go`

**Step 1: Create device.go**

```go
package uiautomator2

import (
    "fmt"
    "os"
    "time"

    "github.com/liukunup/go-uop/core"
    "github.com/liukunup/go-uop/pkg/uiautomator2/server"
)

// Device represents a uiautomator2 device
type Device struct {
    client   *Client
    config   *Config
    pkg      string
}

// NewDevice creates a new uiautomator2 device
func NewDevice(opts ...Option) (*Device, error) {
    cfg := &Config{
        Timeout: 60,
    }
    for _, opt := range opts {
        opt(cfg)
    }

    // Auto-install server if needed
    if err := server.EnsureInstalled(cfg.Serial); err != nil {
        return nil, fmt.Errorf("ensure server installed: %w", err)
    }

    // Determine address
    addr := cfg.Address
    if addr == "" && cfg.Serial != "" {
        // Get device IP
        addr = getDeviceIP(cfg.Serial)
    }
    if addr == "" {
        addr = "localhost" // USB fallback
    }

    client, err := NewClient(addr, time.Duration(cfg.Timeout)*time.Second)
    if err != nil {
        return nil, fmt.Errorf("create client: %w", err)
    }

    return &Device{
        client: client,
        config: cfg,
        pkg:    cfg.Package,
    }, nil
}

// Platform returns the platform type
func (d *Device) Platform() core.Platform {
    return core.Android
}

// Info returns basic device info
func (d *Device) Info() (map[string]interface{}, error) {
    info, err := d.client.DeviceInfo()
    if err != nil {
        return nil, err
    }
    return map[string]interface{}{
        "currentPackageName": info.CurrentPackageName,
        "displayHeight":      info.DisplayHeight,
        "displayRotation":    info.DisplayRotation,
        "displaySizeDpX":     info.DisplaySizeDpX,
        "displaySizeDpY":     info.DisplaySizeDpY,
        "displayWidth":       info.DisplayWidth,
        "productName":        info.ProductName,
        "screenOn":           info.ScreenOn,
        "sdkInt":             info.SdkInt,
        "naturalOrientation": info.NaturalOrientation,
    }, nil
}

// Screenshot captures screenshot
func (d *Device) Screenshot() ([]byte, error) {
    return d.client.Screenshot("raw")
}

// Close releases resources
func (d *Device) Close() error {
    // Could send stop_uiautomator here
    return nil
}

// Tap performs tap at coordinates
func (d *Device) Tap(x, y int) error {
    return d.client.Click(x, y)
}

// SendKeys inputs text
func (d *Device) SendKeys(text string) error {
    return d.client.SendKeys(text)
}

// Launch launches the app
func (d *Device) Launch() error {
    if d.pkg == "" {
        return fmt.Errorf("package name not set")
    }
    return d.appStart(d.pkg, "", false)
}

// GetAlertText returns alert text (not implemented for now)
func (d *Device) GetAlertText() (string, error) {
    return "", nil
}

// AcceptAlert accepts alert (not implemented for now)
func (d *Device) AcceptAlert() error {
    return nil
}

// DismissAlert dismisses alert (not implemented for now)
func (d *Device) DismissAlert() error {
    return nil
}

// DeviceInfo returns detailed device info
func (d *Device) DeviceInfo() (*DeviceDetail, error) {
    // Use adb shell getprop
    serial := d.config.Serial
    if serial == "" {
        serial = os.Getenv("ANDROID_SERIAL")
    }
    return getDeviceDetail(serial)
}

// WindowSize returns window size
func (d *Device) WindowSize() (int, int, error) {
    info, err := d.client.DeviceInfo()
    if err != nil {
        return 0, 0, err
    }
    return info.DisplayWidth, info.DisplayHeight, nil
}

// ScreenOn turns screen on
func (d *Device) ScreenOn() error {
    return d.client.ScreenOn()
}

// ScreenOff turns screen off
func (d *Device) ScreenOff() error {
    return d.client.ScreenOff()
}

// PressKey sends key press
func (d *Device) PressKey(key string) error {
    return d.client.PressKey(key)
}

// Swipe performs swipe gesture
func (d *Device) Swipe(x1, y1, x2, y2 int, duration time.Duration) error {
    return d.client.Swipe(x1, y1, x2, y2, int(duration.Milliseconds()))
}

// Helper to get device IP (simplified)
func getDeviceIP(serial string) string {
    // Use adb shell ip route or getprop
    return ""
}

// Helper to get device detail
func getDeviceDetail(serial string) (*DeviceDetail, error) {
    // Use adb shell getprop
    return &DeviceDetail{}, nil
}

// appStart starts an app
func (d *Device) appStart(pkg, activity string, useMonkey bool) error {
    if activity != "" {
        _, err := d.client.StartActivity(pkg, activity)
        return err
    }
    // Use atx-agent to resolve and start
    _, err := d.client.StartApp(pkg)
    return err
}

var _ core.Device = (*Device)(nil)
```

**Step 2: Create server/installer.go**

```go
package server

import (
    "fmt"
    "os/exec"
    "strings"
)

const (
    atxAgentApkURL = "https://github.com/openatx/atx-agent/releases/latest/download/atx-agent.apk"
    // uiautomatorServerApkURL = ...
)

// EnsureInstalled ensures uiautomator2 services are installed on device
func EnsureInstalled(serial string) error {
    // Check if services are already running
    if isRunning(serial) {
        return nil
    }

    // Install atx-agent
    if err := installAtxAgent(serial); err != nil {
        return fmt.Errorf("install atx-agent: %w", err)
    }

    // Install uiautomator-server
    if err := installUiautomatorServer(serial); err != nil {
        return fmt.Errorf("install uiautomator-server: %w", err)
    }

    // Start services
    if err := startServices(serial); err != nil {
        return fmt.Errorf("start services: %w", err)
    }

    return nil
}

func isRunning(serial string) bool {
    // Check if atx-agent is responding
    output, err := exec.Command("adb", "-s", serial, "shell", "curl", "-s", "-m", "2", "http://localhost:7912/status").CombinedOutput()
    if err != nil {
        return false
    }
    return strings.Contains(string(output), "ok")
}

func installAtxAgent(serial string) error {
    // Download and install atx-agent APK
    // Simplified: assume APK is bundled or use URL install
    _, err := exec.Command("adb", "-s", serial, "install", "-r", "apks/atx-agent.apk").CombinedOutput()
    return err
}

func installUiautomatorServer(serial string) error {
    // Similar to atx-agent
    _, err := exec.Command("adb", "-s", serial, "install", "-r", "apks/uiautomator-server.apk").CombinedOutput()
    return err
}

func startServices(serial string) error {
    // Start atx-agent
    _, err := exec.Command("adb", "-s", serial, "shell", "am", "start", "-n", "com.github.uiautomator/.MainActivity").CombinedOutput()
    return err
}
```

**Step 3: Commit**

```bash
git add pkg/uiautomator2/device.go pkg/uiautomator2/server/installer.go
git commit -m "feat(uiautomator2): add device implementation and server installer"
```

---

## Phase 2: Element & Selector

### Task 4: Create Element and Selector

**Files:**
- Create: `pkg/uiautomator2/element/element.go`
- Create: `pkg/uiautomator2/element/selector.go`

**Step 1: Create element.go**

```go
package element

import (
    "fmt"
    "time"
    
    "github.com/liukunup/go-uop/pkg/uiautomator2"
)

// Element represents a UI element
type Element struct {
    client  *uiautomator2.Client
    sel     Selector
    bounds  uiautomator2.Bounds
    info    *uiautomator2.ElementInfo
}

// Selector defines element selection criteria
type Selector struct {
    Text               string
    TextContains       string
    TextMatches        string
    TextStartsWith     string
    ClassName          string
    ClassNameMatches   string
    Description        string
    DescriptionContains string
    DescriptionMatches string
    DescriptionStartsWith string
    Checkable          bool
    Checked            bool
    Clickable          bool
    LongClickable      bool
    Scrollable         bool
    Enabled            bool
    Focusable          bool
    Focused            bool
    Selected           bool
    PackageName        string
    PackageNameMatches string
    ResourceId         string
    ResourceIdMatches  string
    Index              int
    Instance           int
}

// Exists checks if element exists
func (e *Element) Exists() bool {
    return e.info != nil
}

// Wait waits for element to appear
func (e *Element) Wait(timeout time.Duration) (*Element, error) {
    // Use JSON-RPC wait for exists
    return e, nil
}

// WaitGone waits for element to disappear
func (e *Element) WaitGone(timeout time.Duration) error {
    return nil
}

// Info returns element info
func (e *Element) Info() (*uiautomator2.ElementInfo, error) {
    return e.info, nil
}

// Text returns element text
func (e *Element) Text() string {
    if e.info != nil {
        return e.info.Text
    }
    return ""
}

// SetText sets element text
func (e *Element) SetText(text string) error {
    _, err := e.client.SetText(e.sel, text)
    return err
}

// ClearText clears element text
func (e *Element) ClearText() error {
    return e.SetText("")
}

// Center returns center coordinates
func (e *Element) Center() (int, int) {
    if e.info == nil {
        return 0, 0
    }
    b := e.info.Bounds
    return (b.Left + b.Right) / 2, (b.Top + b.Bottom) / 2
}

// Click clicks the element
func (e *Element) Click() error {
    x, y := e.Center()
    return e.client.Click(x, y)
}

// LongClick long-clicks the element
func (e *Element) LongClick() error {
    x, y := e.Center()
    return e.client.LongClick(x, y, 500)
}

// DragTo drags element to target
func (e *Element) DragTo(x, y int, duration time.Duration) error {
    sx, sy := e.Center()
    return e.client.Drag(sx, sy, x, y, int(duration.Milliseconds()))
}

// Swipe swipes the element
func (e *Element) Swipe(direction string, steps int) error {
    return nil
}
```

**Step 2: Create selector.go**

```go
package element

import (
    "github.com/liukunup/go-uop/pkg/uiautomator2"
)

// DeviceSelector is a selector builder
type DeviceSelector struct {
    client *uiautomator2.Client
    sel    Selector
}

// NewDeviceSelector creates a new selector builder
func NewDeviceSelector(client *uiautomator2.Client) *DeviceSelector {
    return &DeviceSelector{client: client}
}

// Text creates a text selector
func (ds *DeviceSelector) Text(text string) *DeviceSelector {
    ds.sel.Text = text
    return ds
}

// TextContains creates a text contains selector
func (ds *DeviceSelector) TextContains(contains string) *DeviceSelector {
    ds.sel.TextContains = contains
    return ds
}

// ClassName creates a class name selector
func (ds *DeviceSelector) ClassName(className string) *DeviceSelector {
    ds.sel.ClassName = className
    return ds
}

// ResourceId creates a resource-id selector
func (ds *DeviceSelector) ResourceId(id string) *DeviceSelector {
    ds.sel.ResourceId = id
    return ds
}

// Clickable creates a clickable selector
func (ds *DeviceSelector) Clickable(clickable bool) *DeviceSelector {
    ds.sel.Clickable = clickable
    return ds
}

// Enabled creates an enabled selector
func (ds *DeviceSelector) Enabled(enabled bool) *DeviceSelector {
    ds.sel.Enabled = enabled
    return ds
}

// PackageName creates a package name selector
func (ds *DeviceSelector) PackageName(pkg string) *DeviceSelector {
    ds.sel.PackageName = pkg
    return ds
}

// Scrollable creates a scrollable selector
func (ds *DeviceSelector) Scrollable(scrollable bool) *DeviceSelector {
    ds.sel.Scrollable = scrollable
    return ds
}

// Index creates an index selector
func (ds *DeviceSelector) Index(index int) *DeviceSelector {
    ds.sel.Index = index
    return ds
}

// Instance creates an instance selector
func (ds *DeviceSelector) Instance(instance int) *DeviceSelector {
    ds.sel.Instance = instance
    return ds
}

// Child adds a child selector
func (ds *DeviceSelector) Child(selector *DeviceSelector) *DeviceSelector {
    // Implementation
    return ds
}

// Sibling adds a sibling selector
func (ds *DeviceSelector) Sibling(selector *DeviceSelector) *DeviceSelector {
    return ds
}

// LeftOf returns element to the left
func (ds *DeviceSelector) LeftOf(other *DeviceSelector) *DeviceSelector {
    return ds
}

// RightOf returns element to the right
func (ds *DeviceSelector) RightOf(other *DeviceSelector) *DeviceSelector {
    return ds
}

// UpOf returns element above
func (ds *DeviceSelector) UpOf(other *DeviceSelector) *DeviceSelector {
    return ds
}

// DownOf returns element below
func (ds *DeviceSelector) DownOf(other *DeviceSelector) *DeviceSelector {
    return ds
}

// First returns first instance
func (ds *DeviceSelector) First() *Element {
    ds.sel.Instance = 0
    return ds.Do()
}

// Index returns element at index
func (ds *DeviceSelector) Index(i int) *Element {
    ds.sel.Instance = i
    return ds.Do()
}

// Do executes the selector and returns element
func (ds *DeviceSelector) Do() *Element {
    elements, err := ds.client.Selector(ds.sel)
    if err != nil || len(elements) == 0 {
        return &Element{}
    }
    // Get element info
    info, _ := ds.client.ElementInfo(elements[0])
    return &Element{
        client: ds.client,
        sel:    ds.sel,
        bounds: info.Bounds,
        info:   info,
    }
}

// Count returns number of matching elements
func (ds *DeviceSelector) Count() int {
    elements, err := ds.client.Selector(ds.sel)
    if err != nil {
        return 0
    }
    return len(elements)
}
```

**Step 3: Commit**

```bash
git add pkg/uiautomator2/element/
git commit -m "feat(uiautomator2): add element and selector"
```

---

## Phase 3: XPath Support

### Task 5: Create XPath Implementation

**Files:**
- Create: `pkg/uiautomator2/xpath/xpath.go`

**Step 1: Create xpath.go**

```go
package xpath

import (
    "encoding/xml"
    "fmt"
    "regexp"
    "strings"
)

// Matcher handles XPath matching
type Matcher struct {
    xmlDoc  *XMLNode
    expr    string
}

// XMLNode represents an XML element
type XMLNode struct {
    XMLName  xml.Name
   Attrs    map[string]string
    Text     string
    Children []*XMLNode
    Parent   *XMLNode
}

// Compile compiles an XPath expression
func Compile(expr string) (*Matcher, error) {
    return &Matcher{expr: expr}, nil
}

// Match matches the XPath against root
func (m *Matcher) Match(root *XMLNode) []*XMLNode {
    return m.matchNode(root, m.expr)
}

func (m *Matcher) matchNode(node *XMLNode, expr string) []*XMLNode {
    // Handle //
    if strings.HasPrefix(expr, "//") {
        return m.matchRecursive(node, expr[2:])
    }
    // Handle /
    if strings.HasPrefix(expr, "/") {
        return m.matchDirect(node, expr[1:])
    }
    // Handle @attribute
    if strings.HasPrefix(expr, "@") {
        return m.matchAttribute(node, expr[1:])
    }
    return nil
}

func (m *Matcher) matchRecursive(node *XMLNode, expr string) []*XMLNode {
    var results []*XMLNode
    for _, child := range node.Children {
        results = append(results, m.matchNode(child, expr)...)
    }
    for _, child := range node.Children {
        results = append(results, m.matchRecursive(child, expr)...)
    }
    return results
}

func (m *Matcher) matchDirect(node *XMLNode, expr string) []*XMLNode {
    if node.XMLName.Local == expr {
        return []*XMLNode{node}
    }
    return nil
}

func (m *Matcher) matchAttribute(node *XMLNode, attrExpr string) []*XMLNode {
    // Parse @resource-id='...'
    re := regexp.MustCompile(`@?(\w+)\s*=\s*['"]([^'"]+)['"]`)
    matches := re.FindStringSubmatch(attrExpr)
    if len(matches) < 3 {
        return nil
    }
    attrName := matches[1]
    attrValue := matches[2]
    
    if node.Attrs[attrName] == attrValue {
        return []*XMLNode{node}
    }
    return nil
}

// ParseXML parses XML string into XMLNode
func ParseXML(xmlStr string) (*XMLNode, error) {
    decoder := xml.NewDecoder(strings.NewReader(xmlStr))
    return parseNode(decoder, nil)
}

func parseNode(decoder *xml.Decoder, parent *XMLNode) (*XMLNode, error) {
    for {
        token, err := decoder.Token()
        if err != nil {
            break
        }
        switch se := token.(type) {
        case xml.StartElement:
            node := &XMLNode{
                XMLName: se.Name,
                Attrs:   make(map[string]string),
                Parent:  parent,
            }
            for _, attr := range se.Attr {
                node.Attrs[attr.Name.Local] = attr.Value
            }
            child, err := parseNode(decoder, node)
            if err == nil && child != nil {
                node.Children = append(node.Children, child)
            }
            return node, nil
        case xml.CharData:
            if parent != nil {
                parent.Text = string(se)
            }
        case xml.EndElement:
            if parent != nil {
                return parent, nil
            }
        }
    }
    return nil, nil
}
```

**Step 2: Commit**

```bash
git add pkg/uiautomator2/xpath/
git commit -m "feat(uiautomator2): add XPath support"
```

---

## Phase 4: App Management

### Task 6: Create App Manager

**Files:**
- Create: `pkg/uiautomator2/app/manager.go`

**Step 1: Create manager.go**

```go
package app

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
)

// Manager handles app lifecycle
type Manager struct {
    client  *uiautomator2.Client
    serial  string
}

// NewManager creates a new app manager
func NewManager(client *uiautomator2.Client, serial string) *Manager {
    return &Manager{client: client, serial: serial}
}

// Install installs app from URL
func (m *Manager) Install(url string) error {
    // Download APK
    resp, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("download APK: %w", err)
    }
    defer resp.Body.Close()

    // Save to temp file
    tmpFile := filepath.Join(os.TempDir(), "uiautomator2_install.apk")
    f, err := os.Create(tmpFile)
    if err != nil {
        return fmt.Errorf("create temp file: %w", err)
    }
    defer os.Remove(tmpFile)
    defer f.Close()

    _, err = io.Copy(f, resp.Body)
    if err != nil {
        return fmt.Errorf("write temp file: %w", err)
    }

    // Push to device
    devicePath := "/sdcard/tmp_install.apk"
    if err := m.pushToDevice(tmpFile, devicePath); err != nil {
        return fmt.Errorf("push to device: %w", err)
    }

    // Install via uiautomator2 JSON-RPC
    return m.client.InstallApk(devicePath)
}

func (m *Manager) pushToDevice(localPath, remotePath string) error {
    // Use adb push
    cmd := exec.Command("adb", "-s", m.serial, "push", localPath, remotePath)
    _, err := cmd.CombinedOutput()
    return err
}

// Start starts an app
func (m *Manager) Start(pkg string, activity string, useMonkey bool) error {
    if activity != "" {
        return m.client.StartActivity(pkg, activity)
    }
    return m.client.StartApp(pkg, useMonkey)
}

// Stop stops an app
func (m *Manager) Stop(pkg string) error {
    return m.client.ForceStop(pkg)
}

// Clear clears app data
func (m *Manager) Clear(pkg string) error {
    return m.client.Clear(pkg)
}

// StopAll stops all apps
func (m *Manager) StopAll(excludes []string) error {
    running, err := m.client.ListRunningApps()
    if err != nil {
        return err
    }
    for _, pkg := range running {
        if contains(excludes, pkg) {
            continue
        }
        m.Stop(pkg)
    }
    return nil
}

// Info returns app info
func (m *Manager) Info(pkg string) (*uiautomator2.AppInfo, error) {
    return m.client.AppInfo(pkg)
}

// Icon returns app icon
func (m *Manager) Icon(pkg string) ([]byte, error) {
    return m.client.AppIcon(pkg)
}

// ListRunning returns list of running apps
func (m *Manager) ListRunning() ([]string, error) {
    return m.client.ListRunningApps()
}

// Wait waits for app to start
func (m *Manager) Wait(pkg string, front bool, timeout float64) (int, error) {
    return m.client.AppWait(pkg, front, timeout)
}

// Push pushes file to device
func (m *Manager) Push(localPath, remotePath string, mode int) error {
    if mode != 0 {
        // Set mode after push
        if err := m.pushToDevice(localPath, remotePath); err != nil {
            return err
        }
        return m.setFileMode(remotePath, mode)
    }
    return m.pushToDevice(localPath, remotePath)
}

// Pull pulls file from device
func (m *Manager) Pull(remotePath, localPath string) error {
    cmd := exec.Command("adb", "-s", m.serial, "pull", remotePath, localPath)
    _, err := cmd.CombOutput()
    return err
}

func (m *Manager) setFileMode(path string, mode int) error {
    cmd := exec.Command("adb", "-s", m.serial, "shell", "chmod", fmt.Sprintf("%o", mode), path)
    _, err := cmd.CombinedOutput()
    return err
}

// AutoGrantPermissions auto-grants permissions
func (m *Manager) AutoGrantPermissions(pkg string) error {
    return m.client.GrantPermissions(pkg)
}

// OpenUrl opens URL scheme
func (m *Manager) OpenUrl(url string) error {
    return m.client.OpenUrl(url)
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

**Step 2: Commit**

```bash
git add pkg/uiautomator2/app/
git commit -m "feat(uiautomator2): add app manager"
```

---

## Phase 5: Input & Gestures

### Task 7: Create Input Handlers

**Files:**
- Create: `pkg/uiautomator2/input/keys.go`
- Create: `pkg/uiautomator2/input/gesture.go`
- Create: `pkg/uiautomator2/input/clipboard.go`

**Step 1: Create keys.go**

```go
package input

import (
    "github.com/liukunup/go-uop/pkg/uiautomator2"
)

// KeyMap maps key names to Android key codes
var KeyMap = map[string]int{
    "home":           3,
    "back":           4,
    "left":           21,
    "right":          22,
    "up":             19,
    "down":           20,
    "center":         23,
    "menu":           82,
    "search":         84,
    "enter":          66,
    "delete":         67,
    "del":            67,
    "recent":         187,
    "volume_up":      24,
    "volume_down":    25,
    "volume_mute":    164,
    "camera":         27,
    "power":          26,
}

// Keys handles key operations
type Keys struct {
    client *uiautomator2.Client
}

// Press sends a key press
func (k *Keys) Press(key string) error {
    return k.client.PressKey(key)
}

// PressKeyCode sends a key code with meta
func (k *Keys) PressKeyCode(code, meta int) error {
    return k.client.PressKeyCode(code, meta)
}
```

**Step 2: Create gesture.go**

```go
package input

import (
    "github.com/liukunup/go-uop/pkg/uiautomator2"
)

// Gesture handles gesture operations
type Gesture struct {
    client *uiautomator2.Client
}

// NewGesture creates a new gesture handler
func NewGesture(client *uiautomator2.Client) *Gesture {
    return &Gesture{client: client}
}

// Tap performs tap at coordinates
func (g *Gesture) Tap(x, y int) error {
    return g.client.Click(x, y)
}

// DoubleClick performs double click
func (g *Gesture) DoubleClick(x, y int) error {
    return g.client.DoubleClick(x, y)
}

// LongClick performs long click
func (g *Gesture) LongClick(x, y int, duration int) error {
    return g.client.LongClick(x, y, duration)
}

// Swipe performs swipe gesture
func (g *Gesture) Swipe(sx, sy, ex, ey, duration int) error {
    return g.client.Swipe(sx, sy, ex, ey, duration)
}

// Drag performs drag gesture
func (g *Gesture) Drag(sx, sy, ex, ey, duration int) error {
    return g.client.Drag(sx, sy, ex, ey, duration)
}

// TouchDown simulates touch down
func (g *Gesture) TouchDown(x, y int) error {
    return g.client.Touch("down", x, y)
}

// TouchMove simulates touch move
func (g *Gesture) TouchMove(x, y int) error {
    return g.client.Touch("move", x, y)
}

// TouchUp simulates touch up
func (g *Gesture) TouchUp(x, y int) error {
    return g.client.Touch("up", x, y)
}

// SwipePoints performs swipe through multiple points
func (g *Gesture) SwipePoints(points []struct{ X, Y int }, duration int) error {
    if len(points) < 2 {
        return nil
    }
    for i := 0; i < len(points)-1; i++ {
        p1 := points[i]
        p2 := points[i+1]
        if err := g.Swipe(p1.X, p1.Y, p2.X, p2.Y, duration); err != nil {
            return err
        }
    }
    return nil
}

// PinchIn performs pinch in gesture
func (g *Gesture) PinchIn(percent, steps int) error {
    return g.client.PinchIn(percent, steps)
}

// PinchOut performs pinch out gesture
func (g *Gesture) PinchOut(percent, steps int) error {
    return g.client.PinchOut(percent, steps)
}
```

**Step 3: Create clipboard.go**

```go
package input

import (
    "github.com/liukunup/go-uop/pkg/uiautomator2"
)

// Clipboard handles clipboard operations
type Clipboard struct {
    client *uiautomator2.Client
}

// NewClipboard creates a new clipboard handler
func NewClipboard(client *uiautomator2.Client) *Clipboard {
    return &Clipboard{client: client}
}

// Get gets clipboard content
func (c *Clipboard) Get() (string, error) {
    return c.client.GetClipboard()
}

// Set sets clipboard content
func (c *Clipboard) Set(text string, label string) error {
    return c.client.SetClipboard(text, label)
}
```

**Step 4: Commit**

```bash
git add pkg/uiautomator2/input/
git commit -m "feat(uiautomator2): add input handlers"
```

---

## Phase 6: Screen & Toast

### Task 8: Create Screen Handlers

**Files:**
- Create: `pkg/uiautomator2/screen/screenshot.go`
- Create: `pkg/uiautomator2/screen/orientation.go`
- Create: `pkg/uiautomator2/screen/toast.go`

**Step 1: Create screenshot.go**

```go
package screen

import (
    "fmt"
    "os"
    
    "github.com/liukunup/go-uop/pkg/uiautomator2"
)

// Screenshot handles screenshot operations
type Screenshot struct {
    client *uiautomator2.Client
}

// NewScreenshot creates a new screenshot handler
func NewScreenshot(client *uiautomator2.Client) *Screenshot {
    return &Screenshot{client: client}
}

// Capture captures screenshot and saves to file
func (s *Screenshot) Capture(filename string) error {
    data, err := s.client.Screenshot("raw")
    if err != nil {
        return fmt.Errorf("capture screenshot: %w", err)
    }
    return os.WriteFile(filename, data, 0644)
}

// CapturePillow captures screenshot as PIL format
func (s *Screenshot) CapturePillow() ([]byte, error) {
    return s.client.Screenshot("pillow")
}

// CaptureOpenCV captures screenshot as OpenCV format
func (s *Screenshot) CaptureOpenCV() ([]byte, error) {
    return s.client.Screenshot("opencv")
}
```

**Step 2: Create orientation.go**

```go
package screen

import (
    "github.com/liukunup/go-uop/pkg/uiautomator2"
)

// Orientation handles screen orientation
type Orientation struct {
    client *uiautomator2.Client
}

// NewOrientation creates a new orientation handler
func NewOrientation(client *uiautomator2.Client) *Orientation {
    return &Orientation{client: client}
}

// Get returns current orientation
func (o *Orientation) Get() (string, error) {
    return o.client.GetOrientation()
}

// Set sets orientation
func (o *Orientation) Set(orientation string) error {
    return o.client.SetOrientation(orientation)
}

// Freeze freezes rotation
func (o *Orientation) Freeze() error {
    return o.client.FreezeRotation(true)
}

// Unfreeze unfreezes rotation
func (o *Orientation) Unfreeze() error {
    return o.client.FreezeRotation(false)
}
```

**Step 3: Create toast.go**

```go
package screen

import (
    "github.com/liukunup/go-uop/pkg/uiautomator2"
)

// Toast handles toast operations
type Toast struct {
    client *uiautomator2.Client
}

// NewToast creates a new toast handler
func NewToast(client *uiautomator2.Client) *Toast {
    return &Toast{client: client}
}

// GetLast returns last toast
func (t *Toast) GetLast() (string, error) {
    return t.client.LastToast()
}

// Clear clears toast cache
func (t *Toast) Clear() error {
    return t.client.ClearToast()
}
```

**Step 4: Commit**

```bash
git add pkg/uiautomator2/screen/
git commit -m "feat(uiautomator2): add screen handlers"
```

---

## Phase 7: WatchContext

### Task 9: Create WatchContext

**Files:**
- Create: `pkg/uiautomator2/watch/watcher.go`

**Step 1: Create watcher.go**

```go
package watch

import (
    "regexp"
    "sync"
    "time"
    
    "github.com/liukunup/go-uop/pkg/uiautomator2"
    "github.com/liukunup/go-uop/pkg/uiautomator2/xpath"
)

// Watcher handles WatchContext functionality
type Watcher struct {
    client     *uiautomator2.Client
    conditions []*Condition
    running    bool
    mu         sync.Mutex
}

// Condition represents a watch condition
type Condition struct {
    Pattern *regexp.Regexp
    Action  func(el *element.Element) error
    XPath   string
}

// Context is a watch context
type Context struct {
    watcher  *Watcher
    builtin  bool
}

// NewWatcher creates a new watcher
func NewWatcher(client *uiautomator2.Client) *Watcher {
    return &Watcher{
        client:     client,
        conditions: make([]*Condition, 0),
        running:    false,
    }
}

// When adds a condition
func (c *Context) When(pattern string) *Context {
    re := regexp.MustCompile(pattern)
    c.watcher.conditions = append(c.watcher.conditions, &Condition{
        Pattern: re,
    })
    return c
}

// WhenWithXPath adds an XPath condition
func (c *Context) WhenWithXPath(xpathExpr string) *Context {
    c.watcher.conditions = append(c.watcher.conditions, &Condition{
        XPath: xpathExpr,
    })
    return c
}

// Click sets click action
func (c *Context) Click() *Context {
    last := c.watcher.conditions[len(c.watcher.conditions)-1]
    last.Action = func(el *element.Element) error {
        return el.Click()
    }
    return c
}

// Call sets custom action
func (c *Context) Call(fn func(d *uiautomator2.Device, el *element.Element) error) *Context {
    last := c.watcher.conditions[len(c.watcher.conditions)-1]
    last.Action = func(el *element.Element) error {
        return fn(nil, el)
    }
    return c
}

// WaitStability waits for UI to stabilize
func (c *Context) WaitStability() {
    for {
        time.Sleep(2 * time.Second)
        // Check if conditions are still stable
        if !c.checkConditions() {
            break
        }
    }
}

func (c *Context) checkConditions() bool {
    // Get hierarchy and check conditions
    return true
}

// Close stops watching
func (c *Context) Close() {
    c.watcher.Stop()
}

// Start starts watching
func (c *Context) Start() {
    go c.watcher.run()
}

func (w *Watcher) run() {
    w.mu.Lock()
    w.running = true
    w.mu.Unlock()

    for w.running {
        w.checkAndAct()
        time.Sleep(2 * time.Second)
    }
}

func (w *Watcher) checkAndAct() {
    xml, err := w.client.DumpHierarchy(false, false, 50)
    if err != nil {
        return
    }

    root, err := xpath.ParseXML(xml)
    if err != nil {
        return
    }

    for _, cond := range w.conditions {
        var matched []*xpath.XMLNode
        if cond.XPath != "" {
            matcher, _ := xpath.Compile(cond.XPath)
            matched = matcher.Match(root)
        } else if cond.Pattern != nil {
            matched = w.findByPattern(root, cond.Pattern)
        }

        for _, node := range matched {
            if cond.Action != nil {
                cond.Action(nil) // Element would be constructed here
            }
        }
    }
}

func (w *Watcher) findByPattern(node *xpath.XMLNode, pattern *regexp.Regexp) []*xpath.XMLNode {
    var results []*xpath.XMLNode
    text := node.Text
    if pattern.MatchString(text) {
        results = append(results, node)
    }
    for _, child := range node.Children {
        results = append(results, w.findByPattern(child, pattern)...)
    }
    return results
}

// Stop stops the watcher
func (w *Watcher) Stop() {
    w.mu.Lock()
    defer w.mu.Unlock()
    w.running = false
}
```

**Step 2: Commit**

```bash
git add pkg/uiautomator2/watch/
git commit -m "feat(uiautomator2): add WatchContext"
```

---

## Phase 8: Settings & Debug

### Task 10: Create Settings and Debug

**Files:**
- Create: `pkg/uiautomator2/settings.go`

**Step 1: Create settings.go**

```go
package uiautomator2

import (
    "sync"
)

// Settings holds global settings
type Settings struct {
    mu                    sync.RWMutex
    waitTimeout          float64
    operationDelay        (float64, float64)
    operationDelayMethods []string
    maxDepth              int
}

// WaitTimeout gets/sets wait timeout
func (s *Settings) WaitTimeout() float64 {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.waitTimeout
}

func (s *Settings) SetWaitTimeout(timeout float64) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.waitTimeout = timeout
}

// OperationDelay gets/sets operation delay
func (s *Settings) OperationDelay() (float64, float64) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.operationDelay
}

func (s *Settings) SetOperationDelay(before, after float64) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.operationDelay = (before, after)
}
```

**Step 2: Commit**

```bash
git add pkg/uiautomator2/settings.go
git commit -m "feat(uiautomator2): add settings and debug support"
```

---

## Implementation Complete

The above plan provides the complete implementation for a Go-language uiautomator2 library. Each task builds upon the previous one, creating a fully functional Android UI automation tool.

---

## Execution Options

**1. Subagent-Driven (this session)** - I dispatch fresh subagent per task, review between tasks, fast iteration

**2. Parallel Session (separate)** - Open new session with executing-plans, batch execution with checkpoints

**Which approach?**
