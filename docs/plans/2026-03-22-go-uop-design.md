# go-uop 设计文档

> Go SDK Universal Operation — 移动端自动化测试基础库

**版本**: v0.1.0  
**日期**: 2026-03-22  
**状态**: 设计中

---

## 1. 概述

### 1.1 目标

`go-uop` 是一个统一的移动端自动化操作库，支持 iOS (原生 WDA) 和 Android (ADB)，提供：

- **Go API**: 链式调用的 PageObject 风格
- **YAML Runner**: Maestro 风格的脚本解析与执行
- **视觉定位**: OpenCV 模板匹配
- **AI 断言**: 调用 AI API 进行语义判断
- **并行执行**: 多设备并行测试
- **错误恢复**: 自动重试 + 现场保留

### 1.2 使用场景

| 场景 | 说明 |
|------|------|
| 内部测试工具 | QA 团队黑盒测试 |
| 回归测试 | CI/CD 集成 |
| RPA | 流程自动化 |

### 1.3 设计原则

- **完全统一接口**: `Device` 接口自动适配 iOS/Android
- **链式调用**: 所有操作支持 Fluent API
- **零外部依赖 (除必要库)**: iOS 不使用 gwda，原生实现 WDA 协议
- **YAML 优先**: 测试用例以 YAML 定义，支持参数化

---

## 2. 架构

### 2.1 分层结构

```
┌─────────────────────────────────────────────────────┐
│                    User Layer                       │
│   ┌─────────────────┐    ┌──────────────────────┐  │
│   │   Go API        │    │   YAML Runner        │  │
│   │  (链式调用)      │    │  (Maestro 风格)     │  │
│   └────────┬────────┘    └──────────┬───────────┘  │
├────────────┴─────────────────────────┴─────────────┤
│                   Command Layer                     │
│   ┌─────────────────────────────────────────────┐   │
│   │  Command Executor (tap/input/launch/assert)│   │
│   └─────────────────────────────────────────────┘   │
├──────────────────────┬──────────────────────────────┤
│     iOS Driver       │      Android Driver          │
│  (WDA HTTP Native)   │      (ADB Wrapper)           │
├──────────────────────┴──────────────────────────────┤
│                    Core Layer                        │
│   ┌────────┐ ┌────────┐ ┌────────┐ ┌────────────┐ │
│   │ Device │ │ Vision │ │ Report │ │  Parallel   │ │
│   │ Manager│ │ (OpenCV│ │  Gen   │ │  Executor   │ │
│   └────────┘ └────────┘ └────────┘ └────────────┘ │
└─────────────────────────────────────────────────────┘
```

### 2.2 模块说明

| 模块 | 职责 |
|------|------|
| `device` | 设备连接、管理、生命周期 |
| `driver/ios` | iOS WDA 协议实现 |
| `driver/android` | Android ADB 封装 |
| `locator` | 元素定位器 (ByText/ByID/ByXPath 等) |
| `action` | 操作 (Tap/Swipe/SendKeys/Launch 等) |
| `assertion` | 断言 (Text/Visibility/ScreenMatch/AISemantic) |
| `vision` | OpenCV 模板匹配 |
| `yaml` | YAML 解析、表达式、JavaScript 引擎 |
| `parallel` | 并行执行器 |
| `report` | 报告生成器 |
| `retry` | 重试机制 |

---

## 3. Go API 设计

### 3.1 设备连接

```go
// 连接 iOS 设备
dev, err := uop.NewDevice(uop.IOS, uop.WithSerial("00001234"))

// 连接 Android 设备
dev, err := uop.NewDevice(uop.Android, uop.WithSerial("emulator-5554"))

// USB 自动发现
dev, err := uop.NewDevice(uop.IOS) // 无 serial 时自动查找

// WiFi 连接
dev, err := uop.NewDevice(uop.IOS, uop.WithAddress("192.168.1.100:8100"))
```

### 3.2 链式调用

```go
// 基本操作链
dev.
    Launch("com.example.app").
    Find(uop.ByText("登录")).Tap().
    Find(uop.ByID("username")).SendKeys("user").
    Find(uop.ByID("password")).SendKeys("pass").
    Find(uop.ByText("提交")).Tap().
    Wait(uop.WithTimeout(5 * time.Second)).
    Assert().TextContains("欢迎")

// 带滚动和视觉定位
dev.
    SwipeUp().
    Find(uop.ByText("更多")).Tap().
    FindVision("settings_icon.png").Tap() // OpenCV 模板匹配

// 复杂断言
dev.Assert().
    ScreenMatch("expected.png", 0.9).
    TextVisible("登录成功").
    AISemantic("页面应该显示用户头像和用户名", apiKey)
```

