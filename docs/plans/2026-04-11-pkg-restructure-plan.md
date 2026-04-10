# pkg Restructure: ios + wda Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Restructure pkg directory to use adb CLI for Android, go-ios library for iOS device management, and W3C WebDriver for iOS UI automation.

**Architecture:**
- `pkg/adb/` - Android device via adb CLI (unchanged)
- `pkg/ios/` - iOS device management via go-ios library
- `pkg/wda/` - iOS UI automation via W3C WebDriver (moved from pkg/ios/wda, refactored for W3C compliance)
- `pkg/uiautomator2/` - Android UI automation via JSON-RPC (unchanged)

**Tech Stack:**
- go-ios v1.0.206 (already in go.mod as indirect)
- W3C WebDriver protocol
- Go 1.22+

---

## Task 1: Create pkg/wda/ with W3C WebDriver Compliance

**Files:**
- Create: `pkg/wda/client.go` (refactored)
- Create: `pkg/wda/protocol.go`
- Create: `pkg/wda/session.go`
- Create: `pkg/wda/element.go`
- Create: `pkg/wda/app.go`
- Create: `pkg/wda/alert.go`
- Create: `pkg/wda/screenshot.go`
- Create: `pkg/wda/w3c.go` (new - W3C types)
- Delete: `pkg/ios/wda/` (old location)

**Step 1: Create W3C type definitions**

```go
// pkg/wda/w3c.go
package wda

// W3C WebDriver error codes
type ErrorCode string

const (
    ErrElementNotInteractable ErrorCode = "element not interactable"
    ErrInvalidSessionID       ErrorCode = "invalid session id"
    ErrNoSuchElement          ErrorCode = "no such element"
    ErrStaleElementReference   ErrorCode = "stale element reference"
    ErrElementClickIntercepted ErrorCode = "element click intercepted"
)

// W3C Response envelope
type W3CResponse struct {
    Value    any    `json:"value,omitempty"`
    SessionID string `json:"sessionId,omitempty"`
}

// W3C Capabilities (alwaysMatch/firstMatch)
type Capabilities struct {
    AlwaysMatch map[string]any `json:"alwaysMatch,omitempty"`
    FirstMatch  []map[string]any `json:"firstMatch,omitempty"`
}

// NewSessionRequest for POST /session
type NewSessionRequest struct {
    Capabilities Capabilities `json:"capabilities"`
}

// NewSessionResponse from POST /session
type NewSessionResponse struct {
    SessionID  string            `json:"sessionId"`
    Capabilities map[string]any `json:"capabilities"`
}

// Element reference (W3C style)
type ElementID string

// FindElementRequest for POST /session/{id}/element
type FindElementRequest struct {
    Using string `json:"using"`
    Value string `json:"value"`
}

// Element rect response
type ElementRect struct {
    X      float64 `json:"x"`
    Y      float64 `json:"y"`
    Width  float64 `json:"width"`
    Height float64 `json:"height"`
}
```

**Step 2: Refactor client.go with proper W3C session handling**

```go
// pkg/wda/client.go
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

type Client struct {
    BaseURL    *url.URL
    SessionID  string
    Capabilities map[string]any
    HTTPClient *http.Client
}

func NewClient(baseURL string) (*Client, error) {
    u, err := url.Parse(baseURL)
    if err != nil {
        return nil, fmt.Errorf("invalid base URL: %w", err)
    }

    return &Client{
        BaseURL:   u,
        HTTPClient: &http.Client{
            Timeout: 60 * time.Second,
        },
    }, nil
}

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
    if c.SessionID != "" {
        req.Header.Set("X-Session-Id", c.SessionID)
    }

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
        var errResp struct {
            Value struct {
                Error   string `json:"error"`
                Message string `json:"message"`
            } `json:"value"`
        }
        json.Unmarshal(respBody, &errResp)
        return nil, fmt.Errorf("HTTP %d: %s - %s", resp.StatusCode, errResp.Value.Error, errResp.Value.Message)
    }

    return respBody, nil
}

func (c *Client) IsHealthy() bool {
    _, err := c.doRequest("GET", EndpointStatus, nil)
    return err == nil
}

func (c *Client) Close() error {
    if c.SessionID == "" {
        return nil
    }
    path := fmt.Sprintf("%s/%s", EndpointSession, c.SessionID)
    _, err := c.doRequest("DELETE", path, nil)
    c.SessionID = ""
    return err
}
```

