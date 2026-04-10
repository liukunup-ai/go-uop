package wda

const (
	// W3C Standard Endpoints
	EndpointStatus     = "/status"
	EndpointSession    = "/session"
	EndpointScreenshot = "/screenshot"
	EndpointSource     = "/source"
	EndpointAlert      = "/alert"

	// Element endpoints
	EndpointElement   = "/element"    // POST - Find single element
	EndpointElements  = "/elements"   // POST - Find multiple elements
	EndpointElementID = "/element/%s" // GET/POST/DELETE on element

	// App endpoints (WDA extension)
	EndpointAppLaunch    = "/app/launch"
	EndpointAppTerminate = "/app/terminate"
	EndpointAppActivate  = "/app/activate"

	// Touch/Action endpoints (W3C Actions)
	EndpointActions      = "/actions"
	EndpointTouchActions = "/touch/multi/perform"

	// Window/Screen endpoints
	EndpointWindowRect = "/window/rect"

	// WDA Legacy endpoints (backward compatibility)
	EndpointWDATap       = "/wda/tap/0/%d/%d"
	EndpointWDASource    = "/wda/source"
	EndpointWDAKeys      = "/wda/keys"
	EndpointWDAAppLaunch = "/wda/apps/launch"
)

// W3C By strategies
const (
	StrategyCSSSelector     = "css selector"
	StrategyLinkText        = "link text"
	StrategyPartialLinkText = "partial link text"
	StrategyTagName         = "tag name"
	StrategyXPath           = "xpath"
	StrategyClassName       = "class name"
	StrategyID              = "id"
	StrategyName            = "name"
)

// Alert actions
type AlertAction string

const (
	AlertAccept  AlertAction = "accept"
	AlertDismiss AlertAction = "dismiss"
	AlertText    AlertAction = "text"
)

// Key codes (Android compatibility)
const (
	KeyCodeHome       = 3
	KeyCodeBack       = 4
	KeyCodeEnter      = 66
	KeyCodeVolumeUp   = 24
	KeyCodeVolumeDown = 25
	KeyCodePower      = 26
)
