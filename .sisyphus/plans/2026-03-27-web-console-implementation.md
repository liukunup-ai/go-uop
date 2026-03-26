# Web Console Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 构建一个 Web 控制台应用，支持 iOS/Android 设备选择、屏幕预览、命令调试、命令历史和 YAML 导出。

**Architecture:** Go HTTP Server + React SPA，前端使用 embed 包嵌入 Go 二进制。REST API 进行前后端通信。

**Tech Stack:** Go 1.22+, React 18, TypeScript, Vite, Tailwind CSS, Zustand

---

## Phase 1: 项目初始化

### Task 1: 创建 Go 后端入口

**Files:**
- Create: `cmd/console/main.go`
- Modify: `go.mod` (添加依赖)

**Step 1: 创建入口文件**

```go
package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    
    "github.com/liukunup/go-uop/internal/console"
)

var (
    addr string
    openBrowser bool
    devMode bool
)

func main() {
    flag.StringVar(&addr, "addr", ":8080", "HTTP server address")
    flag.BoolVar(&openBrowser, "open", true, "Open browser on start")
    flag.BoolVar(&devMode, "dev", false, "Development mode")
    flag.Parse()
    
    server, err := console.NewServer(addr, devMode)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("🚀 go-uop Console starting on %s\n", addr)
    if openBrowser {
        go func() {
            // 延迟打开浏览器
        }()
    }
    
    log.Fatal(server.Start())
}
```

**Step 2: 运行测试**

Run: `go build -o /dev/null ./cmd/console/`
Expected: 编译错误 (console 包不存在)

**Step 3: 创建 console 包结构**

Create: `internal/console/server.go`

```go
package console

import (
    "net/http"
    "time"
)

type Server struct {
    mux         *http.ServeMux
    addr        string
    devMode     bool
}

func NewServer(addr string, devMode bool) (*Server, error) {
    s := &Server{
        mux:     http.NewServeMux(),
        addr:    addr,
        devMode: devMode,
    }
    s.setupRoutes()
    return s, nil
}

func (s *Server) Start() error {
    srv := &http.Server{
        Addr:         s.addr,
        Handler:      s.mux,
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
    }
    return srv.ListenAndServe()
}
```

**Step 4: 添加路由**

Modify: `internal/console/server.go`

在 `setupRoutes` 中添加:
```go
func (s *Server) setupRoutes() {
    s.mux.HandleFunc("/api/devices", s.handleDevices)
    s.mux.HandleFunc("/api/devices/connect", s.handleConnect)
    s.mux.HandleFunc("/api/devices/", s.handleDeviceOps)
    s.mux.HandleFunc("/api/commands/history", s.handleHistory)
    s.mux.HandleFunc("/api/export/yaml", s.handleYamlExport)
    
    // 前端静态文件
    if !s.devMode {
        s.mux.HandleFunc("/", s.handleFrontend)
    }
}
```

**Step 5: 提交**

```bash
git add cmd/console/main.go internal/console/server.go
git commit -m "feat(console): add server entry point and route structure"
```

---

### Task 2: 初始化 React 前端项目

**Files:**
- Create: `console/package.json`
- Create: `console/vite.config.ts`
- Create: `console/tsconfig.json`
- Create: `console/tailwind.config.js`
- Create: `console/index.html`
- Create: `console/src/main.tsx`
- Create: `console/src/App.tsx`

**Step 1: 创建 package.json**

```json
{
  "name": "go-uop-console",
  "private": true,
  "version": "0.1.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "zustand": "^4.5.0",
    "axios": "^1.6.0"
  },
  "devDependencies": {
    "@types/react": "^18.2.0",
    "@types/react-dom": "^18.2.0",
    "@vitejs/plugin-react": "^4.2.0",
    "autoprefixer": "^10.4.16",
    "postcss": "^8.4.32",
    "tailwindcss": "^3.4.0",
    "typescript": "^5.3.0",
    "vite": "^5.0.0"
  }
}
```

**Step 2: 创建 Vite 配置**

```typescript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  base: './',
  build: {
    outDir: 'dist',
    emptyOutDir: true,
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
```

**Step 3: 创建 Tailwind 配置**

```javascript
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
```

**Step 4: 创建入口文件**

`console/index.html`:
```html
<!DOCTYPE html>
<html lang="zh-CN">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>go-uop Console</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
```

`console/src/main.tsx`:
```tsx
import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import './index.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
```

`console/src/index.css`:
```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

**Step 5: 创建基础 App 组件**

```tsx
function App() {
  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 py-4">
          <h1 className="text-2xl font-bold">go-uop Console</h1>
        </div>
      </header>
      <main className="max-w-7xl mx-auto px-4 py-6">
        <p>Web Console 即将上线...</p>
      </main>
    </div>
  )
}

export default App
```

**Step 6: 运行前端开发服务器验证**

Run: `cd console && npm install && npm run dev`
Expected: 前端开发服务器启动成功

**Step 7: 提交**

```bash
git add console/
git commit -m "feat(console): add React frontend scaffolding"
```

---

## Phase 2: 后端核心实现

### Task 3: 实现 DeviceManager

**Files:**
- Create: `internal/console/device.go`
- Create: `internal/console/types.go`

**Step 1: 定义类型**

```go
// internal/console/types.go
package console

// Device 表示一个移动设备
type Device struct {
    ID      string `json:"id"`
    Platform string `json:"platform"` // "ios" or "android"
    Name    string `json:"name"`
    Serial  string `json:"serial"`
    Status  string `json:"status"` // "available", "connected", "error"
    Model   string `json:"model,omitempty"`
    Address string `json:"address,omitempty"`   // iOS WDA 地址
    PkgName string `json:"packageName,omitempty"` // Android 包名
}

// CommandRecord 表示一条命令记录
type CommandRecord struct {
    ID        string                 `json:"id"`
    Timestamp string                 `json:"timestamp"`
    Type      string                 `json:"command"`
    Params    map[string]interface{} `json:"params"`
    Success   bool                   `json:"success"`
    Output    string                 `json:"output,omitempty"`
    Duration  string                 `json:"duration"`
}

