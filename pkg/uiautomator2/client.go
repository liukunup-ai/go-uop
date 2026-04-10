package uiautomator2

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/liukunup/go-uop/pkg/uiautomator2/jsonrpc"
)

type Client struct {
	rpc         *jsonrpc.Client
	baseURL     string
	httpTimeout time.Duration
}

func NewClient(addr string, timeout time.Duration) (*Client, error) {
	baseURL := fmt.Sprintf("http://%s:7912/jsonrpc/0", addr)

	rpcClient, err := jsonrpc.NewClient(baseURL, timeout)
	if err != nil {
		return nil, fmt.Errorf("create RPC client: %w", err)
	}

	return &Client{
		rpc:         rpcClient,
		baseURL:     baseURL,
		httpTimeout: timeout,
	}, nil
}

func (c *Client) DeviceInfo() (*DeviceInfo, error) {
	result, err := c.rpc.Call("deviceInfo", nil)
	if err != nil {
		return nil, err
	}

	var info DeviceInfo
	if err := json.Unmarshal(result, &info); err != nil {
		return nil, fmt.Errorf("unmarshal device info: %w", err)
	}
	return &info, nil
}

func (c *Client) Screenshot(format string) ([]byte, error) {
	result, err := c.rpc.Call("screenshot", map[string]string{"format": format})
	if err != nil {
		return nil, err
	}

	var value string
	if err := json.Unmarshal(result, &value); err != nil {
		return nil, fmt.Errorf("unmarshal screenshot: %w", err)
	}

	return []byte(value), nil
}

func (c *Client) DumpHierarchy(compressed, pretty bool, maxDepth int) (string, error) {
	params := map[string]interface{}{
		"compressed": compressed,
		"pretty":     pretty,
		"max_depth":  maxDepth,
	}
	result, err := c.rpc.Call("dumpWindowHierarchy", params)
	if err != nil {
		return "", err
	}

	var xml string
	if err := json.Unmarshal(result, &xml); err != nil {
		return "", fmt.Errorf("unmarshal hierarchy: %w", err)
	}
	return xml, nil
}

func (c *Client) Selector(sel Selector) ([]string, error) {
	selMap := make(map[string]interface{})
	if sel.Text != "" {
		selMap["text"] = sel.Text
	}
	if sel.TextContains != "" {
		selMap["textContains"] = sel.TextContains
	}
	if sel.TextMatches != "" {
		selMap["textMatches"] = sel.TextMatches
	}
	if sel.TextStartsWith != "" {
		selMap["textStartsWith"] = sel.TextStartsWith
	}
	if sel.ClassName != "" {
		selMap["className"] = sel.ClassName
	}
	if sel.ClassNameMatches != "" {
		selMap["classNameMatches"] = sel.ClassNameMatches
	}
	if sel.Description != "" {
		selMap["description"] = sel.Description
	}
	if sel.DescriptionContains != "" {
		selMap["descriptionContains"] = sel.DescriptionContains
	}
	if sel.DescriptionMatches != "" {
		selMap["descriptionMatches"] = sel.DescriptionMatches
	}
	if sel.DescriptionStartsWith != "" {
		selMap["descriptionStartsWith"] = sel.DescriptionStartsWith
	}
	if sel.ResourceId != "" {
		selMap["resourceId"] = sel.ResourceId
	}
	if sel.ResourceIdMatches != "" {
		selMap["resourceIdMatches"] = sel.ResourceIdMatches
	}
	if sel.PackageName != "" {
		selMap["packageName"] = sel.PackageName
	}
	if sel.PackageNameMatches != "" {
		selMap["packageNameMatches"] = sel.PackageNameMatches
	}
	if sel.Checkable {
		selMap["checkable"] = sel.Checkable
	}
	if sel.Checked {
		selMap["checked"] = sel.Checked
	}
	if sel.Clickable {
		selMap["clickable"] = sel.Clickable
	}
	if sel.LongClickable {
		selMap["longClickable"] = sel.LongClickable
	}
	if sel.Scrollable {
		selMap["scrollable"] = sel.Scrollable
	}
	if sel.Enabled {
		selMap["enabled"] = sel.Enabled
	}
	if sel.Focusable {
		selMap["focusable"] = sel.Focusable
	}
	if sel.Focused {
		selMap["focused"] = sel.Focused
	}
	if sel.Selected {
		selMap["selected"] = sel.Selected
	}
	if sel.Index > 0 {
		selMap["index"] = sel.Index
	}
	if sel.Instance > 0 {
		selMap["instance"] = sel.Instance
	}

	result, err := c.rpc.Call("selector", selMap)
	if err != nil {
		return nil, err
	}

	var elements []string
	if err := json.Unmarshal(result, &elements); err != nil {
		return nil, fmt.Errorf("unmarshal elements: %w", err)
	}
	return elements, nil
}

