package maestro

type MaestroFlow struct {
	AppID string           `yaml:"appId,omitempty"`
	Name  string           `yaml:"name,omitempty"`
	Tags  []string         `yaml:"tags,omitempty"`
	Steps []MaestroCommand `yaml:"steps,omitempty"`
}

type MaestroCommand struct {
	Name      string `yaml:"name,omitempty"`
	AppId     string `yaml:"appId,omitempty"`
	Launch    string `yaml:"launch,omitempty"`
	Terminate string `yaml:"terminate,omitempty"`
	Install   string `yaml:"install,omitempty"`
	Uninstall string `yaml:"uninstall,omitempty"`

	LaunchApp  *LaunchAppCommand `yaml:"launchApp,omitempty"`
	KillApp    string            `yaml:"killApp,omitempty"`
	StopApp    string            `yaml:"stopApp,omitempty"`
	ClearState string            `yaml:"clearState,omitempty"`

	TapOn     *TapOnCommand     `yaml:"tapOn,omitempty"`
	Tap       *PointCommand     `yaml:"tap,omitempty"`
	DoubleTap *TapOnCommand     `yaml:"doubleTap,omitempty"`
	LongPress *LongPressCommand `yaml:"longPress,omitempty"`

	Swipe      *SwipeCommand `yaml:"swipe,omitempty"`
	SwipeUp    *struct{}     `yaml:"swipeUp,omitempty"`
	SwipeDown  *struct{}     `yaml:"swipeDown,omitempty"`
	SwipeLeft  *struct{}     `yaml:"swipeLeft,omitempty"`
	SwipeRight *struct{}     `yaml:"swipeRight,omitempty"`

	InputText *InputTextCommand `yaml:"inputText,omitempty"`

	Wait        int             `yaml:"wait,omitempty"`
	WaitFor     *WaitForCommand `yaml:"waitFor,omitempty"`
	WaitForGone *WaitForCommand `yaml:"waitForGone,omitempty"`

	AssertVisible    *ElementSelector `yaml:"assertVisible,omitempty"`
	AssertNotVisible *ElementSelector `yaml:"assertNotVisible,omitempty"`
	AssertTrue       *string          `yaml:"assertTrue,omitempty"`
	AssertFalse      *string          `yaml:"assertFalse,omitempty"`

	Back                *struct{}            `yaml:"back,omitempty"`
	PressHome           *struct{}            `yaml:"pressHome,omitempty"`
	PressRecentApps     *struct{}            `yaml:"pressRecentApps,omitempty"`
	WaitForAnimationEnd *WaitForAnimationEnd `yaml:"waitForAnimationToEnd,omitempty"`
	PressKey            *PressKeyCommand     `yaml:"pressKey,omitempty"`

	Screenshot *ScreenshotCommand `yaml:"screenshot,omitempty"`
	Scroll     *ScrollCommand     `yaml:"scroll,omitempty"`

	RunFlow *RunFlowCommand `yaml:"runFlow,omitempty"`
	If      *IfCommand      `yaml:"if,omitempty"`
	Foreach *ForeachCommand `yaml:"foreach,omitempty"`
	While   *WhileCommand   `yaml:"while,omitempty"`
	Repeat  *RepeatCommand  `yaml:"repeat,omitempty"`

	EvalScript *EvalScriptCommand `yaml:"evalScript,omitempty"`

	SetVariable *SetVariableCommand `yaml:"setVariable,omitempty"`

	Log   string `yaml:"log,omitempty"`
	Print string `yaml:"print,omitempty"`

	Comment string `yaml:"comment,omitempty"`
}

type TapOnCommand struct {
	Text     string `yaml:"text,omitempty"`
	ID       string `yaml:"id,omitempty"`
	Index    int    `yaml:"index,omitempty"`
	Optional bool   `yaml:"optional,omitempty"`
	Timeout  string `yaml:"timeout,omitempty"`
}

type PointCommand struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
}

type LaunchAppCommand struct {
	AppID      string `yaml:"appId,omitempty"`
	ClearState bool   `yaml:"clearState,omitempty"`
}

