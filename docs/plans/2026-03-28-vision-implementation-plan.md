# Vision Factory Pattern Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Refactor pkg/vision to use factory pattern supporting multiple matching algorithms (template, multiscale, sift, loftr) with unified interface.

**Architecture:** Factory pattern with interface-based matchers. Each algorithm in separate subpackage. Options pattern for configuration. Debug rendering support.

**Tech Stack:** Go, gocv.io/x/gocv, github.com/yalue/onnxruntime_go

---

## Task 1: Create MatchResult and Helper Methods

**Files:**
- Create: `pkg/vision/result.go`

**Step 1: Write the test**

```go
// result_test.go
package vision

import (
    "image"
    "testing"
)

func TestMatchResult_Center(t *testing.T) {
    r := &MatchResult{
        X: 100, Y: 200, Width: 50, Height: 60,
    }
    cx, cy := r.Center()
    if cx != 125 || cy != 230 {
        t.Errorf("Center() = (%d, %d), want (125, 230)", cx, cy)
    }
}

func TestMatchResult_Rectangle(t *testing.T) {
    r := &MatchResult{
        X: 100, Y: 200, Width: 50, Height: 60,
    }
    rect := r.Rectangle()
    if rect != image.Rect(100, 200, 150, 260) {
        t.Errorf("Rectangle() = %v, want Rect(100, 200, 150, 260)", rect)
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/vision/... -run TestMatchResult -v`
Expected: FAIL - undefined type MatchResult

**Step 3: Write implementation**

```go
// pkg/vision/result.go
package vision

import "image"

type MatchResult struct {
    Found      bool
    X, Y       int
    Width      int
    Height     int
    Score      float64
}

func (r *MatchResult) Center() (int, int) {
    return r.X + r.Width/2, r.Y + r.Height/2
}

func (r *MatchResult) Rectangle() image.Rectangle {
    return image.Rect(r.X, r.Y, r.X+r.Width, r.Y+r.Height)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./pkg/vision/... -run TestMatchResult -v`
Expected: PASS

**Step 5: Commit**

```bash
git add pkg/vision/result.go pkg/vision/result_test.go
git commit -m "feat(vision): add MatchResult with Center and Rectangle helpers"
```

---

## Task 2: Create Matcher Interface and Factory

**Files:**
- Create: `pkg/vision/matcher.go`
- Create: `pkg/vision/matcher_test.go`

**Step 1: Write the test**

```go
// matcher_test.go
package vision

import (
    "testing"
)

func TestNewMatcher_UnknownAlgorithm(t *testing.T) {
    _, err := NewMatcher("unknown")
    if err == nil {
        t.Error("NewMatcher(unknown) should return error")
    }
    if err.Error() != "unknown algorithm: unknown" {
        t.Errorf("error = %q, want %q", err.Error(), "unknown algorithm: unknown")
    }
}

func TestNewMatcher_Template(t *testing.T) {
    m, err := NewMatcher("template")
    if err != nil {
        t.Fatalf("NewMatcher(template) failed: %v", err)
    }
    if m.Name() != "template" {
        t.Errorf("Name() = %q, want %q", m.Name(), "template")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/vision/... -run TestNewMatcher -v`
Expected: FAIL - undefined function NewMatcher

**Step 3: Write implementation**

```go
// pkg/vision/matcher.go
package vision

import (
    "fmt"
)

// Matcher defines the interface for visual matching algorithms
type Matcher interface {
    Find(screenshot, template []byte) ([]*MatchResult, error)
    Name() string
    DebugRender(screenshot []byte, results []*MatchResult) []byte
}

type config struct {
    threshold    float64
    scaleMin     float64
    scaleMax     float64
    scaleStep    float64
    nmsThreshold float64
    debugDir     string
}

var defaultConfig = func() *config {
    return &config{
        threshold:    0.8,
        scaleMin:     0.8,
        scaleMax:     1.2,
        scaleStep:    0.1,
        nmsThreshold: 0.5,
        debugDir:     "",
    }
}

type Option func(*config)

func WithThreshold(t float64) Option {
    return func(c *config) { c.threshold = t }
}

func WithScaleRange(min, max float64) Option {
    return func(c *config) {
        c.scaleMin = min
        c.scaleMax = max
    }
}

func WithScaleStep(step float64) Option {
    return func(c *config) { c.scaleStep = step }
}

func WithNMSThreshold(t float64) Option {
    return func(c *config) { c.nmsThreshold = t }
}

func WithDebug(outputDir string) Option {
    return func(c *config) { c.debugDir = outputDir }
}

func NewMatcher(algo string, opts ...Option) (Matcher, error) {
    cfg := defaultConfig()
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
```

