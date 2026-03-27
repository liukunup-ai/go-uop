# go-uop 架构文档

> 生成时间: 2026-03-27
> 项目: github.com/liukunup/go-uop

---

## 1. 项目概述

**go-uop** 是一个 Go 语言实现的统一移动端自动化测试框架，支持 iOS（通过 WebDriverAgent）和 Android（通过 ADB）两大平台。

### 核心特性

- **统一设备接口**：一套 API 同时支持 iOS 和 Android
- **链式 Fluent API**：可读性强的链式调用构建动作
- **YAML 测试运行器**：支持 Maestro 风格的 YAML 测试流程定义
- **视觉模块**：基于模板匹配的图像识别自动化
- **AI 集成**：OpenAI Provider 智能自动化
- **并行执行**：跨多设备并行测试
- **重试机制**：可配置指数退避重试

---

## 2. 架构分层图

```
┌─────────────────────────────────────────────────────────────────┐
│                        User Layer (用户层)                        │
│  ┌──────────────────────┐    ┌────────────────────────────────┐ │
│  │   Go API (uop*.go)   │    │   YAML Runner (maestro CLI)    │ │
│  └──────────────────────┘    └────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│                     Command Layer (命令层)                        │
│  ┌──────────────┐  ┌──────────────┐  ┌────────────────────────┐ │
│  │   Maestro   │  │     YAML     │  │   Maestro Commands     │ │
│  │  Translator │  │   Evaluator  │  │  (tap, swipe, input)   │ │
│  └──────────────┘  └──────────────┘  └────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│                  Platform Drivers (平台驱动层)                     │
│  ┌──────────────────────────┐    ┌────────────────────────────┐ │
│  │     iOS Device (ios/)    │    │  Android Device (android/) │ │
│  │  ┌────────────────────┐  │    │  ┌─────────────────────┐  │ │
│  │  │  WDA Client (wda/) │  │    │  │   ADB Client (adb/) │  │ │
│  │  │  - HTTP REST API    │  │    │  │   - Shell Commands  │  │ │
│  │  │  - Session Mgmt     │  │    │  │   - Package Mgmt     │  │ │
│  │  └────────────────────┘  │    │  └─────────────────────┘  │ │
│  └──────────────────────────┘    └────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│                     Core Modules (核心模块)                        │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌──────────────┐  │
│  │  Selector  │ │   Action   │ │   Vision   │ │    Retry     │  │
│  │  (元素定位) │ │  (动作定义) │ │ (模板匹配)  │ │ (指数退避)   │  │
│  └────────────┘ └────────────┘ └────────────┘ └──────────────┘  │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐                   │
│  │  Parallel  │ │   Report   │ │     AI     │                   │
│  │  (并行执行) │ │  (报告生成) │ │  (AI集成)  │                   │
│  └────────────┘ └────────────┘ └────────────┘                   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 3. 目录结构

```
go-uop/
├── cmd/
│   └── maestro/main.go          # Maestro CLI 入口
│
├── core/
│   ├── device.go               # Device 接口定义
│   ├── option.go               # 设备选项
│   └── errors.go               # 核心错误定义
│
├── ios/
│   ├── device.go               # iOS 设备实现
│   ├── option.go               # iOS 选项
│   ├── device_test.go
│   └── wda/                    # WebDriverAgent HTTP 客户端
│       ├── client.go           # WDA HTTP 客户端
│       ├── session.go          # 会话管理
│       ├── element.go          # 元素操作
│       ├── alert.go            # 弹窗处理
│       ├── screenshot.go       # 截图
│       ├── app.go              # 应用操作
│       └── protocol.go         # WDA 协议
│
├── android/
│   ├── device.go               # Android 设备实现
│   ├── option.go              # Android 选项
│   ├── device_test.go
│   └── adb/                   # ADB 客户端
│       ├── client.go          # ADB 主客户端
│       ├── shell.go           # Shell 命令
│       ├── input.go           # 输入操作 (tap, swipe, text)
│       ├── screenshot.go      # 截图
│       ├── activity.go        # Activity 管理
│       ├── package.go         # 包管理
│       └── client_test.go
│
├── internal/
│   ├── selector/              # 元素定位器
│   │   ├── selector.go        # Selector 类型和支持方法
│   │   └── selector_test.go
│   │
│   ├── action/                # 动作定义
│   │   ├── action.go          # 各种 Action 类型
│   │   └── action_test.go
│   │
│   ├── vision/                # 视觉模块
│   │   └── template.go        # 模板匹配
│   │
│   ├── retry/                 # 重试机制
│   │   └── retry.go           # 指数退避重试
│   │
│   ├── parallel/              # 并行执行器
│   │   └── executor.go        # Worker Pool 实现
│   │
│   ├── report/                # 测试报告
│   │   └── generator.go
│   │
│   └── ...
│
├── yaml/
│   ├── parser.go              # YAML 解析
│   ├── command.go             # 命令结构定义
│   ├── evaluator.go           # 变量求值器
│   ├── evaluator_test.go
│   ├── parser_test.go
│   └── commands/              # YAML 命令执行
│       ├── control.go         # 控制流 (if, foreach, while)
│       ├── control_test.go
│       ├── tap.go
│       ├── swipe.go
│       ├── input.go
│       ├── navigation.go
│       ├── app.go
│       └── assert.go
│
├── maestro/
│   ├── maestro.go             # Maestro 文件解析
│   ├── types.go               # Maestro 命令类型定义
│   ├── translator.go           # Maestro → Action 翻译器
│   ├── executor.go            # 动作执行器
│   ├── config.go              # 配置
│   ├── parser.go              # Maestro YAML 解析
│   ├── errors.go              # 错误定义
│   ├── integration_test.go
│   └── commands/              # Maestro 命令实现
│       ├── flow.go
│       ├── tap.go
│       ├── swipe.go
│       ├── input.go
│       ├── navigation.go
│       ├── app.go
│       └── assert.go
│
├── ai/
│   ├── provider.go            # AI Provider 接口
│   ├── openai.go              # OpenAI 实现
│   └── bigmodel.go           # BigModel 实现
│
├── uop.go                    # 公共 API 入口
├── uop_fluent.go             # Fluent API (ActionBuilder)
├── uop_assert.go             # 断言 API
├── uop_test.go
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## 4. 核心接口设计

