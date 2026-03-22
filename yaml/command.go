package yaml

type Command struct {
	Name             string             `yaml:"name,omitempty"`
	TapOn            *TapCommand        `yaml:"tapOn,omitempty"`
	Tap              *PointCommand      `yaml:"tap,omitempty"`
	DoubleTap        *TapCommand        `yaml:"doubleTap,omitempty"`
	LongPress        *LongPressCommand  `yaml:"longPress,omitempty"`
	Swipe            *SwipeCommand      `yaml:"swipe,omitempty"`
	SwipeUp          *struct{}          `yaml:"swipeUp,omitempty"`
	SwipeDown        *struct{}          `yaml:"swipeDown,omitempty"`
	SwipeLeft        *struct{}          `yaml:"swipeLeft,omitempty"`
	SwipeRight       *struct{}          `yaml:"swipeRight,omitempty"`
	InputText        *InputCommand      `yaml:"inputText,omitempty"`
	Launch           string             `yaml:"launch,omitempty"`
	Terminate        string             `yaml:"terminate,omitempty"`
	Install          string             `yaml:"install,omitempty"`
	Uninstall        string             `yaml:"uninstall,omitempty"`
	WaitFor          *WaitCommand       `yaml:"waitFor,omitempty"`
	WaitForGone      *WaitCommand       `yaml:"waitForGone,omitempty"`
	Wait             int                `yaml:"wait,omitempty"`
	AssertVisible    *ElementQuery      `yaml:"assertVisible,omitempty"`
	AssertNotVisible *ElementQuery      `yaml:"assertNotVisible,omitempty"`
	AssertTrue       *string            `yaml:"assertTrue,omitempty"`
	Screenshot       *ScreenshotCommand `yaml:"screenshot,omitempty"`
	RunFlow          *RunFlowCommand    `yaml:"runFlow,omitempty"`
	If               *IfCommand         `yaml:"if,omitempty"`
	Foreach          *ForeachCommand    `yaml:"foreach,omitempty"`
	While            *WhileCommand      `yaml:"while,omitempty"`
	EvalScript       *EvalScriptCommand `yaml:"evalScript,omitempty"`
	SetVariable      *SetVarCommand     `yaml:"setVariable,omitempty"`
	Log              string             `yaml:"log,omitempty"`
	Comment          string             `yaml:"comment,omitempty"`
}

type TapCommand struct {
	Text     string `yaml:"text,omitempty"`
	ID       string `yaml:"id,omitempty"`
	XPath    string `yaml:"xpath,omitempty"`
	Index    int    `yaml:"index,omitempty"`
	Optional bool   `yaml:"optional,omitempty"`
	Timeout  string `yaml:"timeout,omitempty"`
}

type PointCommand struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
}

type LongPressCommand struct {
	Text     string `yaml:"text,omitempty"`
	X        int    `yaml:"x,omitempty"`
	Y        int    `yaml:"y,omitempty"`
	Duration int    `yaml:"duration"`
}

type SwipeCommand struct {
	StartX, StartY int `yaml:"startX,omitempty"`
	EndX, EndY     int `yaml:"endX,omitempty"`
	Duration       int `yaml:"duration,omitempty"`
}

type InputCommand struct {
	Text       string        `yaml:"text"`
	Element    *ElementQuery `yaml:"element,omitempty"`
	Secure     bool          `yaml:"secure,omitempty"`
	PressEnter bool          `yaml:"pressEnter,omitempty"`
}

type ElementQuery struct {
	Text  string `yaml:"text,omitempty"`
	ID    string `yaml:"id,omitempty"`
	XPath string `yaml:"xpath,omitempty"`
	Index int    `yaml:"index,omitempty"`
}

type WaitCommand struct {
	Element  *ElementQuery `yaml:"element"`
	Timeout  string        `yaml:"timeout"`
	Optional bool          `yaml:"optional,omitempty"`
}

type ScreenshotCommand struct {
	Name string `yaml:"name"`
}

type RunFlowCommand struct {
	Name   string            `yaml:"name"`
	Params map[string]string `yaml:"params,omitempty"`
}

type IfCommand struct {
	Condition string    `yaml:"condition"`
	Then      []Command `yaml:"then"`
	Else      []Command `yaml:"else,omitempty"`
}

type ForeachCommand struct {
	Variable string    `yaml:"variable"`
	In       string    `yaml:"in"`
	Do       []Command `yaml:"do"`
}

type WhileCommand struct {
	Condition string    `yaml:"condition"`
	MaxIter   int       `yaml:"maxIterations,omitempty"`
	Do        []Command `yaml:"do"`
}

type EvalScriptCommand struct {
	Lang     string   `yaml:"lang"`
	Script   string   `yaml:"script,omitempty"`
	Source   string   `yaml:"source,omitempty"`
	Function string   `yaml:"function,omitempty"`
	Args     []string `yaml:"args,omitempty"`
	SaveTo   string   `yaml:"saveTo,omitempty"`
}

type SetVarCommand struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Flow struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description,omitempty"`
	Platform    string            `yaml:"platform,omitempty"`
	Params      map[string]string `yaml:"params,omitempty"`
	Timeout     string            `yaml:"timeout,omitempty"`
	Steps       []Command         `yaml:"steps"`
}

type TestSuite struct {
	Config    *Config           `yaml:"config,omitempty"`
	Env       map[string]string `yaml:"env,omitempty"`
	Import    []string          `yaml:"import,omitempty"`
	Variables map[string]string `yaml:"variables,omitempty"`
	Flows     []Flow            `yaml:"flows,omitempty"`
	Tests     []TestCase        `yaml:"tests,omitempty"`
}

type Config struct {
	AppID   string `yaml:"appId"`
	Timeout string `yaml:"timeout"`
	Retry   int    `yaml:"retry"`
}

type TestCase struct {
	Name   string            `yaml:"name"`
	Flow   string            `yaml:"flow,omitempty"`
	Params map[string]string `yaml:"params,omitempty"`
	Steps  []Command         `yaml:"steps,omitempty"`
}
