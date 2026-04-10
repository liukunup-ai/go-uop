package uiautomator2

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/liukunup/go-uop/core"
	"github.com/liukunup/go-uop/pkg/uiautomator2/server"
)

type Device struct {
	client *Client
	config *Config
	pkg    string
}

func NewDevice(opts ...Option) (*Device, error) {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	if err := server.EnsureInstalled(cfg.Serial); err != nil {
		return nil, fmt.Errorf("ensure server installed: %w", err)
	}

	addr := cfg.Address
	if addr == "" && cfg.Serial != "" {
		addr = getDeviceIP(cfg.Serial)
	}
	if addr == "" {
		addr = "localhost"
	}

	client, err := NewClient(addr, cfg.HTTPTimeout())
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	return &Device{
		client: client,
		config: cfg,
		pkg:    cfg.Package,
	}, nil
}

func (d *Device) Platform() core.Platform {
	return core.Android
}

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

func (d *Device) Screenshot() ([]byte, error) {
	return d.client.Screenshot("raw")
}

func (d *Device) Close() error {
	return nil
}

func (d *Device) Tap(x, y int) error {
	return d.client.Click(x, y)
}

func (d *Device) SendKeys(text string) error {
	return d.client.SendKeys(text)
}

func (d *Device) Launch() error {
	if d.pkg == "" {
		return fmt.Errorf("package name not set")
	}
	return d.client.StartApp(d.pkg)
}

func (d *Device) GetAlertText() (string, error) {
	return "", nil
}

func (d *Device) AcceptAlert() error {
	return nil
}

func (d *Device) DismissAlert() error {
	return nil
}

func (d *Device) DeviceInfo2() (*DeviceDetail, error) {
	serial := d.config.Serial
	if serial == "" {
		serial = os.Getenv("ANDROID_SERIAL")
	}
	return getDeviceDetail(serial)
}

func (d *Device) WindowSize() (int, int, error) {
	info, err := d.client.DeviceInfo()
	if err != nil {
		return 0, 0, err
	}
	return info.DisplayWidth, info.DisplayHeight, nil
}

func (d *Device) ScreenOn() error {
	return d.client.ScreenOn()
}

func (d *Device) ScreenOff() error {
	return d.client.ScreenOff()
}

func (d *Device) PressKey(key string) error {
	return d.client.PressKey(key)
}

func (d *Device) Swipe(x1, y1, x2, y2 int, duration time.Duration) error {
	return d.client.Swipe(x1, y1, x2, y2, int(duration.Milliseconds()))
}

func (d *Device) LongClick(x, y int, duration time.Duration) error {
	return d.client.LongClick(x, y, int(duration.Milliseconds()))
}

func (d *Device) DoubleClick(x, y int) error {
	return d.client.DoubleClick(x, y)
}

func (d *Device) Drag(x1, y1, x2, y2 int, duration time.Duration) error {
	return d.client.Drag(x1, y1, x2, y2, int(duration.Milliseconds()))
}

func (d *Device) TouchDown(x, y int) error {
	return d.client.Touch("down", x, y)
}

func (d *Device) TouchMove(x, y int) error {
	return d.client.Touch("move", x, y)
}

func (d *Device) TouchUp(x, y int) error {
	return d.client.Touch("up", x, y)
}

func (d *Device) FreezeRotation() error {
	return d.client.FreezeRotation(true)
}

func (d *Device) UnfreezeRotation() error {
	return d.client.FreezeRotation(false)
}

func (d *Device) Orientation() (string, error) {
	return d.client.GetOrientation()
}

func (d *Device) SetOrientation(o string) error {
	return d.client.SetOrientation(o)
}

func (d *Device) OpenNotification() error {
	return d.client.OpenNotification()
}

func (d *Device) OpenQuickSettings() error {
	return d.client.OpenQuickSettings()
}

func (d *Device) GetClipboard() (string, error) {
	return d.client.GetClipboard()
}

func (d *Device) SetClipboard(text, label string) error {
	return d.client.SetClipboard(text, label)
}

func (d *Device) LastToast() (string, error) {
	return d.client.LastToast()
}

func (d *Device) ClearToast() error {
	return d.client.ClearToast()
}

func (d *Device) DumpHierarchy(compressed, pretty bool, maxDepth int) (string, error) {
	return d.client.DumpHierarchy(compressed, pretty, maxDepth)
}

func (d *Device) Wait(activity string, timeout time.Duration) (bool, error) {
	return d.client.Wait(activity, timeout)
}

func (d *Device) WaitForIdle(timeout time.Duration) error {
	return d.client.WaitForIdle(timeout)
}

func (d *Device) WaitForWindowUpdate(pkg string, timeout time.Duration) error {
	return d.client.WaitForWindowUpdate(pkg, timeout)
}