func (c *Client) Touch(action string, x, y int) error {
	_, err := c.rpc.Call("touch", map[string]interface{}{
		"action": action,
		"x":      x,
		"y":      y,
	})
	return err
}

func (c *Client) Click(x, y int) error {
	_, err := c.rpc.Call("click", map[string]interface{}{
		"x": x,
		"y": y,
	})
	return err
}

func (c *Client) LongClick(x, y, duration int) error {
	_, err := c.rpc.Call("longClick", map[string]interface{}{
		"x": x,
		"y": y,
	})
	return err
}

func (c *Client) DoubleClick(x, y int) error {
	_, err := c.rpc.Call("doubleClick", map[string]interface{}{
		"x": x,
		"y": y,
	})
	return err
}

func (c *Client) Swipe(sx, sy, ex, ey, duration int) error {
	_, err := c.rpc.Call("swipe", map[string]interface{}{
		"startX":   sx,
		"startY":   sy,
		"endX":     ex,
		"endY":     ey,
		"duration": duration,
	})
	return err
}

func (c *Client) Drag(sx, sy, ex, ey, duration int) error {
	_, err := c.rpc.Call("drag", map[string]interface{}{
		"startX":   sx,
		"startY":   sy,
		"endX":     ex,
		"endY":     ey,
		"duration": duration,
	})
	return err
}

func (c *Client) SendKeys(text string) error {
	_, err := c.rpc.Call("sendKeys", map[string]interface{}{
		"text": text,
	})
	return err
}

func (c *Client) PressKey(key string) error {
	_, err := c.rpc.Call("pressKey", map[string]interface{}{
		"key": key,
	})
	return err
}

func (c *Client) PressKeyCode(code, meta int) error {
	_, err := c.rpc.Call("pressKeyCode", map[string]interface{}{
		"keycode": code,
		"meta":    meta,
	})
	return err
}

func (c *Client) ScreenOn() error {
	_, err := c.rpc.Call("screenOn", nil)
	return err
}

func (c *Client) ScreenOff() error {
	_, err := c.rpc.Call("screenOff", nil)
	return err
}

func (c *Client) FreezeRotation(freeze bool) error {
	_, err := c.rpc.Call("freezeRotation", freeze)
	return err
}

func (c *Client) GetOrientation() (string, error) {
	result, err := c.rpc.Call("getOrientation", nil)
	if err != nil {
		return "", err
	}
	var orientation string
	if err := json.Unmarshal(result, &orientation); err != nil {
		return "", err
	}
	return orientation, nil
}

func (c *Client) SetOrientation(orientation string) error {
	_, err := c.rpc.Call("setOrientation", orientation)
	return err
}

func (c *Client) SetText(sel Selector, text string) error {
	_, err := c.rpc.Call("setText", map[string]interface{}{
		"text": text,
	})
	return err
}

func (c *Client) ClearText() error {
	_, err := c.rpc.Call("clearText", nil)
	return err
}

func (c *Client) GetClipboard() (string, error) {
	result, err := c.rpc.Call("getClipboard", nil)
	if err != nil {
		return "", err
	}
	var text string
	if err := json.Unmarshal(result, &text); err != nil {
		return "", err
	}
	return text, nil
}

func (c *Client) SetClipboard(text, label string) error {
	_, err := c.rpc.Call("setClipboard", map[string]interface{}{
		"text":  text,
		"label": label,
	})
	return err
}

func (c *Client) GetOrientation2() (string, error) {
	return c.GetOrientation()
}

