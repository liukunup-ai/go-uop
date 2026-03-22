package commands

import (
	"time"

	"github.com/liukunup/go-uop/internal/action"
	"github.com/liukunup/go-uop/maestro"
)

const (
	screenWidth  = 100
	screenHeight = 100
)

type SwipeTranslator struct{}

func NewSwipeTranslator() *SwipeTranslator {
	return &SwipeTranslator{}
}

func (s *SwipeTranslator) TranslateDirection(direction maestro.SwipeDirection) *action.SwipeAction {
	startX, startY, endX, endY := calculateCoordinates(direction)
	return &action.SwipeAction{
		StartX:   startX,
		StartY:   startY,
		EndX:     endX,
		EndY:     endY,
		Duration: 300 * time.Millisecond,
	}
}

func (s *SwipeTranslator) TranslateDirectionWithDuration(direction maestro.SwipeDirection, durationMs int) *action.SwipeAction {
	startX, startY, endX, endY := calculateCoordinates(direction)
	return &action.SwipeAction{
		StartX:   startX,
		StartY:   startY,
		EndX:     endX,
		EndY:     endY,
		Duration: time.Duration(durationMs) * time.Millisecond,
	}
}

func (s *SwipeTranslator) TranslateExtended(cmd *maestro.SwipeCommand) *action.SwipeAction {
	var startX, startY, endX, endY int
	var duration time.Duration

	if cmd.Direction != "" {
		startX, startY, endX, endY = calculateCoordinates(cmd.Direction)
	} else {
		startX, startY = cmd.StartX, cmd.StartY
		endX, endY = cmd.EndX, cmd.EndY
	}

	if cmd.Duration > 0 {
		duration = time.Duration(cmd.Duration) * time.Millisecond
	} else {
		duration = 300 * time.Millisecond
	}

	return &action.SwipeAction{
		StartX:   startX,
		StartY:   startY,
		EndX:     endX,
		EndY:     endY,
		Duration: duration,
	}
}

func (s *SwipeTranslator) TranslateSwipeUntilVisible(direction maestro.SwipeDirection, elem *maestro.ElementSelector) *action.SwipeAction {
	startX, startY, endX, endY := calculateCoordinates(direction)
	return &action.SwipeAction{
		StartX:   startX,
		StartY:   startY,
		EndX:     endX,
		EndY:     endY,
		Duration: 300 * time.Millisecond,
	}
}

func calculateCoordinates(direction maestro.SwipeDirection) (int, int, int, int) {
	centerX := screenWidth / 2
	startY := (screenHeight * 80) / 100
	endY := (screenHeight * 20) / 100

	switch direction {
	case maestro.SwipeUp:
		return centerX, startY, centerX, endY
	case maestro.SwipeDown:
		return centerX, endY, centerX, startY
	case maestro.SwipeLeft:
		return (screenWidth * 70) / 100, endY, (screenWidth * 30) / 100, endY
	case maestro.SwipeRight:
		return (screenWidth * 30) / 100, endY, (screenWidth * 70) / 100, endY
	default:
		return centerX, startY, centerX, endY
	}
}
