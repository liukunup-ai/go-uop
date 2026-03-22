# go-uop Implementation Plan

**Goal:** Build a unified mobile automation library supporting iOS (native WDA) and Android (ADB) with chainable Go API and YAML runner.

**Architecture:** 
- Layered architecture: User Layer (Go API + YAML) → Command Layer → Platform Drivers (iOS/Android) → Core Modules
- Unified `Device` interface with platform-specific implementations
- Chainable fluent API following PageObject pattern

**Tech Stack:**
- Go 1.21+
- gopkg.in/yaml.v3 (YAML parsing)
- github.com/robert-nelsen/otto (JavaScript engine)
- github.com/g不甘于平庸chi/yaegi (Python engine)
- gocv.io/gocv (OpenCV)
- net/http (WDA/ADB protocols)

---

## Phase 1: Core Framework (M1)

### Task 1.1: Initialize Go Module

**Files:**
- Create: `go.mod`
- Create: `uop.go` (main entry)

**Step 1: Create go.mod**

```bash
cd /Users/liukunup/Documents/repo/go-uop
go mod init github.com/yourname/go-uop
```

**Step 2: Create main entry file**

Create `uop.go`:
```go
package uop

import "errors"

// Platform represents target platform
type Platform string

const (
    IOS     Platform = "ios"
    Android Platform = "android"
)

// ErrNotImplemented indicates feature not implemented
var ErrNotImplemented = errors.New("not implemented")

// ErrDeviceNotFound indicates device not found
var ErrDeviceNotFound = errors.New("device not found")
```

**Step 3: Create Device interface**

Create `device.go`:
```go
package uop

import "context"

// Device represents a connected mobile device
type Device interface {
    // Platform returns the device platform
    Platform() Platform
    
    // Info returns device information
    Info() (map[string]interface{}, error)
    
    // Screenshot captures current screen
    Screenshot() ([]byte, error)
    
    // Close releases device resources
    Close() error
}

// NewDevice creates a new device connection
func NewDevice(platform Platform, opts ...DeviceOption) (Device, error) {
    // TODO: implement
    return nil, ErrNotImplemented
}
```

**Step 4: Add DeviceOption pattern**

Create `option.go`:
```go
package uop

import "time"

// DeviceOption configures device creation
type DeviceOption func(*deviceConfig)

type deviceConfig struct {
    serial   string
    address  string
    timeout  time.Duration
}

// WithSerial sets device serial number
func WithSerial(serial string) DeviceOption {
    return func(c *deviceConfig) {
        c.serial = serial
    }
}

// WithAddress sets device address (IP:port for WiFi)
func WithAddress(addr string) DeviceOption {
    return func(c *deviceConfig) {
        c.address = addr
    }
}

// WithTimeout sets connection timeout
func WithTimeout(d time.Duration) DeviceOption {
    return func(c *deviceConfig) {
        c.timeout = d
    }
}
```

**Step 5: Write test**

Create `uop_test.go`:
```go
package uop

import (
    "testing"
    "time"
)

func TestNewDevice_InvalidPlatform(t *testing.T) {
    _, err := NewDevice("unknown")
    if err == nil {
        t.Error("expected error for unknown platform")
    }
}

func TestDeviceOption_WithSerial(t *testing.T) {
    opt := WithSerial("test-123")
    cfg := &deviceConfig{}
    opt(cfg)
    
    if cfg.serial != "test-123" {
        t.Errorf("expected serial 'test-123', got '%s'", cfg.serial)
    }
}

func TestDeviceOption_WithTimeout(t *testing.T) {
    opt := WithTimeout(30 * time.Second)
    cfg := &deviceConfig{}
    opt(cfg)
    
    if cfg.timeout != 30*time.Second {
        t.Errorf("expected 30s timeout, got %v", cfg.timeout)
    }
}
```

**Step 6: Run tests**

```bash
go test -v ./...
```

Expected: 3 tests PASS

**Step 7: Commit**

```bash
git add go.mod uop.go device.go option.go uop_test.go
git commit -m "feat: initialize go-uop module with core types

- Add Platform type (IOS/Android)
- Add Device interface
- Add DeviceOption pattern for configuration
- Add basic error types"
```

---

### Task 1.2: Create Module Structure

**Files:**
- Create: `internal/locator/locator.go`
- Create: `internal/locator/locator_test.go`
- Create: `internal/action/action.go`
- Create: `internal/action/action_test.go`

**Step 1: Create locator package**

Create `internal/locator/locator.go`:
```go
package locator

import (
    "regexp"
    "strings"
)

// SelectorType represents the type of element selector
type SelectorType int

const (
    SelectorTypeText SelectorType = iota
    SelectorTypeID
    SelectorTypeXPath
    SelectorTypeClassName
    SelectorTypePredicate   // iOS only
    SelectorTypeClassChain  // iOS only
)

// Selector describes how to find an element
type Selector struct {
    Type     SelectorType
    Value    string
    Index    int // -1 means default (first/top-left)
    regex    *regexp.Regexp
}

// isRegex checks if value is wrapped in /.../
func isRegex(value string) bool {
    return len(value) >= 2 && 
           value[0] == '/' && 
           value[len(value)-1] == '/'
}

// ParseRegex extracts and compiles regex from /.../ pattern
func ParseRegex(value string) (*regexp.Regexp, bool) {
    if !isRegex(value) {
        return nil, false
    }
    pattern := value[1 : len(value)-1]
    re, err := regexp.Compile(pattern)
    if err != nil {
        return nil, false
    }
    return re, true
}

// NewSelector creates a locator with auto-detected type
func NewSelector(value string) *Selector {
    // Auto-detect regex
    if re, ok := ParseRegex(value); ok {
        return &Selector{
            Type:  SelectorTypeText,
            Value: value,
            regex: re,
        }
    }
    
    return &Selector{
        Type:  SelectorTypeText,
        Value: value,
        Index: -1,
    }
}

// ByText creates a text locator
func ByText(text string) *Selector {
    return NewSelector(text)
}

// ByID creates an ID locator
func ByID(id string) *Selector {
    if re, ok := ParseRegex(id); ok {
        return &Selector{Type: SelectorTypeID, Value: id, regex: re}
    }
    return &Selector{Type: SelectorTypeID, Value: id, Index: -1}
}

// ByXPath creates an XPath locator
func ByXPath(xpath string) *Selector {
    return &Selector{Type: SelectorTypeXPath, Value: xpath, Index: -1}
}

// ByClassName creates a class name locator
func ByClassName(class string) *Selector {
    return &Selector{Type: SelectorTypeClassName, Value: class, Index: -1}
}

// ByPredicate creates a predicate locator (iOS only)
func ByPredicate(predicate string) *Selector {
    return &Selector{Type: SelectorTypePredicate, Value: predicate, Index: -1}
}

// ByClassChain creates a class chain locator (iOS only)
func ByClassChain(chain string) *Selector {
    return &Selector{Type: SelectorTypeClassChain, Value: chain, Index: -1}
}

// Index sets the element index
func (l *Selector) Index(idx int) *Selector {
    l.Index = idx
    return l
}

// Match checks if the locator matches the given text
func (l *Selector) Match(text string) bool {
    if l.regex != nil {
        return l.regex.MatchString(text)
    }
    return strings.Contains(strings.ToLower(text), strings.ToLower(l.Value))
}

// String returns string representation
func (l *Selector) String() string {
    return l.Value
}
```

