# go-uop Web Console 设计文档

**日期**: 2026-03-27
**状态**: 设计完成，待实现

---

## 1. 概述

### 1.1 项目目标

构建一个 Web 控制台应用，用于移动设备调试：

- **前后端分离**：Go HTTP Server + React SPA
- **单文件部署**：Go 二进制包含 React 前端（使用 `embed` 包）
- **设备管理**：支持 iOS (WebDriverAgent) 和 Android (ADB) 设备
- **命令调试**：支持点击、输入、滑动、按键等完整操作
- **命令历史**：记录执行过的命令，支持回放和重新执行
- **YAML 导出**：导出操作步骤为 Maestro 兼容的 YAML 脚本

### 1.2 技术栈

| 层级 | 技术 | 说明 |
|------|------|------|
| 前端 | React 18 + TypeScript | 单页应用 |
| UI 框架 | Tailwind CSS | 快速布局 |
| 状态管理 | Zustand | 轻量级状态管理 |
| 后端 | Go 1.22+ | HTTP 服务 |
| 设备驱动 | 现有 ios/android 包 | 复用现有代码 |
| 构建工具 | Vite | 前端构建 |
| 嵌入打包 | Go embed | 前端资源打包 |

---

## 2. 系统架构

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│                        用户浏览器                                │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                    React SPA                              │  │
│  │  ┌─────────┐  ┌─────────────────┐  ┌──────────────────┐   │  │
│  │  │ 设备列表 │  │    屏幕预览      │  │    命令面板      │   │  │
│  │  │         │  │                 │  │                  │   │  │
│  │  │ ○ iPhone│  │   ┌─────────┐   │  │ [Tap] [Input]   │   │  │
│  │  │ ○ Pixel │  │   │  实际   │   │  │ [Swipe][Key]   │   │  │
│  │  │         │  │   │  屏幕   │   │  │ [Launch][Home] │   │  │
│  │  └─────────┘  │   │  截图   │   │  ├──────────────────┤   │  │
│  │               │   └─────────┘   │  │    命令历史      │   │  │
│  │               │                 │  │  1. Tap(100,200)│   │  │
│  │  [连接设备]    │  [刷新截图]      │  │  2. Input("abc")│   │  │
│  │               │                 │  │  [导出YAML]     │   │  │
│  │  设备信息:     │  坐标: (X, Y)    │  └──────────────────┘   │  │
│  │  - 平台: iOS  │  点击屏幕可获取   │                         │  │
│  │  - 型号: ...  │                 │                         │  │
│  └───────────────┴─────────────────┴─────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              │ HTTP/REST API
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Go HTTP Server                             │
│                         :8080                                    │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │                    API Handlers                            │ │
│  │  GET  /api/devices          - 列出设备                      │ │
│  │  POST /api/devices/connect  - 连接设备                      │ │
│  │  GET /api/devices/:id/info  - 设备信息                      │ │
│  │  GET /api/devices/:id/screenshot - 获取截图                │ │
│  │  POST /api/devices/:id/commands - 执行命令                 │ │
│  │  GET /api/devices/:id/source  - 获取UI源码                 │ │
│  │  GET /api/commands/history   - 命令历史                     │ │
│  │  GET /api/export/yaml       - 导出YAML                     │ │
│  └────────────────────────────────────────────────────────────┘ │
│                              │                                  │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │                   Device Manager                           │ │
│  │  ┌──────────────┐    ┌──────────────┐                      │ │
│  │  │ iOS Driver   │    │ Android Driver│                     │ │
│  │  │ (WebDriver)  │    │ (ADB)        │                      │ │
│  │  └──────────────┘    └──────────────┘                      │ │
│  └────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 目录结构

```
go-uop/
├── cmd/
│   └── console/              # 新增：Web Console 入口
│       └── main.go
├── console/                  # 新增：Web Console 前端
│   ├── src/
│   │   ├── components/       # React 组件
│   │   │   ├── DeviceList.tsx
│   │   │   ├── ScreenPreview.tsx
│   │   │   ├── CommandPanel.tsx
│   │   │   └── CommandHistory.tsx
│   │   ├── hooks/            # 自定义 Hooks
│   │   │   └── useDevice.ts
│   │   ├── stores/           # Zustand Store
│   │   │   └── deviceStore.ts
│   │   ├── types/            # TypeScript 类型
│   │   │   └── index.ts
│   │   ├── api/              # API 客户端
│   │   │   └── client.ts
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── index.html
│   ├── package.json
│   ├── vite.config.ts
│   └── tailwind.config.js
├── internal/                 # 新增：内部包
│   └── console/
│       ├── server.go         # HTTP Server
│       ├── handler.go        # API Handlers
│       ├── device.go         # 设备管理
│       ├── history.go        # 命令历史
│       └── yaml.go           # YAML 导出
├── embed/                   # 新增：前端资源嵌入
│   └── embed.go
└── ... (existing files)
```

