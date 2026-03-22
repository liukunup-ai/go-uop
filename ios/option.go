package ios

type config struct {
	address string
}

type Option func(*config)

func WithAddress(addr string) Option {
	return func(c *config) {
		c.address = addr
	}
}
