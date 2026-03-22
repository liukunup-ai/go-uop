package vision

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

type MatchResult struct {
	X      int
	Y      int
	Width  int
	Height int
	Score  float64
}

type TemplateMatcher struct {
	screenshot []byte
	template   []byte
}

func NewTemplateMatcher(screenshot []byte) *TemplateMatcher {
	return &TemplateMatcher{
		screenshot: screenshot,
	}
}

func (tm *TemplateMatcher) SetTemplate(templatePath string) error {
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("read template: %w", err)
	}
	tm.template = data
	return nil
}

func (tm *TemplateMatcher) SetTemplateBytes(data []byte) {
	tm.template = data
}

func (tm *TemplateMatcher) FindBestMatch() (*MatchResult, error) {
	if tm.screenshot == nil || tm.template == nil {
		return nil, fmt.Errorf("screenshot or template not set")
	}

	screenImg, err := png.Decode(newBytesReader(tm.screenshot))
	if err != nil {
		return nil, fmt.Errorf("decode screenshot: %w", err)
	}

	templateImg, err := png.Decode(newBytesReader(tm.template))
	if err != nil {
		return nil, fmt.Errorf("decode template: %w", err)
	}

	bounds := screenImg.Bounds()
	templateBounds := templateImg.Bounds()

	if bounds.Dx() < templateBounds.Dx() || bounds.Dy() < templateBounds.Dy() {
		return nil, fmt.Errorf("template larger than screenshot")
	}

	result := simpleMatch(screenImg, templateImg)
	return result, nil
}

func simpleMatch(screen, template image.Image) *MatchResult {
	return &MatchResult{
		X:      0,
		Y:      0,
		Width:  template.Bounds().Dx(),
		Height: template.Bounds().Dy(),
		Score:  0.0,
	}
}

type byteReader struct {
	data []byte
	pos  int
}

func newBytesReader(data []byte) *byteReader {
	return &byteReader{data: data}
}

func (r *byteReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, fmt.Errorf("end of data")
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

func SaveMatchResult(img image.Image, result *MatchResult, outputPath string) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	return png.Encode(f, img)
}

func LoadImageFromFile(path string) (image.Image, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	return LoadImageFromBytes(data)
}

func LoadImageFromBytes(data []byte) (image.Image, error) {
	return png.Decode(newBytesReader(data))
}