func (c *Client) GetRotation() (int, error) {
	result, err := c.rpc.Call("getRotation", nil)
	if err != nil {
		return 0, err
	}
	var rotation int
	if err := json.Unmarshal(result, &rotation); err != nil {
		return 0, err
	}
	return rotation, nil
}

func (c *Client) GetDisplaySizeDp() (int, int, error) {
	result, err := c.rpc.Call("getDisplaySizeDp", nil)
	if err != nil {
		return 0, 0, err
	}
	var dp struct {
		X int `json:"x"`
		Y int `json:"y"`
	}
	if err := json.Unmarshal(result, &dp); err != nil {
		return 0, 0, err
	}
	return dp.X, dp.Y, nil
}

func (c *Client) GetDisplayWidthHeight() (int, int, error) {
	info, err := c.DeviceInfo()
	if err != nil {
		return 0, 0, err
	}
	return info.DisplayWidth, info.DisplayHeight, nil
}

func (c *Client) WaitForIdle(timeout time.Duration) error {
	_, err := c.rpc.Call("waitForIdle", map[string]interface{}{
		"timeout": int(timeout / time.Millisecond),
	})
	return err
}

func (c *Client) WaitForWindowUpdate(pkg string, timeout time.Duration) error {
	_, err := c.rpc.Call("waitForWindowUpdate", map[string]interface{}{
		"packageName": pkg,
		"timeout":     int(timeout / time.Millisecond),
	})
	return err
}

func (c *Client) OpenNotification() error {
	_, err := c.rpc.Call("openNotification", nil)
	return err
}

func (c *Client) OpenQuickSettings() error {
	_, err := c.rpc.Call("openQuickSettings", nil)
	return err
}

func (c *Client) LastToast() (string, error) {
	result, err := c.rpc.Call("getLastToast", map[string]interface{}{
		"timeout": 10000,
	})
	if err != nil {
		return "", err
	}
	var toast string
	if err := json.Unmarshal(result, &toast); err != nil {
		return "", err
	}
	return toast, nil
}

func (c *Client) ClearToast() error {
	_, err := c.rpc.Call("clearToast", nil)
	return err
}

func (c *Client) RegisterWatcher(name string, cond Selector) error {
	_, err := c.rpc.Call("registerWatcher", map[string]interface{}{
		"name": name,
		"cond": cond,
	})
	return err
}

func (c *Client) UnregisterWatcher(name string) error {
	_, err := c.rpc.Call("unregisterWatcher", map[string]interface{}{
		"name": name,
	})
	return err
}

func (c *Client) UnregisterAllWatchers() error {
	_, err := c.rpc.Call("unregisterAllWatchers", nil)
	return err
}

func (c *Client) RunWatchers() error {
	_, err := c.rpc.Call("runWatchers", nil)
	return err
}

