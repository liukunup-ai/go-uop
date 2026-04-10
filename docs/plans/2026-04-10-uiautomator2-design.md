# uiautomator2 Go Implementation Design

## Overview

This document describes the design for implementing a Go-language version of [uiautomator2](https://github.com/openatx/uiautomator2) - a Python library for Android UI automation. The Go implementation will be located at `pkg/uiautomator2` and provide the same functionality as the Python original.

## Architecture

```
pkg/uiautomator2/
├── client.go           # HTTP client + JSON-RPC protocol
├── device.go           # Device implementation (core.Device interface)
├── option.go           # Configuration options
├── types.go            # Type definitions
├── server/
│   └── installer.go    # Auto-install atx-agent/uiautomator-server
├── element/
│   ├── selector.go     # UiSelector implementation
│   └── element.go      # UIElement representation
├── xpath/
│   └── xpath.go        # XPath parsing and matching
├── app/
│   └── manager.go      # App lifecycle management
├── input/
│   ├── clipboard.go    # Clipboard operations
│   ├── keys.go         # Key events
│   └── gesture.go      # Gesture operations (touch, swipe, drag)
├── screen/
│   ├── screenshot.go   # Screenshot
│   ├── orientation.go  # Screen orientation
│   └── toast.go        # Toast handling
├── watch/
│   └── watcher.go      # WatchContext popup monitoring
└── jsonrpc/
    └── rpc.go          # JSON-RPC protocol wrapper
```

## Core Design Decisions

### 1. JSON-RPC Protocol

uiautomator2 uses JSON-RPC over HTTP to communicate with a service running on the Android device (default port 7912). All calls follow the JSON-RPC 2.0 spec:

```json
// Request
{"jsonrpc": "2.0", "id": "R10001", "method": "deviceInfo", "params": {}}

// Response
{"jsonrpc": "2.0", "id": "R10001", "result": {...}}
```

### 2. Device-side Service Management

The implementation will automatically:
- Detect if uiautomator2 services are installed on the device
- Install/update atx-agent APK if needed
- Install/update uiautomator-server APK if needed  
- Start the services automatically

APKs can be bundled or downloaded from GitHub releases.

### 3. Element Location

Dual-mode element location:
- **UiSelector Mode**: Use `dump_hierarchy()` to get XML, match elements by properties (text, resourceId, className, etc.)
- **XPath Mode**: Parse XML hierarchy, apply XPath expressions

Both modes use the same underlying XML hierarchy.

### 4. Multi-Device Support

- Device identified by serial number
- Environment variable `ANDROID_SERIAL` support
- WiFi connection by IP:port

## Complete API Specification

### Connection

```go
// Connect by serial number
d := uiautomator2.Connect("serial")  // connect_usb
d := uiautomator2.ConnectWiFi("ip:port")

// Environment variable ANDROID_SERIAL
d := uiautomator2.Connect()  // uses ANDROID_SERIAL
```

### Device Info

```go
d.Info()          // deviceInfo() - current package, display, screen status
d.DeviceInfo()    // detailed device info - arch, brand, model, sdk, version
d.WindowSize()    // (width, height)
d.AppCurrent()    // {package, activity, pid}
d.Serial()        // device serial
d.WlanIP()        // WLAN IP or nil
d.WaitActivity(act, timeout)  // wait for activity
```

### Key Events

```go
d.ScreenOn()
d.ScreenOff()
d.PressKey("home")  // home, back, left, right, up, down, center, menu, search, enter, delete/del, recent, volume_up, volume_down, volume_mute, camera, power
d.Unlock()
```

### Gestures

```go
d.Tap(x, y)
d.TapPct(xPct, yPct)  // percentage 0.0-1.0
d.DoubleClick(x, y)
d.LongClick(x, y, duration)  // duration in seconds, default 0.5
d.Swipe(sx, sy, ex, ey, duration)
d.SwipeExt(dir, scale, box)  // dir: "left"|"right"|"up"|"down"
d.Drag(sx, sy, ex, ey, duration)
d.SwipePoints([]Point, duration)
d.TouchDown(x, y)
d.TouchMove(x, y)
d.TouchUp(x, y)
```

### Screen Operations

```go
d.Screenshot(filename string, format string)  // format: "pillow", "opencv", "raw"
d.DumpHierarchy(compressed, pretty, maxDepth)
d.Orientation()           // returns "natural"|"left"|"right"|"upsidedown"
d.SetOrientation(o)
d.FreezeRotation()
d.FreezeRotation(false)  // unfreeze
d.OpenNotification()
d.OpenQuickSettings()
```

### Selector

```go
d(text="Settings")
d(className="android.widget.TextView")
d(resourceId="com.example:id/button")
d(textContains="Set")
d(textMatches("^Set.*")
d(textStartsWith="Set")
d(description="Submit")
d(checkable=true)
d(checked=false)
d(clickable=true)
d(scrollable=true)
d(enabled=true)
d(packageName="com.example")

// Chained
d(className="android.widget.ListView").child(text="Bluetooth")
d(text="Google").sibling(className="android.widget.ImageView")

// Relative positioning
d(A).LeftOf(B)
d(A).RightOf(B)
d(A).UpOf(B)
d(A).DownOf(B)

// Instance
d(text="Add")[0]
d(text="Add").Count

// Child by
d(className="android.widget.ListView").ChildByText("Bluetooth", className="android.widget.LinearLayout", allowScrollSearch=true)
d(className="android.widget.ListView").ChildByDescription(...)
d(className="android.widget.ListView").ChildByInstance(...)
```

### Element Operations

```go
el := d(text="Settings")
el.Exists(timeout)       // wait for element
el.Wait(timeout)         // wait for element, returns *Element or nil
el.WaitGone(timeout)     // wait for element to disappear
el.Info()                // element info struct

el.Text()                // get text
el.SetText("hello")      // set text
el.ClearText()           // clear text

x, y := el.Center()      // center coordinates
x, y := el.Center(offsetX, offsetY)

el.Screenshot()           // screenshot of element

el.Click()
el.Click(timeout)
el.Click(offsetX, offsetY)
el.ClickExists(timeout)  // click if exists, return bool
el.ClickGone(maxretry, interval)

el.LongClick()
el.DragTo(x, y, duration)
el.DragTo(el2, duration)

el.Swipe("left"|"right"|"up"|"down", steps)
el.Gesture((sx1, sy1), (sx2, sy2), (ex1, ey1), (ex2, ey2))
el.PinchIn(percent, steps)
el.PinchOut()

el.Fling()                    // default: vert forward
el.Fling.Horiz.Forward()
el.Fling.Vert.Backward()
el.Fling.Horiz.ToBeginning(maxSwipes)
el.Fling.ToEnd()

el.Scroll()
el.Scroll.Forward()
el.Scroll.Up()
el.Scroll.Horiz.ToBeginning()
el.Scroll.ToEnd()
el.Scroll.To(text="Security")  // scroll until element appears
```

### XPath

```go
d.XPath("//*[@text='Settings']")
d.XPath("@personal-fm")  // syntax sugar for resource-id

sl := d.Xpath("//*[@text='Submit']")
sl.Click()
sl.Click(timeout)
sl.ClickExists()
sl.Get()       // wait and get element
sl.GetText()
sl.SetText("hello")
sl.Wait()
sl.WaitGone()
```

### Input

```go
d.SendKeys("hello world")
d.SendKeys("hello", clear=true)
d.ClearText()
d.SendAction("search"|"go"|"send"|"next"|"done"|"previous")
d.HideKeyboard()
d.CurrentIme()
```

### Clipboard

```go
d.SetClipboard("hello")
d.SetClipboard("hello", label)
d.Clipboard()  // get clipboard
```

### Toast

```go
d.LastToast()   // returns string or nil
d.ClearToast()
```

### App Management

```go
d.AppInstall("http://example.com/app.apk")
d.AppStart("com.example.app")
d.AppStart("com.example.app", ".MainActivity")
d.AppStart("com.example.app", use_monkey=true)
d.AppStop("com.example.app")
d.AppClear("com.example.app")
d.AppStopAll()
d.AppStopAll(excludes=["com.examples.demo"])

d.AppInfo("com.example.app")
d.AppIcon("com.example.app")  // returns image bytes

d.AppListRunning()  // []string of package names

d.AppWait("com.example.app")           // returns pid
d.AppWait("com.example.app", front=true)
d.AppWait("com.example.app", timeout=20)

d.Push("local/file.txt", "/sdcard/")
d.Push("local/file.txt", "/sdcard/renamed.txt")
d.Push(fileObj, "/sdcard/")
d.Push("foo.sh", "/data/local/tmp/", mode=0o755)

d.Pull("/sdcard/tmp.txt", "local.txt")

d.AppAutoGrantPermissions("io.appium.android.apis")
d.OpenUrl("appname://appnamehost")
```

### Session

```go
sess := d.Session("com.netease.cloudmusic")
sess.Close()
sess.Restart()
sess := d.Session("com.netease.cloudmusic", attach=true)
sess.Running()  // bool

// with context manager
with d.Session("com.netease.cloudmusic") as sess:
    sess(text="Play").click()
```

### WatchContext

```go
ctx := d.WatchContext()
ctx.When("立即下载").When("取消").Click()
ctx.When("同意").Click()
ctx.When("确定").Click()
ctx.When("仲夏之夜").Call(func(d, el) { d.Press("back") })
ctx.WaitStability()
ctx.Close()

with d.WatchContext(builtin=True) as ctx:
    ctx.When("@tb:id/jview_view").When("//*[@content-desc='图片']").Click()
```

### Settings

```go
d.ImplicitlyWait(10.0)
d.Settings["wait_timeout"] = 20.0
d.Settings["operation_delay"] = (0.5, 1.0)
d.Settings["operation_delay_methods"] = ["click", "swipe", "drag", "press"]
d.Settings["max_depth"] = 50
```

### Other

```go
d.StopUiautomator()
d.Debug = true
```

## Implementation Phases

### Phase 1: Core Infrastructure
- `client.go` - HTTP client + JSON-RPC
- `jsonrpc/rpc.go` - JSON-RPC types and helpers
- `device.go` - Basic device structure
- `option.go` - Configuration

### Phase 2: Server Installation
- `server/installer.go` - APK install/update/launch

### Phase 3: Basic Operations
- Screen: screenshot, dump_hierarchy
- Input: tap, swipe, keys, press
- Device info: info, device_info

### Phase 4: Element & Selector
- `element/element.go` - UIElement
- `element/selector.go` - UiSelector builder
- XPath: `xpath/xpath.go`

### Phase 5: Advanced Gestures
- Double click, long click
- Drag, swipe_points
- Touch down/move/up
- Pinch, gesture

### Phase 6: App Management
- `app/manager.go` - Full app lifecycle
- Push/pull file transfer

### Phase 7: Extended Features
- Clipboard, toast, IME
- Orientation, freeze rotation
- WatchContext

### Phase 8: Settings & Debug
- Global settings
- Debug mode

## References

- [uiautomator2 Python](https://github.com/openatx/uiautomator2)
- [Android UiAutomator](https://developer.android.com/training/testing/ui-automator)
- [JSON-RPC 2.0 Spec](https://www.jsonrpc.org/specification)