**Step 2: Write locator tests**

Create `internal/locator/locator_test.go`:
```go
package locator

import (
    "testing"
)

func TestByText_ExactMatch(t *testing.T) {
    l := ByText("登录")
    
    if !l.Match("登录") {
        t.Error("expected exact match")
    }
    if l.Match("登录按钮") {
        t.Error("should not match partial")
    }
}

func TestByText_RegexMatch(t *testing.T) {
    l := ByText("/登.*/")
    
    if !l.Match("登录") {
        t.Error("expected regex match '登录'")
    }
    if !l.Match("登录页") {
        t.Error("expected regex match '登录页'")
    }
    if l.Match("取消") {
        t.Error("should not match '取消'")
    }
}

func TestByText_RegexUserId(t *testing.T) {
    l := ByText("/^用户\\d+$/")
    
    if !l.Match("用户1") {
        t.Error("expected match '用户1'")
    }
    if !l.Match("用户123") {
        t.Error("expected match '用户123'")
    }
    if l.Match("用户名") {
        t.Error("should not match '用户名'")
    }
}

func TestByID_Regex(t *testing.T) {
    l := ByID("/btn_.*_confirm/")
    
    if !l.Match("btn_submit_confirm") {
        t.Error("expected match")
    }
    if !l.Match("btn_cancel_confirm") {
        t.Error("expected match")
    }
}

func TestSelector_Index(t *testing.T) {
    l := ByText("确定").Index(2)
    
    if l.Index != 2 {
        t.Errorf("expected index 2, got %d", l.Index)
    }
}

func TestSelector_DefaultIndex(t *testing.T) {
    l := ByText("确定")
    
    if l.Index != -1 {
        t.Errorf("expected default index -1, got %d", l.Index)
    }
}

func TestByXPath(t *testing.T) {
    l := ByXPath("//Button[@text='OK']")
    
    if l.Type != SelectorTypeXPath {
        t.Errorf("expected XPath type, got %d", l.Type)
    }
    if l.Value != "//Button[@text='OK']" {
        t.Errorf("unexpected value: %s", l.Value)
    }
}

func TestByClassName(t *testing.T) {
    l := ByClassName("android.widget.Button")
    
    if l.Type != SelectorTypeClassName {
        t.Errorf("expected ClassName type, got %d", l.Type)
    }
}

func TestByPredicate(t *testing.T) {
    l := ByPredicate(`name == "test"`)
    
    if l.Type != SelectorTypePredicate {
        t.Errorf("expected Predicate type, got %d", l.Type)
    }
}

func TestByClassChain(t *testing.T) {
    l := ByClassChain(`**/Button[*]`)
    
    if l.Type != SelectorTypeClassChain {
        t.Errorf("expected ClassChain type, got %d", l.Type)
    }
}
```

**Step 3: Run tests**

```bash
go test -v ./internal/locator/...
```

**Step 4: Commit**

```bash
git add internal/locator/
git commit -m "feat: add locator package with regex auto-detection

- ByText, ByID, ByXPath, ByClassName, ByPredicate, ByClassChain
- Auto-detect regex patterns wrapped in /.../
- Index support for multi-element selection
- Default index -1 (first/top-left)"
```

---

### Task 1.3: Create Action Package

**Files:**
- Create: `internal/action/action.go`
- Create: `internal/action/action_test.go`

**Step 1: Create action package**

Create `internal/action/action.go`:
```go
package action

import (
    "time"
    
    "github.com/yourname/go-uop/internal/locator"
)

// Action represents a device action
type Action interface {
    Do() error
}

// TapAction taps on coordinates or element
type TapAction struct {
    X, Y     int
    Element  *locator.Selector
}

// SwipeAction performs swipe gesture
type SwipeAction struct {
    StartX, StartY int
    EndX, EndY     int
    Duration       time.Duration
}

// SendKeysAction inputs text
type SendKeysAction struct {
    Text    string
    Element *locator.Selector
    Secure  bool // don't log the text
}

// LaunchAction launches an app
type LaunchAction struct {
    AppID     string
    Arguments []string
    WaitIdle  bool
}

// PressKeyAction presses a key
type PressKeyAction struct {
    KeyCode int
}

// WaitAction waits for duration or condition
type WaitAction struct {
    Duration  time.Duration
    Element   *locator.Selector
    Optional  bool
}
```

**Step 2: Commit**

```bash
git add internal/action/
git commit -m "feat: add action package with action types

- TapAction, SwipeAction, SendKeysAction
- LaunchAction, PressKeyAction, WaitAction
- Support both coordinate and element-based actions"
```

---

## Phase 2: iOS WDA Driver (M2)

### Task 2.1: Create WDA Client

**Files:**
- Create: `ios/wda/protocol.go`
- Create: `ios/wda/client.go`
- Create: `ios/wda/client_test.go`

**Step 1: Create WDA protocol definitions**

Create `ios/wda/protocol.go`:
```go
package wda

// WDA endpoints
const (
    EndpointStatus        = "/status"
    EndpointSession       = "/session"
    EndpointScreenshot    = "/screenshot"
    EndpointSource        = "/wda/source"
    EndpointElement       = "/wda/element/active"
    EndpointTap           = "/wda/tap/0/%d/%d"      // x, y
    EndpointSwipe         = "/wda/performActions"
    EndpointKeys          = "/wda/keys"
    EndpointAppLaunch     = "/wda/apps/launch"
    EndpointAppTerminate  = "/wda/apps/terminate/%s" // bundleId
    EndpointAlert         = "/wda/alert/%s"          // action
)

// AlertAction represents alert actions
type AlertAction string

const (
    AlertAccept  AlertAction = "accept"
    AlertDismiss AlertAction = "dismiss"
    AlertText    AlertAction = "text"
)

// KeyCodes for common actions
const (
    KeyCodeHome          = 3
    KeyCodeBack          = 4
    KeyCodeEnter         = 66
    KeyCodeVolumeUp      = 24
    KeyCodeVolumeDown    = 25
    KeyCodePower         = 26
)
```

