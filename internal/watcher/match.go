package watcher

import (
	"context"
	"os"
	"regexp"

	"github.com/liukunup/go-uop/core"
	"github.com/liukunup/go-uop/pkg/vision"
)

type MatchCondition interface {
	Match(ctx context.Context, device core.Device) (bool, error)
}

type ImageMatch struct {
	TemplatePath string
	Threshold    float64
}

func NewImageMatch(templatePath string, threshold float64) *ImageMatch {
	return &ImageMatch{
		TemplatePath: templatePath,
		Threshold:    threshold,
	}
}

func (m *ImageMatch) Match(ctx context.Context, device core.Device) (bool, error) {
	screenshot, err := device.Screenshot()
	if err != nil {
		return false, err
	}

	template, err := os.ReadFile(m.TemplatePath)
	if err != nil {
		return false, err
	}

	matcher, err := vision.NewMatcher("template", vision.WithThreshold(m.Threshold))
	if err != nil {
		return false, err
	}

	results, err := matcher.Find(screenshot, template)
	if err != nil {
		return false, err
	}

	if len(results) == 0 {
		return false, nil
	}

	return results[0].Score >= m.Threshold, nil
}

type TextMatch struct {
	Text string
}

func NewTextMatch(text string) *TextMatch {
	return &TextMatch{Text: text}
}

func (m *TextMatch) Match(ctx context.Context, device core.Device) (bool, error) {
	alertText, err := device.GetAlertText()
	if err != nil {
		return false, nil
	}
	return alertText == m.Text, nil
}

type RegexMatch struct {
	Pattern string
}

func NewRegexMatch(pattern string) *RegexMatch {
	return &RegexMatch{Pattern: pattern}
}

func (m *RegexMatch) Match(ctx context.Context, device core.Device) (bool, error) {
	alertText, err := device.GetAlertText()
	if err != nil {
		return false, nil
	}
	matched, _ := regexp.MatchString(m.Pattern, alertText)
	return matched, nil
}

type CompoundMatch struct {
	Operator   string
	Conditions []MatchCondition
}

func NewCompoundMatch(operator string, conditions []MatchCondition) *CompoundMatch {
	return &CompoundMatch{
		Operator:   operator,
		Conditions: conditions,
	}
}

func (m *CompoundMatch) Match(ctx context.Context, device core.Device) (bool, error) {
	for _, cond := range m.Conditions {
		matched, err := cond.Match(ctx, device)
		if err != nil {
			return false, err
		}
		if m.Operator == "or" && matched {
			return true, nil
		}
		if m.Operator == "and" && !matched {
			return false, nil
		}
	}
	if m.Operator == "and" {
		return true, nil
	}
	return false, nil
}
