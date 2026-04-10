package ios

// config holds iOS device configuration options
type config struct {
	udid        string
	address     string // For WDA connection
	skipSession bool
	bundleID    string
}

// Option is a function that modifies config
type Option func(*config)

// WithUDID sets the device UDID for go-ios connection
func WithUDID(udid string) Option {
	return func(c *config) {
		c.udid = udid
	}
}

// WithAddress sets the WDA server address (default: http://localhost:8100)
func WithAddress(addr string) Option {
	return func(c *config) {
		c.address = addr
	}
}

// SkipSession creates client without starting a WDA session
func SkipSession() Option {
	return func(c *config) {
		c.skipSession = true
	}
}