---

## 3. API 设计

### 3.1 设备相关 API

#### GET /api/devices
列出所有可用设备

**响应**:
```json
{
  "devices": [
    {
      "id": "ios-123456",
      "platform": "ios",
      "name": "iPhone 15 Pro",
      "serial": "123456",
      "status": "connected",
      "model": "iPhone 15 Pro",
      "address": "http://localhost:8100"
    },
    {
      "id": "android-emulator-5554",
      "platform": "android",
      "name": "Pixel 6",
      "serial": "emulator-5554",
      "status": "available",
      "model": "Pixel 6"
    }
  ]
}
```

#### POST /api/devices/connect
连接指定设备

**请求**:
```json
{
  "platform": "ios",
  "serial": "123456",
  "address": "http://localhost:8100"
}
```

**响应**:
```json
{
  "success": true,
  "device": {
    "id": "ios-123456",
    "platform": "ios",
    "name": "iPhone 15 Pro",
    "status": "connected"
  }
}
```

#### GET /api/devices/:id/screenshot
获取设备屏幕截图

**响应**: PNG 图像数据
**Content-Type**: `image/png`

#### GET /api/devices/:id/info
获取设备信息

**响应**:
```json
{
  "platform": "ios",
  "bundleId": "com.example.app",
  "model": "iPhone 15 Pro",
  "screenSize": {"width": 393, "height": 852}
}
```

### 3.2 命令执行 API

#### POST /api/devices/:id/commands
执行命令

**请求**:
```json
{
  "command": "tap",
  "params": {"x": 100, "y": 200}
}
```

**响应**:
```json
{
  "success": true,
  "output": null,
  "duration": "45ms"
}
```

**支持的命令**:

| 命令 | 参数 | 说明 |
|------|------|------|
| `tap` | x, y | 点击坐标 |
| `double_tap` | x, y | 双击坐标 |
| `long_press` | x, y, duration | 长按 |
| `swipe` | x1, y1, x2, y2, duration | 滑动 |
| `input` | text | 输入文本 |
| `launch` | - | 启动应用 |
| `terminate` | - | 终止应用 |
| `press_key` | key | 按键 (Android) |
| `screenshot` | - | 截图 |
| `get_source` | - | 获取UI源码 |

#### GET /api/commands/history
获取命令历史

**响应**:
```json
{
  "history": [
    {
      "id": "cmd-001",
      "timestamp": "2026-03-27T10:00:00Z",
      "command": "tap",
      "params": {"x": 100, "y": 200},
      "success": true
    },
    {
      "id": "cmd-002",
      "timestamp": "2026-03-27T10:00:05Z",
      "command": "input",
      "params": {"text": "hello"},
      "success": true
    }
  ]
}
```

### 3.3 导出 API

#### GET /api/export/yaml
导出命令历史为 YAML

**响应**:
```yaml
name: debug-session
steps:
  - tapOn:
      x: 100
      y: 200
  - inputText:
      text: "hello"
  - tapOn:
      text: "登录"
```

---

## 4. 前端组件设计

### 4.1 布局结构

```
┌────────────────────────────────────────────────────────────────────┐
│  Header: go-uop Console                              [导出 YAML]  │
├────────────┬─────────────────────────────┬────────────────────────┤
│            │                             │                        │
│  设备列表   │      屏幕预览区域            │    命令面板            │
│            │                             │                        │
│ ┌────────┐ │  ┌─────────────────────┐   │  ┌──────────────────┐  │
│ │ ○ iOS  │ │  │                     │   │  │ 快捷命令          │  │
│ │   设备1│ │  │    设备屏幕截图      │   │  │                  │  │
│ │   设备2│ │  │    (可点击)         │   │  │ [📱 Tap]        │  │
│ │         │  │                     │   │  │ [⌨️ Input]      │  │
│ ├────────┤ │  │                     │   │  │ [👆 Swipe]      │  │
│ │ ○ Anrd │ │  └─────────────────────┘   │  │ [🔘 KeyEvent]   │  │
│ │   设备3│ │                             │  │ [🚀 Launch]     │  │
│ └────────┘ │  [📸 刷新截图] [📋 获取源码] │  │ [❌ Terminate]  │  │
│            │                             │  ├──────────────────┤  │
│ [连接设备]  │  坐标: (X: 100, Y: 200)    │  │ 命令历史         │  │
│ [断开设备]  │  (点击屏幕显示)            │  │                  │  │
│            │                             │  │ 1. tap 100,200  │  │
│ ┌────────┐ │                             │  │ 2. input "abc"  │  │
│ │设备信息 │ │                             │  │ 3. swipe ...    │  │
│ │平台:iOS │ │                             │  │                  │  │
│ │型号:XXX │ │                             │  │ [清空] [回放]   │  │
│ └────────┘ │                             │  └──────────────────┘  │
│            │                             │                        │
└────────────┴─────────────────────────────┴────────────────────────┘
```