**Step 4: Run test to verify it passes**

Run: `go test ./pkg/vision/... -run TestNewMatcher -v`
Expected: PASS

**Step 5: Commit**

```bash
git add pkg/vision/matcher.go pkg/vision/matcher_test.go
git commit -m "feat(vision): add Matcher interface and factory"
```

---

## Task 3: Create matchers directory and Template Matcher

**Files:**
- Create: `pkg/vision/matchers/template.go`
- Create: `pkg/vision/matchers/template_test.go`

**Step 1: Write the test**

```go
package matchers

import (
    "testing"
)

func TestTemplateMatcher_Name(t *testing.T) {
    m := newTemplateMatcher(nil)
    if m.Name() != "template" {
        t.Errorf("Name() = %q, want %q", m.Name(), "template")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/vision/matchers/... -run TestTemplateMatcher -v`
Expected: FAIL - undefined function newTemplateMatcher

**Step 3: Write stub implementation**

```go
// pkg/vision/matchers/template.go
package matchers

import (
    "github.com/liukunup/go-uop/pkg/vision"
)

type templateMatcher struct {
    config *config
}

func newTemplateMatcher(cfg *config) *templateMatcher {
    return &templateMatcher{config: cfg}
}

func (m *templateMatcher) Name() string {
    return "template"
}

func (m *templateMatcher) Find(screenshot, template []byte) ([]*vision.MatchResult, error) {
    // TODO: implement OpenCV template matching
    return nil, nil
}

func (m *templateMatcher) DebugRender(screenshot []byte, results []*vision.MatchResult) []byte {
    return nil
}
```

**Note:** Config type is not yet defined. We'll add it next as a shared types file.

**Step 4: Fix compilation error - create shared config**

Run: `go build ./pkg/vision/...` to check errors

**Step 5: Create config.go in pkg/vision**

```go
// pkg/vision/config.go
package vision

type config struct {
    threshold    float64
    scaleMin     float64
    scaleMax     float64
    scaleStep    float64
    nmsThreshold float64
    debugDir     string
}
```

**Step 6: Run test to verify it passes**

Run: `go test ./pkg/vision/matchers/... -run TestTemplateMatcher -v`
Expected: PASS

**Step 7: Commit**

```bash
git add pkg/vision/config.go pkg/vision/matchers/template.go pkg/vision/matchers/template_test.go
git commit -m "feat(vision): add template matcher stub"
```

---

## Task 4: Implement Template Matcher with OpenCV

**Files:**
- Modify: `pkg/vision/matchers/template.go`

**Step 1: Write the failing test with real OpenCV behavior**

```go
// pkg/vision/matchers/template_test.go
package matchers

import (
    "os"
    "testing"
)

func TestTemplateMatcher_Find(t *testing.T) {
    // Skip if OpenCV not available
    if os.Getenv("TEST_OPENCV") != "1" {
        t.Skip("Skipping OpenCV test")
    }
    
    m := newTemplateMatcher(nil)
    results, err := m.Find(nil, nil)
    if err != nil {
        t.Fatalf("Find() error = %v", err)
    }
    if len(results) != 0 {
        t.Errorf("Find() with nil input should return empty, got %d results", len(results))
    }
}
```

**Step 2: Run test to verify it passes with current stub**

Run: `go test ./pkg/vision/matchers/... -run TestTemplateMatcher_Name -v`
Expected: PASS

**Step 3: Implement OpenCV template matching**