### 4.1 Device 接口 (core/device.go)

```go
type Device interface {
    Platform() Platform           // 返回平台类型 (ios/android)
    Info() (map[string]interface{}, error)  // 设备信息
    Screenshot() ([]byte, error)  // 截图
    Tap(x, y int) error           // 点击坐标
    SendKeys(text string) error   // 输入文本
    Launch() error                // 启动应用
    Close() error                 // 关闭连接
}
```

### 4.2 Provider 接口 (ai/provider.go)

```go
type Provider interface {
    Name() string
    // Assert 判断实际文本是否符合预期描述
    // 返回: (是否匹配, 置信度[0-1], 错误)
    Assert(ctx context.Context, actual, expected string) (bool, float64, error)
    // Rerank 根据匹配度对文本重新排序
    // 返回: 按分数降序排列的文本列表
    Rerank(ctx context.Context, texts []string, expected string) ([]RankedText, error)
}

type RankedText struct {
    Text  string
    Score float64
}

type Config struct {
    APIKey      string
    BaseURL     string
    Model       string
    TopP        float64
    Temperature float64
    MaxTokens   int
}
```

### 4.3 Action 接口 (internal/action/action.go)

```go
type Action interface {
    Do() error
}

// 支持的 Action 类型
type TapAction struct {
    X, Y    int
    Element *selector.Selector
}

type SwipeAction struct {
    StartX, StartY int
    EndX, EndY     int
    Duration       time.Duration
}

type SendKeysAction struct {
    Text              string
    Element           *selector.Selector
    Secure            bool
    Enter             bool
    ClearExistingText bool
}

type LaunchAction struct {
    AppID      string
    Arguments  []string
    WaitIdle   bool
    ClearState bool
}

type PressKeyAction struct {
    KeyCode int
}

type WaitAction struct {
    Duration time.Duration
    Element  *selector.Selector
    Optional bool
}

type AssertAction struct {
    Element   *selector.Selector
    MustExist bool
    Timeout   time.Duration
}

type ScreenshotAction struct {
    Name string
    Path string
}

type RunFlowAction struct {
    SubflowPath string
    EnvVars     map[string]string
    Depth       int
}
```

