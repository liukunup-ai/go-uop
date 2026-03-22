package commands

import (
	"time"

	"github.com/liukunup/go-uop/internal/action"
	"github.com/liukunup/go-uop/maestro"
)

const (
	// DefaultTimeout is the default timeout for wait operations (15 seconds)
	DefaultTimeout = 15 * time.Second
)

type NavigationTranslator struct{}

func NewNavigationTranslator() *NavigationTranslator {
	return &NavigationTranslator{}
}

// TranslateWaitForAnimationToEnd translates WaitForAnimationEnd to WaitAction
// If timeout is specified, uses that duration in milliseconds
// Otherwise uses DefaultTimeout
func (t *NavigationTranslator) TranslateWaitForAnimationToEnd(cmd *maestro.WaitForAnimationEnd) *action.WaitAction {
	if cmd == nil {
		return &action.WaitAction{
			Duration: DefaultTimeout,
		}
	}

	if cmd.Timeout > 0 {
		return &action.WaitAction{
			Duration: time.Duration(cmd.Timeout) * time.Millisecond,
		}
	}

	return &action.WaitAction{
		Duration: DefaultTimeout,
	}
}

// TranslatePressKey translates PressKeyCommand to PressKeyAction
// Maps common key names to Android key codes
func (t *NavigationTranslator) TranslatePressKey(cmd *maestro.PressKeyCommand) *action.PressKeyAction {
	if cmd == nil {
		return nil
	}

	return &action.PressKeyAction{
		KeyCode: MapKeyToCode(cmd.Key),
	}
}

// TranslateBack translates the Back command to PressKeyAction with back key
func (t *NavigationTranslator) TranslateBack() *action.PressKeyAction {
	return &action.PressKeyAction{
		KeyCode: MapKeyToCode("back"),
	}
}

// TranslatePressHome translates the PressHome command to PressKeyAction with home key
func (t *NavigationTranslator) TranslatePressHome() *action.PressKeyAction {
	return &action.PressKeyAction{
		KeyCode: MapKeyToCode("home"),
	}
}

// TranslatePressRecentApps translates the PressRecentApps command to PressKeyAction
func (t *NavigationTranslator) TranslatePressRecentApps() *action.PressKeyAction {
	return &action.PressKeyAction{
		KeyCode: MapKeyToCode("recent_apps"),
	}
}

// MapKeyToCode maps a key name to its Android key code
// Reference: https://developer.android.com/reference/android/view/KeyEvent
func MapKeyToCode(key string) int {
	switch key {
	case "home":
		return 3 // KEYCODE_HOME
	case "back":
		return 4 // KEYCODE_BACK
	case "enter":
		return 66 // KEYCODE_ENTER
	case "recent_apps", "recent":
		return 187 // KEYCODE_APP_SWITCH
	case "volume_up":
		return 24 // KEYCODE_VOLUME_UP
	case "volume_down":
		return 25 // KEYCODE_VOLUME_DOWN
	case "power":
		return 26 // KEYCODE_POWER
	case "delete":
		return 67 // KEYCODE_DEL
	case "tab":
		return 61 // KEYCODE_TAB
	case "escape", "esc":
		return 111 // KEYCODE_ESCAPE
	case "space":
		return 62 // KEYCODE_SPACE
	case "menu":
		return 82 // KEYCODE_MENU
	default:
		return 0
	}
}