```go
// pkg/vision/matchers/template.go
package matchers

import (
    "bytes"
    "image"
    "image/png"
    "os"

    "gocv.io/x/gocv"
    "github.com/liukunup/go-uop/pkg/vision"
)

type templateMatcher struct {
    config *config
}

func newTemplateMatcher(cfg *config) *templateMatcher {
    return &templateMatcher{config: cfg}
}

func (m *templateMatcher) Name() string {
    return "template"
}

func (m *templateMatcher) Find(screenshot, templateImg []byte) ([]*vision.MatchResult, error) {
    if screenshot == nil || templateImg == nil {
        return nil, nil
    }

    // Decode images
    screen, err := gocv.IMDecode(screenshot, gocv.IMReadGrayScale)
    if err != nil {
        return nil, err
    }
    defer screen.Close()

    tmpl, err := gocv.IMDecode(templateImg, gocv.IMReadGrayScale)
    if err != nil {
        return nil, err
    }
    defer tmpl.Close()

    // Ensure template is smaller than screen
    if screen.Cols() < tmpl.Cols() || screen.Rows() < tmpl.Rows() {
        return nil, nil
    }

    // Perform template matching
    result := gocv.NewMat()
    defer result.Close()

    gocv.MatchTemplate(screen, tmpl, &result, gocv.TmCcoeffNormed, gocv.NewMat())

    // Find best match
    _, maxVal, _, maxLoc := gocv.MinMaxLoc(result)
    
    threshold := 0.8
    if m.config != nil {
        threshold = m.config.threshold
    }

    if maxVal < threshold {
        return nil, nil
    }

    return []*vision.MatchResult{
        {
            Found:  true,
            X:      maxLoc.X,
            Y:      maxLoc.Y,
            Width:  tmpl.Cols(),
            Height: tmpl.Rows(),
            Score:  maxVal,
        },
    }, nil
}

func (m *templateMatcher) DebugRender(screenshot []byte, results []*vision.MatchResult) []byte {
    if m.config == nil || m.config.debugDir == "" || len(results) == 0 {
        return nil
    }

    img, err := gocv.IMDecode(screenshot, gocv.IMReadColor)
    if err != nil {
        return nil
    }
    defer img.Close()

    for i, r := range results {
        if !r.Found {
            continue
        }

        // Draw rectangle
        rect := image.Rect(r.X, r.Y, r.X+r.Width, r.Y+r.Height)
        gocv.Rectangle(&img, rect, color.RGBA{0, 255, 0, 0}, 2)

        // Draw center cross
        cx, cy := r.Center()
        gocv.Line(&img, image.Point{cx - 10, cy}, image.Point{cx + 10, cy}, color.RGBA{0, 255, 0, 0}, 2)
        gocv.Line(&img, image.Point{cx, cy - 10}, image.Point{cx, cy + 10}, color.RGBA{0, 255, 0, 0}, 2)

        // Draw label
        label := fmt.Sprintf("[%d] Score: %.2f Pos: (%d,%d)", i, r.Score, r.X, r.Y)
        gocv.PutText(&img, label, image.Point{r.X, r.Y - 10}, gocv.FontHersheySimplex, 0.5, color.RGBA{255, 255, 255, 0}, 2)
    }

    buf, _ := gocv.IMEncode(".png", img)
    return buf
}
```

**Step 4: Verify compilation**

Run: `go build ./pkg/vision/matchers/...`
Expected: Need to add missing imports (bytes, fmt, image/color, os)

**Step 5: Fix imports and commit**

```bash
git add pkg/vision/matchers/template.go
git commit -m "feat(vision): implement OpenCV template matching"
```

---

## Task 5: Create Multiscale Matcher

**Files:**
- Create: `pkg/vision/matchers/multiscale.go`
- Create: `pkg/vision/matchers/multiscale_test.go`

**Step 1: Write the test**

```go
package matchers

import (
    "testing"
)

func TestMultiscaleMatcher_Name(t *testing.T) {
    m := newMultiscaleMatcher(nil)
    if m.Name() != "multiscale" {
        t.Errorf("Name() = %q, want %q", m.Name(), "multiscale")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/vision/matchers/... -run TestMultiscaleMatcher -v`
