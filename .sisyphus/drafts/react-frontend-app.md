# Draft: React 前端 + Go 后端桌面应用

## 项目背景

**现有代码库**:
- Go SDK for mobile automation (github.com/liukunup/go-uop)
- 支持 iOS (WebDriverAgent) 和 Android (ADB)
- 已有 `core.Device` 接口定义
- 已有 `ios/device.go` 和 `android/device.go` 实现
- 已有 `cmd/maestro/main.go` CLI 工具

**目标**: 构建一个带有 React 前端的桌面应用，可以编译成单个二进制可执行文件。

---

## 需求 (待确认)

### 技术栈选择

**推荐方案**: Wails (Go backend + React frontend)

**备选方案**:
- Tauri (Rust backend - 不推荐，因为这是 Go 项目)
- Electron (需要额外运行时)
- Next.js + 独立 Go server (两个进程)

### 功能需求 (待确认)

1. **设备选择**
   - 列出可用的 iOS/Android 设备
   - 支持按平台筛选
   - 显示设备信息 (serial, model, status)

2. **设备查看**
   - 实时屏幕截图
   - 查看 UI 层级结构 (source tree)

3. **命令调试**
   - 点击坐标
   - 发送文本输入
   - 启动/关闭应用
   - 按键事件
   - 滑动操作

4. **其他**
   - 命令历史记录
   - 日志输出

---

## 待确认问题

1. 是否需要支持同时连接多个设备?
2. 命令调试的具体用例是什么? (自动化测试? 手动调试?)
3. 是否需要保存/加载命令脚本?
4. UI 风格偏好? (简洁/功能丰富)
5. 是否需要 OAuth 登录或用户管理?