**Step 3: Refactor session.go with W3C capabilities**

```go
// pkg/wda/session.go
package wda

import (
    "encoding/json"
    "fmt"
)

// StartSession creates a new W3C WebDriver session
func (c *Client) StartSession(bundleID string) error {
    body := NewSessionRequest{
        Capabilities: Capabilities{
            AlwaysMatch: map[string]any{
                "platformName":     "ios",
                "browserName":      "Safari",
                "appium:automationName": "XCUITest",
                "appium:bundleId":  bundleID,
            },
            FirstMatch: []map[string]any{
                {},
            },
        },
    }

    respBody, err := c.doRequest("POST", EndpointSession, body)
    if err != nil {
        return fmt.Errorf("start session: %w", err)
    }

    var w3cResp W3CResponse
    if err := json.Unmarshal(respBody, &w3cResp); err != nil {
        return fmt.Errorf("parse session response: %w", err)
    }

    // Extract session ID from value object (W3C style)
    if w3cResp.SessionID != "" {
        c.SessionID = w3cResp.SessionID
    } else if m, ok := w3cResp.Value.(map[string]any); ok {
        if sid, ok := m["sessionId"].(string); ok {
            c.SessionID = sid
        }
        if caps, ok := m["capabilities"].(map[string]any); ok {
            c.Capabilities = caps
        }
    }

    if c.SessionID == "" {
        return fmt.Errorf("no sessionId in response")
    }
    return nil
}

// StopSession terminates the current session
func (c *Client) StopSession() error {
    return c.Close()
}
```

**Step 4: Update protocol.go with W3C endpoints**

```go
// pkg/wda/protocol.go
package wda

// W3C WebDriver endpoints
const (
    EndpointStatus       = "/status"
    EndpointSession     = "/session"
    EndpointScreenshot  = "/screenshot"
    EndpointElement     = "/element"           // W3C style
    EndpointElements    = "/elements"          // W3C style
    EndpointSource      = "/source"
    EndpointTap         = "/touch/perform"     // W3C touch
    EndpointClick       = "/element/{elementId}/click"  // W3C
    EndpointSendKeys    = "/element/{elementId}/value"  // W3C
    EndpointRect        = "/element/{elementId}/rect"   // W3C
    EndpointAppLaunch   = "/session/{sessionId}/app/launch"  // Appium
    EndpointAppTerminate= "/session/{sessionId}/app/terminate"
    EndpointAlert       = "/alert"
    EndpointTimeouts    = "/timeouts"
)

// Location strategies for Find Element
const (
    UsingCSSSelector    = "css selector"
    UsingXPath          = "xpath"
    UsingTagName        = "tag name"
    UsingLinkText       = "link text"
    UsingPartialLinkText= "partial link text"
    UsingClassName      = "class name"
    UsingID             = "id"
    UsingName           = "name"
)
```

**Step 5: Implement element.go with FindElement/FindElements**