### 4.4 Selector 类型 (internal/selector/selector.go)

支持 6 种元素定位方式：

```go
type SelectorType int

const (
    SelectorTypeText SelectorType = iota      // 文本匹配
    SelectorTypeID                            // ID 匹配
    SelectorTypeXPath                         // XPath 定位
    SelectorTypeClassName                     // 类名定位
    SelectorTypePredicate                     // 谓词定位 (iOS 专用)
    SelectorTypeClassChain                    // 类链定位 (iOS 专用)
)

// Selector 构造器
func ByText(text string) *Selector
func ByID(id string) *Selector
func ByXPath(xpath string) *Selector
func ByClassName(class string) *Selector
func ByPredicate(predicate string) *Selector   // iOS only
func ByClassChain(chain string) *Selector     // iOS only

// 支持正则表达式，格式: /regex_pattern/
func NewSelector(value string) *Selector
```

---

## 5. 平台驱动详解

### 5.1 iOS 驱动 (ios/)

通过 **WebDriverAgent (WDA)** 实现 iOS 自动化。WDA 是 Facebook 开源的 iOS 自动化框架。

**WDA Client** (ios/wda/client.go)：
- 基于 HTTP REST API 与 WDA 通信
- 默认地址：`http://localhost:8100`
- 支持 Session 管理

**主要能力**：
- 元素查找与操作
- 手势操作 (tap, swipe, pinch)
- 屏幕截图
- 应用启停
- Alert 处理

**设备初始化**：
```go
device, err := ios.NewDevice("com.example.app",
    ios.WithAddress("http://localhost:8100"))
```

### 5.2 Android 驱动 (android/)

通过 **Android Debug Bridge (ADB)** 实现 Android 自动化。

**ADB Client** (android/adb/client.go)：
- 通过执行 `adb` 命令与设备通信
- 支持多设备并行

**主要能力**：
| 方法 | 说明 |
|------|------|
| `Tap(x, y)` | 点击坐标 |
| `SendText(text)` | 输入文本 |
| `Swipe(x1, y1, x2, y2, duration)` | 滑动 |
| `StartActivity(pkg/activity)` | 启动 Activity |
| `StopPackage(pkg)` | 停止应用 |
| `Screenshot()` | 屏幕截图 |
| `PressKey(keyCode)` | 按键 |
| `Shell(cmd)` | 执行 Shell 命令 |

**设备初始化**：
```go
device, err := android.NewDevice(
    android.WithSerial("emulator-5554"),
    android.WithPackage("com.example.app"))
```

---

## 6. YAML 命令系统

### 6.1 标准 YAML 命令 (yaml/command.go)

支持 Maestro 风格命令：

```yaml
name: login flow
platform: ios
steps:
  - launch: com.example.app
  - tapOn:
      text: "登录"
  - inputText:
      text: "user@example.com"
      element:
        id: "email_input"
  - tapOn:
      text: "登录"
```

**支持的命令**：