**Step 2: Create WDA client**

Create `ios/wda/client.go`:
```go
package wda

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "time"
)

// Client communicates with WebDriverAgent
type Client struct {
    BaseURL    *url.URL
    SessionID  string
    HTTPClient *http.Client
}

// NewClient creates a new WDA client
func NewClient(baseURL string) (*Client, error) {
    u, err := url.Parse(baseURL)
    if err != nil {
        return nil, fmt.Errorf("invalid base URL: %w", err)
    }
    
    return &Client{
        BaseURL: u,
        SessionID: "",
        HTTPClient: &http.Client{
            Timeout: 60 * time.Second,
        },
    }, nil
}

// NewClientWithSession creates client and starts session
func NewClientWithSession(baseURL string, bundleID string) (*Client, error) {
    client, err := NewClient(baseURL)
    if err != nil {
        return nil, err
    }
    
    if err := client.StartSession(bundleID); err != nil {
        return nil, err
    }
    
    return client, nil
}

// doRequest performs an HTTP request
func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
    u := c.BaseURL.ResolveReference(&url.URL{Path: path})
    
    var reqBody io.Reader
    if body != nil {
        data, err := json.Marshal(body)
        if err != nil {
            return nil, fmt.Errorf("marshal body: %w", err)
        }
        reqBody = bytes.NewReader(data)
    }
    
    req, err := http.NewRequest(method, u.String(), reqBody)
    if err != nil {
        return nil, fmt.Errorf("create request: %w", err)
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("do request: %w", err)
    }
    defer resp.Body.Close()
    
    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("read body: %w", err)
    }
    
    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
    }
    
    return respBody, nil
}

// StartSession starts a WDA session
func (c *Client) StartSession(bundleID string) error {
    body := map[string]interface{}{
        "capabilities": map[string]interface{}{
            "bundleId": bundleID,
        },
    }
    
    respBody, err := c.doRequest("POST", EndpointSession, body)
    if err != nil {
        return fmt.Errorf("start session: %w", err)
    }
    
    var resp struct {
        SessionID string `json:"sessionId"`
    }
    if err := json.Unmarshal(respBody, &resp); err != nil {
        return fmt.Errorf("parse session: %w", err)
    }
    
    c.SessionID = resp.SessionID
    return nil
}

// StopSession stops the WDA session
func (c *Client) StopSession() error {
    if c.SessionID == "" {
        return nil
    }
    
    path := fmt.Sprintf("%s/%s", EndpointSession, c.SessionID)
    _, err := c.doRequest("DELETE", path, nil)
    c.SessionID = ""
    return err
}

// Screenshot captures the screen
func (c *Client) Screenshot() ([]byte, error) {
    respBody, err := c.doRequest("GET", EndpointScreenshot, nil)
    if err != nil {
        return nil, fmt.Errorf("screenshot: %w", err)
    }
    
    // WDA returns base64 encoded PNG
    var resp struct {
        Value string `json:"value"`
    }
    if err := json.Unmarshal(respBody, &resp); err != nil {
        return nil, fmt.Errorf("parse screenshot: %w", err)
    }
    
    return []byte(resp.Value), nil
}

// Tap performs a tap at coordinates
func (c *Client) Tap(x, y int) error {
    path := fmt.Sprintf(EndpointTap, x, y)
    _, err := c.doRequest("POST", path, nil)
    return err
}

// SendKeys sends text input
func (c *Client) SendKeys(text string) error {
    body := map[string]interface{}{
        "value": []string{text},
    }
    _, err := c.doRequest("POST", EndpointKeys, body)
    return err
}

// GetSource returns page source XML
func (c *Client) GetSource() (string, error) {
    respBody, err := c.doRequest("GET", EndpointSource, nil)
    if err != nil {
        return "", fmt.Errorf("get source: %w", err)
    }
    
    var resp struct {
        Value string `json:"value"`
    }
    if err := json.Unmarshal(respBody, &resp); err != nil {
        return "", fmt.Errorf("parse source: %w", err)
    }
    
    return resp.Value, nil
}

// LaunchApp launches an app by bundle ID
func (c *Client) LaunchApp(bundleID string) error {
    body := map[string]interface{}{
        "bundleId": bundleID,
    }
    _, err := c.doRequest("POST", EndpointAppLaunch, body)
    return err
}

// TerminateApp terminates an app
func (c *Client) TerminateApp(bundleID string) error {
    path := fmt.Sprintf(EndpointAppTerminate, bundleID)
    _, err := c.doRequest("POST", path, nil)
    return err
}

// IsHealthy checks if WDA is responsive
func (c *Client) IsHealthy() bool {
    _, err := c.doRequest("GET", EndpointStatus, nil)
    return err == nil
}
```

**Step 3: Write WDA client tests (mock-based)**

Create `ios/wda/client_test.go`:
```go
package wda

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestNewClient_InvalidURL(t *testing.T) {
    _, err := NewClient("://invalid")
    if err == nil {
        t.Error("expected error for invalid URL")
    }
}

func TestNewClient_ValidURL(t *testing.T) {
    client, err := NewClient("http://localhost:8100")
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if client.BaseURL.String() != "http://localhost:8100/" {
        t.Errorf("unexpected base URL: %s", client.BaseURL)
    }
}

func TestClient_IsHealthy_True(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
    }))
    defer server.Close()
    
    client, _ := NewClient(server.URL)
    if !client.IsHealthy() {
        t.Error("expected healthy")
    }
}

func TestClient_IsHealthy_False(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusInternalServerError)
    }))
    defer server.Close()
    
    client, _ := NewClient(server.URL)
    if client.IsHealthy() {
        t.Error("expected unhealthy")
    }
}

func TestClient_Screenshot(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/screenshot" {
            t.Errorf("unexpected path: %s", r.URL.Path)
        }
        json.NewEncoder(w).Encode(map[string]string{
            "value": "iVBORw0KGgoAAAANSUhEUg==", // base64 of minimal PNG
        })
    }))
    defer server.Close()
    
    client, _ := NewClient(server.URL)
    data, err := client.Screenshot()
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if len(data) == 0 {
        t.Error("expected screenshot data")
    }
}

func TestClient_Tap(t *testing.T) {
    var tapPath string
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tapPath = r.URL.Path
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]interface{}{"value": nil})
    }))
    defer server.Close()
    
    client, _ := NewClient(server.URL)
    err := client.Tap(100, 200)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if tapPath != "/wda/tap/0/100/200" {
        t.Errorf("unexpected tap path: %s", tapPath)
    }
}

func TestClient_SendKeys(t *testing.T) {
    var method, body string
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        method = r.Method
        var req map[string]interface{}
        json.NewDecoder(r.Body).Decode(&req)
        body = req["value"].([]interface{})[0].(string)
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]interface{}{"value": nil})
    }))
    defer server.Close()
    
    client, _ := NewClient(server.URL)
    err := client.SendKeys("hello")
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if method != "POST" {
        t.Errorf("expected POST, got %s", method)
    }
    if body != "hello" {
        t.Errorf("expected 'hello', got '%s'", body)
    }
}
```