Expected: FAIL - undefined newMultiscaleMatcher

**Step 3: Write implementation**

```go
// pkg/vision/matchers/multiscale.go
package matchers

import (
    "sort"

    "gocv.io/x/gocv"
    "github.com/liukunup/go-uop/pkg/vision"
)

type multiscaleMatcher struct {
    config *config
}

func newMultiscaleMatcher(cfg *config) *multiscaleMatcher {
    return &multiscaleMatcher{config: cfg}
}

func (m *multiscaleMatcher) Name() string {
    return "multiscale"
}

func (m *multiscaleMatcher) Find(screenshot, templateImg []byte) ([]*vision.MatchResult, error) {
    if screenshot == nil || templateImg == nil {
        return nil, nil
    }

    cfg := defaultConfig()
    if m.config != nil {
        cfg = m.config
    }

    // Decode images
    screen, err := gocv.IMDecode(screenshot, gocv.IMReadGrayScale)
    if err != nil {
        return nil, err
    }
    defer screen.Close()

    tmpl, err := gocv.IMDecode(templateImg, gocv.IMReadGrayScale)
    if err != nil {
        return nil, err
    }
    defer tmpl.Close()

    // Multi-scale search
    var allMatches []matchCandidate

    for scale := cfg.scaleMin; scale <= cfg.scaleMax; scale += cfg.scaleStep {
        scaledTmpl := gocv.NewMat()
        gocv.Resize(tmpl, &scaledTmpl, image.Point{
            X: int(float64(tmpl.Cols()) * scale),
            Y: int(float64(tmpl.Rows()) * scale),
        })
        defer scaledTmpl.Close()

        if screen.Cols() < scaledTmpl.Cols() || screen.Rows() < scaledTmpl.Rows() {
            continue
        }

        result := gocv.NewMat()
        gocv.MatchTemplate(screen, scaledTmpl, &result, gocv.TmCcoeffNormed, gocv.NewMat())

        // Find all matches above threshold
        for y := 0; y < result.Rows(); y++ {
            for x := 0; x < result.Cols(); x++ {
                val := result.GetFloatAt(y, x)
                if val >= cfg.threshold {
                    allMatches = append(allMatches, matchCandidate{
                        x: x, y: y,
                        width: scaledTmpl.Cols(), height: scaledTmpl.Rows(),
                        score: val, scale: scale,
                    })
                }
            }
        }
        result.Close()
    }

    if len(allMatches) == 0 {
        return nil, nil
    }

    // Sort by score descending
    sort.Slice(allMatches, func(i, j int) bool {
        return allMatches[i].score > allMatches[j].score
    })

    // Apply NMS
    var results []*vision.MatchResult
    used := make([]bool, len(allMatches))

    for i := 0; i < len(allMatches); i++ {
        if used[i] {
            continue
        }

        cand := allMatches[i]
        results = append(results, &vision.MatchResult{
            Found:  true,
            X:      cand.x,
            Y:      cand.y,
            Width:  cand.width,
            Height: cand.height,
            Score:  cand.score,
        })

        // Suppress overlapping candidates
        for j := i + 1; j < len(allMatches); j++ {
            if used[j] {
                continue
            }
            if overlap(cand, allMatches[j]) > cfg.nmsThreshold {
                used[j] = true
            }
        }
    }

    return results, nil
}

func (m *multiscaleMatcher) DebugRender(screenshot []byte, results []*vision.MatchResult) []byte {
    // Same as templateMatcher
    return debugRender(m.config, screenshot, results)
}

type matchCandidate struct {
    x, y   int
    width  int
    height int
    score  float32
    scale  float64
}

func overlap(a, b matchCandidate) float64 {
    x1 := max(a.x, b.x)
    y1 := max(a.y, b.y)
    x2 := min(a.x+a.width, b.x+b.width)
    y2 := min(a.y+a.height, b.y+b.height)

    if x2 <= x1 || y2 <= y1 {
        return 0
    }

    inter := float64((x2 - x1) * (y2 - y1))
    union := float64(a.width*a.height + b.width*b.height - inter)
    return inter / union
}
```