type LongPressCommand struct {
	Text     string `yaml:"text,omitempty"`
	ID       string `yaml:"id,omitempty"`
	X        int    `yaml:"x,omitempty"`
	Y        int    `yaml:"y,omitempty"`
	Duration int    `yaml:"duration"`
}

// SwipeDirection defines the direction for swipe actions
type SwipeDirection string

const (
	SwipeUp    SwipeDirection = "UP"
	SwipeDown  SwipeDirection = "DOWN"
	SwipeLeft  SwipeDirection = "LEFT"
	SwipeRight SwipeDirection = "RIGHT"
)

type SwipeCommand struct {
	StartX, StartY    int              `yaml:"startX,omitempty"`
	EndX, EndY        int              `yaml:"endX,omitempty"`
	Duration          int              `yaml:"duration,omitempty"`
	Direction         SwipeDirection   `yaml:"direction,omitempty"`
	SwipeUntilVisible *ElementSelector `yaml:"swipeUntilVisible,omitempty"`
}

type InputTextCommand struct {
	Text              string `yaml:"text"`
	ID                string `yaml:"id,omitempty"`
	Optional          bool   `yaml:"optional,omitempty"`
	Enter             bool   `yaml:"enter,omitempty"`
	ClearExistingText bool   `yaml:"clearExistingText,omitempty"`
}

type WaitForCommand struct {
	Text    string `yaml:"text,omitempty"`
	ID      string `yaml:"id,omitempty"`
	Index   int    `yaml:"index,omitempty"`
	Timeout string `yaml:"timeout,omitempty"`
}

type ElementSelector struct {
	Text     string `yaml:"text,omitempty"`
	ID       string `yaml:"id,omitempty"`
	Index    int    `yaml:"index,omitempty"`
	Optional bool   `yaml:"optional,omitempty"`
}

type ScreenshotCommand struct {
	Name string `yaml:"name"`
}

type TakeScreenshotCommand struct {
	Name string `yaml:"name,omitempty"`
}

type ScrollDirection string

const (
	ScrollUp    ScrollDirection = "UP"
	ScrollDown  ScrollDirection = "DOWN"
	ScrollLeft  ScrollDirection = "LEFT"
	ScrollRight ScrollDirection = "RIGHT"
)

type ScrollCommand struct {
	Direction ScrollDirection `yaml:"direction,omitempty"`
	Duration  int             `yaml:"duration,omitempty"`
}

type RunFlowCommand struct {
	Name   string            `yaml:"name"`
	Path   string            `yaml:"path,omitempty"`
	Params map[string]string `yaml:"params,omitempty"`
}

type IfCommand struct {
	Condition string           `yaml:"condition"`
	Then      []MaestroCommand `yaml:"then"`
	Else      []MaestroCommand `yaml:"else,omitempty"`
}

type ForeachCommand struct {
	Variable string           `yaml:"variable"`
	In       string           `yaml:"in"`
	Do       []MaestroCommand `yaml:"do"`
}

type WhileCommand struct {
	Condition string           `yaml:"condition"`
	MaxIter   int              `yaml:"maxIterations,omitempty"`
	Do        []MaestroCommand `yaml:"do"`
}

type RepeatCommand struct {
	Times int              `yaml:"times"`
	Do    []MaestroCommand `yaml:"do"`
}

type EvalScriptCommand struct {
	Lang   string   `yaml:"lang"`
	Script string   `yaml:"script,omitempty"`
	Source string   `yaml:"source,omitempty"`
	Args   []string `yaml:"args,omitempty"`
	SaveTo string   `yaml:"saveTo,omitempty"`
}

type SetVariableCommand struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type TextSelector struct {
	Text     string `yaml:"text"`
	Optional bool   `yaml:"optional,omitempty"`
	Index    int    `yaml:"index,omitempty"`
}

type IDSelector struct {
	ID       string `yaml:"id"`
	Optional bool   `yaml:"optional,omitempty"`
	Index    int    `yaml:"index,omitempty"`
}

type IndexSelector struct {
	Index    int  `yaml:"index"`
	Optional bool `yaml:"optional,omitempty"`
}

type PointSelector struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
}

type WaitForAnimationEnd struct {
	Timeout int `yaml:"timeout,omitempty"`
}

type PressKeyCommand struct {
	Key string `yaml:"key"`
}