```go
// pkg/wda/element.go
package wda

import (
    "encoding/json"
    "fmt"
    "strings"
)

// FindElement locates a single element using W3C strategy
func (c *Client) FindElement(strategy, selector string) (string, error) {
    if c.SessionID == "" {
        return "", fmt.Errorf("no session")
    }
    
    body := FindElementRequest{
        Using: strategy,
        Value: selector,
    }
    
    path := fmt.Sprintf("%s/%s/element", EndpointSession, c.SessionID)
    respBody, err := c.doRequest("POST", path, body)
    if err != nil {
        return "", err
    }
    
    var w3cResp struct {
        Value map[string]string `json:"value"`
    }
    if err := json.Unmarshal(respBody, &w3cResp); err != nil {
        return "", fmt.Errorf("parse element response: %w", err)
    }
    
    // W3C returns {"element-6066-11e4-a012-44a0-8f36-f38267b8d19": "id"}
    for _, v := range w3cResp.Value {
        return v, nil
    }
    
    return "", fmt.Errorf("element not found")
}

// FindElements locates multiple elements
func (c *Client) FindElements(strategy, selector string) ([]string, error) {
    if c.SessionID == "" {
        return nil, fmt.Errorf("no session")
    }
    
    body := FindElementRequest{
        Using: strategy,
        Value: selector,
    }
    
    path := fmt.Sprintf("%s/%s/elements", EndpointSession, c.SessionID)
    respBody, err := c.doRequest("POST", path, body)
    if err != nil {
        return nil, err
    }
    
    var w3cResp struct {
        Value []map[string]string `json:"value"`
    }
    if err := json.Unmarshal(respBody, &w3cResp); err != nil {
        return nil, fmt.Errorf("parse elements response: %w", err)
    }
    
    var ids []string
    for _, elem := range w3cResp.Value {
        for _, v := range elem {
            ids = append(ids, v)
        }
    }
    
    return ids, nil
}

// GetSource returns page source
func (c *Client) GetSource() (string, error) {
    if c.SessionID == "" {
        return "", fmt.Errorf("no session")
    }
    
    path := fmt.Sprintf("%s/%s/source", EndpointSession, c.SessionID)
    respBody, err := c.doRequest("GET", path, nil)
    if err != nil {
        return "", fmt.Errorf("get source: %w", err)
    }
    
    var w3cResp W3CResponse
    if err := json.Unmarshal(respBody, &w3cResp); err != nil {
        return "", fmt.Errorf("parse source: %w", err)
    }
    
    if s, ok := w3cResp.Value.(string); ok {
        return s, nil
    }
    return "", nil
}

// Click an element (W3C)
func (c *Client) Click(elementID string) error {
    if c.SessionID == "" {
        return fmt.Errorf("no session")
    }
    
    path := fmt.Sprintf("%s/%s/element/%s/click", EndpointSession, c.SessionID, elementID)
    _, err := c.doRequest("POST", path, nil)
    return err
}

// SendKeys to an element (W3C)
func (c *Client) SendKeys(elementID, text string) error {
    if c.SessionID == "" {
        return fmt.Errorf("no session")
    }
    
    path := fmt.Sprintf("%s/%s/element/%s/value", EndpointSession, c.SessionID, elementID)
    body := map[string]any{
        "text": text,
    }
    _, err := c.doRequest("POST", path, body)
    return err
}

// GetElementRect returns element bounds
func (c *Client) GetElementRect(elementID string) (*ElementRect, error) {
    if c.SessionID == "" {
        return nil, fmt.Errorf("no session")
    }
    
    path := fmt.Sprintf("%s/%s/element/%s/rect", EndpointSession, c.SessionID, elementID)
    respBody, err := c.doRequest("GET", path, nil)
    if err != nil {
        return nil, err
    }
    
    var w3cResp struct {
        Value ElementRect `json:"value"`
    }
    if err := json.Unmarshal(respBody, &w3cResp); err != nil {
        return nil, fmt.Errorf("parse rect: %w", err)
    }
    
    return &w3cResp.Value, nil
}

// Tap at coordinates (WDA-specific)
func (c *Client) Tap(x, y int) error {
    if c.SessionID == "" {
        return fmt.Errorf("no session")
    }
    
    path := fmt.Sprintf("%s/%s/wda/touch/perform", EndpointSession, c.SessionID)
    body := map[string]any{
        "actions": []map[string]any{
            {
                "action": "tap",
                "options": map[string]int{
                    "x": x,
                    "y": y,
                },
            },
        },
    }
    _, err := c.doRequest("POST", path, body)
    return err
}

// Legacy tap for backward compat (WDA-style)
func (c *Client) TapLegacy(x, y int) error {
    path := fmt.Sprintf("/wda/tap/0/%d/%d", x, y)
    _, err := c.doRequest("POST", path, nil)
    return err
}

// SendKeysLegacy types text (WDA-style)
func (c *Client) SendKeysLegacy(text string) error {
    body := map[string]interface{}{
        "value": []string{text},
    }
    _, err := c.doRequest("POST", "/wda/keys", body)
    return err
}
```

**Step 6: Update app.go**