**Step 4: Verify compilation and run test**

Run: `go test ./pkg/vision/matchers/... -run TestMultiscaleMatcher -v`
Expected: PASS

**Step 5: Commit**

```bash
git add pkg/vision/matchers/multiscale.go pkg/vision/matchers/multiscale_test.go
git commit -m "feat(vision): add multiscale template matching with NMS"
```

---

## Task 6: Create SIFT Matcher

**Files:**
- Create: `pkg/vision/sift/matcher.go`
- Create: `pkg/vision/sift/matcher_test.go`

**Step 1: Write the test stub**

```go
package sift

import (
    "testing"
)

func TestSIFTMatcher_Name(t *testing.T) {
    m := newSIFTMatcher(nil)
    if m.Name() != "sift" {
        t.Errorf("Name() = %q, want %q", m.Name(), "sift")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/vision/sift/... -v`
Expected: FAIL - package sift does not exist

**Step 3: Create implementation**

```go
// pkg/vision/sift/matcher.go
package sift

import (
    "gocv.io/x/gocv"
    "github.com/liukunup/go-uop/pkg/vision"
)

type siftMatcher struct {
    config *config
}

func newSIFTMatcher(cfg *config) *siftMatcher {
    return &siftMatcher{config: cfg}
}

func (m *siftMatcher) Name() string {
    return "sift"
}

func (m *siftMatcher) Find(screenshot, templateImg []byte) ([]*vision.MatchResult, error) {
    if screenshot == nil || templateImg == nil {
        return nil, nil
    }

    // Decode images
    img1, err := gocv.IMDecode(screenshot, gocv.IMReadGrayScale)
    if err != nil {
        return nil, err
    }
    defer img1.Close()

    img2, err := gocv.IMDecode(templateImg, gocv.IMReadGrayScale)
    if err != nil {
        return nil, err
    }
    defer img2.Close()

    // Create SIFT detector
    sift := gocv.NewSIFT()
    defer sift.Close()

    // Detect keypoints and compute descriptors
    kp1, des1 := sift.DetectAndCompute(img1, gocv.NewMat())
    defer kp1.Close()
    defer des1.Close()

    kp2, des2 := sift.DetectAndCompute(img2, gocv.NewMat())
    defer kp2.Close()
    defer des2.Close()

    if len(kp1) == 0 || len(kp2) == 0 {
        return nil, nil
    }

    // BFMatcher with kNN
    bf := gocv.NewBFMatcher(gocv.NormHamming, false)
    defer bf.Close()

    matches := bf.KnnMatch(des1, des2, 2)

    // Ratio test
    var goodMatches []gocv.DMatch
    for _, m := range matches {
        if len(m) > 1 {
            if m[0].Distance < 0.75*m[1].Distance {
                goodMatches = append(goodMatches, m[0])
            }
        }
    }

    if len(goodMatches) < 4 {
        return nil, nil
    }

    // RANSAC for geometric verification
    // Extract point coordinates
    srcPoints := make([][]byte, len(goodMatches))
    dstPoints := make([][]byte, len(goodMatches))
    for i, match := range goodMatches {
        srcPoints[i] = []byte{float32(kp1[match.QueryIdx].X), float32(kp1[match.QueryIdx].Y)}
        dstPoints[i] = []byte{float32(kp2[match.TrainIdx].X), float32(kp2[match.TrainIdx].Y)}
    }

    // Find homography with RANSAC
    homography, mask := gocv.FindHomography(srcPoints, dstPoints, gocv.Ransac)
    defer homography.Close()
    defer mask.Close()

    // Count inliers
    threshold := 0.8
    if m.config != nil {
        threshold = m.config.threshold
    }

    inliers := 0
    for i := 0; i < mask.Total(); i++ {
        if mask.GetIntAt(i, 0) == 1 {
            inliers++
        }
    }

    if float64(inliers)/float64(len(goodMatches)) < threshold {
        return nil, nil
    }

    // Estimate bounding box
    h := img1.Rows()
    w := img1.Cols()

    // Apply homography to corners
    corners := []image.Point{
        {0, 0}, {w, 0}, {w, h}, {0, h},
    }

    minX, minY := corners[0].X, corners[0].Y
    maxX, maxY := corners[0].X, corners[0].Y
    for _, p := range corners[1:] {
        if p.X < minX {
            minX = p.X
        }
        if p.X > maxX {
            maxX = p.X
        }
        if p.Y < minY {
            minY = p.Y
        }
        if p.Y > maxY {
            maxY = p.Y
        }
    }

    score := float64(inliers) / float64(len(goodMatches))

    return []*vision.MatchResult{
        {
            Found:  true,
            X:      minX,
            Y:      minY,
            Width:  maxX - minX,
            Height: maxY - minY,
            Score:  score,
        },
    }, nil
}

func (m *siftMatcher) DebugRender(screenshot []byte, results []*vision.MatchResult) []byte {
    return debugRender(m.config, screenshot, results)
}
```