| 命令 | 说明 |
|------|------|
| `launch` | 启动应用 |
| `terminate` | 终止应用 |
| `install` | 安装应用 |
| `uninstall` | 卸载应用 |
| `tapOn` | 点击元素 |
| `tap` | 点击坐标 |
| `doubleTap` | 双击 |
| `longPress` | 长按 |
| `swipe` | 滑动 |
| `swipeUp/Down/Left/Right` | 方向滑动 |
| `inputText` | 输入文本 |
| `waitFor` | 等待元素出现 |
| `waitForGone` | 等待元素消失 |
| `assertVisible` | 断言元素可见 |
| `assertNotVisible` | 断言元素不可见 |
| `assertTrue` | 断言条件为真 |
| `screenshot` | 截图 |
| `runFlow` | 运行子流程 |
| `if/then/else` | 条件判断 |
| `foreach` | 循环 |
| `while` | 条件循环 |
| `setVariable` | 设置变量 |
| `log` | 日志 |

### 6.2 Maestro YAML (maestro/types.go)

与标准 YAML 类似但有细微差异，专门用于 Maestro 兼容模式。

**Maestro 特有命令**：
```yaml
- back: {}                    # 按返回键
- pressHome: {}              # 按 Home 键
- pressRecentApps: {}        # 按最近应用键
- waitForAnimationEnd: {}    # 等待动画结束
- pressKey: {key: "HOME"}    # 按指定键
- repeat:                    # 重复执行
    times: 3
    do:
      - tapOn: {text: "Next"}
```

---

## 7. 核心模块详解

### 7.1 Selector 模块 (internal/selector/)

元素定位模块，支持多种定位策略和正则匹配：

```go
// 文本定位
selector := selector.ByText("登录")

// ID 定位
selector := selector.ByID("submit_button")

// XPath 定位
selector := selector.ByXPath("//button[@text='提交']")

// iOS 专用定位
selector := selector.ByPredicate("type == 'XCUIElementTypeButton'")
selector := selector.ByClassChain("**/Button[`name == "submit"`]")

// 正则表达式定位
selector := selector.NewSelector("/^登录.*$/")

// 链式设置索引
selector := selector.ByText("取消").SetIndex(1)
```

### 7.2 Retry 模块 (internal/retry/)

指数退避重试策略：

```go
type Config struct {
    MaxAttempts int           // 最大重试次数 (默认 3)
    Delay       time.Duration // 初始延迟 (默认 100ms)
    MaxDelay    time.Duration // 最大延迟 (默认 30s)
    Backoff     float64       // 退避系数 (默认 2.0)
}

// 使用示例
err := retry.Do(func() error {
    return someFlakyOperation()
}, retry.DefaultConfig())

// 自定义配置
err := retry.Do(fn, Config{
    MaxAttempts: 5,
    Delay:       200 * time.Millisecond,
    MaxDelay:    10 * time.Second,
    Backoff:     1.5,
})

// 带 Context 支持 (可取消/超时)
err := retry.DoWithContext(ctx, func(ctx context.Context) error {
    // ...
}, config)

// 直至成功 (有最大时长限制)
err := retry.DoUntilSuccess(fn, 30*time.Second)

// 带结果的版本
result, err := retry.DoWithResult(func() (T, error) {
    return someOperation()
}, config)
```

### 7.3 Parallel 模块 (internal/parallel/)

并行任务执行器：

```go
// 并行任务执行器
executor := parallel.NewExecutor(maxWorkers)
results := executor.Execute([]parallel.Task{
    func() error { return task1() },
    func() error { return task2() },
    func() error { return task3() },
})

// Worker Pool 模式 (推荐用于大量任务)
pool := parallel.NewPool(workers, queueSize)
pool.Submit(task1)
pool.Submit(task2)
result := pool.SubmitAndWait(task3)
pool.Close()

// 收集结果
err := parallel.CollectResults(results)
```

### 7.4 Vision 模块 (internal/vision/)

模板匹配实现：

```go
// 创建模板匹配器
matcher := vision.NewTemplateMatcher(screenshot)

// 设置模板文件
err := matcher.SetTemplate("button_template.png")

// 设置模板字节数据
matcher.SetTemplateBytes(templateData)

// 查找最佳匹配
result, err := matcher.FindBestMatch()
if result != nil {
    fmt.Printf("Found at (%d, %d), size: %dx%d, score: %.2f\n",
        result.X, result.Y, result.Width, result.Height, result.Score)
}

// 图像加载辅助
img, err := vision.LoadImageFromFile("screenshot.png")
img, err := vision.LoadImageFromBytes(imageData)
```