### 3.3 定位器 (Selector)

定位器支持**自动识别正则表达式**，以 `/` 包裹的字符串自动识别为正则。

```go
// 文本匹配 (自动识别正则)
uop.ByText("登录")                    // 精确匹配
uop.ByText("/登.*/")                  // 正则匹配: 匹配"登录"、"登录页"等
uop.ByText("/^用户\\d+$/")            // 正则: "用户1", "用户2"...

// ID 匹配 (自动识别正则)
uop.ByID("submit")                    // 精确匹配
uop.ByID("/btn_.*_confirm/")          // 正则匹配

// XPath (原样使用，不做正则处理)
uop.ByXPath("//Button[@text='登录']")

// 类名
uop.ByClassName("android.widget.Button")

// iOS 特有
uop.ByPredicate(`name == "test"`)    // Predicate 语法
uop.ByClassChain(`**/Button[`*`]`)   // Class Chain 语法

// 组合定位
uop.ByAnd(uop.ByText("用户"), uop.ByID("user_avatar"))

// 索引定位 (多元素时的选择策略)
// 默认: 匹配多个元素时，选择屏幕位置左上角第一个
// 显式: 指定 index 选择第 N 个 (从 0 开始)
uop.ByText("取消")                    // 默认选左上角第一个
uop.ByText("取消").Index(0)           // 同上，显式指定第一个
uop.ByText("取消").Index(1)           // 选择第二个
uop.ByText("取消").Index(2)           // 选择第三个

// 多元素场景示例
// 假设屏幕上有 3 个 "确定" 按钮，分别在 (100,100), (200,200), (300,300)
// - dev.Find(uop.ByText("确定")).Tap()  // 默认点击 (100,100)
// - dev.Find(uop.ByText("确定").Index(1)).Tap()  // 点击 (200,200)
```

**多元素选择策略:**

| 策略 | 说明 | 使用方式 |
|------|------|----------|
| 默认 (左上角优先) | 匹配多个时，选择 Y 坐标最小 → X 坐标最小的元素 | `uop.ByText("确定")` |
| 显式序号 | 指定 index 选择第 N 个 (0-based) | `uop.ByText("确定").Index(2)` |

### 3.4 操作 (Action)

```go
// 触摸操作
dev.Tap(x, y)
dev.Tap(uop.ByText("确定"))         // 定位后点击
dev.DoubleTap(x, y)
dev.LongPress(x, y, 2 * time.Second)

// 滑动
dev.Swipe(x1, y1, x2, y2)
dev.SwipeUp()
dev.SwipeDown()
dev.Drag(x1, y1, x2, y2, 0.5*time.Second)

// 文本输入
dev.SendKeys("hello world")
dev.SendKeys("password", uop.WithSecure()) // 隐藏日志

// 应用管理
dev.Launch("com.example.app")
dev.Terminate("com.example.app")
dev.Install("/path/to/app.ipa")
dev.Uninstall("com.example.app")

// 系统操作
dev.PressHome()
dev.PressBack()
dev.PressButton(uop.ButtonVolumeUp)
dev.Screenshot("/path/to/save.png")
dev.RecordScreen("/path/to/video.mp4", 30*time.Second)

// 等待
dev.Wait(uop.WithTimeout(10 * time.Second))
dev.WaitForElement(uop.ByText("加载完成"))
dev.WaitForGone(uop.ByID("loading"))

// 页面源码
source := dev.GetSource()           // XML 源码
accessible := dev.GetAccessibleSource()
```

### 3.5 断言 (Assertion)

```go
// 元素断言
dev.Assert().
    TextVisible("登录成功").
    TextContains("欢迎").
    ElementVisible(uop.ByID("avatar")).
    ElementGone(uop.ByID("loading"))

// 屏幕断言
dev.Assert().
    ScreenMatch("expected.png", 0.9).
    ScreenDiff("expected.png", "diff.png", 0.1) // 容许 10% 差异

// AI 语义断言
dev.Assert().
    AISemantic(
        "页面应该显示用户头像在右上角，用户名在头像下方",
        uop.WithVisionAPI(visionAPIKey),
    )

// 视觉匹配 + 点击
matched := dev.Assert().
    VisionExists("button.png", 0.85)
if matched {
    dev.VisionTap("button.png", 0.85)
}
```

