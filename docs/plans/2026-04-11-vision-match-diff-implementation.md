# Vision Match & Diff Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 在 `pkg/vision` 下创建 `match` 和 `diff` 两个子包，分别封装图像匹配和图像差异对比功能。

**Architecture:** 
- `match` 包封装现有 Matcher 接口，提供更简洁的 API
- `diff` 包实现图像差异对比，支持区域框选和可视化输出
- 保留现有 `pkg/vision/matcher.go` 作为兼容层（type alias）

**Tech Stack:** Go, image, imaging library (if needed)

---

## Task 1: 创建 match 包基础结构

**Files:**
- Create: `pkg/vision/match/match.go`
- Create: `pkg/vision/match/result.go`

**Step 1: 创建 match.go - Matcher 接口和工厂函数**

```go
package match

import (
	"fmt"
	"image"
)

type Matcher interface {
	Find(screenshot, template []byte) ([]*MatchResult, error)
	Name() string
}

type Config struct {
	Threshold    float64
	ScaleMin     float64
	ScaleMax     float64
	ScaleStep    float64
	NMSThreshold float64
	DebugDir     string
}

var defaultConfig = &Config{
	Threshold:    0.8,
	ScaleMin:     0.8,
	ScaleMax:     1.2,
	ScaleStep:    0.1,
	NMSThreshold: 0.5,
	DebugDir:     "",
}

type Option func(*Config)

func WithThreshold(t float64) Option {
	return func(c *Config) { c.Threshold = t }
}

func WithScaleRange(min, max float64) Option {
	return func(c *Config) { c.ScaleMin = min; c.ScaleMax = max }
}

func WithScaleStep(step float64) Option {
	return func(c *Config) { c.ScaleStep = step }
}

func WithNMSThreshold(t float64) Option {
	return func(c *Config) { c.NMSThreshold = t }
}

func WithDebug(outputDir string) Option {
	return func(c *Config) { c.DebugDir = outputDir }
}

// New 创建 Matcher 实例
func New(algo string, opts ...Option) (Matcher, error) {
	cfg := defaultConfig.clone()
	for _, opt := range opts {
		opt(cfg)
	}

	switch algo {
	case "template":
		return newTemplateMatcher(cfg), nil
	case "multiscale":
		return newMultiscaleMatcher(cfg), nil
	case "sift":
		return newSIFTMatcher(cfg), nil
	case "loftr":
		return newLoFTRMatcher(cfg), nil
	default:
		return nil, fmt.Errorf("unknown algorithm: %s", algo)
	}
}

func (c *Config) clone() *Config {
	return &Config{
		Threshold:    c.Threshold,
		ScaleMin:     c.ScaleMin,
		ScaleMax:     c.ScaleMax,
		ScaleStep:    c.ScaleStep,
		NMSThreshold: c.NMSThreshold,
		DebugDir:     c.DebugDir,
	}
}
```

**Step 2: 创建 result.go - MatchResult 类型**

```go
package match

import "image"

type MatchResult struct {
	Found  bool
	X, Y   int
	Width  int
	Height int
	Score  float64
}

func (r *MatchResult) Center() (int, int) {
	return r.X + r.Width/2, r.Y + r.Height/2
}

func (r *MatchResult) Rectangle() image.Rectangle {
	return image.Rect(r.X, r.Y, r.X+r.Width, r.Y+r.Height)
}
```

---

## Task 2: 实现 match 包算法文件

**Files:**
- Create: `pkg/vision/match/template.go`
- Create: `pkg/vision/match/multiscale.go`
- Create: `pkg/vision/match/sift.go`
- Create: `pkg/vision/match/loftr.go`

**Step 1: 从现有 `pkg/vision/matcher_*.go` 复制实现到各算法文件**

从现有文件复制内容到新位置，并修改 package 名称为 `match`。

---

## Task 3: 创建 diff 包基础结构

**Files:**
- Create: `pkg/vision/diff/diff.go`
- Create: `pkg/vision/diff/result.go`
- Create: `pkg/vision/diff/differ.go`

**Step 1: 创建 diff.go - Differ 接口和工厂函数**

```go
package diff

import (
	"fmt"
)

type Differ interface {
	Compare(img1, img2 []byte, cfg *Config) (*DiffResult, error)
	Name() string
}

type Config struct {
	Threshold float64  // 差异阈值 (0.0-1.0)
	Region    *Rect    // 可选：框选区域
	OutputDir string   // 差异图输出目录
}

type Rect struct {
	X, Y, Width, Height int
}

var defaultConfig = &Config{
	Threshold: 0.1,
	Region:    nil,
	OutputDir: "",
}

type Option func(*Config)

func WithThreshold(t float64) Option {
	return func(c *Config) { c.Threshold = t }
}

func WithRegion(x, y, w, h int) Option {
	return func(c *Config) { c.Region = &Rect{X: x, Y: y, Width: w, Height: h} }
}

func WithOutputDir(dir string) Option {
	return func(c *Config) { c.OutputDir = dir }
}

// New 创建 Differ 实例
func New(algo string, opts ...Option) (Differ, error) {
	cfg := defaultConfig.clone()
	for _, opt := range opts {
		opt(cfg)
	}

	switch algo {
	case "pixel":
		return newPixelDiffer(cfg), nil
	default:
		return nil, fmt.Errorf("unknown algorithm: %s", algo)
	}
}

func (c *Config) clone() *Config {
	n := &Config{
		Threshold: c.Threshold,
		OutputDir: c.OutputDir,
	}
	if c.Region != nil {
		n.Region = &Rect{X: c.Region.X, Y: c.Region.Y, Width: c.Region.Width, Height: c.Region.Height}
	}
	return n
}
```