// CommandRequest 表示命令请求
type CommandRequest struct {
    Command string                 `json:"command"`
    Params  map[string]interface{} `json:"params"`
}
```

**Step 2: 实现 DeviceManager**

```go
// internal/console/device.go
package console

import (
    "sync"
    "time"

    "github.com/liukunup/go-uop/android"
    "github.com/liukunup/go-uop/android/adb"
    "github.com/liukunup/go-uop/core"
    "github.com/liukunup/go-uop/ios"
)

type DeviceManager struct {
    mu        sync.RWMutex
    devices   map[string]core.Device
    info      map[string]*Device
}

func NewDeviceManager() *DeviceManager {
    return &DeviceManager{
        devices: make(map[string]core.Device),
        info:    make(map[string]*Device),
    }
}

// ListDevices 列出所有可用设备
func (m *DeviceManager) ListDevices() ([]*Device, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    var result []*Device

    // 列出 Android 设备
    androidDevices, err := adb.Devices()
    if err == nil {
        for _, d := range androidDevices {
            status := "available"
            if _, connected := m.devices["android-"+d.Serial]; connected {
                status = "connected"
            }
            result = append(result, &Device{
                ID:       "android-" + d.Serial,
                Platform: "android",
                Name:     d.Model,
                Serial:   d.Serial,
                Status:   status,
                Model:    d.Model,
            })
        }
    }

    // iOS 设备需要用户手动提供地址
    // 目前只返回已连接的 iOS 设备
    for id, info := range m.info {
        if info.Platform == "ios" {
            status := "available"
            if _, connected := m.devices[id]; connected {
                status = "connected"
            }
            result = append(result, &Device{
                ID:      id,
                Platform: "ios",
                Name:    info.Name,
                Serial:  info.Serial,
                Status:  status,
                Address: info.Address,
            })
        }
    }

    return result, nil
}

// ConnectDevice 连接设备
func (m *DeviceManager) ConnectDevice(d *Device) (core.Device, error) {
    m.mu.Lock()
    defer m.mu.Unlock()

    var device core.Device
    var err error

    switch d.Platform {
    case "android":
        device, err = android.NewDevice(
            android.WithSerial(d.Serial),
            android.WithPackage(d.PkgName),
        )
    case "ios":
        device, err = ios.NewDevice(
            d.Serial, // bundleID for iOS
            ios.WithAddress(d.Address),
        )
    default:
        err = ErrUnsupportedPlatform
    }

    if err != nil {
        return nil, err
    }

    m.devices[d.ID] = device
    m.info[d.ID] = d
    d.Status = "connected"

    return device, nil
}

// DisconnectDevice 断开设备
func (m *DeviceManager) DisconnectDevice(id string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    device, ok := m.devices[id]
    if !ok {
        return ErrDeviceNotFound
    }

    if err := device.Close(); err != nil {
        return err
    }

    delete(m.devices, id)
    if info, ok := m.info[id]; ok {
        info.Status = "available"
    }

    return nil
}

// GetConnected 获取已连接的设备
func (m *DeviceManager) GetConnected(id string) (core.Device, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    device, ok := m.devices[id]
    if !ok {
        return nil, ErrDeviceNotConnected
    }

    return device, nil
}

// ExecuteCommand 执行命令
func (m *DeviceManager) ExecuteCommand(id string, cmd string, params map[string]interface{}) (*CommandRecord, error) {
    device, err := m.GetConnected(id)
    if err != nil {
        return nil, err
    }

    record := &CommandRecord{
        ID:        generateID(),
        Timestamp: time.Now().Format(time.RFC3339),
        Type:      cmd,
        Params:    params,
    }

    start := time.Now()
    defer func() {
        record.Duration = time.Since(start).String()
    }()

    switch cmd {
    case "tap":
        x, _ := toInt(params["x"])
        y, _ := toInt(params["y"])
        err = device.Tap(x, y)
    case "input":
        text, _ := toString(params["text"])
        err = device.SendKeys(text)
    case "launch":
        err = device.Launch()
    case "terminate":
        err = device.Terminate()
    case "screenshot":
        // 截图单独处理
    default:
        err = ErrUnknownCommand
    }

    record.Success = err == nil
    if err != nil {
        record.Output = err.Error()
    }

    return record, err
}
```

**Step 3: 添加错误定义**

Create: `internal/console/errors.go`

```go
package console

import "errors"

var (
    ErrDeviceNotFound      = errors.New("device not found")
    ErrDeviceNotConnected  = errors.New("device not connected")
    ErrUnsupportedPlatform = errors.New("unsupported platform")
    ErrUnknownCommand      = errors.New("unknown command")
)
```

**Step 4: 添加辅助函数**

Add to: `internal/console/device.go`

```go
func generateID() string {
    return fmt.Sprintf("cmd-%d", time.Now().UnixNano())
}

func toInt(v interface{}) (int, bool) {
    switch val := v.(type) {
    case int:
        return val, true
    case float64:
        return int(val), true
    case string:
        i, err := strconv.Atoi(val)
        return i, err == nil
    }
    return 0, false
}

func toString(v interface{}) (string, bool) {
    if s, ok := v.(string); ok {
        return s, true
    }
    return "", false
}
```

**Step 5: 提交**

```bash
git add internal/console/
git commit -m "feat(console): add DeviceManager and types"
```

---

### Task 4: 实现 HTTP Handlers

**Files:**
- Create: `internal/console/handler.go`

**Step 1: 实现设备列表处理**

```go
// handleDevices 返回所有设备列表
func (s *Server) handleDevices(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
        return
    }

    devices, err := s.deviceMgr.ListDevices()
    if err != nil {
        writeError(w, http.StatusInternalServerError, "LIST_DEVICES_FAILED")
        return
    }

    writeJSON(w, map[string]interface{}{
        "devices": devices,
    })
}
```

**Step 2: 实现设备连接处理**

```go
// handleConnect 连接设备
func (s *Server) handleConnect(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
        return
    }

    var req Device
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "INVALID_REQUEST")
        return
    }

    device, err := s.deviceMgr.ConnectDevice(&req)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "CONNECT_FAILED")
        return
    }

    info, _ := device.Info()
    writeJSON(w, map[string]interface{}{
        "success": true,
        "device":  req,
        "info":    info,
    })
}
```

**Step 3: 实现设备操作处理**

```go
// handleDeviceOps 处理设备相关操作
func (s *Server) handleDeviceOps(w http.ResponseWriter, r *http.Request) {
    path := strings.TrimPrefix(r.URL.Path, "/api/devices/")
    
    // 解析路径: :id/screenshot, :id/info, :id/commands
    parts := strings.Split(path, "/")
    if len(parts) < 2 {
        writeError(w, http.StatusBadRequest, "INVALID_PATH")
        return
    }

    deviceID := parts[0]
    operation := parts[1]

    switch operation {
    case "screenshot":
        s.handleScreenshot(w, r, deviceID)
    case "info":
        s.handleDeviceInfo(w, r, deviceID)
    case "commands":
        s.handleCommands(w, r, deviceID)
    default:
        writeError(w, http.StatusNotFound, "OPERATION_NOT_FOUND")
    }
}

