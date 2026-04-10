# Vision 模块重构设计

## 概述

将 `pkg/vision` 下的匹配功能拆分为 `match` 和 `diff` 两个独立子包，提高代码组织清晰度。

## 目标

1. `match` - 封装现有 Matcher 接口，提供更简洁的 API
2. `diff` - 新增图像差异对比功能，支持区域框选和可视化输出

## 目录结构

```
pkg/vision/
├── match/                      # 图像匹配功能
│   ├── match.go                # Matcher 接口和工厂函数
│   ├── result.go               # MatchResult 类型定义
│   ├── template.go             # 模板匹配实现
│   ├── multiscale.go           # 多尺度匹配实现
│   ├── sift.go                 # SIFT 特征匹配实现
│   └── loftr.go                # LoFTR 匹配实现
├── diff/                       # 图像差异对比功能
│   ├── diff.go                 # Diff 接口和工厂函数
│   ├── result.go               # DiffResult, DiffRegion 类型
│   └── differ.go               # 差异对比核心算法
├── matcher.go                  # 现有 Matcher 接口（保留，alias）
├── result.go                   # 现有 MatchResult（保留，alias）
├── matcher_*.go                # 现有实现（保留，迁移后移除）
└── ...
```

## API 设计

### match 模块

```go
package match

// Matcher 图像匹配接口
type Matcher interface {
    Find(screenshot, template []byte) ([]*MatchResult, error)
    Name() string
}

// MatchResult 匹配结果
type MatchResult struct {
    Found  bool
    X, Y   int
    Width  int
    Height int
    Score  float64
}

func (r *MatchResult) Center() (int, int) { ... }
func (r *MatchResult) Rectangle() image.Rectangle { ... }

// New 创建 Matcher
func New(algo string, opts ...Option) (Matcher, error)

// Config 匹配配置
type Config struct {
    Threshold    float64
    ScaleMin     float64
    ScaleMax     float64
    ScaleStep    float64
    NMSThreshold float64
    DebugDir     string
}

// Option 配置选项
type Option func(*Config)
func WithThreshold(t float64) Option { ... }
func WithScaleRange(min, max float64) Option { ... }
...
```

### diff 模块

```go
package diff

// Config 差异对比配置
type Config struct {
    Threshold float64  // 差异阈值 (0.0-1.0)
    Region    *Rect    // 可选：框选区域
    OutputDir string   // 差异图输出目录
}

// Rect 区域定义
type Rect struct {
    X, Y, Width, Height int
}

// DiffResult 差异对比结果
type DiffResult struct {
    HasDiff    bool          // 是否有差异
    Diffs      []DiffRegion  // 差异区域列表
    OutputPath string        // 可视化差异图路径
    Similarity float64       // 相似度 (0.0-1.0)
}

// DiffRegion 差异区域
type DiffRegion struct {
    X, Y, Width, Height int
    Score        float64  // 该区域差异分数
    PixelCount   int      // 差异像素数量
}

// Differ 差异对比接口
type Differ interface {
    Compare(img1, img2 []byte, cfg *Config) (*DiffResult, error)
}

// New 创建 Differ
func New(algo string, opts ...Option) (Differ, error)

// Option 配置选项
type Option func(*Config)
func WithThreshold(t float64) Option { ... }
func WithRegion(x, y, w, h int) Option { ... }
func WithOutputDir(dir string) Option { ... }
```

## diff 算法实现

### 像素级差异对比

1. 将两张图片缩放到相同尺寸
2. 按区域（或全图）计算像素差异
3. 差异超过阈值的位置标记为差异点
4. 使用聚类算法将相邻差异点合并为差异区域
5. 生成可视化差异图（高亮差异区域）

### 可视化输出

生成的差异图命名格式：`diff_<timestamp>_<hash>.png`

差异图特性：
- 原始图片作为背景
- 差异区域用红色框标注
- 支持多区域同时显示

## 实现计划

### Phase 1: match 模块
1. 创建 `pkg/vision/match/` 目录结构
2. 实现 `match.go` - 接口和工厂函数
3. 实现 `result.go` - 结果类型
4. 实现各算法文件（可从现有文件迁移）

### Phase 2: diff 模块
1. 创建 `pkg/vision/diff/` 目录结构
2. 实现 `diff.go` - 接口和工厂函数
3. 实现 `result.go` - 结果类型
4. 实现 `differ.go` - 核心差异算法
5. 实现可视化输出

### Phase 3: 清理和兼容
1. 更新 `pkg/vision/matcher.go` 作为兼容层（alias）
2. 添加单元测试
3. 更新文档

## 迁移策略

保留现有 `pkg/vision/matcher.go` 作为兼容层：
```go
// Deprecated: 使用 pkg/vision/match 替代
type Matcher = match.Matcher
type MatchResult = match.MatchResult
```

## 验收标准

1. ✅ `match.New("template")` 等价于原有 `vision.NewMatcher("template")`
2. ✅ `diff.Compare()` 能正确检测两张图片的差异区域
3. ✅ 支持区域框选，只对比指定区域
4. ✅ 生成的可视化差异图清晰标注差异区域
5. ✅ 所有现有测试通过
