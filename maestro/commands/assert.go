package commands

import (
	"strconv"
	"time"

	"github.com/liukunup/go-uop/internal/action"
	"github.com/liukunup/go-uop/internal/selector"
	"github.com/liukunup/go-uop/maestro"
)

type AssertTranslator struct{}

func NewAssertTranslator() *AssertTranslator {
	return &AssertTranslator{}
}

func (t *AssertTranslator) TranslateAssertVisibleShorthand(text string) *action.AssertAction {
	return &action.AssertAction{
		Element:   selector.ByText(text),
		MustExist: true,
	}
}

func (t *AssertTranslator) TranslateAssertNotVisibleShorthand(text string) *action.AssertAction {
	return &action.AssertAction{
		Element:   selector.ByText(text),
		MustExist: false,
	}
}

func (t *AssertTranslator) TranslateAssertVisible(cmd *maestro.ElementSelector) *action.AssertAction {
	var elem *selector.Selector

	if cmd.ID != "" {
		elem = selector.ByID(cmd.ID)
	} else if cmd.Text != "" {
		elem = selector.ByText(cmd.Text)
	}

	if cmd.Index > 0 && elem != nil {
		elem.SetIndex(cmd.Index)
	}

	return &action.AssertAction{
		Element:   elem,
		MustExist: true,
	}
}

func (t *AssertTranslator) TranslateAssertNotVisible(cmd *maestro.ElementSelector) *action.AssertAction {
	var elem *selector.Selector

	if cmd.ID != "" {
		elem = selector.ByID(cmd.ID)
	} else if cmd.Text != "" {
		elem = selector.ByText(cmd.Text)
	}

	if cmd.Index > 0 && elem != nil {
		elem.SetIndex(cmd.Index)
	}

	return &action.AssertAction{
		Element:   elem,
		MustExist: false,
	}
}

type AssertWithTimeout struct {
	*action.AssertAction
	Timeout time.Duration
}

func (t *AssertTranslator) TranslateAssertVisibleWithTimeout(cmd *maestro.TapOnCommand) *action.AssertAction {
	var elem *selector.Selector

	if cmd.ID != "" {
		elem = selector.ByID(cmd.ID)
	} else if cmd.Text != "" {
		elem = selector.ByText(cmd.Text)
	}

	if cmd.Index > 0 && elem != nil {
		elem.SetIndex(cmd.Index)
	}

	act := &action.AssertAction{
		Element:   elem,
		MustExist: true,
	}

	if cmd.Timeout != "" {
		if timeoutMs, err := strconv.ParseInt(cmd.Timeout, 10, 64); err == nil {
			act.Timeout = time.Duration(timeoutMs) * time.Millisecond
		}
	}

	return act
}

func (t *AssertTranslator) TranslateAssertNotVisibleWithTimeout(cmd *maestro.TapOnCommand) *action.AssertAction {
	var elem *selector.Selector

	if cmd.ID != "" {
		elem = selector.ByID(cmd.ID)
	} else if cmd.Text != "" {
		elem = selector.ByText(cmd.Text)
	}

	if cmd.Index > 0 && elem != nil {
		elem.SetIndex(cmd.Index)
	}

	act := &action.AssertAction{
		Element:   elem,
		MustExist: false,
	}

	if cmd.Timeout != "" {
		if timeoutMs, err := strconv.ParseInt(cmd.Timeout, 10, 64); err == nil {
			act.Timeout = time.Duration(timeoutMs) * time.Millisecond
		}
	}

	return act
}
