package android

type config struct {
	udid        string
	packageName string
}

type Option func(*config)

func WithUDID(udid string) Option {
	return func(c *config) {
		c.udid = udid
	}
}

func WithPackage(pkg string) Option {
	return func(c *config) {
		c.packageName = pkg
	}
}