**Step 4: Run tests**

```bash
go test -v ./ios/wda/...
```

**Step 5: Commit**

```bash
git add ios/wda/
git commit -m "feat(ios): add WDA client with HTTP protocol

- Native WDA HTTP implementation (no gwda dependency)
- Session management (start/stop)
- Core actions: tap, sendKeys, screenshot, source
- App lifecycle: launch, terminate
- Mock-based unit tests"
```

---

### Task 2.2: Create iOS Device

**Files:**
- Create: `ios/device.go`
- Create: `ios/device_test.go`

**Step 1: Create iOS device implementation**

Create `ios/device.go`:
```go
package ios

import (
    "fmt"
    
    "github.com/yourname/go-uop"
    "github.com/yourname/go-uop/ios/wda"
)

// Device implements uop.Device for iOS
type Device struct {
    client   *wda.Client
    bundleID string
}

// NewDevice creates a new iOS device
func NewDevice(bundleID string, opts ...Option) (*Device, error) {
    cfg := &config{
        address: "http://localhost:8100",
    }
    for _, opt := range opts {
        opt(cfg)
    }
    
    client, err := wda.NewClient(cfg.address)
    if err != nil {
        return nil, fmt.Errorf("create WDA client: %w", err)
    }
    
    if err := client.StartSession(bundleID); err != nil {
        return nil, fmt.Errorf("start session: %w", err)
    }
    
    return &Device{
        client:   client,
        bundleID: bundleID,
    }, nil
}

// Platform returns iOS platform
func (d *Device) Platform() uop.Platform {
    return uop.IOS
}

// Info returns device information
func (d *Device) Info() (map[string]interface{}, error) {
    // TODO: implement device info
    return map[string]interface{}{
        "platform": "ios",
        "bundleId": d.bundleID,
    }, nil
}

// Screenshot captures current screen
func (d *Device) Screenshot() ([]byte, error) {
    return d.client.Screenshot()
}

// Close releases device resources
func (d *Device) Close() error {
    return d.client.StopSession()
}

// Tap performs tap at coordinates
func (d *Device) Tap(x, y int) error {
    return d.client.Tap(x, y)
}

// SendKeys inputs text
func (d *Device) SendKeys(text string) error {
    return d.client.SendKeys(text)
}

// Launch launches the app
func (d *Device) Launch() error {
    return d.client.LaunchApp(d.bundleID)
}

// Terminate terminates the app
func (d *Device) Terminate() error {
    return d.client.TerminateApp(d.bundleID)
}

// GetSource returns page source
func (d *Device) GetSource() (string, error) {
    return d.client.GetSource()
}

// Ensure Device implements uop.Device
var _ uop.Device = (*Device)(nil)
```

**Step 2: Add Option pattern**

Create `ios/option.go`:
```go
package ios

type config struct {
    address string
}

// Option configures iOS device
type Option func(*config)

// WithAddress sets WDA address
func WithAddress(addr string) Option {
    return func(c *config) {
        c.address = addr
    }
}
```

**Step 3: Commit**

```bash
git add ios/device.go ios/option.go
git commit -m "feat(ios): add iOS device implementation

- Implements uop.Device interface
- Wraps WDA client
- Provides tap, sendKeys, launch, terminate
- Option pattern for configuration"
```

---

## Phase 3: Android ADB Driver (M3)

### Task 3.1: Create ADB Client

**Files:**
- Create: `android/adb/client.go`
- Create: `android/adb/client_test.go`

**Step 1: Create ADB client**