func (s *Server) handleScreenshot(w http.ResponseWriter, r *http.Request, deviceID string) {
    device, err := s.deviceMgr.GetConnected(deviceID)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "DEVICE_NOT_CONNECTED")
        return
    }

    img, err := device.Screenshot()
    if err != nil {
        writeError(w, http.StatusInternalServerError, "SCREENSHOT_FAILED")
        return
    }

    w.Header().Set("Content-Type", "image/png")
    w.Write(img)
}
```

**Step 4: 实现命令执行处理**

```go
// handleCommands 执行命令
func (s *Server) handleCommands(w http.ResponseWriter, r *http.Request, deviceID string) {
    if r.Method != http.MethodPost {
        writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
        return
    }

    var req CommandRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "INVALID_REQUEST")
        return
    }

    record, err := s.deviceMgr.ExecuteCommand(deviceID, req.Command, req.Params)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "COMMAND_FAILED")
        return
    }

    // 添加到历史
    s.historyMgr.Add(record)

    writeJSON(w, record)
}
```

**Step 5: 添加 JSON 响应辅助函数**

Add to: `internal/console/handler.go`

```go
func writeJSON(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, code string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "error": map[string]interface{}{
            "code": code,
        },
    })
}
```

**Step 6: 提交**

```bash
git add internal/console/handler.go
git commit -m "feat(console): add HTTP handlers"
```

---

### Task 5: 实现 HistoryManager 和 YAML 导出

**Files:**
- Create: `internal/console/history.go`
- Create: `internal/console/yaml.go`

**Step 1: 实现 HistoryManager**

```go
// internal/console/history.go
package console

type HistoryManager struct {
    mu      sync.RWMutex
    history []CommandRecord
    maxSize int
}

func NewHistoryManager(maxSize int) *HistoryManager {
    if maxSize <= 0 {
        maxSize = 100
    }
    return &HistoryManager{
        history: make([]CommandRecord, 0, maxSize),
        maxSize: maxSize,
    }
}

func (m *HistoryManager) Add(record *CommandRecord) {
    m.mu.Lock()
    defer m.mu.Unlock()

    m.history = append(m.history, *record)
    
    // 限制历史大小
    if len(m.history) > m.maxSize {
        m.history = m.history[len(m.history)-m.maxSize:]
    }
}

func (m *HistoryManager) GetAll() []CommandRecord {
    m.mu.RLock()
    defer m.mu.RUnlock()

    result := make([]CommandRecord, len(m.history))
    copy(result, m.history)
    return result
}

func (m *HistoryManager) Clear() {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.history = m.history[:0]
}

func (m *HistoryManager) GetSelected(ids []string) []CommandRecord {
    m.mu.RLock()
    defer m.mu.RUnlock()

    idSet := make(map[string]bool)
    for _, id := range ids {
        idSet[id] = true
    }

    var result []CommandRecord
    for _, record := range m.history {
        if idSet[record.ID] {
            result = append(result, record)
        }
    }
    return result
}
```

**Step 2: 实现 YAML 导出**

```go
// internal/console/yaml.go
package console

import (
    "fmt"
    "strings"

    "gopkg.in/yaml.v3"
)

type YamlFlow struct {
    Name  string        `yaml:"name"`
    AppID string        `yaml:"appId,omitempty"`
    Steps []YamlCommand `yaml:"steps"`
}

type YamlCommand struct {
    TapOn     *TapCommand     `yaml:"tapOn,omitempty"`
    InputText *InputCommand   `yaml:"inputText,omitempty"`
    Swipe     *SwipeCommand   `yaml:"swipe,omitempty"`
    Launch    *struct{}       `yaml:"launch,omitempty"`
    Terminate *struct{}       `yaml:"terminate,omitempty"`
}

type TapCommand struct {
    X    int    `yaml:"x,omitempty"`
    Y    int    `yaml:"y,omitempty"`
    Text string `yaml:"text,omitempty"`
}

type InputCommand struct {
    Text string `yaml:"text"`
}

type SwipeCommand struct {
    X1       int `yaml:"x1"`
    Y1       int `yaml:"y1"`
    X2       int `yaml:"x2"`
    Y2       int `yaml:"y2"`
    Duration int `yaml:"duration,omitempty"`
}

func ExportToYaml(records []CommandRecord, name string) ([]byte, error) {
    flow := YamlFlow{
        Name:  name,
        Steps: make([]YamlCommand, 0, len(records)),
    }

    for _, record := range records {
        if !record.Success {
            continue
        }

        cmd := YamlCommand{}
        switch record.Type {
        case "tap":
            x, _ := toInt(record.Params["x"])
            y, _ := toInt(record.Params["y"])
            cmd.TapOn = &TapCommand{X: x, Y: y}
        case "input":
            text, _ := toString(record.Params["text"])
            cmd.InputText = &InputCommand{Text: text}
        case "swipe":
            x1, _ := toInt(record.Params["x1"])
            y1, _ := toInt(record.Params["y1"])
            x2, _ := toInt(record.Params["x2"])
            y2, _ := toInt(record.Params["y2"])
            dur, _ := toInt(record.Params["duration"])
            cmd.Swipe = &SwipeCommand{X1: x1, Y1: y1, X2: x2, Y2: y2, Duration: dur}
        case "launch":
            cmd.Launch = &struct{}{}
        case "terminate":
            cmd.Terminate = &struct{}{}
        }

        flow.Steps = append(flow.Steps, cmd)
    }

    return yaml.Marshal(&flow)
}
```

**Step 3: 添加历史和导出 handler**

Modify: `internal/console/handler.go`

```go
// handleHistory 获取命令历史
func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
    history := s.historyMgr.GetAll()
    writeJSON(w, map[string]interface{}{
        "history": history,
    })
}