```go
// pkg/wda/app.go
package wda

import "fmt"

// LaunchApp launches an app by bundle ID
func (c *Client) LaunchApp(bundleID string) error {
    if c.SessionID == "" {
        return fmt.Errorf("no session")
    }
    
    path := fmt.Sprintf("%s/%s/app/launch", EndpointSession, c.SessionID)
    body := map[string]interface{}{
        "bundleId": bundleID,
    }
    _, err := c.doRequest("POST", path, body)
    return err
}

// TerminateApp terminates an app
func (c *Client) TerminateApp(bundleID string) error {
    if c.SessionID == "" {
        return fmt.Errorf("no session")
    }
    
    path := fmt.Sprintf("%s/%s/app/terminate", EndpointSession, c.SessionID)
    body := map[string]interface{}{
        "bundleId": bundleID,
    }
    _, err := c.doRequest("POST", path, body)
    return err
}
```

**Step 7: Update alert.go**

```go
// pkg/wda/alert.go
package wda

import (
    "encoding/json"
    "fmt"
)

type AlertAction string

const (
    AlertAccept  AlertAction = "accept"
    AlertDismiss AlertAction = "dismiss"
    AlertText    AlertAction = "getText"
)

func (c *Client) AcceptAlert() error {
    if c.SessionID == "" {
        return fmt.Errorf("no session")
    }
    
    path := fmt.Sprintf("%s/%s/alert/accept", EndpointSession, c.SessionID)
    _, err := c.doRequest("POST", path, nil)
    return err
}

func (c *Client) DismissAlert() error {
    if c.SessionID == "" {
        return fmt.Errorf("no session")
    }
    
    path := fmt.Sprintf("%s/%s/alert/dismiss", EndpointSession, c.SessionID)
    _, err := c.doRequest("POST", path, nil)
    return err
}

func (c *Client) GetAlertText() (string, error) {
    if c.SessionID == "" {
        return "", fmt.Errorf("no session")
    }
    
    path := fmt.Sprintf("%s/%s/alert/text", EndpointSession, c.SessionID)
    respBody, err := c.doRequest("GET", path, nil)
    if err != nil {
        return "", fmt.Errorf("get alert text: %w", err)
    }
    
    var w3cResp W3CResponse
    if err := json.Unmarshal(respBody, &w3cResp); err != nil {
        return "", fmt.Errorf("parse alert text: %w", err)
    }
    
    if s, ok := w3cResp.Value.(string); ok {
        return s, nil
    }
    return "", nil
}
```

**Step 8: Update screenshot.go**

```go
// pkg/wda/screenshot.go
package wda

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
)

func (c *Client) Screenshot() ([]byte, error) {
    if c.SessionID == "" {
        return nil, fmt.Errorf("no session")
    }
    
    path := fmt.Sprintf("%s/%s/screenshot", EndpointSession, c.SessionID)
    respBody, err := c.doRequest("GET", path, nil)
    if err != nil {
        return nil, fmt.Errorf("screenshot: %w", err)
    }
    
    var w3cResp struct {
        Value string `json:"value"`
    }
    if err := json.Unmarshal(respBody, &w3cResp); err != nil {
        return nil, fmt.Errorf("parse screenshot: %w", err)
    }
    
    return base64.StdEncoding.DecodeString(w3cResp.Value)
}
```

**Step 9: Delete old pkg/ios/wda/ directory**

```bash
rm -rf pkg/ios/wda/
```

**Step 10: Commit**

```bash
git add -A
git commit -m "refactor(wda): move to pkg/wda and add W3C WebDriver compliance

- Move pkg/ios/wda to pkg/wda
- Add W3C session creation with alwaysMatch/firstMatch capabilities
- Add proper response envelope parsing (value.sessionId)
- Add FindElement/FindElements with W3C location strategies
- Add X-Session-Id header for authenticated requests
- Add GetElementRect, Click, Element SendKeys
- Update all endpoints to include session ID in path
- Add W3C error response parsing"
```

---

## Task 2: Create pkg/ios/ using go-ios library