### 7.5 AI 模块 (ai/)

AI 驱动的智能自动化：

```go
// 创建 Provider
provider, err := ai.NewProvider("openai", ai.Config{
    APIKey:      "sk-...",
    BaseURL:     "https://api.openai.com/v1",
    Model:       "gpt-4",
    Temperature: 0.7,
    MaxTokens:   1000,
})

// 断言判断 - 判断实际文本是否符合预期描述
isMatch, confidence, err := provider.Assert(ctx,
    actual:   "错误: 网络连接失败",
    expected: "错误信息应该包含失败原因")

// 文本排序 - 根据匹配度对候选文本排序
ranked, err := provider.Rerank(ctx,
    texts:   []string{"登录成功", "密码错误", "网络超时"},
    expected: "登录失败的原因")
// 返回: [{"密码错误": 0.95}, {"网络超时": 0.80}, {"登录成功": 0.10}]

// 支持多个 Provider
provider, _ = ai.NewProvider("bigmodel", config)
```

---

## 8. Fluent API (uop_fluent.go)

链式调用构建复杂动作序列：

```go
err := uop.NewActionBuilder(device).
    // 基础操作
    Tap(100, 200).                             // 点击坐标
    TapElement(selector.ByID("submit")).       // 点击元素
    SendKeys("hello").                         // 输入文本
    
    // 滑动操作
    Swipe(100, 400, 100, 200).                 // 从 (100,400) 滑到 (100,200)
    SwipeWithDuration(100, 400, 100, 200, 500*time.Millisecond).
    SwipeUp().                                // 上滑 (默认 20% 屏幕高度)
    SwipeUpDistance(0.3).                      // 上滑 30% 屏幕高度
    
    // 应用操作
    Launch().                                 // 启动应用
    
    // 等待
    Wait("1s").                               // 等待 1 秒
    Wait("500").                              // 等待 500 毫秒
    
    // 执行所有动作
    Do()
```

**Fluent API 错误处理**：
- 每个方法都会检查前置错误，遇到错误后后续操作会被跳过
- 最终通过 `Do()` 获取执行结果

---

## 9. Maestro CLI

入口：`cmd/maestro/main.go`

### 9.1 命令行用法

```bash
# 验证 YAML 语法
maestro validate flow.yaml

# 执行测试流程
maestro test flow.yaml --device ios --output ./screenshots

# 指定平台
maestro test flow.yaml --device android

# 指定配置文件
maestro test flow.yaml --config config.yaml

# 查看帮助
maestro --help
maestro -h
```

### 9.2 执行流程

```
用户执行 maestro test flow.yaml --device ios
    ↓
1. 解析命令行参数和 flags
    ↓
2. 读取并解析 YAML 文件 → Flow 结构体
    ↓
3. 根据 platform 创建 Device 连接
    ↓
4. 创建 Translator，将 Flow 转换为 []Action
    ↓
5. 创建 Executor，执行每个 Action
    ↓
6. 失败时自动截图保存
    ↓
7. 输出执行结果
```

---

## 10. 数据流详解

### 10.1 iOS 自动化流程

```
┌─────────────────────────────────────────────────────────────────┐
│  Go Code: ios.NewDevice("com.example.app")                      │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│  ios/device.go → wda.NewClient("http://localhost:8100")          │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│  ios/wda/client.go → HTTP REST API Calls                         │
│  POST /session/{sessionId}/element                             │
│  POST /session/{sessionId}/execute                             │
│  GET  /session/{sessionId}/screenshot                          │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│  WebDriverAgent (设备上运行)                                     │
│  - REST API Server (Port 8100)                                  │
│  - XCUI Test Framework                                          │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│  iOS Device                                                     │
│  - XCUITest APIs                                                │
│  - Accessibility Tree                                           │
└─────────────────────────────────────────────────────────────────┘
```