### 3.6 视觉 (Vision)

```go
// 模板匹配
pos, found := dev.Vision().Find("button.png", 0.85)
if found {
    dev.Vision().Tap(pos.X, pos.Y)
}

// 等待视觉元素出现
dev.Vision().WaitFor("loading.gif", 0.8, 10*time.Second)

// 多模板匹配
positions := dev.Vision().FindAll("item.png", 0.8)
for _, pos := range positions {
    dev.Tap(pos.X, pos.Y)
}
```

---

## 4. YAML Runner 设计

### 4.1 文件结构

```
test/
├── flows/
│   ├── login.yaml        # 登录流程
│   └── checkout.yaml     # 结账流程
├── data/
│   └── users.yaml        # 测试数据
└── smoke.yaml            # 主入口脚本
```

### 4.2 流程定义

YAML 流程文件由两部分组成：

**1. 全局定义 (可选)**

```yaml
# flows/login.yaml
name: login                              # 流程名称
description: 用户登录流程                  # 流程描述
platform: both                           # 平台: ios | android | both (默认)

# 参数定义
params:
  username: string   (required)
  password: string   (required)
  rememberMe: bool   (default: false)

# 全局超时 (可被单步覆盖)
timeout: 30s
```

**2. 命令列表 (必需)**

```yaml
# flows/login.yaml
steps:
  - launch: com.example.app
  
  - tapOn: {text: "登录", timeout: 5s}
  
  - inputText:
      text: ${params.username}
      element: {id: "username_input"}
  
  - inputText:
      text: ${params.password}
      element: {id: "password_input"}
      secure: true
  
  - tapOn: {text: "提交"}
  
  - waitFor: {element: {text: "我的"}, timeout: 15s}
  
  - assertVisible: {text: "我的"}
  
  - takeScreenshot: {name: "after_login"}
```

### 4.3 主入口脚本

```yaml
# smoke.yaml
config:
  appId: com.example.app
  timeout: 30s
  retry: 3

env:
  baseUrl: "https://staging.example.com"
  apiKey: ${ENV.API_KEY}

import:
  - flows/login.yaml
  - flows/checkout.yaml

variables:
  testUser: "test_${RANDOM(6)}@example.com"

tests:
  - name: "用户登录"
    flow: login
    params:
      username: ${variables.testUser}
      password: "Test123456"

  - name: "完整购买"
    steps:
      - runFlow: login
        params:
          username: "buyer@example.com"
          password: "Test123456"
      - tapOn: {text: "商品"}
      - runFlow: checkout
```

### 4.4 完整命令列表

#### 4.4.1 应用控制

```yaml
# 启动应用
- launch: com.example.app
- launch:
    appId: com.example.app
    arguments:
      - "-env"
      - "staging"
    waitForIdle: true

# 终止应用
- terminate: com.example.app

# 安装应用
- install: /path/to/app.ipa
- install:
    path: /path/to/app.apk
    grantPermissions: true

# 卸载应用
- uninstall: com.example.app
```

#### 4.4.2 触摸操作

```yaml
# 点击元素
- tapOn:
    text: "确定"
- tapOn:
    id: "submit_button"
- tapOn:
    xpath: "//Button[@text='OK']"
- tapOn:
    index: 2
    text: "取消"

# 点击坐标
- tap:
    x: 500
    y: 300

# 双击
- doubleTap:
    text: "编辑"

# 长按
- longPress:
    text: "删除"
    duration: 2000  # 毫秒

# 滑动
- swipe:
    startX: 500
    startY: 1000
    endX: 500
    endY: 500
    duration: 500

# 快速滑动
- swipeUp: {}
- swipeDown: {}
- swipeLeft: {}
- swipeRight: {}

# 拖拽
- drag:
    from: {text: "item1"}
    to: {text: "target"}
    duration: 1000
```

#### 4.4.3 文本输入

```yaml
# 输入文本
- inputText:
    text: "hello world"
    element:
      id: "input_field"

# 清空并输入
- clearInput:
    element:
      id: "input_field"
- inputText:
    text: "new value"
    element:
      id: "input_field"

# 密码输入 (隐藏日志)
- inputText:
    text: ${params.password}
    element:
      id: "password"
    secure: true

# 按键操作
- pressKey: ENTER
- pressKey: BACK
- pressKey: HOME
- pressKey: 66  # keycode
```

#### 4.4.4 等待与暂停