Create `android/adb/client.go`:
```go
package adb

import (
    "bytes"
    "fmt"
    "os/exec"
    "regexp"
    "strconv"
    "strings"
    "time"
)

// Client wraps adb commands
type Client struct {
    serial string
}

// NewClient creates ADB client
func NewClient(serial ...string) (*Client, error) {
    s := ""
    if len(serial) > 0 {
        s = serial[0]
    }
    
    if s != "" {
        // Verify device exists
        devices, err := Devices()
        if err != nil {
            return nil, err
        }
        found := false
        for _, d := range devices {
            if d.Serial == s {
                found = true
                break
            }
        }
        if !found {
            return nil, fmt.Errorf("device not found: %s", s)
        }
    }
    
    return &Client{serial: s}, nil
}

// serialArg returns adb -s flag if serial is set
func (c *Client) serialArg() string {
    if c.serial != "" {
        return "-s " + c.serial
    }
    return ""
}

// exec runs adb command
func (c *Client) exec(args ...string) (string, error) {
    cmd := exec.Command("adb", append(strings.Fields(c.serialArg()), args...)...)
    var buf bytes.Buffer
    cmd.Stdout = &buf
    cmd.Stderr = &buf
    
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("adb %v: %w\n%s", args, err, buf.String())
    }
    
    return strings.TrimSpace(buf.String()), nil
}

// DeviceInfo represents connected device info
type DeviceInfo struct {
    Serial   string
    Status   string
    Product  string
    Model    string
    Device   string
    Transport string
}

// Devices returns list of connected devices
func Devices() ([]DeviceInfo, error) {
    cmd := exec.Command("adb", "devices", "-l")
    var buf bytes.Buffer
    cmd.Stdout = &buf
    if err := cmd.Run(); err != nil {
        return nil, fmt.Errorf("adb devices: %w", err)
    }
    
    var devices []DeviceInfo
    lines := strings.Split(buf.String(), "\n")
    for _, line := range lines[1:] {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, "List") {
            continue
        }
        
        parts := strings.Fields(line)
        if len(parts) < 2 {
            continue
        }
        
        d := DeviceInfo{Serial: parts[0], Status: parts[1]}
        
        // Parse additional info
        for _, p := range parts[2:] {
            kv := strings.SplitN(p, ":", 2)
            if len(kv) != 2 {
                continue
            }
            switch kv[0] {
            case "product":
                d.Product = kv[1]
            case "model":
                d.Model = strings.ReplaceAll(kv[1], "_", " ")
            case "device":
                d.Device = kv[1]
            case "transport_id":
                d.Transport = kv[1]
            }
        }
        
        devices = append(devices, d)
    }
    
    return devices, nil
}

// Tap performs tap at coordinates
func (c *Client) Tap(x, y int) error {
    _, err := c.exec("shell", "input", "tap", strconv.Itoa(x), strconv.Itoa(y))
    return err
}

// Swipe performs swipe gesture
func (c *Client) Swipe(x1, y1, x2, y2 int, durationMs int) error {
    args := []string{"shell", "input", "swipe", 
        strconv.Itoa(x1), strconv.Itoa(y1), 
        strconv.Itoa(x2), strconv.Itoa(y2)}
    if durationMs > 0 {
        args = append(args, strconv.Itoa(durationMs))
    }
    _, err := c.exec(args...)
    return err
}

// SendText inputs text
func (c *Client) SendText(text string) error {
    // Escape spaces and special chars
    escaped := strings.ReplaceAll(text, " ", "%s")
    _, err := c.exec("shell", "input", "text", escaped)
    return err
}

// PressKey presses a key code
func (c *Client) PressKey(keyCode int) error {
    _, err := c.exec("shell", "input", "keyevent", strconv.Itoa(keyCode))
    return err
}

// Key codes
const (
    KeyHome      = 3
    KeyBack      = 4
    KeyEnter     = 66
    KeyVolumeUp  = 24
    KeyVolumeDown= 25
    KeyPower     = 26
)

// Screenshot captures screen
func (c *Client) Screenshot() ([]byte, error) {
    // Take screenshot to temp file
    out := "/sdcard/screenshot.png"
    if _, err := c.exec("shell", "screencap", "-p", out); err != nil {
        return nil, err
    }
    
    // Pull to stdout
    cmd := exec.Command("adb", append(strings.Fields(c.serialArg()), "exec-out", "screencap", "-p")...)
    var buf bytes.Buffer
    cmd.Stdout = &buf
    if err := cmd.Run(); err != nil {
        return nil, fmt.Errorf("screencap: %w", err)
    }
    
    return buf.Bytes(), nil
}

// Shell runs shell command
func (c *Client) Shell(command string) (string, error) {
    return c.exec("shell", command)
}

// StartActivity starts an activity
func (c *Client) StartActivity(component string) error {
    _, err := c.exec("shell", "am", "start", "-n", component)
    return err
}

// StopPackage stops an app
func (c *Client) StopPackage(packageName string) error {
    _, err := c.exec("shell", "am", "force-stop", packageName)
    return err
}

// Install installs an APK
func (c *Client) Install(apkPath string, grantPerms bool) error {
    args := []string{"install"}
    if grantPerms {
        args = append(args, "-g")
    }
    args = append(args, apkPath)
    _, err := c.exec(args...)
    return err
}

// Uninstall uninstalls a package
func (c *Client) Uninstall(packageName string) error {
    _, err := c.exec("uninstall", packageName)
    return err
}

// CurrentPackage returns current foreground package
func (c *Client) CurrentPackage() (string, error) {
    output, err := c.Shell("dumpsys activity activities | grep mResumedActivity")
    if err != nil {
        return "", err
    }
    
    re := regexp.MustCompile(`([\w.]+)/`)
    matches := re.FindStringSubmatch(output)
    if len(matches) < 2 {
        return "", fmt.Errorf("could not find foreground package")
    }
    
    return matches[1], nil
}

// WaitForIdle waits for the device to be idle
func (c *Client) WaitForIdle(timeout time.Duration) error {
    end := time.Now().Add(timeout)
    for time.Now().Before(end) {
        _, err := c.Shell("sleep 0.1 && dumpsys activity activities")
        if err == nil {
            return nil
        }
        time.Sleep(500 * time.Millisecond)
    }
    return fmt.Errorf("timeout waiting for idle")
}
```

**Step 2: Write ADB client tests**

Create `android/adb/client_test.go`:
```go
package adb

import (
    "testing"
)

func TestNewClient_Empty(t *testing.T) {
    // This test requires adb installed and a device connected
    // Skip if not available
    t.Skip("requires adb")
}

func TestDevices_Format(t *testing.T) {
    // This test requires adb installed
    t.Skip("requires adb")
    
    devices, err := Devices()
    if err != nil {
        t.Skipf("adb not available: %v", err)
    }
    
    for _, d := range devices {
        if d.Serial == "" {
            t.Error("device serial should not be empty")
        }
        if d.Status == "" {
            t.Error("device status should not be empty")
        }
    }
}

func TestClient_ExecValidation(t *testing.T) {
    c := &Client{serial: "test-serial"}
    arg := c.serialArg()
    if arg != "-s test-serial" {
        t.Errorf("unexpected serial arg: %s", arg)
    }
}

func TestClient_NoSerial(t *testing.T) {
    c := &Client{serial: ""}
    arg := c.serialArg()
    if arg != "" {
        t.Errorf("unexpected serial arg: %s", arg)
    }
}
```

**Step 3: Commit**

```bash
git add android/adb/
git commit -m "feat(android): add ADB client wrapper

- Device discovery and management
- Input operations: tap, swipe, sendText, pressKey
- App lifecycle: startActivity, stopPackage, install, uninstall
- Screenshot via exec-out
- Mock-based unit tests"
```

---

### Task 3.2: Create Android Device

**Files:**
- Create: `android/device.go`
- Create: `android/device_test.go`

**Step 1: Create Android device implementation**