**Files:**
- Create: `pkg/ios/device.go` (new implementation)
- Create: `pkg/ios/option.go`
- Create: `pkg/ios/wda_client.go` (delegates UI ops to wda)
- Create: `pkg/ios/screenshot.go`
- Create: `pkg/ios/app.go`
- Delete: `pkg/ios/device.go` (old)
- Delete: `pkg/ios/option.go` (old)
- Delete: `pkg/ios/device_test.go` (old)

**Step 1: Create option.go for iOS device options**

```go
// pkg/ios/option.go
package ios

type config struct {
    udid        string
    address     string  // For WDA connection
    skipSession bool
    bundleID    string
}

type Option func(*config)

// WithUDID sets the device UDID for go-ios connection
func WithUDID(udid string) Option {
    return func(c *config) {
        c.udid = udid
    }
}

// WithAddress sets the WDA server address (default: http://localhost:8100)
func WithAddress(addr string) Option {
    return func(c *config) {
        c.address = addr
    }
}

// SkipSession creates client without starting a WDA session
func SkipSession() Option {
    return func(c *config) {
        c.skipSession = true
    }
}
```

**Step 2: Create go-ios based screenshot.go**

```go
// pkg/ios/screenshot.go
package ios

import (
    "fmt"
    
    "github.com/danielpaulus/go-ios/ios"
    "github.com/danielpaulus/go-ios/ios/instruments"
)

func (d *Device) screenshotGoIOS() ([]byte, error) {
    // Find device by UDID
    device, err := d.findDevice()
    if err != nil {
        return nil, fmt.Errorf("find device: %w", err)
    }
    
    // Use go-ios screenshot service
    screenshotSvc, err := instruments.NewScreenshotService(device)
    if err != nil {
        return nil, fmt.Errorf("create screenshot service: %w", err)
    }
    defer screenshotSvc.Close()
    
    return screenshotSvc.TakeScreenshot()
}

func (d *Device) findDevice() (ios.DeviceEntry, error) {
    deviceList, err := ios.ListDevices()
    if err != nil {
        return ios.DeviceEntry{}, fmt.Errorf("list devices: %w", err)
    }
    
    if d.config.udid != "" {
        for _, dev := range deviceList.DeviceList {
            if dev.Properties.SerialNumber == d.config.udid {
                return dev, nil
            }
        }
        return ios.DeviceEntry{}, fmt.Errorf("device not found: %s", d.config.udid)
    }
    
    // Return first device if no UDID specified
    if len(deviceList.DeviceList) > 0 {
        return deviceList.DeviceList[0], nil
    }
    
    return ios.DeviceEntry{}, fmt.Errorf("no devices found")
}
```

**Step 3: Create app.go using go-ios**

```go
// pkg/ios/app.go
package ios

import (
    "fmt"
    
    "github.com/danielpaulus/go-ios/ios"
    "github.com/danielpaulus/go-ios/ios/appservice"
    "github.com/danielpaulus/go-ios/ios/instruments"
)

func (d *Device) launchAppGoIOS(bundleID string) error {
    device, err := d.findDevice()
    if err != nil {
        return fmt.Errorf("find device: %w", err)
    }
    
    // Try appservice first (iOS 17+)
    appsvc, err := appservice.New(device)
    if err == nil {
        defer appsvc.Close()
        _, err = appsvc.LaunchApp(bundleID, nil, nil, nil, true)
        if err == nil {
            return nil
        }
    }
    
    // Fall back to instruments (pre-iOS 17)
    pctrl, err := instruments.NewProcessControl(device)
    if err != nil {
        return fmt.Errorf("create process control: %w", err)
    }
    defer pctrl.Close()
    
    _, err = pctrl.LaunchApp(bundleID, nil)
    return err
}

func (d *Device) terminateAppGoIOS(bundleID string) error {
    device, err := d.findDevice()
    if err != nil {
        return fmt.Errorf("find device: %w", err)
    }
    
    // Try appservice first (iOS 17+)
    appsvc, err := appservice.New(device)
    if err == nil {
        defer appsvc.Close()
        processes, err := appsvc.ListProcesses()
        if err != nil {
            return err
        }
        for _, p := range processes {
            if p.BundleID == bundleID {
                return appsvc.KillProcess(p.Pid)
            }
        }
    }
    
    // Fall back to instruments
    pctrl, err := instruments.NewProcessControl(device)
    if err != nil {
        return fmt.Errorf("create process control: %w", err)
    }
    defer pctrl.Close()
    
    // Can't easily kill by bundle ID in instruments, so return not implemented
    return fmt.Errorf("terminate app not fully supported via go-ios")
}
```