```yaml
# 等待元素出现
- waitFor:
    element:
      text: "加载完成"
    timeout: 30s
    optional: true  # 不出现也不报错

# 等待元素消失
- waitForGone:
    element:
      id: "loading"
    timeout: 10s

# 等待一段时间
- wait: 2000  # 毫秒

# 等待_activity
- waitForActivity:
    name: ".MainActivity"
    timeout: 10s

# 等待网络空闲
- waitForIdle: {}
```

#### 4.4.5 断言

```yaml
# 元素可见
- assertVisible:
    text: "登录成功"
- assertVisible:
    id: "welcome_text"

# 元素不可见
- assertNotVisible:
    text: "错误"

# 元素存在 (不一定可见)
- assertExist:
    id: "loading_spinner"

# 文本包含
- assertText:
    element:
      id: "status"
    contains: "完成"

# 表达式断言
- assertTrue: ${SCREEN_MATCH("success.png", 0.9)}
- assertTrue: ${TEXT_CONTAINS("用户名: test")}
- assertTrue: ${ELEMENT_VISIBLE("submit_button")}

# AI 语义断言
- assertTrue: ${AI_SEMANTIC("页面应该显示绿色对勾表示成功")}
- assertTrue: ${AI_SEMANTIC("用户头像应该在右上角")}
```

#### 4.4.6 视觉操作

```yaml
# 视觉点击
- visionTap:
    image: "button.png"
    threshold: 0.85
    timeout: 10s

# 等待视觉元素
- visionWait:
    image: "loading.png"
    threshold: 0.8
    timeout: 30s

# 视觉存在断言
- visionAssert:
    image: "success_icon.png"
    threshold: 0.9
    exists: true
```

#### 4.4.7 控制流

```yaml
# 循环
- foreach:
    variable: item
    in: ${LIST("item1", "item2", "item3")}
    do:
      - tapOn:
          text: ${item}
      - swipeUp: {}

# 条件执行
- if:
    condition: ${TEXT_CONTAINS("VIP用户")}
    then:
      - tapOn: {text: "VIP专属"}
    else:
      - tapOn: {text: "普通用户"}

# 条件循环
- while:
    condition: ${ELEMENT_VISIBLE("load_more")}
    maxIterations: 10
    do:
      - tapOn: {text: "加载更多"}

# 重试块
- retry:
    times: 3
    interval: 1000  # ms
    on:
      - text: "网络错误"
      - text: "加载失败"
    do:
      - swipeDown: {}
      - wait: 2000
```

#### 4.4.8 子流程

```yaml
# 运行子流程
- runFlow: login
- runFlow:
    name: login
    params:
      username: "test@example.com"
      password: "pass123"

# 嵌入流程
- include:
    - tapOn: {text: "设置"}
    - assertVisible: {text: "账户"}
```

#### 4.4.9 变量与表达式

```yaml
# 变量赋值
- setVariable:
    name: userId
    value: "user_${RANDOM(8)}"

# 提取变量
- extractVariable:
    name: balance
    from: ${TEXT("balance_value")}

# JavaScript 表达式
- eval: |
    const user = variables.testUser;
    return user.split('@')[0];

# HTTP 请求 (可选)
- http:
    method: GET
    url: "${config.baseUrl}/api/user"
    headers:
      Authorization: "Bearer ${env.token}"
    saveTo: apiResponse

# 截图
- takeScreenshot:
    name: "step_${STEP_INDEX}"

# 录屏
- startRecording:
    name: "test_recording"
- stopRecording: {}

# 输出日志
- log: "用户 ${variables.userId} 登录成功"

# 注释
- comment: "以下开始结账流程"
```

#### 4.4.10 脚本执行 (evalScript)

支持执行 JavaScript 和 Python 脚本，可访问设备上下文和变量。