### 10.2 Android 自动化流程

```
┌─────────────────────────────────────────────────────────────────┐
│  Go Code: android.NewDevice(WithSerial("emulator-5554"))        │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│  android/device.go → adb.NewClient("emulator-5554")            │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│  android/adb/client.go → Shell Commands                         │
│  adb -s emulator-5554 shell input tap x y                       │
│  adb -s emulator-5554 shell input text "hello"                 │
│  adb -s emulator-5554 shell screencap -p > screen.png          │
│  adb -s emulator-5554 shell am start -n pkg/activity          │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│  Android Device                                                 │
│  - uiautomator2 (API 18+)                                       │
│  - shell input (基础输入)                                        │
│  - am/pm (应用管理)                                              │
└─────────────────────────────────────────────────────────────────┘
```

### 10.3 YAML 执行流程

```
┌─────────────────────────────────────────────────────────────────┐
│  flow.yaml                                                      │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ name: login flow                                        │   │
│  │ steps:                                                  │   │
│  │   - tapOn: {text: "Login"}                             │   │
│  │   - inputText: {text: "user@example.com"}              │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│  maestro.ParseFlow() / yaml.ParseFlow()                         │
│  → MaestroFlow / Flow 结构体                                    │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│  Translator.TranslateFlow()                                    │
│  - 遍历 Steps                                                   │
│  - 根据命令类型创建对应的 Action Wrapper                         │
│  - tapOn → TapOnWrapper                                        │
│  - inputText → SendKeysWrapper                                 │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│  []Action (动作列表)                                            │
│  [TapOnWrapper{element: ByText("Login")},                       │
│   SendKeysWrapper{text: "user@example.com"}]                    │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│  Executor.Execute()                                            │
│  - 按顺序执行每个 Action                                         │
│  - 打印进度 [STEP 1/2] tapOn                                     │
│  - 失败时自动截图                                               │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│  Device.Tap() / Device.SendKeys()                               │
│  → iOS: WDA API / Android: ADB Shell                            │
└─────────────────────────────────────────────────────────────────┘
```

---

## 11. 依赖关系

```
go-uop
│
├── gopkg.in/yaml.v3 (v3.0.1)           # YAML 解析
│
└── github.com/openai/openai-go/v3 (v3.29.0)  # OpenAI API
    │
    └── github.com/tidwall/*            # 间接依赖
        ├── gjson (v1.18.0)            # JSON 解析
        ├── match (v1.1.1)             # 模式匹配
        ├── pretty (v1.2.1)             # JSON 格式化
        └── sjson (v1.2.5)              # JSON 构建
```

---

## 12. 扩展点

| 扩展点 | 接口/类型 | 说明 |
|--------|-----------|------|
| 新平台支持 | `core.Device` | 实现 Device 接口即可接入统一 API |
| 新 AI Provider | `ai.Provider` | 实现 Assert/Rerank 方法 |
| 新 Selector | `internal/selector` | 扩展 SelectorType 枚举 |
| 新 Action | `internal/action` | 实现 Action 接口 |
| 新 YAML 命令 | `yaml/commands/` | 添加新的命令处理器 |
| 新 Maestro 命令 | `maestro/commands/` | 添加 Maestro 兼容命令 |

### 12.1 实现新平台示例

```go
package myplatform

type Device struct {
    client *MyClient
}

func (d *Device) Platform() core.Platform {
    return "myplatform"
}

func (d *Device) Info() (map[string]interface{}, error) {
    return map[string]interface{}{"platform": "myplatform"}, nil
}

func (d *Device) Screenshot() ([]byte, error) {
    return d.client.CaptureScreen()
}

func (d *Device) Tap(x, y int) error {
    return d.client.Touch(x, y)
}

func (d *Device) SendKeys(text string) error {
    return d.client.Type(text)
}

func (d *Device) Launch() error {
    return d.client.StartApp()
}

func (d *Device) Close() error {
    return d.client.Disconnect()
}

var _ core.Device = (*Device)(nil)
```