### 4.2 组件规格

#### DeviceList (设备列表)
- 显示 iOS/Android 两组设备
- 每组下列出可用设备（从 API 获取）
- 选中状态高亮显示
- 连接/断开按钮
- 设备信息面板

**状态**:
```typescript
interface DeviceState {
  devices: Device[];
  selectedDevice: Device | null;
  connectedDevice: Device | null;
  isConnecting: boolean;
}
```

#### ScreenPreview (屏幕预览)
- 显示设备截图（定时刷新或手动刷新）
- 支持点击屏幕获取坐标
- 点击坐标后自动填充到命令面板
- 显示当前鼠标位置坐标

**交互**:
- 点击屏幕 → 获取坐标 → 高亮显示
- 刷新按钮 → 重新获取截图
- 源码按钮 → 获取 UI 层级结构

#### CommandPanel (命令面板)
- 快捷命令按钮组
- 命令参数输入表单
- 执行按钮
- 实时结果显示

**命令表单**:
```typescript
interface CommandForm {
  type: 'tap' | 'input' | 'swipe' | 'press_key' | 'launch' | 'terminate';
  params: Record<string, any>;
}
```

#### CommandHistory (命令历史)
- 时间线形式展示历史命令
- 每条记录显示：时间、命令、参数、状态
- 支持清空历史
- 支持选中命令重新执行
- 勾选多条命令后可导出为 YAML

---

## 5. 后端设计

### 5.1 Server 结构

```go
type ConsoleServer struct {
    mux         *http.ServeMux
    deviceMgr   *DeviceManager
    historyMgr  *HistoryManager
    addr        string
}

type DeviceManager struct {
    mu          sync.RWMutex
    devices     map[string]core.Device  // deviceID -> Device
    connected   map[string]core.Device  // 当前连接的设备
}

type HistoryManager struct {
    mu       sync.RWMutex
    history  []CommandRecord
    maxSize  int
}
```

### 5.2 命令执行流程

```
1. 接收命令请求
   └─> POST /api/devices/:id/commands

2. 验证设备连接
   └─> deviceMgr.GetConnected(id)

3. 创建命令记录
   └─> historyMgr.StartRecord(command)

4. 执行命令
   ├─> tap     → device.Tap(x, y)
   ├─> input   → device.SendKeys(text)
   ├─> swipe   → device.Swipe(x1, y1, x2, y2, duration)
   ├─> launch  → device.Launch()
   ├─> press_key → device.PressKey(code)
   └─> screenshot → device.Screenshot()

5. 更新记录状态
   └─> historyMgr.FinishRecord(success, output)

6. 返回结果
   └─> JSON response
```

### 5.3 前端资源嵌入

```go
//go:embed all:_console_dist
var consoleFS embed.FS

func (s *ConsoleServer) serveFrontend(w http.ResponseWriter, r *http.Request) {
    // 尝试从 embed.FS 读取文件
    // 如果文件不存在，返回 index.html (SPA 路由处理)
}
```

---

## 6. 数据模型

### 6.1 设备模型

```typescript
interface Device {
  id: string;
  platform: 'ios' | 'android';
  name: string;
  serial: string;
  status: 'available' | 'connected' | 'error';
  model?: string;
  address?: string;    // iOS WDA 地址
  packageName?: string; // Android 包名
}

interface DeviceInfo {
  platform: string;
  screenSize: { width: number; height: number };
  [key: string]: any;
}
```

### 6.2 命令模型

```typescript
interface Command {
  id: string;
  timestamp: string;
  type: CommandType;
  params: Record<string, any>;
  success: boolean;
  output?: string;
  duration: string;
}

type CommandType = 
  | 'tap' 
  | 'double_tap' 
  | 'long_press' 
  | 'swipe' 
  | 'input' 
  | 'launch' 
  | 'terminate'
  | 'press_key'
  | 'screenshot'
  | 'get_source';
```

### 6.3 YAML 导出格式

```yaml
name: debug-session-20260327
appId: com.example.app
steps:
  - tapOn:
      x: 100
      y: 200
  - inputText:
      text: "hello world"
  - tapOn:
      text: "登录"
```

