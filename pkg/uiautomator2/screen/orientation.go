package screen

import (
	"os"
)

type Screenshot struct {
	screenshot func(format string) ([]byte, error)
}

func NewScreenshot(screenshot func(format string) ([]byte, error)) *Screenshot {
	return &Screenshot{screenshot: screenshot}
}

func (s *Screenshot) Capture(filename string) error {
	data, err := s.screenshot("raw")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (s *Screenshot) CaptureRaw() ([]byte, error) {
	return s.screenshot("raw")
}

func (s *Screenshot) CapturePillow() ([]byte, error) {
	return s.screenshot("pillow")
}

func (s *Screenshot) CaptureOpenCV() ([]byte, error) {
	return s.screenshot("opencv")
}

func (s *Screenshot) CaptureBase64() (string, error) {
	data, err := s.screenshot("raw")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type Orientation struct {
	getOrientation func() (string, error)
	setOrientation func(o string) error
	freezeRotation func(freeze bool) error
}

func NewOrientation(
	getOrientation func() (string, error),
	setOrientation func(o string) error,
	freezeRotation func(freeze bool) error,
) *Orientation {
	return &Orientation{
		getOrientation: getOrientation,
		setOrientation: setOrientation,
		freezeRotation: freezeRotation,
	}
}

func (o *Orientation) Get() (string, error) {
	return o.getOrientation()
}

func (o *Orientation) Set(orientation string) error {
	return o.setOrientation(orientation)
}

func (o *Orientation) Freeze() error {
	return o.freezeRotation(true)
}

func (o *Orientation) Unfreeze() error {
	return o.freezeRotation(false)
}

func (o *Orientation) Natural() error {
	return o.setOrientation("natural")
}

func (o *Orientation) Left() error {
	return o.setOrientation("left")
}

func (o *Orientation) Right() error {
	return o.setOrientation("right")
}
