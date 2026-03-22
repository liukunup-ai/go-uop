package commands

import (
	"github.com/liukunup/go-uop/internal/action"
	"github.com/liukunup/go-uop/maestro"
)

type AppCommandTranslator struct{}

func NewAppCommandTranslator() *AppCommandTranslator {
	return &AppCommandTranslator{}
}

func (t *AppCommandTranslator) TranslateLaunchShorthand(appID string) *action.LaunchAction {
	return &action.LaunchAction{
		AppID: appID,
	}
}

func (t *AppCommandTranslator) TranslateLaunchExtended(cmd *maestro.LaunchAppCommand) *action.LaunchAction {
	act := &action.LaunchAction{
		AppID: cmd.AppID,
	}
	if cmd.ClearState {
		act.ClearState = true
	}
	return act
}

func (t *AppCommandTranslator) TranslateKill(appID string) *action.KillAction {
	return &action.KillAction{
		AppID: appID,
	}
}

func (t *AppCommandTranslator) TranslateStop(appID string) *action.StopAction {
	return &action.StopAction{
		AppID:    appID,
		Graceful: true,
	}
}

func (t *AppCommandTranslator) TranslateClearState(appID string) *action.ClearStateAction {
	return &action.ClearStateAction{
		AppID: appID,
	}
}