**Step 2: 创建 result.go - DiffResult, DiffRegion 类型**

```go
package diff

type DiffResult struct {
	HasDiff    bool          // 是否有差异
	Diffs      []DiffRegion  // 差异区域列表
	OutputPath string        // 可视化差异图路径
	Similarity float64       // 相似度 (0.0-1.0)
}

type DiffRegion struct {
	X, Y, Width, Height int
	Score        float64  // 该区域差异分数
	PixelCount   int      // 差异像素数量
}
```

**Step 3: 创建 differ.go - 像素级差异对比实现**

```go
package diff

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"time"
)

type pixelDiffer struct {
	cfg *Config
}

func newPixelDiffer(cfg *Config) Differ {
	return &pixelDiffer{cfg: cfg}
}

func (d *pixelDiffer) Name() string {
	return "pixel"
}

func (d *pixelDiffer) Compare(img1, img2 []byte, cfg *Config) (*DiffResult, error) {
	if cfg == nil {
		cfg = d.cfg
	}

	m1, err := decodeImage(img1)
	if err != nil {
		return nil, err
	}
	m2, err := decodeImage(img2)
	if err != nil {
		return nil, err
	}

	// 统一尺寸
	m1 = normalizeSize(m1, m2.Bounds())
	m2 = normalizeSize(m2, m1.Bounds())

	bounds := m1.Bounds()
	result := &DiffResult{
		Diffs: make([]DiffRegion, 0),
	}

	// 如果有区域限制，使用区域
	if cfg.Region != nil {
		result.Diffs = d.compareRegion(m1, m2, cfg.Region, cfg.Threshold)
	} else {
		result.Diffs = d.compareFullImage(m1, m2, cfg.Threshold)
	}

	// 计算相似度
	totalPixels := bounds.Dx() * bounds.Dy()
	diffPixels := 0
	for _, r := range result.Diffs {
		diffPixels += r.PixelCount
	}
	result.Similarity = 1.0 - float64(diffPixels)/float64(totalPixels)
	result.HasDiff = len(result.Diffs) > 0

	// 生成差异图
	if cfg.OutputDir != "" && result.HasDiff {
		path, err := d.renderDiff(m1, m2, result.Diffs, cfg.OutputDir)
		if err == nil {
			result.OutputPath = path
		}
	}

	return result, nil
}

func (d *pixelDiffer) compareRegion(m1, m2 *image.RGBA, region *Rect, threshold float64) []DiffRegion {
	// 实现区域对比逻辑
	// ...
}

func (d *pixelDiffer) compareFullImage(m1, m2 *image.RGBA, threshold float64) []DiffRegion {
	// 实现全图对比逻辑
	// ...
}

func (d *pixelDiffer) renderDiff(m1, m2 *image.RGBA, diffs []DiffRegion, outputDir string) (string, error) {
	// 创建差异图，标注差异区域
	// ...
}
```

---

## Task 4: 添加单元测试

**Files:**
- Create: `pkg/vision/match/match_test.go`
- Create: `pkg/vision/diff/diff_test.go`

**Step 1: 编写 match_test.go**

```go
package match

import (
	"testing"
)

func TestNew(t *testing.T) {
	m, err := New("template")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if m == nil {
		t.Fatal("New() returned nil")
	}
	if m.Name() != "template" {
		t.Errorf("Name() = %s, want template", m.Name())
	}
}

func TestNew_UnknownAlgorithm(t *testing.T) {
	_, err := New("unknown")
	if err == nil {
		t.Fatal("New() expected error for unknown algorithm")
	}
}

func TestOptions(t *testing.T) {
	m, err := New("template",
		WithThreshold(0.9),
		WithScaleRange(0.5, 1.5),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	_ = m // 使用 matcher
}
```

**Step 2: 编写 diff_test.go**

```go
package diff

import (
	"testing"
)

func TestNew(t *testing.T) {
	d, err := New("pixel")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if d == nil {
		t.Fatal("New() returned nil")
	}
	if d.Name() != "pixel" {
		t.Errorf("Name() = %s, want pixel", d.Name())
	}
}

func TestNew_UnknownAlgorithm(t *testing.T) {
	_, err := New("unknown")
	if err == nil {
		t.Fatal("New() expected error for unknown algorithm")
	}
}

func TestOptions(t *testing.T) {
	d, err := New("pixel",
		WithThreshold(0.05),
		WithRegion(0, 0, 100, 100),
		WithOutputDir("/tmp"),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	_ = d // 使用 differ
}
```

---

## Task 5: 添加兼容层

**Files:**
- Modify: `pkg/vision/matcher.go`

**Step 1: 更新 matcher.go 添加类型别名**

```go
package vision

// Deprecated: Use pkg/vision/match instead.
type Matcher = match.Matcher

// Deprecated: Use pkg/vision/match instead.
type MatchResult = match.MatchResult

// MatchResult 的方法需要委托
func (r *MatchResult) Center() (int, int) {
	return r.X + r.Width/2, r.Y + r.Height/2
}

func (r *MatchResult) Rectangle() image.Rectangle {
	return image.Rect(r.X, r.Y, r.X+r.Width, r.Y+r.Height)
}
```

---

## Task 6: 验证和清理

**Step 1: 运行所有测试**

```bash
go test ./pkg/vision/... -v
```

**Step 2: 检查 LSP 错误**

```bash
go build ./pkg/vision/...
```

**Step 3: 更新文档（如需要）**