```yaml
# 执行 JavaScript 脚本
- evalScript:
    lang: javascript
    script: |
      // 访问内置对象
      const { device, variables, utils } = scriptContext;
      
      // 计算随机 ID
      const randomId = 'user_' + Math.random().toString(36).substr(2, 8);
      return randomId;

# 执行 Python 脚本
- evalScript:
    lang: python
    script: |
      import random
      import string
      
      # 访问上下文
      device = script_context['device']
      variables = script_context['variables']
      
      # 计算随机 ID
      random_id = 'user_' + ''.join(random.choices(string.ascii_lowercase, k=8))
      return random_id

# 脚本赋值给变量
- evalScript:
    lang: javascript
    script: |
      const timestamp = Date.now();
      const orderId = 'ORD_' + timestamp;
      return orderId;
    saveTo: orderId

# 条件判断中使用脚本
- if:
    condition: |
        // JavaScript 判断条件
        const count = parseInt(variables.itemCount);
        const isVip = variables.userLevel === 'VIP';
        return count > 10 && isVip;
    then:
      - tapOn: {text: "批量操作"}

# 循环中使用脚本
- foreach:
    variable: index
    in: ${LIST("0", "1", "2", "3")}
    do:
      - evalScript:
          lang: javascript
          script: |
            const i = parseInt(variables.index);
            return i * i;  // 计算平方
          saveTo: square

# 复杂业务逻辑
- evalScript:
    lang: python
    script: |
      # 生成测试数据
      users = []
      for i in range(5):
          users.append({
              'name': f'user_{i}',
              'email': f'user_{i}@test.com',
              'age': 20 + i
          })
      return users
    saveTo: testUsers

# 访问设备信息
- evalScript:
    lang: javascript
    script: |
      // 获取设备信息
      const info = device.info();
      return {
          platform: info.platform,
          version: info.version,
          model: info.model
      };

# 访问截图并进行图像处理
- evalScript:
    lang: python
    script: |
      import base64
      import json
      
      # 获取当前截图
      screenshot = device.screenshot()
      
      # 可以对截图进行处理
      # 这里只是返回截图大小作为示例
      return {
          'size': len(screenshot),
          'hasContent': len(screenshot) > 1000
      };

# 调用其他脚本文件
- evalScript:
    lang: javascript
    source: scripts/helpers.js
    function: generateOrderId
    args:
      - prefix: "ORD"
      - length: 12
```

**内置上下文对象:**

| 对象 | 说明 | JavaScript 示例 | Python 示例 |
|------|------|-----------------|-------------|
| `device` | 设备操作接口 | `device.info()` | `device.info()` |
| `variables` | 当前变量 | `variables.userId` | `variables['userId']` |
| `utils` | 工具函数 | `utils.uuid()` | `utils.uuid()` |
| `console` | 日志输出 | `console.log()` | `print()` |

**device 对象可用方法:**

```go
// 设备信息
device.Info() map[string]interface{}

// 设备交互
device.Screenshot() ([]byte, error)
device.GetSource() (string, error)
device.ElementExists(locator string) (bool, error)
device.GetElementText(locator string) (string, error)
device.GetElementAttribute(locator, attr string) (string, error)

// 变量操作
variables.Set(key, value)
variables.Get(key) interface{}
variables.All() map[string]interface{}
```

**JavaScript 引擎:**

- 使用 `otto` 库 (纯 Go 实现)
- 内置 Math, Date, JSON, String, Array 等标准对象
- 可访问 scriptContext 中的 device, variables, utils

**Python 引擎:**

- 使用 `gpython` 或 `go-python` (需要 CGO)
- 或者使用 `yaegi` (纯 Go 实现，但功能受限)
- 内置 os, sys, json, re, random, datetime 等标准库

### 4.5 表达式系统

#### 4.5.1 内置函数

| 函数 | 说明 | 示例 |
|------|------|------|
| `${ENV(key)}` | 读取环境变量 | `${ENV.API_KEY}` |
| `${RANDOM(n)}` | 生成随机字符串 | `${RANDOM(6)}` |
| `${LIST(a,b,c)}` | 创建列表 | `${LIST("a","b")}` |
| `${TEXT(id)}` | 获取元素文本 | `${TEXT("username")}` |
| `${ATTR(id, name)}` | 获取元素属性 | `${ATTR("img","src")}` |
| `${SCREEN_MATCH(img, threshold)}` | 视觉匹配 | `${SCREEN_MATCH("ok.png", 0.9)}` |
| `${AI_SEMANTIC(description)}` | AI 语义判断 | `${AI_SEMANTIC("页面显示成功")}` |
| `${ELEMENT_VISIBLE(locator)}` | 元素可见性 | `${ELEMENT_VISIBLE("submit")}` |
| `${NOW()}` | 当前时间戳 | `${NOW()}` |
| `${UUID()}` | UUID 生成 | `${UUID()}` |

#### 4.5.2 JavaScript 支持

```yaml
- eval: |
    // 计算随机用户
    const random = Math.floor(Math.random() * 1000);
    return "user" + random + "@test.com";

- if:
    condition: |
        // JavaScript 条件
        const count = parseInt(variables.itemCount);
        return count > 10;
    then:
      - tapOn: {text: "批量操作"}
```

### 4.6 参数化设计