---

## 7. 命令详情

### 7.1 Tap (点击)

**参数**: `{ x: number, y: number }`
**UI**: 点击屏幕预览获取坐标，或手动输入

### 7.2 Double Tap (双击)

**参数**: `{ x: number, y: number }`

### 7.3 Long Press (长按)

**参数**: `{ x: number, y: number, duration: number }` (duration in ms)

### 7.4 Swipe (滑动)

**参数**: `{ x1: number, y1: number, x2: number, y2: number, duration: number }`

### 7.5 Input (输入文本)

**参数**: `{ text: string }`

### 7.6 Launch (启动应用)

**参数**: 无 (使用连接时指定的包名/Bundle ID)

### 7.7 Terminate (终止应用)

**参数**: 无

### 7.8 Press Key (按键 - 仅 Android)

**参数**: `{ key: number }` (Android keycode)

常用 keycode:
- 3 = HOME
- 4 = BACK
- 82 = MENU
- 26 = POWER

### 7.9 Screenshot (截图)

**参数**: 无
**返回**: PNG 图像

### 7.10 Get Source (获取UI源码)

**参数**: 无
**返回**: XML/JSON 格式的 UI 层级结构

---

## 8. 错误处理

### 8.1 错误响应格式

```json
{
  "error": {
    "code": "DEVICE_NOT_CONNECTED",
    "message": "设备未连接，请先连接设备",
    "details": {}
  }
}
```

### 8.2 错误码

| 错误码 | 说明 |
|--------|------|
| `DEVICE_NOT_FOUND` | 设备不存在 |
| `DEVICE_NOT_CONNECTED` | 设备未连接 |
| `DEVICE_CONNECTION_FAILED` | 设备连接失败 |
| `COMMAND_EXECUTION_FAILED` | 命令执行失败 |
| `INVALID_PARAMETERS` | 参数无效 |
| `SCREENSHOT_FAILED` | 截图失败 |

---

## 9. 实现计划

### Phase 1: 项目初始化
- [ ] 创建 `cmd/console/main.go` 入口
- [ ] 创建 `internal/console/` 包结构
- [ ] 初始化 React 项目 (`console/` 目录)
- [ ] 配置 Vite + Tailwind

### Phase 2: 后端 API
- [ ] 实现 DeviceManager
- [ ] 实现 HistoryManager
- [ ] 实现 HTTP Handlers
- [ ] 测试 API 端点

### Phase 3: 前端组件
- [ ] DeviceList 组件
- [ ] ScreenPreview 组件
- [ ] CommandPanel 组件
- [ ] CommandHistory 组件
- [ ] 状态管理 (Zustand)

### Phase 4: 集成与打包
- [ ] 前端嵌入 Go 二进制
- [ ] 构建脚本
- [ ] 测试完整流程

---

## 10. 依赖

### Go 依赖
```go
// 新增依赖
github.com/gorilla/mux      // HTTP 路由 (或用标准库)
github.com/rivo/uniseg     // Unicode 处理
```

### 前端依赖
```json
{
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "zustand": "^4.5.0",
    "axios": "^1.6.0"
  },
  "devDependencies": {
    "typescript": "^5.3.0",
    "vite": "^5.0.0",
    "tailwindcss": "^3.4.0",
    "@types/react": "^18.2.0",
    "@types/react-dom": "^18.2.0"
  }
}
```

---

## 11. 构建流程

### 开发模式

```bash
# Terminal 1: 运行 Go 后端
cd /Users/liukunup/Documents/repo/go-uop
go run cmd/console/main.go -dev

# Terminal 2: 运行 React 开发服务器
cd /Users/liukunup/Documents/repo/go-uop/console
npm run dev
```

访问 `http://localhost:5173` 进行开发。

### 生产构建

```bash
# 构建 React 前端
cd console
npm run build

# 构建 Go 二进制 (包含前端资源)
cd ..
go build -o uop-console cmd/console/main.go

# 运行
./uop-console
```

默认监听 `http://localhost:8080`，自动打开浏览器。

---

## 12. 命令行参数

```bash
./uop-console [选项]

选项:
  -addr string
        HTTP 服务地址 (默认 ":8080")
  -open
        启动后自动打开浏览器 (默认 true)
  -dev
        开发模式，前端使用独立服务器
```

---

## 13. 后续扩展 (TODO)

- [ ] WebSocket 实时推送截图 (减少轮询)
- [ ] 批量命令执行
- [ ] 脚本管理 (保存/加载 YAML)
- [ ] 设备屏幕录制
- [ ] 元素定位器 (基于 XML source)
- [ ] AI 辅助命令生成