// handleYamlExport 导出为 YAML
func (s *Server) handleYamlExport(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()
    ids := query["ids"]
    
    var records []CommandRecord
    if len(ids) > 0 {
        records = s.historyMgr.GetSelected(ids)
    } else {
        records = s.historyMgr.GetAll()
    }

    name := query.Get("name")
    if name == "" {
        name = "debug-session"
    }

    yamlData, err := ExportToYaml(records, name)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "EXPORT_FAILED")
        return
    }

    w.Header().Set("Content-Type", "text/yaml")
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.yaml", name))
    w.Write(yamlData)
}
```

**Step 4: 提交**

```bash
git add internal/console/history.go internal/console/yaml.go
git commit -m "feat(console): add HistoryManager and YAML export"
```

---

## Phase 3: 前端组件实现

### Task 6: 实现 API 客户端和状态管理

**Files:**
- Create: `console/src/api/client.ts`
- Create: `console/src/types/index.ts`
- Create: `console/src/stores/deviceStore.ts`

**Step 1: 定义类型**

```typescript
// console/src/types/index.ts
export interface Device {
  id: string;
  platform: 'ios' | 'android';
  name: string;
  serial: string;
  status: 'available' | 'connected' | 'error';
  model?: string;
  address?: string;
  packageName?: string;
}

export interface CommandRecord {
  id: string;
  timestamp: string;
  command: string;
  params: Record<string, any>;
  success: boolean;
  output?: string;
  duration: string;
}

export interface CommandRequest {
  command: string;
  params: Record<string, any>;
}
```

**Step 2: 实现 API 客户端**

```typescript
// console/src/api/client.ts
import axios from 'axios';
import type { Device, CommandRecord, CommandRequest } from '../types';

const client = axios.create({
  baseURL: '/api',
  timeout: 30000,
});

export const api = {
  // 获取设备列表
  listDevices: async (): Promise<Device[]> => {
    const res = await client.get<{ devices: Device[] }>('/devices');
    return res.data.devices;
  },

  // 连接设备
  connectDevice: async (device: Partial<Device>): Promise<{ device: Device }> => {
    const res = await client.post<{ device: Device }>('/devices/connect', device);
    return res.data;
  },

  // 获取截图
  getScreenshot: async (deviceId: string): Promise<Blob> => {
    const res = await client.get(`/devices/${deviceId}/screenshot`, {
      responseType: 'blob',
    });
    return res.data;
  },

  // 获取设备信息
  getDeviceInfo: async (deviceId: string): Promise<any> => {
    const res = await client.get(`/devices/${deviceId}/info`);
    return res.data;
  },

  // 执行命令
  executeCommand: async (
    deviceId: string,
    cmd: CommandRequest
  ): Promise<CommandRecord> => {
    const res = await client.post<CommandRecord>(
      `/devices/${deviceId}/commands`,
      cmd
    );
    return res.data;
  },

  // 获取命令历史
  getHistory: async (): Promise<CommandRecord[]> => {
    const res = await client.get<{ history: CommandRecord[] }>('/commands/history');
    return res.data.history;
  },

  // 导出 YAML
  exportYaml: async (ids?: string[], name?: string): Promise<string> => {
    const params = new URLSearchParams();
    if (ids && ids.length > 0) {
      ids.forEach((id) => params.append('ids', id));
    }
    if (name) {
      params.append('name', name);
    }
    const res = await client.get(`/export/yaml?${params.toString()}`, {
      responseType: 'text',
    });
    return res.data;
  },
};
```

**Step 3: 实现 Zustand Store**

```typescript
// console/src/stores/deviceStore.ts
import { create } from 'zustand';
import type { Device, CommandRecord } from '../types';
import { api } from '../api/client';

interface DeviceStore {
  // 状态
  devices: Device[];
  selectedDevice: Device | null;
  connectedDevice: Device | null;
  history: CommandRecord[];
  screenshot: string | null;
  isLoading: boolean;
  error: string | null;

  // 操作
  fetchDevices: () => Promise<void>;
  selectDevice: (device: Device | null) => void;
  connectDevice: (device: Device) => Promise<void>;
  disconnectDevice: () => Promise<void>;
  fetchScreenshot: () => Promise<void>;
  executeCommand: (command: string, params: Record<string, any>) => Promise<void>;
  fetchHistory: () => Promise<void>;
  clearHistory: () => void;
  exportYaml: (ids?: string[]) => Promise<void>;
}