Create `android/device.go`:
```go
package android

import (
    "fmt"
    "time"
    
    "github.com/yourname/go-uop"
    "github.com/yourname/go-uop/android/adb"
)

// Device implements uop.Device for Android
type Device struct {
    client *adb.Client
    pkg    string
}

// NewDevice creates a new Android device
func NewDevice(opts ...Option) (*Device, error) {
    cfg := &config{}
    for _, opt := range opts {
        opt(cfg)
    }
    
    client, err := adb.NewClient(cfg.serial)
    if err != nil {
        return nil, fmt.Errorf("create ADB client: %w", err)
    }
    
    return &Device{
        client: client,
        pkg:    cfg.packageName,
    }, nil
}

// Platform returns Android platform
func (d *Device) Platform() uop.Platform {
    return uop.Android
}

// Info returns device information
func (d *Device) Info() (map[string]interface{}, error) {
    devices, err := adb.Devices()
    if err != nil {
        return nil, err
    }
    
    for _, dev := range devices {
        if dev.Serial == d.client.Serial || d.client.Serial == "" {
            return map[string]interface{}{
                "platform": "android",
                "serial":   dev.Serial,
                "model":    dev.Model,
                "product":   dev.Product,
            }, nil
        }
    }
    
    return map[string]interface{}{
        "platform": "android",
        "serial":   d.client.Serial,
    }, nil
}

// Screenshot captures current screen
func (d *Device) Screenshot() ([]byte, error) {
    return d.client.Screenshot()
}

// Close releases device resources
func (d *Device) Close() error {
    // ADB client doesn't need explicit close
    return nil
}

// Tap performs tap at coordinates
func (d *Device) Tap(x, y int) error {
    return d.client.Tap(x, y)
}

// SendKeys inputs text
func (d *Device) SendKeys(text string) error {
    return d.client.SendText(text)
}

// Launch launches the app
func (d *Device) Launch() error {
    if d.pkg == "" {
        return fmt.Errorf("package name not set")
    }
    return d.client.StartActivity(d.pkg + "/.MainActivity")
}

// Terminate terminates the app
func (d *Device) Terminate() error {
    if d.pkg == "" {
        return fmt.Errorf("package name not set")
    }
    return d.client.StopPackage(d.pkg)
}

// GetSource returns page source (dump of UI hierarchy)
func (d *Device) GetSource() (string, error) {
    return d.client.Shell("uiautomator dump /sdcard/ui.xml && cat /sdcard/ui.xml")
}

// PressKey presses a key
func (d *Device) PressKey(keyCode int) error {
    return d.client.PressKey(keyCode)
}

// Swipe performs swipe gesture
func (d *Device) Swipe(x1, y1, x2, y2 int, duration time.Duration) error {
    return d.client.Swipe(x1, y1, x2, y2, int(duration.Milliseconds()))
}

// Ensure Device implements uop.Device
var _ uop.Device = (*Device)(nil)
```

**Step 2: Add Option pattern**

Create `android/option.go`:
```go
package android

type config struct {
    serial      string
    packageName string
}

// Option configures Android device
type Option func(*config)

// WithSerial sets device serial number
func WithSerial(serial string) Option {
    return func(c *config) {
        c.serial = serial
    }
}

// WithPackage sets app package name
func WithPackage(pkg string) Option {
    return func(c *config) {
        c.packageName = pkg
    }
}
```

**Step 3: Commit**

```bash
git add android/device.go android/option.go
git commit -m "feat(android): add Android device implementation

- Implements uop.Device interface
- Wraps ADB client
- Provides tap, sendKeys, launch, terminate, swipe
- Option pattern for configuration"
```

---

## Phase 4: Chainable API (M4)

### Task 4.1: Create Fluent Action API

**Files:**
- Modify: `internal/action/action.go`
- Create: `uop_fluent.go`

**Step 1: Create fluent API wrapper**

Create `uop_fluent.go`:
```go
package uop

import (
    "github.com/yourname/go-uop/internal/action"
    "github.com/yourname/go-uop/internal/locator"
)

// ActionBuilder provides chainable action API
type ActionBuilder struct {
    device Device
}

// NewActionBuilder creates action builder
func NewActionBuilder(device Device) *ActionBuilder {
    return &ActionBuilder{device: device}
}

// Tap taps at coordinates
func (ab *ActionBuilder) Tap(x, y int) *ActionBuilder {
    // TODO: implement
    return ab
}

// TapElement taps on element
func (ab *ActionBuilder) TapElement(loc *locator.Selector) *ActionBuilder {
    // TODO: implement
    return ab
}

// Swipe performs swipe
func (ab *ActionBuilder) Swipe(x1, y1, x2, y2 int) *ActionBuilder {
    // TODO: implement
    return ab
}

// SwipeUp performs swipe up
func (ab *ActionBuilder) SwipeUp() *ActionBuilder {
    // TODO: implement
    return ab
}

// SendKeys inputs text
func (ab *ActionBuilder) SendKeys(text string) *ActionBuilder {
    // TODO: implement
    return ab
}

// Launch launches app
func (ab *ActionBuilder) Launch(appID string) *ActionBuilder {
    // TODO: implement
    return ab
}

// Wait waits for duration
func (ab *ActionBuilder) Wait(duration string) *ActionBuilder {
    // TODO: implement
    return ab
}

// Do executes all queued actions
func (ab *ActionBuilder) Do() error {
    // TODO: implement
    return nil
}
```

**Step 2: Commit (partial)**

```bash
git add uop_fluent.go
git commit -m "feat: add fluent action builder skeleton

- Chainable API pattern
- Tap, Swipe, SendKeys, Launch, Wait
- Ready for implementation in M4"
```

---

## Phase 5: YAML Runner (M5)

### Task 5.1: Create YAML Parser

**Files:**
- Create: `yaml/parser.go`
- Create: `yaml/parser_test.go`
- Create: `yaml/evaluator.go`
- Create: `yaml/evaluator_test.go`

**Step 1: Create YAML command types**