---

## 13. 设计模式

| 模式 | 应用场景 |
|------|----------|
| **策略模式** | Selector 支持多种定位策略 |
| **建造者模式** | Fluent API 使用链式调用 |
| **适配器模式** | 不同平台统一通过 Device 接口 |
| **命令模式** | Action 将操作封装为对象 |
| **装饰器模式** | Retry 对操作进行重试包装 |
| **代理模式** | WDA/ADB Client 封装平台通信 |
| **工厂模式** | `ai.NewProvider()` 创建不同 AI Provider |

---

## 14. 错误处理

### 14.1 核心错误 (core/errors.go)

```go
var (
    ErrNotImplemented = errors.New("not implemented")
    ErrDeviceNotFound = errors.New("device not found")
)
```

### 14.2 Maestro 错误 (maestro/errors.go)

```go
var (
    ErrUnsupportedCommand = errors.New("unsupported command")
    ErrParsingFailed      = errors.New("parsing failed")
    ErrExecutionFailed    = errors.New("execution failed")
    ErrTranslationFailed  = errors.New("translation failed")
)
```

### 14.3 错误包装

所有错误都通过 `fmt.Errorf` 进行上下文包装：

```go
if err := ab.device.Tap(x, y); err != nil {
    return fmt.Errorf("tap: %w", err)
}
```

---

## 15. 配置选项

### 15.1 iOS 选项 (ios/option.go)

```go
type Option func(*config)

func WithAddress(addr string) Option
// 示例: ios.WithAddress("http://localhost:8100")

func WithTimeout(timeout time.Duration) Option
// 示例: ios.WithTimeout(30 * time.Second)
```

### 15.2 Android 选项 (android/option.go)

```go
type Option func(*config)

func WithSerial(serial string) Option
// 示例: android.WithSerial("emulator-5554")

func WithPackage(pkg string) Option
// 示例: android.WithPackage("com.example.app")

func WithTimeout(timeout time.Duration) Option
```

### 15.3 核心选项 (core/option.go)

```go
type DeviceOption func(*deviceConfig)

func WithTimeout(timeout time.Duration) DeviceOption
func WithSerial(serial string) DeviceOption
func WithAddress(address string) DeviceOption
```

---

## 16. 测试

### 16.1 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./ios/...
go test ./android/...
go test ./maestro/...

# 带详细输出
go test -v ./...

# 带覆盖率
go test -cover ./...

# 运行特定测试
go test -v -run TestSelector ./internal/selector/...
```

### 16.2 测试文件结构

```
uop_test.go              # 公共 API 测试
ios/device_test.go       # iOS 设备测试
android/device_test.go   # Android 设备测试
internal/selector/selector_test.go
internal/action/action_test.go
maestro/translator_test.go
maestro/config_test.go
maestro/integration_test.go
yaml/evaluator_test.go
yaml/parser_test.go
```

---

## 17. 构建和发布

### 17.1 Makefile 常用命令

```bash
make build        # 构建所有二进制文件
make test        # 运行测试
make lint        # 代码检查
make clean       # 清理构建产物
make deps        # 下载依赖
```

### 17.2 发布版本

```bash
# 标签版本
git tag v0.1.0
git push origin v0.1.0

# 构建发布包
make release
```

---

## 18. 参考链接

- [WebDriverAgent](https://github.com/appium/WebDriverAgent)
- [go-ios](https://github.com/danielpaulus/go-ios)
- [Android ADB](https://developer.android.com/studio/command-line/adb)
- [Maestro](https://docs.maestro.dev)
- [OpenAI Go SDK](https://github.com/openai/openai-go)
- [YAML v3](https://pkg.go.dev/gopkg.in/yaml.v3)