func (c *Client) HasWatcher(name string) (bool, error) {
	result, err := c.rpc.Call("hasWatcher", name)
	if err != nil {
		return false, err
	}
	var exists bool
	if err := json.Unmarshal(result, &exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (c *Client) GetWatcherStats() (map[string]int, error) {
	result, err := c.rpc.Call("getWatcherStats", nil)
	if err != nil {
		return nil, err
	}
	var stats map[string]int
	if err := json.Unmarshal(result, &stats); err != nil {
		return nil, err
	}
	return stats, nil
}

func (c *Client) StartActivity(pkg, activity string) error {
	_, err := c.rpc.Call("startActivity", map[string]interface{}{
		"pkg":      pkg,
		"activity": activity,
	})
	return err
}

func (c *Client) StartApp(pkg string) error {
	_, err := c.rpc.Call("startApp", map[string]interface{}{
		"pkg": pkg,
	})
	return err
}

func (c *Client) ForceStop(pkg string) error {
	_, err := c.rpc.Call("forceStop", map[string]interface{}{
		"pkg": pkg,
	})
	return err
}

func (c *Client) Clear(pkg string) error {
	_, err := c.rpc.Call("clear", map[string]interface{}{
		"pkg": pkg,
	})
	return err
}

func (c *Client) InstallApk(path string) error {
	_, err := c.rpc.Call("installApk", map[string]interface{}{
		"path": path,
	})
	return err
}

func (c *Client) ListRunningApps() ([]string, error) {
	result, err := c.rpc.Call("listRunningApps", nil)
	if err != nil {
		return nil, err
	}
	var apps []string
	if err := json.Unmarshal(result, &apps); err != nil {
		return nil, err
	}
	return apps, nil
}

func (c *Client) AppInfo(pkg string) (*AppInfo, error) {
	result, err := c.rpc.Call("appInfo", map[string]interface{}{
		"pkg": pkg,
	})
	if err != nil {
		return nil, err
	}
	var info AppInfo
	if err := json.Unmarshal(result, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func (c *Client) AppIcon(pkg string) ([]byte, error) {
	result, err := c.rpc.Call("appIcon", map[string]interface{}{
		"pkg": pkg,
	})
	if err != nil {
		return nil, err
	}
	var icon string
	if err := json.Unmarshal(result, &icon); err != nil {
		return nil, err
	}
	return []byte(icon), nil
}

func (c *Client) AppWait(pkg string, front bool, timeout float64) (int, error) {
	result, err := c.rpc.Call("appWait", map[string]interface{}{
		"pkg":     pkg,
		"front":   front,
		"timeout": timeout,
	})
	if err != nil {
		return 0, err
	}
	var pid int
	if err := json.Unmarshal(result, &pid); err != nil {
		return 0, err
	}
	return pid, nil
}

func (c *Client) GrantPermissions(pkg string) error {
	_, err := c.rpc.Call("grantPermissions", map[string]interface{}{
		"pkg": pkg,
	})
	return err
}

func (c *Client) OpenUrl(url string) error {
	_, err := c.rpc.Call("openUrl", map[string]interface{}{
		"url": url,
	})
	return err
}

func (c *Client) PushFile(localPath, remotePath string) error {
	_, err := c.rpc.Call("pushFile", map[string]interface{}{
		"localPath":  localPath,
		"remotePath": remotePath,
	})
	return err
}

func (c *Client) PullFile(remotePath, localPath string) error {
	_, err := c.rpc.Call("pullFile", map[string]interface{}{
		"remotePath": remotePath,
		"localPath":  localPath,
	})
	return err
}

func (c *Client) PinchIn(percent, steps int) error {
	_, err := c.rpc.Call("pinchIn", map[string]interface{}{
		"percent": percent,
		"steps":   steps,
	})
	return err
}

func (c *Client) PinchOut(percent, steps int) error {
	_, err := c.rpc.Call("pinchOut", map[string]interface{}{
		"percent": percent,
		"steps":   steps,
	})
	return err
}

func (c *Client) SwipeExt(dir string, scale, box float64) error {
	_, err := c.rpc.Call("swipeExt", map[string]interface{}{
		"dir":   dir,
		"scale": scale,
		"box":   box,
	})
	return err
}

func (c *Client) Gesture(gesture string) error {
	_, err := c.rpc.Call("gesture", gesture)
	return err
}

func (c *Client) Wait(activity string, timeout time.Duration) (bool, error) {
	result, err := c.rpc.Call("wait", map[string]interface{}{
		"activity": activity,
		"timeout":  int(timeout / time.Millisecond),
	})
	if err != nil {
		return false, err
	}
	var ok bool
	if err := json.Unmarshal(result, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (c *Client) WaitForExists(sel Selector, timeout time.Duration) (bool, error) {
	result, err := c.rpc.Call("waitForExists", map[string]interface{}{
		"selector": sel,
		"timeout":  int(timeout / time.Millisecond),
	})
	if err != nil {
		return false, err
	}
	var ok bool
	if err := json.Unmarshal(result, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (c *Client) WaitUntilExists(sel Selector, timeout time.Duration) (bool, error) {
	return c.WaitForExists(sel, timeout)
}

func (c *Client) Exists(sel Selector) bool {
	elements, err := c.Selector(sel)
	if err != nil || len(elements) == 0 {
		return false
	}
	return true
}

func (c *Client) ElementInfo(elementID string) (*ElementInfo, error) {
	result, err := c.rpc.Call("elementInfo", map[string]interface{}{
		"elementId": elementID,
	})
	if err != nil {
		return nil, err
	}
	var info ElementInfo
	if err := json.Unmarshal(result, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func (c *Client) StopUiautomator() error {
	_, err := c.rpc.Call("stopUiautomator", nil)
	return err
}
