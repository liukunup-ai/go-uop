# go-uop

Go SDK for unified mobile automation supporting iOS (via WebDriverAgent) and Android (via ADB).

## Features

- **Unified Device Interface**: Single API for both iOS and Android
- **Chainable Fluent API**: Build actions with readable chain calls
- **YAML Test Runner**: Define test flows in YAML with Maestro-style commands
- **Vision Module**: Template matching for image-based automation
- **AI Integration**: OpenAI provider for intelligent automation
- **Parallel Execution**: Run tests across multiple devices
- **Retry Mechanisms**: Configurable retry with exponential backoff
- **Serial Communication**: Event-driven serial port with command table management

## Installation

```bash
go get github.com/liukunup/go-uop
```

## Quick Start

### iOS

```go
device, err := ios.NewDevice("com.example.app",
    ios.WithAddress("http://localhost:8100"))
if err != nil {
    log.Fatal(err)
}
defer device.Close()

// Take screenshot
screenshot, err := device.Screenshot()
```

### Android

```go
device, err := android.NewDevice(
    android.WithSerial("emulator-5554"),
    android.WithPackage("com.example.app"))
if err != nil {
    log.Fatal(err)
}
defer device.Close()

// Tap at coordinates
err = device.Tap(100, 200)
```

### Serial

```go
// Load command table from YAML
ct := serial.NewCommandTable()
ct.LoadFile("commands.yaml")

// Open serial port with command table
s, err := serial.NewSerial(&serial.Config{
    Name:    "/dev/ttyUSB0",
    Baud:    115200,
    Commands: ct,
})

// Send command by ID (async with callback)
s.SendByID("reset", func(result *serial.SendResult) {
    if result.Success {
        fmt.Println("Command executed")
    }
})

// Or send by name
s.SendByName("设备复位", callback)
```

### Fluent API

```go
err := uop.NewActionBuilder(device).
    Tap(100, 200).
    SendKeys("hello").
    Launch("com.example.app").
    Do()
```

## Architecture

```
User Layer (Go API + YAML)
    ↓
Command Layer
    ↓
Platform Drivers (iOS/Android)
    ↓
Core Modules (locator, action, vision, retry)
```

## Packages

| Package | Description |
|---------|-------------|
| `core` | Core types and interfaces |
| `ios` | iOS device implementation |
| `ios/wda` | WebDriverAgent client |
| `android` | Android device implementation |
| `android/adb` | ADB client |
| `serial` | Serial port with event-driven monitoring |
| `internal/locator` | Element locator types |
| `internal/action` | Action types |
| `internal/vision` | Vision/image processing |
| `internal/retry` | Retry utilities |
| `internal/parallel` | Parallel execution |
| `internal/report` | Test report generation |
| `yaml` | YAML parsing and command types |
| `yaml/commands` | Command execution |
| `ai` | AI provider integration |

## YAML Commands

```yaml
name: login flow
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

## License

MIT

## Links

- https://github.com/openai/openai-go
- https://github.com/ElectricBubble/gwda
- https://github.com/danielpaulus/go-ios
- https://github.com/appium/WebDriverAgent
- https://playwright.cn/docs/intro
- https://docs.maestro.dev