**Step 4: Create wda_client.go for WDA delegation**

```go
// pkg/ios/wda_client.go
package ios

import (
    "github.com/liukunup/go-uop/pkg/wda"
)

// wdaClient wraps WDA client for iOS UI automation
// This handles Tap, SendKeys, GetSource, Alert operations
type wdaClient struct {
    client   *wda.Client
    bundleID string
}

func newWDAClient(addr, bundleID string) (*wdaClient, error) {
    client, err := wda.NewClient(addr)
    if err != nil {
        return nil, err
    }
    
    return &wdaClient{
        client:   client,
        bundleID: bundleID,
    }, nil
}

func (c *wdaClient) Tap(x, y int) error {
    // Try W3C tap first
    if err := c.client.Tap(x, y); err != nil {
        // Fall back to legacy tap
        return c.client.TapLegacy(x, y)
    }
    return nil
}

func (c *wdaClient) SendKeys(text string) error {
    // Use legacy keys endpoint
    return c.client.SendKeysLegacy(text)
}

func (c *wdaClient) GetSource() (string, error) {
    return c.client.GetSource()
}

func (c *wdaClient) GetAlertText() (string, error) {
    return c.client.GetAlertText()
}

func (c *wdaClient) AcceptAlert() error {
    return c.client.AcceptAlert()
}

func (c *wdaClient) DismissAlert() error {
    return c.client.DismissAlert()
}

func (c *wdaClient) Close() error {
    return c.client.Close()
}
```

**Step 5: Create new device.go combining go-ios + WDA**

```go
// pkg/ios/device.go
package ios

import (
    "fmt"
    "sync"
    
    "github.com/liukunup/go-uop/core"
    "github.com/liukunup/go-uop/pkg/wda"
)

type Device struct {
    mu         sync.Mutex
    config     *config
    wda        *wdaClient
    goIOSDev   interface{}  // ios.DeviceEntry for go-ios operations
}

func NewDevice(bundleID string, opts ...Option) (*Device, error) {
    cfg := &config{
        address: "http://localhost:8100",
    }
    for _, opt := range opts {
        opt(cfg)
    }
    cfg.bundleID = bundleID
    
    // Initialize WDA client for UI operations
    wdaClient, err := newWDAClient(cfg.address, bundleID)
    if err != nil {
        return nil, fmt.Errorf("create WDA client: %w", err)
    }
    
    if !cfg.skipSession {
        if err := wdaClient.client.StartSession(bundleID); err != nil {
            return nil, fmt.Errorf("start WDA session: %w", err)
        }
    }
    
    return &Device{
        config: cfg,
        wda:    wdaClient,
    }, nil
}

func (d *Device) Platform() core.Platform {
    return core.IOS
}

func (d *Device) Info() (map[string]interface{}, error) {
    return map[string]interface{}{
        "platform": "ios",
        "bundleId": d.config.bundleID,
        "wda":     d.config.address,
    }, nil
}

// Screenshot uses go-ios for high-quality screenshot via instruments
// Falls back to WDA screenshot if go-ios fails
func (d *Device) Screenshot() ([]byte, error) {
    // Try go-ios screenshot first (higher quality)
    if img, err := d.screenshotGoIOS(); err == nil {
        return img, nil
    }
    
    // Fall back to WDA screenshot
    return d.wda.client.Screenshot()
}

func (d *Device) Close() error {
    if d.wda != nil {
        return d.wda.Close()
    }
    return nil
}

// Tap uses WDA for UI tap
func (d *Device) Tap(x, y int) error {
    return d.wda.Tap(x, y)
}

// SendKeys uses WDA for text input
func (d *Device) SendKeys(text string) error {
    return d.wda.SendKeys(text)
}

// Launch uses go-ios for app launch
func (d *Device) Launch() error {
    return d.launchAppGoIOS(d.config.bundleID)
}

// GetSource uses WDA
func (d *Device) GetSource() (string, error) {
    return d.wda.GetSource()
}

// Alert operations use WDA
func (d *Device) GetAlertText() (string, error) {
    return d.wda.GetAlertText()
}

func (d *Device) AcceptAlert() error {
    return d.wda.AcceptAlert()
}

func (d *Device) DismissAlert() error {
    return d.wda.DismissAlert()
}

var _ core.Device = (*Device)(nil)
```

