package diff

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
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

	bounds := m1.Bounds()
	result := &DiffResult{
		Diffs: make([]DiffRegion, 0),
	}

	if cfg.Region != nil {
		result.Diffs = d.compareRegion(m1, m2, cfg.Region, cfg.Threshold)
	} else {
		result.Diffs = d.compareFullImage(m1, m2, bounds, cfg.Threshold)
	}

	totalPixels := bounds.Dx() * bounds.Dy()
	diffPixels := 0
	for _, r := range result.Diffs {
		diffPixels += r.PixelCount
	}
	result.Similarity = 1.0 - float64(diffPixels)/float64(totalPixels)
	result.HasDiff = len(result.Diffs) > 0

	if cfg.OutputDir != "" && result.HasDiff {
		path, err := d.renderDiff(m1, m2, result.Diffs, cfg.OutputDir)
		if err == nil {
			result.OutputPath = path
		}
	}

	return result, nil
}

func (d *pixelDiffer) compareRegion(m1, m2 *image.RGBA, region *Rect, threshold float64) []DiffRegion {
	r := image.Rect(region.X, region.Y, region.X+region.Width, region.Y+region.Height)
	return d.compareFullImage(m1, m2, r, threshold)
}

func (d *pixelDiffer) compareFullImage(m1, m2 *image.RGBA, bounds image.Rectangle, threshold float64) []DiffRegion {
	diffMap := make(map[string]bool)
	var diffRegions []DiffRegion

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c1 := m1.RGBAAt(x, y)
			c2 := m2.RGBAAt(x, y)

			dr := float64(c1.R) - float64(c2.R)
			dg := float64(c1.G) - float64(c2.G)
			db := float64(c1.B) - float64(c2.B)
			da := float64(c1.A) - float64(c2.A)

			diff := (dr*dr + dg*dg + db*db + da*da) / (255.0 * 255.0 * 4.0)

			if diff > threshold {
				key := fmt.Sprintf("%d,%d", x, y)
				diffMap[key] = true
			}
		}
	}

	if len(diffMap) == 0 {
		return diffRegions
	}

	visited := make(map[string]bool)
	for key := range diffMap {
		if visited[key] {
			continue
		}

		region := d.floodFill(diffMap, visited, key)
		parts := splitKey(key)
		diffRegions = append(diffRegions, DiffRegion{
			X:          parts.x,
			Y:          parts.y,
			Width:      region.maxX - region.minX,
			Height:     region.maxY - region.minY,
			Score:      float64(region.count) / float64(len(diffMap)),
			PixelCount: region.count,
		})
	}

	return diffRegions
}

type point struct{ x, y int }
type regionInfo struct {
	minX, maxX int
	minY, maxY int
	count      int
}

func (d *pixelDiffer) floodFill(diffMap, visited map[string]bool, startKey string) regionInfo {
	parts := splitKey(startKey)
	queue := []point{{x: parts.x, y: parts.y}}
	info := regionInfo{minX: parts.x, maxX: parts.x, minY: parts.y, maxY: parts.y}

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		key := fmt.Sprintf("%d,%d", p.x, p.y)

		if visited[key] {
			continue
		}
		if !diffMap[key] {
			continue
		}

		visited[key] = true
		info.count++

		if p.x < info.minX {
			info.minX = p.x
		}
		if p.x > info.maxX {
			info.maxX = p.x
		}
		if p.y < info.minY {
			info.minY = p.y
		}
		if p.y > info.maxY {
			info.maxY = p.y
		}

		queue = append(queue,
			point{x: p.x + 1, y: p.y},
			point{x: p.x - 1, y: p.y},
			point{x: p.x, y: p.y + 1},
			point{x: p.x, y: p.y - 1},
		)
	}

	return info
}

func splitKey(key string) point {
	var x, y int
	fmt.Sscanf(key, "%d,%d", &x, &y)
	return point{x: x, y: y}
}

func (d *pixelDiffer) renderDiff(m1, m2 *image.RGBA, diffs []DiffRegion, outputDir string) (string, error) {
	img := image.NewRGBA(m1.Bounds())

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.SetRGBA(x, y, m2.RGBAAt(x, y))
		}
	}

	red := color.RGBA{R: 255, G: 0, B: 0, A: 128}
	for _, r := range diffs {
		rect := image.Rect(r.X, r.Y, r.X+r.Width, r.Y+r.Height)
		for py := rect.Min.Y; py < rect.Max.Y; py++ {
			for px := rect.Min.X; px < rect.Max.X; px++ {
				if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
					img.SetRGBA(px, py, red)
				}
			}
		}
	}

	os.MkdirAll(outputDir, 0755)
	filename := fmt.Sprintf("diff_%d.png", time.Now().UnixNano())
	outputPath := filepath.Join(outputDir, filename)

	f, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}

func decodeImage(data []byte) (*image.RGBA, error) {
	m, _, err := image.Decode(&stdImageReader{data: data})
	if err != nil {
		return nil, err
	}

	bounds := m.Bounds()
	rgba := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.SetRGBA(x, y, m.At(x, y).(color.RGBA))
		}
	}

	return rgba, nil
}

type stdImageReader struct {
	data []byte
	pos  int
}

func (r *stdImageReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, fmt.Errorf("EOF")
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}