export const useDeviceStore = create<DeviceStore>((set, get) => ({
  // 初始状态
  devices: [],
  selectedDevice: null,
  connectedDevice: null,
  history: [],
  screenshot: null,
  isLoading: false,
  error: null,

  // 获取设备列表
  fetchDevices: async () => {
    try {
      set({ isLoading: true, error: null });
      const devices = await api.listDevices();
      set({ devices, isLoading: false });
    } catch (err: any) {
      set({ error: err.message, isLoading: false });
    }
  },

  // 选择设备
  selectDevice: (device) => {
    set({ selectedDevice: device });
  },

  // 连接设备
  connectDevice: async (device) => {
    try {
      set({ isLoading: true, error: null });
      await api.connectDevice(device);
      set({ connectedDevice: device, isLoading: false });
      // 连接成功后获取截图
      get().fetchScreenshot();
    } catch (err: any) {
      set({ error: err.message, isLoading: false });
    }
  },

  // 断开设备
  disconnectDevice: async () => {
    set({ connectedDevice: null, screenshot: null });
  },

  // 获取截图
  fetchScreenshot: async () => {
    const { connectedDevice } = get();
    if (!connectedDevice) return;

    try {
      const blob = await api.getScreenshot(connectedDevice.id);
      const url = URL.createObjectURL(blob);
      set({ screenshot: url });
    } catch (err) {
      console.error('Screenshot failed:', err);
    }
  },

  // 执行命令
  executeCommand: async (command, params) => {
    const { connectedDevice } = get();
    if (!connectedDevice) {
      set({ error: 'No device connected' });
      return;
    }

    try {
      set({ isLoading: true, error: null });
      const record = await api.executeCommand(connectedDevice.id, { command, params });
      set((state) => ({
        history: [...state.history, record],
        isLoading: false,
      }));
      // 如果是截图命令，刷新截图
      if (command === 'screenshot' || command === 'tap' || command === 'swipe') {
        get().fetchScreenshot();
      }
    } catch (err: any) {
      set({ error: err.message, isLoading: false });
    }
  },

  // 获取历史
  fetchHistory: async () => {
    try {
      const history = await api.getHistory();
      set({ history });
    } catch (err) {
      console.error('Fetch history failed:', err);
    }
  },

  // 清空历史
  clearHistory: () => {
    set({ history: [] });
  },

  // 导出 YAML
  exportYaml: async (ids) => {
    try {
      const yaml = await api.exportYaml(ids);
      // 创建下载
      const blob = new Blob([yaml], { type: 'text/yaml' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'debug-session.yaml';
      a.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      console.error('Export YAML failed:', err);
    }
  },
}));
```

**Step 4: 提交**

```bash
git add console/src/api/ console/src/types/ console/src/stores/
git commit -m "feat(console): add API client and Zustand store"
```

---

### Task 7: 实现前端组件

**Files:**
- Create: `console/src/components/DeviceList.tsx`
- Create: `console/src/components/ScreenPreview.tsx`
- Create: `console/src/components/CommandPanel.tsx`
- Create: `console/src/components/CommandHistory.tsx`

**Step 1: 实现 DeviceList 组件**

```tsx
// console/src/components/DeviceList.tsx
import React from 'react';
import { useDeviceStore } from '../stores/deviceStore';

export const DeviceList: React.FC = () => {
  const {
    devices,
    selectedDevice,
    connectedDevice,
    selectDevice,
    connectDevice,
    disconnectDevice,
    fetchDevices,
    isLoading,
  } = useDeviceStore();

  React.useEffect(() => {
    fetchDevices();
    const interval = setInterval(fetchDevices, 5000);
    return () => clearInterval(interval);
  }, []);

  const iosDevices = devices.filter((d) => d.platform === 'ios');
  const androidDevices = devices.filter((d) => d.platform === 'android');

  const handleConnect = async () => {
    if (!selectedDevice) return;
    await connectDevice(selectedDevice);
  };

  return (
    <div className="bg-white rounded-lg shadow p-4">
      <h2 className="text-lg font-semibold mb-4">设备列表</h2>

      {/* iOS 设备 */}
      <div className="mb-4">
        <h3 className="text-sm font-medium text-gray-500 mb-2">iOS</h3>
        <div className="space-y-2">
          {iosDevices.map((device) => (
            <div
              key={device.id}
              onClick={() => selectDevice(device)}
              className={`p-3 rounded border cursor-pointer ${
                selectedDevice?.id === device.id
                  ? 'border-blue-500 bg-blue-50'
                  : 'border-gray-200 hover:border-gray-300'
              }`}
            >
              <div className="flex items-center">
                <span className="text-lg mr-2">📱</span>
                <div>
                  <div className="font-medium">{device.name}</div>
                  <div className="text-xs text-gray-500">{device.serial}</div>
                </div>
                {connectedDevice?.id === device.id && (
                  <span className="ml-auto text-green-500 text-sm">已连接</span>
                )}
              </div>
            </div>
          ))}
          {iosDevices.length === 0 && (
            <p className="text-sm text-gray-400">未发现 iOS 设备</p>
          )}
        </div>
      </div>

      {/* Android 设备 */}
      <div className="mb-4">
        <h3 className="text-sm font-medium text-gray-500 mb-2">Android</h3>
        <div className="space-y-2">
          {androidDevices.map((device) => (
            <div
              key={device.id}
              onClick={() => selectDevice(device)}
              className={`p-3 rounded border cursor-pointer ${
                selectedDevice?.id === device.id
                  ? 'border-blue-500 bg-blue-50'
                  : 'border-gray-200 hover:border-gray-300'
              }`}
            >
              <div className="flex items-center">
                <span className="text-lg mr-2">📱</span>
                <div>
                  <div className="font-medium">{device.name}</div>
                  <div className="text-xs text-gray-500">{device.serial}</div>
                </div>
                {connectedDevice?.id === device.id && (
                  <span className="ml-auto text-green-500 text-sm">已连接</span>
                )}
              </div>
            </div>
          ))}
          {androidDevices.length === 0 && (
            <p className="text-sm text-gray-400">未发现 Android 设备</p>
          )}
        </div>
      </div>

      {/* 操作按钮 */}
      <div className="flex gap-2">
        <button
          onClick={handleConnect}
          disabled={!selectedDevice || isLoading}
          className="flex-1 px-4 py-2 bg-blue-500 text-white rounded disabled:bg-gray-300"
        >
          连接设备
        </button>
        <button
          onClick={disconnectDevice}
          disabled={!connectedDevice || isLoading}
          className="flex-1 px-4 py-2 bg-red-500 text-white rounded disabled:bg-gray-300"
        >
          断开
        </button>
      </div>
    </div>
  );
};
```

**Step 2: 实现 ScreenPreview 组件**

```tsx
// console/src/components/ScreenPreview.tsx
import React, { useRef, useState, useEffect } from 'react';
import { useDeviceStore } from '../stores/deviceStore';

export const ScreenPreview: React.FC = () => {
  const { connectedDevice, screenshot, fetchScreenshot } = useDeviceStore();
  const [mousePos, setMousePos] = useState<{ x: number; y: number } | null>(null);
  const imgRef = useRef<HTMLImageElement>(null);
  const [imgSize, setImgSize] = useState({ width: 0, height: 0 });

  useEffect(() => {
    const interval = setInterval(() => {
      if (connectedDevice) {
        fetchScreenshot();
      }
    }, 3000);
    return () => clearInterval(interval);
  }, [connectedDevice]);

  const handleMouseMove = (e: React.MouseEvent) => {
    if (!imgRef.current) return;
    const rect = imgRef.current.getBoundingClientRect();
    const scaleX = imgRef.current.naturalWidth / rect.width;
    const scaleY = imgRef.current.naturalHeight / rect.height;
    setMousePos({
      x: Math.round((e.clientX - rect.left) * scaleX),
      y: Math.round((e.clientY - rect.top) * scaleY),
    });
  };

  return (
    <div className="bg-white rounded-lg shadow p-4">
      <h2 className="text-lg font-semibold mb-4">屏幕预览</h2>

      {!connectedDevice ? (
        <div className="aspect-[9/16] bg-gray-100 rounded flex items-center justify-center">
          <p className="text-gray-400">请先连接设备</p>
        </div>
      ) : screenshot ? (
        <div className="relative">
          <img
            ref={imgRef}
            src={screenshot}
            alt="Device Screen"
            className="max-w-full mx-auto rounded border"
            onMouseMove={handleMouseMove}
            onMouseLeave={() => setMousePos(null)}
          />
          {mousePos && (
            <div className="absolute top-2 right-2 bg-black/70 text-white px-2 py-1 rounded text-sm">
              ({mousePos.x}, {mousePos.y})
            </div>
          )}
        </div>
      ) : (
        <div className="aspect-[9/16] bg-gray-100 rounded flex items-center justify-center">
          <p className="text-gray-400">加载中...</p>
        </div>
      )}

      {/* 操作按钮 */}
      <div className="flex gap-2 mt-4">
        <button
          onClick={() => fetchScreenshot()}
          disabled={!connectedDevice}
          className="flex-1 px-4 py-2 bg-gray-100 border rounded disabled:bg-gray-50"
        >
          📸 刷新截图
        </button>
      </div>

      {/* 坐标信息 */}
      {mousePos && (
        <div className="mt-2 text-sm text-gray-600">
          点击坐标: X={mousePos.x}, Y={mousePos.y}
        </div>
      )}
    </div>
  );
};
```

**Step 3: 实现 CommandPanel 组件**

```tsx
// console/src/components/CommandPanel.tsx
import React, { useState } from 'react';
import { useDeviceStore } from '../stores/deviceStore';

interface CommandButton {
  label: string;
  command: string;
  icon: string;
  fields?: { name: string; label: string; type: string; default?: any }[];
}

const commands: CommandButton[] = [
  { label: 'Tap', command: 'tap', icon: '👆', fields: [
    { name: 'x', label: 'X', type: 'number' },
    { name: 'y', label: 'Y', type: 'number' },
  ]},
  { label: 'Input', command: 'input', icon: '⌨️', fields: [
    { name: 'text', label: '文本', type: 'text' },
  ]},
  { label: 'Swipe', command: 'swipe', icon: '👆', fields: [
    { name: 'x1', label: 'X1', type: 'number' },
    { name: 'y1', label: 'Y1', type: 'number' },
    { name: 'x2', label: 'X2', type: 'number' },
    { name: 'y2', label: 'Y2', type: 'number' },
  ]},
  { label: 'Launch', command: 'launch', icon: '🚀', fields: [] },
  { label: 'Terminate', command: 'terminate', icon: '❌', fields: [] },
];

export const CommandPanel: React.FC = () => {
  const { connectedDevice, executeCommand, error } = useDeviceStore();
  const [selectedCommand, setSelectedCommand] = useState<string | null>(null);
  const [params, setParams] = useState<Record<string, any>>({});

  const handleCommandSelect = (cmd: CommandButton) => {
    setSelectedCommand(cmd.command);
    const defaults: Record<string, any> = {};
    cmd.fields?.forEach((f) => {
      defaults[f.name] = f.default || '';
    });
    setParams(defaults);
  };

  const handleExecute = async () => {
    if (!selectedCommand) return;
    await executeCommand(selectedCommand, params);
    setSelectedCommand(null);
    setParams({});
  };

  const selected = commands.find((c) => c.command === selectedCommand);

  return (
    <div className="bg-white rounded-lg shadow p-4">
      <h2 className="text-lg font-semibold mb-4">命令面板</h2>

      {/* 命令按钮 */}
      <div className="grid grid-cols-2 gap-2 mb-4">
        {commands.map((cmd) => (
          <button
            key={cmd.command}
            onClick={() => handleCommandSelect(cmd)}
            disabled={!connectedDevice}
            className={`px-4 py-2 rounded border text-left disabled:opacity-50 ${
              selectedCommand === cmd.command
                ? 'border-blue-500 bg-blue-50'
                : 'border-gray-200 hover:border-gray-300'
            }`}
          >
            <span className="mr-2">{cmd.icon}</span>
            {cmd.label}
          </button>
        ))}
      </div>

      {/* 参数表单 */}
      {selected && selected.fields && selected.fields.length > 0 && (
        <div className="space-y-3 mb-4">
          {selected.fields.map((field) => (
            <div key={field.name}>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                {field.label}
              </label>
              <input
                type={field.type}
                value={params[field.name] || ''}
                onChange={(e) =>
                  setParams({
                    ...params,
                    [field.name]: field.type === 'number' ? Number(e.target.value) : e.target.value,
                  })
                }
                className="w-full px-3 py-2 border rounded"
              />
            </div>
          ))}
          <button
            onClick={handleExecute}
            disabled={!connectedDevice}
            className="w-full px-4 py-2 bg-green-500 text-white rounded disabled:bg-gray-300"
          >
            执行
          </button>
        </div>
      )}

      {/* 无参数命令 */}
      {selected && (!selected.fields || selected.fields.length === 0) && (
        <button
          onClick={handleExecute}
          disabled={!connectedDevice}
          className="w-full px-4 py-2 bg-green-500 text-white rounded disabled:bg-gray-300"
        >
          执行 {selected.label}
        </button>
      )}

      {/* 错误提示 */}
      {error && (
        <div className="mt-4 p-3 bg-red-50 border border-red-200 rounded text-red-600 text-sm">
          {error}
        </div>
      )}
    </div>
  );
};
```

**Step 4: 实现 CommandHistory 组件**

```tsx
// console/src/components/CommandHistory.tsx
import React, { useState, useEffect } from 'react';
import { useDeviceStore } from '../stores/deviceStore';

export const CommandHistory: React.FC = () => {
  const { history, fetchHistory, clearHistory, exportYaml, connectedDevice } = useDeviceStore();
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());

  useEffect(() => {
    if (connectedDevice) {
      fetchHistory();
    }
  }, [connectedDevice]);

  const toggleSelect = (id: string) => {
    const newSet = new Set(selectedIds);
    if (newSet.has(id)) {
      newSet.delete(id);
    } else {
      newSet.add(id);
    }
    setSelectedIds(newSet);
  };

  const handleExport = () => {
    exportYaml(selectedIds.size > 0 ? Array.from(selectedIds) : undefined);
  };

  return (
    <div className="bg-white rounded-lg shadow p-4">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold">命令历史</h2>
        <div className="flex gap-2">
          <button
            onClick={handleExport}
            disabled={history.length === 0}
            className="px-3 py-1 text-sm bg-blue-100 text-blue-600 rounded disabled:bg-gray-100"
          >
            导出 YAML
          </button>
          <button
            onClick={clearHistory}
            disabled={history.length === 0}
            className="px-3 py-1 text-sm bg-red-100 text-red-600 rounded disabled:bg-gray-100"
          >
            清空
          </button>
        </div>
      </div>

      {/* 历史列表 */}
      <div className="space-y-2 max-h-96 overflow-y-auto">
        {history.length === 0 ? (
          <p className="text-gray-400 text-center py-4">暂无命令历史</p>
        ) : (
          [...history].reverse().map((record) => (
            <div
              key={record.id}
              onClick={() => toggleSelect(record.id)}
              className={`p-3 rounded border cursor-pointer ${
                selectedIds.has(record.id)
                  ? 'border-blue-500 bg-blue-50'
                  : 'border-gray-200 hover:border-gray-300'
              }`}
            >
              <div className="flex items-center justify-between">
                <span className={`font-medium ${record.success ? 'text-green-600' : 'text-red-600'}`}>
                  {record.command.toUpperCase()}
                </span>
                <span className="text-xs text-gray-400">
                  {new Date(record.timestamp).toLocaleTimeString()}
                </span>
              </div>
              <div className="text-sm text-gray-600 mt-1">
                {formatParams(record.params)}
              </div>
              {record.duration && (
                <div className="text-xs text-gray-400 mt-1">
                  耗时: {record.duration}
                </div>
              )}
            </div>
          ))
        )}
      </div>

      {selectedIds.size > 0 && (
        <div className="mt-4 p-2 bg-blue-50 rounded text-sm text-blue-600">
          已选择 {selectedIds.size} 条命令
        </div>
      )}
    </div>
  );
};

function formatParams(params: Record<string, any>): string {
  const entries = Object.entries(params);
  if (entries.length === 0) return '-';
  return entries.map(([k, v]) => `${k}: ${v}`).join(', ');
}
```

**Step 5: 更新 App.tsx**

```tsx
// console/src/App.tsx
import React from 'react';
import { DeviceList } from './components/DeviceList';
import { ScreenPreview } from './components/ScreenPreview';
import { CommandPanel } from './components/CommandPanel';
import { CommandHistory } from './components/CommandHistory';
import { useDeviceStore } from './stores/deviceStore';

function App() {
  const { exportYaml } = useDeviceStore();

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
          <h1 className="text-2xl font-bold text-gray-800">go-uop Console</h1>
          <button
            onClick={() => exportYaml()}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            导出全部 YAML
          </button>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 py-6">
        <div className="grid grid-cols-12 gap-6">
          {/* 左侧: 设备列表 */}
          <div className="col-span-3">
            <DeviceList />
          </div>

          {/* 中间: 屏幕预览 */}
          <div className="col-span-5">
            <ScreenPreview />
          </div>

          {/* 右侧: 命令面板 */}
          <div className="col-span-4 space-y-4">
            <CommandPanel />
            <CommandHistory />
          </div>
        </div>
      </main>
    </div>
  );
}