**Step 4: Verify compilation and run test**

Run: `go test ./pkg/vision/sift/... -v`
Expected: PASS (with OpenCV)

**Step 5: Commit**

```bash
git add pkg/vision/sift/matcher.go pkg/vision/sift/matcher_test.go
git commit -m "feat(vision): add SIFT + RANSAC matcher"
```

---

## Task 7: Create LoFTR Matcher

**Files:**
- Create: `pkg/vision/loftr/matcher.go`
- Create: `pkg/vision/loftr/matcher_test.go`

**Step 1: Write the test stub**

```go
package loftr

import (
    "testing"
)

func TestLoFTRMatcher_Name(t *testing.T) {
    m := newLoFTRMatcher(nil)
    if m.Name() != "loftr" {
        t.Errorf("Name() = %q, want %q", m.Name(), "loftr")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/vision/loftr/... -v`
Expected: FAIL - package loftr does not exist

**Step 3: Create implementation**

```go
// pkg/vision/loftr/matcher.go
package loftr

import (
    "github.com/yalue/onnxruntime_go"
    "github.com/liukunup/go-uop/pkg/vision"
)

type loftrMatcher struct {
    config     *config
    session    *onnxruntime_go.AdvancedSession
    inputShape []int64
}

func newLoFTRMatcher(cfg *config) *loftrMatcher {
    return &loftrMatcher{config: cfg}
}

func (m *loftrMatcher) Name() string {
    return "loftr"
}

func (m *loftrMatcher) Find(screenshot, templateImg []byte) ([]*vision.MatchResult, error) {
    if screenshot == nil || templateImg == nil {
        return nil, nil
    }

    // TODO: Implement LoFTR ONNX inference
    // 1. Initialize ONNX session with LoFTR model (from oooooha/loftr2onnx)
    // 2. Preprocess images (resize, normalize)
    // 3. Run inference
    // 4. Parse matches and apply confidence threshold
    // 5. Return bounding boxes

    return nil, nil
}

func (m *loftrMatcher) DebugRender(screenshot []byte, results []*vision.MatchResult) []byte {
    return debugRender(m.config, screenshot, results)
}
```

**Step 4: Verify compilation**

Run: `go build ./pkg/vision/loftr/...`
Expected: PASS (stub compiles)

**Step 5: Commit**

```bash
git add pkg/vision/loftr/matcher.go pkg/vision/loftr/matcher_test.go
git commit -m "feat(vision): add LoFTR matcher stub"
```

---

## Task 8: Add Debug Render Helper

**Files:**
- Create: `pkg/vision/matchers/debug.go`

**Step 1: Write debug render helper**

