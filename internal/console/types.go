package console

// Device represents a mobile device
type Device struct {
	ID       string `json:"id"`
	Platform string `json:"platform"` // "ios" or "android"
	Name     string `json:"name"`
	Serial   string `json:"serial"`
	Status   string `json:"status"` // "available", "connected", "error"
	Model    string `json:"model,omitempty"`
	Address  string `json:"address,omitempty"`     // iOS WDA address
	PkgName  string `json:"packageName,omitempty"` // Android package name
}

// CommandRecord represents a command execution record
type CommandRecord struct {
	ID        string                 `json:"id"`
	Timestamp string                 `json:"timestamp"`
	Type      string                 `json:"command"`
	Params    map[string]interface{} `json:"params"`
	Success   bool                   `json:"success"`
	Output    string                 `json:"output,omitempty"`
	Duration  string                 `json:"duration"`
}

// CommandRequest represents a command request
type CommandRequest struct {
	Command string                 `json:"command"`
	Params  map[string]interface{} `json:"params"`
}