```yaml
# 参数定义
params:
  - name: username
    type: string
    required: true
  - name: password
    type: string
    required: true
    secure: true  # 日志隐藏
  - name: timeout
    type: int
    default: 30

# 参数验证
validate:
  username:
    pattern: "^\\w+@\\w+\\.\\w+$"
    message: "邮箱格式不正确"
  password:
    minLength: 6
    message: "密码至少6位"
```

### 4.7 数据驱动

```yaml
# 数据文件
data:
  file: data/users.csv
  delimiter: ","

# 内联数据
data:
  - username: user1@test.com
    password: pass123
    role: buyer
  - username: user2@test.com
    password: pass456
    role: seller

# 遍历数据
tests:
  - name: "测试用户: ${user.username}"
    flow: login
    with:
      - ${data[0]}
      - ${data[1]}
      - ${data[2]}
```

---

## 5. AI 断言设计

### 5.1 接口定义

```go
// VisionProvider 定义视觉 AI 接口
type VisionProvider interface {
    // Analyze 分析截图，返回语义描述
    Analyze(image []byte, prompt string) (string, error)
    
    // Judge 判断截图是否符合描述
    Judge(image []byte, description string) (bool, float64, error)
}

// 内置实现
type OpenAIProvider struct {
    APIKey string
    Model  string // "gpt-4o", "gpt-4o-mini"
}

```

### 5.2 使用方式

```go
// Go API
dev.Assert().AISemantic(
    "页面应该显示绿色对勾和'操作成功'文字",
    uop.WithProvider("openai", "sk-xxx"),
)

// YAML
- assertTrue: ${AI_SEMANTIC("应该显示用户头像在右上角")}
```

---

## 6. 并行执行设计

### 6.1 使用场景

并行执行主要用于 **多设备覆盖测试** 或 **加速执行**：

| 场景 | 说明 |
|------|------|
| 多设备测试 | 同一套用例同时在 3 台手机上跑，验证兼容性 |
| 多用例并发 | 多个独立测试用例并行执行，缩短总执行时间 |
| 数据驱动 | 同一用例使用不同数据并行执行 |

### 6.2 Go API

```go
// 方式 1: 加载 YAML 测试套件
suite, err := uop.LoadSuite("tests/smoke.yaml")

// 方式 2: 编程方式添加测试
suite := uop.NewSuite()
suite.AddTest("iOS 登录", uop.IOS, loginFlow)
suite.AddTest("Android 登录", uop.Android, loginFlow)

// 创建运行器
runner := uop.NewRunner(suite,
    uop.WithWorkers(4),           // 最多 4 个并发
    uop.WithRetry(2),            // 失败重试 2 次
    uop.WithReport("./reports"), // 报告输出目录
)

// 连接多台设备
devices, _ := uop.DiscoverDevices(uop.IOS)  // 查找 iOS 设备
devices, _ = uop.DiscoverDevices(uop.Android) // 查找 Android 设备

// 执行
results := runner.Run(devices...)

// 生成报告
report := uop.NewReport(results)
report.SaveHTML("report.html")
report.SaveJSON("report.json")
report.SaveJUnitXML("junit.xml")  // CI 集成
```

### 6.3 设备发现与分配

```go
// 发现所有可用设备
allDevices, _ := uop.DiscoverAll()

// 发现特定平台设备
iosDevices, _ := uop.DiscoverDevices(uop.IOS)
androidDevices, _ := uop.DiscoverDevices(uop.Android)

// 按序列号指定设备
dev, _ := uop.NewDevice(uop.IOS, uop.WithSerial("00001234"))

// WiFi 连接
dev, _ := uop.NewDevice(uop.IOS, uop.WithAddress("192.168.1.100:8100"))
```

### 6.4 报告结构

```go
// 测试结果结构
type Result struct {
    Name     string
    Device   string      // 设备标识
    Platform string      // ios / android
    Status   string      // passed / failed / skipped
    Duration time.Duration
    Steps    []StepResult
    Error    string      // 失败原因
    Screenshots []string // 截图路径
}

// 步骤结果
type StepResult struct {
    Action   string
    Params   map[string]interface{}
    Result   string      // success / failed
    Duration time.Duration
    Error    string
    Screenshot string
}
```

### 6.5 报告输出

```go
// HTML 报告 (详细，带截图)
report.SaveHTML("report.html")

// JSON 报告 (程序化使用)
report.SaveJSON("report.json")

// JUnit XML (CI/CD 集成)
report.SaveJUnitXML("junit.xml")

// 控制台输出
report.PrintSummary(os.Stdout)
```

