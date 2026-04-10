package uiautomator2

// DeviceInfo represents basic device information (matches d.info)
type DeviceInfo struct {
	CurrentPackageName string `json:"currentPackageName"`
	DisplayHeight      int    `json:"displayHeight"`
	DisplayRotation    int    `json:"displayRotation"`
	DisplaySizeDpX     int    `json:"displaySizeDpX"`
	DisplaySizeDpY     int    `json:"displaySizeDpY"`
	DisplayWidth       int    `json:"displayWidth"`
	ProductName        string `json:"productName"`
	ScreenOn           bool   `json:"screenOn"`
	SdkInt             int    `json:"sdkInt"`
	NaturalOrientation bool   `json:"naturalOrientation"`
}

// DeviceDetail represents detailed device info (matches d.device_info)
type DeviceDetail struct {
	Arch    string `json:"arch"`
	Brand   string `json:"brand"`
	Model   string `json:"model"`
	Sdk     int    `json:"sdk"`
	Serial  string `json:"serial"`
	Version int    `json:"version"`
}

// AppInfo represents application information
type AppInfo struct {
	MainActivity string `json:"mainActivity"`
	Label        string `json:"label"`
	VersionName  string `json:"versionName"`
	VersionCode  int    `json:"versionCode"`
	Size         int64  `json:"size"`
}

// Selector represents UI selector criteria
type Selector struct {
	Text                  string `json:"text,omitempty"`
	TextContains          string `json:"textContains,omitempty"`
	TextMatches           string `json:"textMatches,omitempty"`
	TextStartsWith        string `json:"textStartsWith,omitempty"`
	ClassName             string `json:"className,omitempty"`
	ClassNameMatches      string `json:"classNameMatches,omitempty"`
	Description           string `json:"description,omitempty"`
	DescriptionContains   string `json:"descriptionContains,omitempty"`
	DescriptionMatches    string `json:"descriptionMatches,omitempty"`
	DescriptionStartsWith string `json:"descriptionStartsWith,omitempty"`
	Checkable             bool   `json:"checkable,omitempty"`
	Checked               bool   `json:"checked,omitempty"`
	Clickable             bool   `json:"clickable,omitempty"`
	LongClickable         bool   `json:"longClickable,omitempty"`
	Scrollable            bool   `json:"scrollable,omitempty"`
	Enabled               bool   `json:"enabled,omitempty"`
	Focusable             bool   `json:"focusable,omitempty"`
	Focused               bool   `json:"focused,omitempty"`
	Selected              bool   `json:"selected,omitempty"`
	PackageName           string `json:"packageName,omitempty"`
	PackageNameMatches    string `json:"packageNameMatches,omitempty"`
	ResourceId            string `json:"resourceId,omitempty"`
	ResourceIdMatches     string `json:"resourceIdMatches,omitempty"`
	Index                 int    `json:"index,omitempty"`
	Instance              int    `json:"instance,omitempty"`
}

// Point represents coordinates
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Bounds represents rectangular bounds
type Bounds struct {
	Left   int `json:"left"`
	Top    int `json:"top"`
	Right  int `json:"right"`
	Bottom int `json:"bottom"`
}

// ElementInfo represents UI element information
type ElementInfo struct {
	ContentDescription string `json:"contentDescription"`
	Checked            bool   `json:"checked"`
	Scrollable         bool   `json:"scrollable"`
	Text               string `json:"text"`
	PackageName        string `json:"packageName"`
	Selected           bool   `json:"selected"`
	Enabled            bool   `json:"enabled"`
	Bounds             Bounds `json:"bounds"`
	ClassName          string `json:"className"`
	Focused            bool   `json:"focused"`
	Focusable          bool   `json:"focusable"`
	Clickable          bool   `json:"clickable"`
	ChildCount         int    `json:"childCount"`
	LongClickable      bool   `json:"longClickable"`
	VisibleBounds      Bounds `json:"visibleBounds"`
	Checkable          bool   `json:"checkable"`
}

// SessionInfo represents current app session
type SessionInfo struct {
	Activity string `json:"activity"`
	Package  string `json:"package"`
	Pid      int    `json:"pid,omitempty"`
}

// SelectorResult represents selector query result
type SelectorResult struct {
	Elements []string `json:"elements"`
	Count    int      `json:"count"`
}