export default App;
```

**Step 6: 提交**

```bash
git add console/src/components/ console/src/App.tsx
git commit -m "feat(console): add React components"
```

---

## Phase 4: 集成与打包

### Task 8: 实现前端资源嵌入

**Files:**
- Create: `embed/embed.go`

**Step 1: 创建 embed 目录和文件**

```go
// embed/embed.go
package embed

import (
    "embed"
    "io/fs"
    "net/http"
)

// ConsoleFS holds the embedded console frontend files.
//
//go:embed all:_out
var ConsoleFS embed.FS

// GetFS returns the console filesystem for use with http.FileSystem.
// It strips the "all:" prefix from the embed directive.
func GetFS() (fs.FS, error) {
    sub, err := fs.Sub(ConsoleFS, "_out")
    if err != nil {
        return nil, err
    }
    return sub, nil
}

// HTTPFileSystem returns an http.FileSystem that serves the console frontend.
func HTTPFileSystem() http.FileSystem {
    f, err := GetFS()
    if err != nil {
        // Return a dummy filesystem that will show an error
        return http.Dir(".")
    }
    return http.FS(f)
}
```

**Step 2: 修改 server.go 添加前端服务**

Modify: `internal/console/server.go`

```go
package console

import (
    "embed"
    "io/fs"
    "net/http"
)