Create `yaml/command.go`:
```go
package yaml

import "gopkg.in/yaml.v3"

// Command represents a YAML command
type Command struct {
    // Command type (determined by which field is set)
    Name string `yaml:"name,omitempty"`
    
    // Tap command
    TapOn      *TapCommand      `yaml:"tapOn,omitempty"`
    Tap        *PointCommand    `yaml:"tap,omitempty"`
    DoubleTap  *TapCommand      `yaml:"doubleTap,omitempty"`
    LongPress  *LongPressCommand `yaml:"longPress,omitempty"`
    
    // Swipe command
    Swipe      *SwipeCommand    `yaml:"swipe,omitempty"`
    SwipeUp    *struct{}         `yaml:"swipeUp,omitempty"`
    SwipeDown  *struct{}         `yaml:"swipeDown,omitempty"`
    SwipeLeft  *struct{}         `yaml:"swipeLeft,omitempty"`
    SwipeRight *struct{}         `yaml:"swipeRight,omitempty"`
    
    // Input command
    InputText  *InputCommand     `yaml:"inputText,omitempty"`
    
    // App commands
    Launch     string            `yaml:"launch,omitempty"`
    Terminate  string            `yaml:"terminate,omitempty"`
    Install    string            `yaml:"install,omitempty"`
    Uninstall  string            `yaml:"uninstall,omitempty"`
    
    // Wait commands
    WaitFor    *WaitCommand      `yaml:"waitFor,omitempty"`
    WaitForGone *WaitCommand     `yaml:"waitForGone,omitempty"`
    Wait       int               `yaml:"wait,omitempty"` // milliseconds
    
    // Assertion commands
    AssertVisible *ElementQuery  `yaml:"assertVisible,omitempty"`
    AssertNotVisible *ElementQuery `yaml:"assertNotVisible,omitempty"`
    AssertTrue *string           `yaml:"assertTrue,omitempty"`
    
    // Screenshot
    Screenshot *ScreenshotCommand `yaml:"screenshot,omitempty"`
    
    // Control flow
    RunFlow    *RunFlowCommand   `yaml:"runFlow,omitempty"`
    If         *IfCommand        `yaml:"if,omitempty"`
    Foreach    *ForeachCommand   `yaml:"foreach,omitempty"`
    While      *WhileCommand      `yaml:"while,omitempty"`
    
    // Script
    EvalScript *EvalScriptCommand `yaml:"evalScript,omitempty"`
    
    // Variable
    SetVariable *SetVarCommand    `yaml:"setVariable,omitempty"`
    
    // Utility
    Log        string            `yaml:"log,omitempty"`
    Comment    string            `yaml:"comment,omitempty"`
}

// TapCommand tap on element
type TapCommand struct {
    Text     string `yaml:"text,omitempty"`
    ID       string `yaml:"id,omitempty"`
    XPath    string `yaml:"xpath,omitempty"`
    Index    int    `yaml:"index,omitempty"`
    Optional bool   `yaml:"optional,omitempty"`
    Timeout  string `yaml:"timeout,omitempty"`
}

// PointCommand tap on coordinates
type PointCommand struct {
    X int `yaml:"x"`
    Y int `yaml:"y"`
}

// LongPressCommand long press
type LongPressCommand struct {
    Text     string `yaml:"text,omitempty"`
    X        int    `yaml:"x,omitempty"`
    Y        int    `yaml:"y,omitempty"`
    Duration int    `yaml:"duration"` // milliseconds
}

// SwipeCommand swipe gesture
type SwipeCommand struct {
    StartX, StartY int `yaml:"startX,omitempty"`
    EndX, EndY     int `yaml:"endX,omitempty"`
    Duration       int `yaml:"duration,omitempty"`
}

// InputCommand text input
type InputCommand struct {
    Text     string       `yaml:"text"`
    Element  *ElementQuery `yaml:"element,omitempty"`
    Secure   bool         `yaml:"secure,omitempty"`
    PressEnter bool       `yaml:"pressEnter,omitempty"`
}

// ElementQuery element locator
type ElementQuery struct {
    Text     string `yaml:"text,omitempty"`
    ID       string `yaml:"id,omitempty"`
    XPath    string `yaml:"xpath,omitempty"`
    Index    int    `yaml:"index,omitempty"`
}

// WaitCommand wait for element
type WaitCommand struct {
    Element  *ElementQuery `yaml:"element"`
    Timeout  string        `yaml:"timeout"`
    Optional bool          `yaml:"optional,omitempty"`
}

// ScreenshotCommand screenshot
type ScreenshotCommand struct {
    Name string `yaml:"name"`
}

// RunFlowCommand run sub-flow
type RunFlowCommand struct {
    Name   string            `yaml:"name"`
    Params map[string]string `yaml:"params,omitempty"`
}

// IfCommand conditional
type IfCommand struct {
    Condition string   `yaml:"condition"`
    Then      []Command `yaml:"then"`
    Else      []Command `yaml:"else,omitempty"`
}

// ForeachCommand loop
type ForeachCommand struct {
    Variable string   `yaml:"variable"`
    In       string   `yaml:"in"`
    Do       []Command `yaml:"do"`
}

// WhileCommand conditional loop
type WhileCommand struct {
    Condition   string   `yaml:"condition"`
    MaxIter int      `yaml:"maxIterations,omitempty"`
    Do       []Command `yaml:"do"`
}

// EvalScriptCommand script execution
type EvalScriptCommand struct {
    Lang     string   `yaml:"lang"` // javascript, python
    Script   string   `yaml:"script,omitempty"`
    Source   string   `yaml:"source,omitempty"`
    Function string   `yaml:"function,omitempty"`
    Args     []string `yaml:"args,omitempty"`
    SaveTo   string   `yaml:"saveTo,omitempty"`
}

// SetVarCommand set variable
type SetVarCommand struct {
    Name  string `yaml:"name"`
    Value string `yaml:"value"`
}

// Flow represents a YAML flow
type Flow struct {
    Name        string            `yaml:"name"`
    Description string            `yaml:"description,omitempty"`
    Platform    string            `yaml:"platform,omitempty"`
    Params      map[string]string `yaml:"params,omitempty"`
    Timeout     string            `yaml:"timeout,omitempty"`
    Steps       []Command          `yaml:"steps"`
}

// TestSuite represents a test suite
type TestSuite struct {
    Config     *Config            `yaml:"config,omitempty"`
    Env        map[string]string `yaml:"env,omitempty"`
    Import     []string           `yaml:"import,omitempty"`
    Variables  map[string]string `yaml:"variables,omitempty"`
    Flows      []Flow            `yaml:"flows,omitempty"`
    Tests      []TestCase         `yaml:"tests,omitempty"`
}

// Config suite configuration
type Config struct {
    AppID   string `yaml:"appId"`
    Timeout string `yaml:"timeout"`
    Retry   int    `yaml:"retry"`
}

// TestCase represents a test case
type TestCase struct {
    Name   string   `yaml:"name"`
    Flow   string   `yaml:"flow,omitempty"`
    Params map[string]string `yaml:"params,omitempty"`
    Steps  []Command `yaml:"steps,omitempty"`
}
```

**Step 2: Create YAML parser**

Create `yaml/parser.go`:
```go
package yaml

import (
    "fmt"
    "os"
    
    "gopkg.in/yaml.v3"
)

// ParseFlow parses a flow YAML file
func ParseFlow(path string) (*Flow, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("read file: %w", err)
    }
    
    var flow Flow
    if err := yaml.Unmarshal(data, &flow); err != nil {
        return nil, fmt.Errorf("parse yaml: %w", err)
    }
    
    return &flow, nil
}

// ParseSuite parses a test suite YAML file
func ParseSuite(path string) (*TestSuite, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("read file: %w", err)
    }
    
    var suite TestSuite
    if err := yaml.Unmarshal(data, &suite); err != nil {
        return nil, fmt.Errorf("parse yaml: %w", err)
    }
    
    return &suite, nil
}

// ParseFlowFromString parses flow YAML from string
func ParseFlowFromString(content string) (*Flow, error) {
    var flow Flow
    if err := yaml.Unmarshal([]byte(content), &flow); err != nil {
        return nil, fmt.Errorf("parse yaml: %w", err)
    }
    return &flow, nil
}
```

