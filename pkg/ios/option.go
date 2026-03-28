package ios

type config struct {
	address     string
	udid        string
	skipSession bool
}

type Option func(*config)

func WithAddress(addr string) Option {
	return func(c *config) {
		c.address = addr
	}
}

func WithUDID(udid string) Option {
	return func(c *config) {
		c.udid = udid
	}
}

func SkipSession() Option {
	return func(c *config) {
		c.skipSession = true
	}
}