**控制台输出示例:**
```
✅ 用户登录 (iPhone 15 Pro)       12s  PASSED
❌ 用户登录 (Pixel 7)            8s   FAILED: 元素未找到 "提交"
✅ 购买流程 (iPhone 15 Pro)       45s  PASSED
⏭️  支付测试 (Android Emulator)  --   SKIPPED

总计: 4 | 通过: 2 | 失败: 1 | 跳过: 1 | 耗时: 1m5s
```

---

## 7. 错误恢复设计

### 7.1 重试策略

```go
// 内置重试
uop.WithRetry(3),
uop.WithRetryInterval(1 * time.Second),

// 自定义重试条件
uop.WithRetryPolicy(func(err error) bool {
    return strings.Contains(err.Error(), "network")
})

// YAML 配置
config:
  retry: 3
  retryInterval: 1s
  retryOn:
    - "网络错误"
    - "加载失败"
```

### 7.2 现场保留

```go
// 失败时自动截图和录屏
uop.WithFailureCapture(true)

// 自定义失败处理
uop.WithOnFailure(func(ctx *uop.Context) {
    ctx.Screenshot("failure.png")
    ctx.RecordStop()
    ctx.Log("Test failed, artifacts saved")
})
```

---

## 8. iOS WDA 原生实现

### 8.1 协议端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/status` | GET | WDA 状态 |
| `/session` | POST | 创建会话 |
| `/session/{id}` | DELETE | 销毁会话 |
| `/wda/tap/0/{x}/{y}` | POST | 点击坐标 |
| `/wda/element/active` | GET | 获取激活元素 |
| `/wda/element/{id}/click` | POST | 点击元素 |
| `/wda/keys` | POST | 发送文本 |
| `/screenshot` | GET | 截图 |
| `/wda/source` | GET | 页面源码 |
| `/wda/app/launch` | POST | 启动应用 |
| `/wda/app/terminate/{bundleId}` | POST | 终止应用 |

### 8.2 实现要点

```go
// WDA HTTP 客户端
type WDAClient struct {
    BaseURL    string // http://localhost:8100
    SessionID  string
    HTTPClient *http.Client
}

// USB 模式 (通过 go-ios 的设备连接获取地址)
// WiFi 模式 (直接连接 IP:Port)
```

---

## 9. 技术选型

| 组件 | 选择 | 理由 |
|------|------|------|
| 模板匹配 | OpenCV (go-opencv/go-cv) | 成熟稳定的视觉库 |
| YAML 解析 | gopkg.in/yaml.v3 | 标准库，稳定 |
| JavaScript | otto | 纯 Go 实现，无外部依赖 |
| Python | yaegi | 纯 Go 实现，支持 Python 3 |
| HTTP Client | net/http | 标准库 |
| JSON | encoding/json | 标准库 |
| 图像处理 | image + golang.org/x/image | 标准库 + 扩展 |

---

## 10. 目录结构

```
go-uop/
├── go.mod
├── go.sum
├── README.md
├── docs/
│   └── plans/
│       └── 2026-03-22-go-uop-design.md
├── examples/
│   ├── basic_test.go
│   └── parallel_test.go
├── internal/
│   ├── assert/
│   │   ├── assertion.go
│   │   └── assertion_test.go    # 断言单元测试
│   ├── command/
│   │   ├── executor.go
│   │   └── executor_test.go     # 命令执行单元测试
│   ├── device/
│   │   ├── manager.go
│   │   └── manager_test.go     # 设备管理单元测试
│   ├── locator/
│   │   ├── locator.go
│   │   └── locator_test.go     # 定位器单元测试 (含正则匹配)
│   ├── parallel/
│   │   ├── executor.go
│   │   └── executor_test.go    # 并行执行单元测试
│   ├── report/
│   │   ├── generator.go
│   │   └── generator_test.go   # 报告生成单元测试
│   ├── retry/
│   │   ├── retry.go
│   │   └── retry_test.go      # 重试机制单元测试
│   ├── script/
│   │   ├── context.go
│   │   ├── javascript/
│   │   │   ├── engine.go
│   │   │   └── engine_test.go  # JS 引擎单元测试
│   │   └── python/
│   │       ├── engine.go
│   │       └── engine_test.go  # Python 引擎单元测试
│   └── vision/
│       ├── template.go
│       └── template_test.go    # 视觉匹配单元测试
├── ios/
│   ├── device.go
│   ├── device_test.go         # iOS 设备单元测试
│   ├── driver.go
│   ├── wda/
│   │   ├── client.go
│   │   ├── client_test.go     # WDA 客户端单元测试
│   │   └── protocol.go
│   └── locator.go
├── android/
│   ├── device.go
│   ├── device_test.go         # Android 设备单元测试
│   ├── driver.go
│   ├── adb/
│   │   ├── client.go
│   │   └── client_test.go     # ADB 客户端单元测试
│   └── locator.go
├── yaml/
│   ├── parser.go
│   ├── parser_test.go         # YAML 解析单元测试
│   ├── evaluator.go
│   └── commands/
│       ├── tap.go
│       ├── input.go
│       ├── assert.go
│       ├── flow.go
│       ├── control.go
│       └── script.go
├── ai/
│   ├── provider.go
│   ├── openai.go
│   └── provider_test.go       # AI Provider 单元测试
└── uop.go  # 主入口，导出统一 API
```