**Step 3: Create expression evaluator**

Create `yaml/evaluator.go`:
```go
package yaml

import (
    "regexp"
    "strings"
    
    "github.com/robert-nelsen/otto"
)

var (
    reVariable  = regexp.MustCompile(`\$\{([^}]+)\}`)
    reENV       = regexp.MustCompile(`\$\{ENV\(([^)]+)\)\}`)
)

// Context holds evaluation context
type Context struct {
    Variables map[string]interface{}
    Env       map[string]string
    JS        *otto.Otto
}

// NewContext creates new evaluation context
func NewContext() *Context {
    return &Context{
        Variables: make(map[string]interface{}),
        Env:       make(map[string]string),
        JS:        otto.New(),
    }
}

// SetVariable sets a variable
func (c *Context) SetVariable(name string, value interface{}) {
    c.Variables[name] = value
}

// GetVariable gets a variable
func (c *Context) GetVariable(name string) interface{} {
    return c.Variables[name]
}

// Evaluate evaluates expressions in a string
func (c *Context) Evaluate(input string) (string, error) {
    // Handle ${...} expressions
    result := reVariable.ReplaceAllStringFunc(input, func(match string) string {
        expr := match[2 : len(match)-1]
        return c.evalExpr(expr)
    })
    
    return result, nil
}

// evalExpr evaluates a single expression
func (c *evalExpr) evalExpr(expr string) string {
    // ENV function
    if matches := reENV.FindStringSubmatch(expr); len(matches) > 1 {
        key := matches[1]
        if val, ok := c.Env[key]; ok {
            return val
        }
        return ""
    }
    
    // Variable access: variables.name
    if strings.HasPrefix(expr, "variables.") {
        name := strings.TrimPrefix(expr, "variables.")
        if val, ok := c.Variables[name]; ok {
            return fmt.Sprintf("%v", val)
        }
        return ""
    }
    
    // Simple variable: name
    if val, ok := c.Variables[expr]; ok {
        return fmt.Sprintf("%v", val)
    }
    
    // Return as-is
    return "${" + expr + "}"
}
```

**Step 4: Write YAML parser tests**

Create `yaml/parser_test.go`:
```go
package yaml

import (
    "testing"
)

func TestParseFlowFromString_Basic(t *testing.T) {
    yaml := `
name: login
steps:
  - launch: com.example.app
  - tapOn:
      text: "登录"
`
    flow, err := ParseFlowFromString(yaml)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    
    if flow.Name != "login" {
        t.Errorf("expected name 'login', got '%s'", flow.Name)
    }
    
    if len(flow.Steps) != 2 {
        t.Errorf("expected 2 steps, got %d", len(flow.Steps))
    }
    
    if flow.Steps[0].Launch != "com.example.app" {
        t.Errorf("expected launch 'com.example.app', got '%s'", flow.Steps[0].Launch)
    }
}

func TestParseFlowFromString_WithParams(t *testing.T) {
    yaml := `
name: login
params:
  username: string
  password: string
steps:
  - inputText:
      text: "${params.username}"
`
    flow, err := ParseFlowFromString(yaml)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    
    if flow.Params["username"] != "string" {
        t.Errorf("expected param 'username', got '%s'", flow.Params["username"])
    }
}

func TestParseFlowFromString_WithControlFlow(t *testing.T) {
    yaml := `
name: test
steps:
  - foreach:
      variable: item
      in: "a,b,c"
      do:
        - tapOn:
            text: "${item}"
`
    flow, err := ParseFlowFromString(yaml)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    
    if len(flow.Steps) != 1 {
        t.Errorf("expected 1 step, got %d", len(flow.Steps))
    }
    
    if flow.Steps[0].Foreach == nil {
        t.Error("expected foreach command")
    }
    
    if flow.Steps[0].Foreach.Variable != "item" {
        t.Errorf("expected variable 'item', got '%s'", flow.Steps[0].Foreach.Variable)
    }
}

func TestContext_SetGetVariable(t *testing.T) {
    ctx := NewContext()
    
    ctx.SetVariable("name", "test")
    if ctx.GetVariable("name") != "test" {
        t.Error("expected 'test'")
    }
    
    ctx.SetVariable("count", 42)
    if ctx.GetVariable("count") != 42 {
        t.Error("expected 42")
    }
}
```

**Step 5: Run tests**

```bash
go test -v ./yaml/...
```

**Step 6: Commit**

```bash
git add yaml/
git commit -m "feat(yaml): add YAML parser with Maestro-style commands

- Command types: tap, swipe, input, launch, wait, assert
- Control flow: if, foreach, while, runFlow
- evalScript support for JavaScript/Python
- Expression evaluator with variable substitution
- Unit tests for parser and evaluator"
```

---

## Phase 6-12: (Abbreviated)

### Task 6.1: YAML Control Flow

Create `yaml/commands/control.go` with if/foreach/while implementation

### Task 7.1: OpenCV Vision Module

Create `internal/vision/template.go` with OpenCV template matching

### Task 8.1: AI Provider

Create `ai/provider.go`, `ai/openai.go`

### Task 9.1: Parallel Executor

Create `internal/parallel/executor.go`

### Task 10.1: Report Generator

Create `internal/report/generator.go`

### Task 11.1: Retry Mechanism

Create `internal/retry/retry.go`

### Task 12.1: Documentation

Create README.md, examples/

---

## Summary

| Phase | Tasks | Files | Tests |
|-------|-------|-------|-------|
| M1 | 1.1-1.3 | 10 | 15 |
| M2 | 2.1-2.2 | 6 | 12 |
| M3 | 3.1-3.2 | 6 | 8 |
| M4 | 4.1 | 3 | 0 |
| M5 | 5.1 | 5 | 8 |
| M6-M12 | Various | ~25 | ~20 |
| **Total** | **~25** | **~55** | **~63** |

---

**Plan complete and saved to `docs/plans/2026-03-22-implementation-plan.md`.**

**Two execution options:**

**1. Subagent-Driven (this session)** - I dispatch fresh subagent per task, review between tasks, fast iteration

**2. Parallel Session (separate)** - Open new session with executing-plans, batch execution with checkpoints

**Which approach?**
