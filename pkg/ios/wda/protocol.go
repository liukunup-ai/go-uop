package wda

const (
	EndpointStatus       = "/status"
	EndpointSession      = "/session"
	EndpointScreenshot   = "/screenshot"
	EndpointSource       = "/wda/source"
	EndpointElement      = "/wda/element/active"
	EndpointTap          = "/wda/tap/0/%d/%d"
	EndpointSwipe        = "/wda/performActions"
	EndpointKeys         = "/wda/keys"
	EndpointAppLaunch    = "/wda/apps/launch"
	EndpointAppTerminate = "/wda/apps/terminate/%s"
	EndpointAlert        = "/wda/alert/%s"
)

type AlertAction string

const (
	AlertAccept  AlertAction = "accept"
	AlertDismiss AlertAction = "dismiss"
	AlertText    AlertAction = "text"
)

const (
	KeyCodeHome       = 3
	KeyCodeBack       = 4
	KeyCodeEnter      = 66
	KeyCodeVolumeUp   = 24
	KeyCodeVolumeDown = 25
	KeyCodePower      = 26
)