func (d *Device) Unlock() error {
	if err := d.client.ScreenOn(); err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)
	return d.client.Swipe(0, 0, 100, 100, 500)
}

func (d *Device) AppStart(pkg, activity string, useMonkey bool) error {
	if activity != "" {
		return d.client.StartActivity(pkg, activity)
	}
	return d.client.StartApp(pkg)
}

func (d *Device) AppStop(pkg string) error {
	return d.client.ForceStop(pkg)
}

func (d *Device) AppClear(pkg string) error {
	return d.client.Clear(pkg)
}

func (d *Device) AppStopAll(excludes []string) error {
	running, err := d.client.ListRunningApps()
	if err != nil {
		return err
	}
	for _, p := range running {
		if contains(excludes, p) {
			continue
		}
		d.client.ForceStop(p)
	}
	return nil
}

func (d *Device) AppInfo(pkg string) (*AppInfo, error) {
	return d.client.AppInfo(pkg)
}

func (d *Device) AppIcon(pkg string) ([]byte, error) {
	return d.client.AppIcon(pkg)
}

func (d *Device) AppListRunning() ([]string, error) {
	return d.client.ListRunningApps()
}

func (d *Device) AppWait(pkg string, front bool, timeout float64) (int, error) {
	return d.client.AppWait(pkg, front, timeout)
}

func (d *Device) AppInstall(apkURL string) error {
	resp, err := httpGet(apkURL)
	if err != nil {
		return err
	}
	tmpFile := os.TempDir() + "/uiautomator2_install.apk"
	if err := os.WriteFile(tmpFile, resp, 0644); err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	serial := d.config.Serial
	if serial == "" {
		serial = os.Getenv("ANDROID_SERIAL")
	}
	devicePath := "/sdcard/tmp_install.apk"
	if err := adbPush(serial, tmpFile, devicePath); err != nil {
		return err
	}
	return d.client.InstallApk(devicePath)
}

func (d *Device) Push(localPath, remotePath string, mode int) error {
	serial := d.config.Serial
	if serial == "" {
		serial = os.Getenv("ANDROID_SERIAL")
	}
	if err := adbPush(serial, localPath, remotePath); err != nil {
		return err
	}
	if mode != 0 {
		return adbChmod(serial, remotePath, mode)
	}
	return nil
}

func (d *Device) Pull(remotePath, localPath string) error {
	serial := d.config.Serial
	if serial == "" {
		serial = os.Getenv("ANDROID_SERIAL")
	}
	return adbPull(serial, remotePath, localPath)
}

func (d *Device) AppAutoGrantPermissions(pkg string) error {
	return d.client.GrantPermissions(pkg)
}

func (d *Device) OpenUrl(url string) error {
	return d.client.OpenUrl(url)
}

func (d *Device) Selector() *SelectorBuilder {
	return &SelectorBuilder{client: d.client, sel: Selector{}}
}

func (d *Device) XPath(expr string) *XPathSelector {
	return &XPathSelector{device: d, expr: expr}
}

func (d *Device) StopUiautomator() error {
	return d.client.StopUiautomator()
}

func (d *Device) SetSettings(settings map[string]interface{}) error {
	return nil
}

func (d *Device) GetSettings() map[string]interface{} {
	return map[string]interface{}{
		"wait_timeout":            20.0,
		"operation_delay":         [2]float64{0, 0},
		"operation_delay_methods": []string{"click", "swipe"},
		"max_depth":               50,
	}
}

func (d *Device) ImplicitlyWait(timeout time.Duration) {
}

func getDeviceIP(serial string) string {
	cmd := exec.Command("adb", "-s", serial, "shell", "ip", "route", "get", "1")
	output, _ := cmd.CombinedOutput()
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "src") {
			parts := strings.Fields(line)
			for i, p := range parts {
				if p == "src" && i+1 < len(parts) {
					return parts[i+1]
				}
			}
		}
	}
	return ""
}

func getDeviceDetail(serial string) (*DeviceDetail, error) {
	getprop := func(prop string) string {
		cmd := exec.Command("adb", "-s", serial, "shell", "getprop", prop)
		output, _ := cmd.CombinedOutput()
		return strings.TrimSpace(string(output))
	}

	return &DeviceDetail{
		Arch:    getprop("ro.arch"),
		Brand:   getprop("ro.product.brand"),
		Model:   getprop("ro.product.model"),
		Sdk:     0,
		Serial:  serial,
		Version: 0,
	}, nil
}

func httpGet(url string) ([]byte, error) {
	resp, err := exec.Command("curl", "-s", "-m", "60", "-L", "-o", "/dev/stdout", url).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	return resp, nil
}

func adbPush(serial, local, remote string) error {
	cmd := exec.Command("adb", "-s", serial, "push", local, remote)
	_, err := cmd.CombinedOutput()
	return err
}