```go
// pkg/vision/matchers/debug.go
package matchers

import (
    "fmt"
    "image"
    "image/color"
    "os"
    "path/filepath"

    "gocv.io/x/gocv"
    "github.com/liukunup/go-uop/pkg/vision"
)

func debugRender(cfg *config, screenshot []byte, results []*vision.MatchResult) []byte {
    if cfg == nil || cfg.debugDir == "" || len(results) == 0 {
        return nil
    }

    img, err := gocv.IMDecode(screenshot, gocv.IMReadColor)
    if err != nil {
        return nil
    }
    defer img.Close()

    for i, r := range results {
        if !r.Found {
            continue
        }

        // Draw rectangle
        rect := image.Rect(r.X, r.Y, r.X+r.Width, r.Y+r.Height)
        gocv.Rectangle(&img, rect, color.RGBA{0, 255, 0, 0}, 2)

        // Draw center cross
        cx, cy := r.Center()
        gocv.Line(&img, image.Point{cx - 10, cy}, image.Point{cx + 10, cy}, color.RGBA{0, 255, 0, 0}, 2)
        gocv.Line(&img, image.Point{cx, cy - 10}, image.Point{cx, cy + 10}, color.RGBA{0, 255, 0, 0}, 2)

        // Draw label
        label := fmt.Sprintf("[%d] Score: %.2f Pos: (%d,%d)", i, r.Score, r.X, r.Y)
        gocv.PutText(&img, label, image.Point{r.X, r.Y - 10}, gocv.FontHersheySimplex, 0.5, color.RGBA{255, 255, 255, 0}, 2)
    }

    // Save to debug directory
    filename := fmt.Sprintf("debug_%d.png", len(results))
    outputPath := filepath.Join(cfg.debugDir, filename)
    os.MkdirAll(cfg.debugDir, 0755)
    gocv.IMWrite(outputPath, img)

    buf, _ := gocv.IMEncode(".png", img)
    return buf
}
```

**Step 2: Verify compilation**

Run: `go build ./pkg/vision/matchers/...`
Expected: PASS

**Step 3: Commit**

```bash
git add pkg/vision/matchers/debug.go
git commit -m "feat(vision): add debug render helper"
```

---

## Task 9: Delete Old Template File

**Files:**
- Delete: `pkg/vision/template.go` (old implementation)

**Step 1: Remove old file**

```bash
rm pkg/vision/template.go
```

**Step 2: Verify no references remain**

Run: `grep -r "simpleMatch\|byteReader\|NewTemplateMatcher\|SetTemplate" pkg/vision/`
Expected: No matches

**Step 3: Commit**

```bash
git rm pkg/vision/template.go
git commit -m "refactor(vision): remove old template.go"
```

---

## Task 10: Integration Test

**Files:**
- Create: `pkg/vision/integration_test.go`

**Step 1: Write integration test**

```go
package vision

import (
    "os"
    "testing"
)

func TestIntegration_AllMatchers(t *testing.T) {
    if os.Getenv("TEST_OPENCV") != "1" {
        t.Skip("Skipping integration test")
    }

    // Load test images
    screenshot, _ := os.ReadFile("testdata/screenshot.png")
    template, _ := os.ReadFile("testdata/button.png")

    matchers := []string{"template", "multiscale", "sift", "loftr"}

    for _, algo := range matchers {
        t.Run(algo, func(t *testing.T) {
            m, err := NewMatcher(algo,
                WithThreshold(0.7),
                WithScaleRange(0.8, 1.2),
                WithNMSThreshold(0.3),
            )
            if err != nil {
                t.Fatalf("NewMatcher(%s) failed: %v", algo, err)
            }

            results, err := m.Find(screenshot, template)
            if err != nil {
                t.Fatalf("Find() failed: %v", err)
            }

            t.Logf("Algorithm %s found %d matches", algo, len(results))
        })
    }
}
```

**Step 2: Run integration test**

Run: `TEST_OPENCV=1 go test ./pkg/vision/... -run TestIntegration -v`
Expected: Tests run (may skip if images not found)

**Step 3: Commit**

```bash
git add pkg/vision/integration_test.go
git commit -m "test(vision): add integration test"
```

---

## Summary

| Task | Description |
|------|-------------|
| 1 | MatchResult with Center and Rectangle helpers |
| 2 | Matcher interface and factory |
| 3 | matchers directory and template matcher stub |
| 4 | Implement OpenCV template matching |
| 5 | Multiscale matcher with NMS |
| 6 | SIFT + RANSAC matcher |
| 7 | LoFTR matcher stub |
| 8 | Debug render helper |
| 9 | Remove old template.go |
| 10 | Integration test |
