package commands

import (
	"fmt"
	"time"

	"github.com/liukunup/go-uop/internal/action"
	"github.com/liukunup/go-uop/maestro"
)

const DefaultScreenshotPath = "./maestro-screenshots"

type MediaCommandTranslator struct {
	screenshotPath string
}

func NewMediaCommandTranslator() *MediaCommandTranslator {
	return &MediaCommandTranslator{
		screenshotPath: DefaultScreenshotPath,
	}
}

func NewMediaCommandTranslatorWithPath(path string) *MediaCommandTranslator {
	return &MediaCommandTranslator{
		screenshotPath: path,
	}
}

func (m *MediaCommandTranslator) TranslateScrollShorthand() *action.SwipeAction {
	return m.TranslateScrollExtended(&maestro.ScrollCommand{
		Direction: maestro.ScrollUp,
	})
}

func (m *MediaCommandTranslator) TranslateScrollExtended(cmd *maestro.ScrollCommand) *action.SwipeAction {
	direction := maestro.ScrollUp
	if cmd.Direction != "" {
		direction = cmd.Direction
	}

	swipeTranslator := NewSwipeTranslator()
	swipeDir := m.convertScrollToSwipeDirection(direction)

	duration := 300 * time.Millisecond
	if cmd.Duration > 0 {
		duration = time.Duration(cmd.Duration) * time.Millisecond
	}

	act := swipeTranslator.TranslateDirection(swipeDir)
	act.Duration = duration

	return act
}

func (m *MediaCommandTranslator) TranslateTakeScreenshot(cmd *maestro.TakeScreenshotCommand) *action.ScreenshotAction {
	name := m.generateScreenshotName()
	if cmd != nil && cmd.Name != "" {
		name = cmd.Name
	}

	return &action.ScreenshotAction{
		Name: name,
		Path: m.screenshotPath,
	}
}

func (m *MediaCommandTranslator) TranslateTakeScreenshotShorthand() *action.ScreenshotAction {
	return m.TranslateTakeScreenshot(nil)
}

func (m *MediaCommandTranslator) convertScrollToSwipeDirection(direction maestro.ScrollDirection) maestro.SwipeDirection {
	switch direction {
	case maestro.ScrollUp:
		return maestro.SwipeUp
	case maestro.ScrollDown:
		return maestro.SwipeDown
	case maestro.ScrollLeft:
		return maestro.SwipeLeft
	case maestro.ScrollRight:
		return maestro.SwipeRight
	default:
		return maestro.SwipeUp
	}
}

func (m *MediaCommandTranslator) generateScreenshotName() string {
	return fmt.Sprintf("screenshot_%d", time.Now().UnixNano())
}