func adbPull(serial, remote, local string) error {
	cmd := exec.Command("adb", "-s", serial, "pull", remote, local)
	_, err := cmd.CombinedOutput()
	return err
}

func adbChmod(serial, path string, mode int) error {
	cmd := exec.Command("adb", "-s", serial, "shell", "chmod", fmt.Sprintf("%o", mode), path)
	_, err := cmd.CombinedOutput()
	return err
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

type SelectorBuilder struct {
	client *Client
	sel    Selector
}

func (sb *SelectorBuilder) Text(text string) *SelectorBuilder {
	sb.sel.Text = text
	return sb
}

func (sb *SelectorBuilder) TextContains(text string) *SelectorBuilder {
	sb.sel.TextContains = text
	return sb
}

func (sb *SelectorBuilder) ClassName(className string) *SelectorBuilder {
	sb.sel.ClassName = className
	return sb
}

func (sb *SelectorBuilder) ResourceId(id string) *SelectorBuilder {
	sb.sel.ResourceId = id
	return sb
}

func (sb *SelectorBuilder) Clickable(clickable bool) *SelectorBuilder {
	sb.sel.Clickable = clickable
	return sb
}

func (sb *SelectorBuilder) Enabled(enabled bool) *SelectorBuilder {
	sb.sel.Enabled = enabled
	return sb
}

func (sb *SelectorBuilder) Scrollable(scrollable bool) *SelectorBuilder {
	sb.sel.Scrollable = scrollable
	return sb
}

func (sb *SelectorBuilder) PackageName(pkg string) *SelectorBuilder {
	sb.sel.PackageName = pkg
	return sb
}

func (sb *SelectorBuilder) Instance(instance int) *SelectorBuilder {
	sb.sel.Instance = instance
	return sb
}

func (sb *SelectorBuilder) First() *Element {
	sb.sel.Instance = 0
	return sb.Do()
}

func (sb *SelectorBuilder) Index(i int) *Element {
	sb.sel.Instance = i
	return sb.Do()
}

func (sb *SelectorBuilder) Do() *Element {
	elements, err := sb.client.Selector(sb.sel)
	if err != nil || len(elements) == 0 {
		return &Element{}
	}
	info, _ := sb.client.ElementInfo(elements[0])
	return &Element{
		client: sb.client,
		sel:    sb.sel,
		info:   info,
	}
}

func (sb *SelectorBuilder) Count() int {
	elements, err := sb.client.Selector(sb.sel)
	if err != nil {
		return 0
	}
	return len(elements)
}

type Element struct {
	client *Client
	sel    Selector
	info   *ElementInfo
	bounds Bounds
}

func (e *Element) Exists() bool {
	return e.info != nil
}

func (e *Element) Wait(timeout time.Duration) *Element {
	return e
}

func (e *Element) WaitGone(timeout time.Duration) error {
	return nil
}

func (e *Element) Info() *ElementInfo {
	return e.info
}

func (e *Element) Text() string {
	if e.info != nil {
		return e.info.Text
	}
	return ""
}

func (e *Element) SetText(text string) error {
	return e.client.SetText(e.sel, text)
}

func (e *Element) ClearText() error {
	return e.client.ClearText()
}

func (e *Element) Center() (int, int) {
	if e.info == nil {
		return 0, 0
	}
	b := e.info.Bounds
	return (b.Left + b.Right) / 2, (b.Top + b.Bottom) / 2
}

func (e *Element) Click() error {
	x, y := e.Center()
	return e.client.Click(x, y)
}

func (e *Element) LongClick() error {
	x, y := e.Center()
	return e.client.LongClick(x, y, 500)
}

func (e *Element) DragTo(x, y int, duration time.Duration) error {
	sx, sy := e.Center()
	return e.client.Drag(sx, sy, x, y, int(duration.Milliseconds()))
}

func (e *Element) Swipe(direction string, steps int) error {
	return nil
}

func (e *Element) PinchIn(percent, steps int) error {
	return e.client.PinchIn(percent, steps)
}

func (e *Element) PinchOut(percent, steps int) error {
	return e.client.PinchOut(percent, steps)
}

func (e *Element) Screenshot() ([]byte, error) {
	return nil, nil
}

type XPathSelector struct {
	device *Device
	expr   string
}

func (xs *XPathSelector) Click() error {
	return nil
}

func (xs *XPathSelector) Get() (*Element, error) {
	return nil, nil
}

func (xs *XPathSelector) GetText() (string, error) {
	return "", nil
}

func (xs *XPathSelector) SetText(text string) error {
	return nil
}

func (xs *XPathSelector) Wait() *Element {
	return nil
}

func (xs *XPathSelector) WaitGone() {
}

var _ core.Device = (*Device)(nil)