// ConsoleAssets holds the embedded console frontend files.
//
//go:embed all:_out
var ConsoleAssets embed.FS

func (s *Server) serveFrontend(w http.ResponseWriter, r *http.Request) {
    // 尝试获取文件
    path := r.URL.Path
    if path == "/" {
        path = "/index.html"
    }

    // 去掉前导斜杠，因为 embed.FS 不带前导斜杠
    path = strings.TrimPrefix(path, "/")

    data, err := fs.ReadFile(ConsoleAssets, "_out/"+path)
    if err != nil {
        // 文件不存在，返回 index.html (SPA 路由处理)
        data, err = fs.ReadFile(ConsoleAssets, "_out/index.html")
        if err != nil {
            http.NotFound(w, r)
            return
        }
    }

    // 设置 Content-Type
    contentType := getContentType(path)
    w.Header().Set("Content-Type", contentType)
    w.Write(data)
}

func getContentType(path string) string {
    ext := strings.ToLower(path[strings.LastIndex(path, "."):])
    switch ext {
    case ".html":
        return "text/html"
    case ".js":
        return "application/javascript"
    case ".css":
        return "text/css"
    case ".png":
        return "image/png"
    case ".svg":
        return "image/svg+xml"
    case ".json":
        return "application/json"
    default:
        return "text/plain"
    }
}
```

**Step 3: 更新路由**

Modify: `internal/console/server.go`

```go
func (s *Server) setupRoutes() {
    // API 路由
    s.mux.HandleFunc("/api/devices", s.handleDevices)
    s.mux.HandleFunc("/api/devices/connect", s.handleConnect)
    s.mux.HandleFunc("/api/devices/", s.handleDeviceOps)
    s.mux.HandleFunc("/api/commands/history", s.handleHistory)
    s.mux.HandleFunc("/api/export/yaml", s.handleYamlExport)

    // 前端路由 (SPA)
    s.mux.HandleFunc("/", s.serveFrontend)
}
```

**Step 4: 更新 main.go 添加构建命令**

Modify: `cmd/console/main.go`

```go
package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"

    "github.com/liukunup/go-uop/internal/console"
)

