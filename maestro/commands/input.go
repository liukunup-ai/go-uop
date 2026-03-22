package commands

import (
	"github.com/liukunup/go-uop/internal/action"
	"github.com/liukunup/go-uop/internal/selector"
	"github.com/liukunup/go-uop/maestro"
)

type InputTextTranslator struct{}

func NewInputTextTranslator() *InputTextTranslator {
	return &InputTextTranslator{}
}

// TranslateShorthand translates a simple text string to a SendKeysAction
func (t *InputTextTranslator) TranslateShorthand(text string) *action.SendKeysAction {
	return &action.SendKeysAction{
		Text: text,
	}
}

// TranslateExtended translates an InputTextCommand to a SendKeysAction
func (t *InputTextTranslator) TranslateExtended(cmd *maestro.InputTextCommand) *action.SendKeysAction {
	var elem *selector.Selector

	if cmd.ID != "" {
		elem = selector.ByID(cmd.ID)
	}

	return &action.SendKeysAction{
		Text:              cmd.Text,
		Element:           elem,
		Enter:             cmd.Enter,
		ClearExistingText: cmd.ClearExistingText,
	}
}