**Step 6: Delete old ios files**

```bash
rm pkg/ios/device.go pkg/ios/option.go pkg/ios/device_test.go
```

**Step 7: Commit**

```bash
git add -A
git commit -m "refactor(ios): use go-ios for device mgmt, WDA for UI automation

- Add go-ios based screenshot via instruments
- Add go-ios based app launch (appservice/instruments)
- WDA client moved to pkg/wda with W3C compliance
- Device orchestrates go-ios (screenshot, launch) + WDA (tap, keys, alerts)
- Remove old WDA wrapper from pkg/ios/"
```

---

## Task 3: Update go.mod and verify builds

**Step 1: Update go.mod to promote go-ios to direct dependency**

```bash
go get github.com/danielpaulus/go-ios@v1.0.206
go mod tidy
```

**Step 2: Verify build**

```bash
go build ./...
```

**Step 3: Run tests**

```bash
go test ./pkg/... -v -short
```

**Step 4: Commit**

```bash
git add go.mod go.sum
git commit -m "deps: promote go-ios to direct dependency"
```

---

## Final Structure

```
pkg/
├── adb/                    # ✅ adb CLI wrapper (unchanged)
│   ├── client.go
│   ├── input.go
│   ├── screenshot.go
│   ├── shell.go
│   ├── package.go
│   └── activity.go
│
├── ios/                   # ⚡ NEW: go-ios + WDA orchestration
│   ├── device.go          # Core device implementation
│   ├── option.go          # Configuration options
│   ├── screenshot.go      # go-ios screenshot
│   ├── app.go             # go-ios app launch
│   └── wda_client.go      # WDA UI operations delegator
│
├── wda/                   # ⚡ MOVED & REFACTORED: W3C WebDriver
│   ├── client.go          # HTTP client + W3C session
│   ├── protocol.go        # W3C endpoints
│   ├── session.go         # Session lifecycle
│   ├── element.go         # FindElement, Click, SendKeys
│   ├── app.go             # App launch/terminate
│   ├── alert.go           # Alert handling
│   ├── screenshot.go      # Screenshot
│   └── w3c.go             # W3C types & error codes
│
└── uiautomator2/          # ✅ unchanged
    └── ...
```

---

## Dependencies Flow

```
core.Device (interface)
    ↑
    │
    ├─ android.Device (pkg/adb) ──────→ adb CLI
    │
    ├─ ios.Device (pkg/ios) ─────────┐
    │                                │
    │         ┌─────────────────────┴──────────┐
    │         ↓                           ↓
    │   go-ios DeviceEntry            wda.Client (pkg/wda)
    │   - Screenshot (instruments)    - Tap, SendKeys
    │   - App Launch                  - FindElement
    │   - Pair                        - Alerts
    │   - ListDevices                 - Screenshot
    │
    └─ uiautomator2.Device (pkg/uiautomator2) ──→ JSON-RPC → device server
```

---

## Notes

1. **go-ios limitations**: Does NOT provide tap/sendText - UI automation still requires WDA
2. **iOS 17+**: Requires `sudo ios tunnel start` for full go-ios functionality
3. **W3C compliance**: Current implementation is a best-effort W3C-style wrapper around WDA's actual REST API
4. **Backward compatibility**: Legacy tap/keys endpoints maintained for compatibility

---

## Execution Options

**Plan complete and saved to `docs/plans/2026-04-11-pkg-restyle-plan.md`. Two execution options:**

**1. Subagent-Driven (this session)** - I dispatch fresh subagent per task, review between tasks, fast iteration

**2. Parallel Session (separate)** - Open new session with executing-plans, batch execution with checkpoints

**Which approach?**