var (
    addr        string
    openBrowser bool
    devMode     bool
    buildOnly   bool
)

func main() {
    flag.StringVar(&addr, "addr", ":8080", "HTTP server address")
    flag.BoolVar(&openBrowser, "open", true, "Open browser on start")
    flag.BoolVar(&devMode, "dev", false, "Development mode")
    flag.BoolVar(&buildOnly, "build", false, "Only build frontend")
    flag.Parse()

    // 构建前端
    if err := buildFrontend(); err != nil {
        log.Fatal(err)
    }

    if buildOnly {
        fmt.Println("Frontend built successfully")
        return
    }

    // 启动服务器
    server, err := console.NewServer(addr, devMode)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("🚀 go-uop Console starting on %s\n", addr)

    if openBrowser {
        go openBrowserFunc(addr)
    }

    log.Fatal(server.Start())
}

func buildFrontend() error {
    consoleDir, err := findConsoleDir()
    if err != nil {
        return fmt.Errorf("console directory not found: %w", err)
    }

    outDir := filepath.Join(consoleDir, "_out")

    // 检查是否需要构建
    needBuild := true
    if _, err := os.Stat(outDir); err == nil {
        // 检查前端源码是否更新
        statOut, _ := os.Stat(outDir)
        statSrc, _ := os.Stat(filepath.Join(consoleDir, "src"))

        if statOut != nil && statSrc != nil && statOut.ModTime().After(statSrc.ModTime()) {
            needBuild = false
        }
    }

    if !needBuild && !devMode {
        return nil
    }

    // 移除旧构建
    os.RemoveAll(outDir)

    // 安装依赖
    fmt.Println("Installing npm dependencies...")
    if err := runCmd("npm", "install"); err != nil {
        return fmt.Errorf("npm install failed: %w", err)
    }

    // 构建
    fmt.Println("Building frontend...")
    if err := runCmd("npm", "run", "build"); err != nil {
        return fmt.Errorf("npm run build failed: %w", err)
    }

    // 移动 dist 到 _out
    distDir := filepath.Join(consoleDir, "dist")
    if err := os.Rename(distDir, outDir); err != nil {
        return fmt.Errorf("failed to move dist: %w", err)
    }

    fmt.Println("Frontend built successfully")
    return nil
}

func runCmd(name string, args ...string) error {
    cmd := exec.Command(name, args...)
    cmd.Dir = findConsoleDir()
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

func findConsoleDir() (string, error) {
    _, filename, _, _ := runtime.Caller(0)
    root := filepath.Dir(filepath.Dir(filename))
    consoleDir := filepath.Join(root, "console")
    if _, err := os.Stat(consoleDir); os.IsNotExist(err) {
        return "", err
    }
    return consoleDir, nil
}

func openBrowserFunc(addr string) {
    // 延迟打开浏览器
}
```

**Step 5: 提交**

```bash
git add embed/embed.go internal/console/server.go cmd/console/main.go
git commit -m "feat(console): add frontend embedding"
```

---

## Final Verification Wave

### Task F1: 编译验证

**Step 1: 运行编译**

Run: `go build -o uop-console ./cmd/console/`
Expected: 编译成功，生成 uop-console 二进制

### Task F2: 启动验证

**Step 1: 启动应用**

Run: `./uop-console -open=false`
Expected: 服务启动在 :8080

**Step 2: 测试 API**

```bash
# 获取设备列表
curl http://localhost:8080/api/devices

# 获取页面
curl http://localhost:8080/
```

### Task F3: 前端功能测试

**Step 1: 打开浏览器访问**

URL: http://localhost:8080

**Step 2: 验证组件渲染**
- 设备列表显示
- 屏幕预览区域
- 命令面板按钮
- 命令历史区域

### Task F4: 集成测试

**Step 1: 连接 Android 设备**

```bash
# 确保有 Android 设备连接
adb devices
```

**Step 2: 执行命令测试**

```bash
curl -X POST http://localhost:8080/api/devices/android-<serial>/commands \
  -H "Content-Type: application/json" \
  -d '{"command":"tap","params":{"x":100,"y":200}}'
```

---

## Commit Strategy

- Task 1: `feat(console): add server entry point and route structure`
- Task 2: `feat(console): add React frontend scaffolding`
- Task 3: `feat(console): add DeviceManager and types`
- Task 4: `feat(console): add HTTP handlers`
- Task 5: `feat(console): add HistoryManager and YAML export`
- Task 6: `feat(console): add API client and Zustand store`
- Task 7: `feat(console): add React components`
- Task 8: `feat(console): add frontend embedding`
- Final: `feat(console): complete web console implementation`

---

## Success Criteria

1. ✅ Go 后端编译成功
2. ✅ 前端资源正确嵌入
3. ✅ API 端点正常工作
4. ✅ 设备列表显示
5. ✅ 屏幕截图显示
6. ✅ 命令执行成功
7. ✅ 命令历史记录
8. ✅ YAML 导出成功
