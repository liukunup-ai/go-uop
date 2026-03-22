package commands

import (
	"github.com/liukunup/go-uop/internal/action"
	"github.com/liukunup/go-uop/internal/selector"
	"github.com/liukunup/go-uop/maestro"
)

type TapOnTranslator struct{}

func NewTapOnTranslator() *TapOnTranslator {
	return &TapOnTranslator{}
}

func (t *TapOnTranslator) TranslateShorthand(text string) *action.TapAction {
	return &action.TapAction{
		Element: selector.ByText(text),
	}
}

func (t *TapOnTranslator) TranslateExtended(cmd *maestro.TapOnCommand) *action.TapAction {
	var elem *selector.Selector

	if cmd.ID != "" {
		elem = selector.ByID(cmd.ID)
	} else if cmd.Text != "" {
		elem = selector.ByText(cmd.Text)
	}

	if cmd.Index > 0 && elem != nil {
		elem.SetIndex(cmd.Index)
	}

	return &action.TapAction{
		Element: elem,
	}
}

func (t *TapOnTranslator) TranslatePoint(x, y int) *action.TapAction {
	return &action.TapAction{
		X: x,
		Y: y,
	}
}