### 单元测试说明

每个模块都配有 `_test.go` 文件，支持：

| 测试类型 | 说明 |
|----------|------|
| Mock 测试 | 使用 mock 设备/mock 服务，避免依赖真机 |
| 表驱动测试 | Go 标准表驱动模式，易于扩展 |
| 基准测试 | 性能敏感模块 (`vision`, `locator`) 包含 `Benchmark_*` |
| 示例测试 | `Example_*` 函数验证文档正确性 |

```bash
# 运行所有测试
go test ./...

# 运行特定模块测试
go test ./internal/locator/...

# 运行带覆盖率
go test -cover ./...

# 运行基准测试
go test -bench=. ./internal/vision/
```
go-uop/
├── go.mod
├── go.sum
├── README.md
├── docs/
│   └── plans/
│       └── 2026-03-22-go-uop-design.md
├── examples/
│   ├── basic_test.go
│   ├── parallel_test.go
│   └── yaml/
│       └── smoke_test.yaml
├── internal/
│   ├── assert/
│   │   └── assertion.go
│   ├── command/
│   │   └── executor.go
│   ├── device/
│   │   └── manager.go
│   ├── locator/
│   │   └── locator.go
│   ├── parallel/
│   │   └── executor.go
│   ├── report/
│   │   └── generator.go
│   ├── retry/
│   │   └── retry.go
│   ├── script/
│   │   ├── context.go        # 脚本上下文 (device, variables, utils)
│   │   ├── javascript/
│   │   │   └── engine.go     # otto 引擎封装
│   │   └── python/
│   │       └── engine.go     # yaegi 引擎封装
│   └── vision/
│       └── template.go
├── ios/
│   ├── device.go
│   ├── driver.go
│   ├── wda/
│   │   ├── client.go
│   │   └── protocol.go
│   └── locator.go
├── android/
│   ├── device.go
│   ├── driver.go
│   ├── adb/
│   │   └── client.go
│   └── locator.go
├── yaml/
│   ├── parser.go
│   ├── evaluator.go
│   └── commands/
│       ├── tap.go
│       ├── input.go
│       ├── assert.go
│       ├── flow.go
│       ├── control.go
│       └── script.go          # evalScript 命令
├── ai/
│   ├── provider.go
│   ├── openai.go
└── uop.go  # 主入口，导出统一 API
```

---

## 11. 里程碑

| 阶段 | 内容 | 优先级 |
|------|------|--------|
| **M1** | 核心框架骨架、设备连接 | P0 |
| **M2** | iOS WDA 协议实现 | P0 |
| **M3** | Android ADB 实现 | P0 |
| **M4** | 链式 API (Action/Selector/Assertion) | P0 |
| **M5** | YAML 解析器 (基础命令) | P0 |
| **M6** | YAML 控制流 (if/for/while) | P1 |
| **M7** | OpenCV 视觉定位 | P1 |
| **M8** | AI 断言 (OpenAI) | P2 |
| **M9** | 并行执行器 | P1 |
| **M10** | 报告生成器 | P1 |
| **M11** | 错误恢复、重试机制 | P1 |
| **M12** | 文档、示例、CI | P2 |

---

## 12. 风险与决策

| 风险 | 缓解措施 |
|------|----------|
| WDA 协议兼容性 | 优先实现核心端点，后续迭代 |
| OpenCV Go 绑定稳定性 | 考虑使用 cgo 调用或寻找更稳定的绑定 |
| YAML 解析复杂度 | 使用成熟的解析库，表达式计算独立模块 |
| AI API 成本/延迟 | 支持配置 API Provider，结果可缓存 |
