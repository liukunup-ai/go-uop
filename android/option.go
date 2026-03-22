package android

type config struct {
	serial      string
	packageName string
}

type Option func(*config)

func WithSerial(serial string) Option {
	return func(c *config) {
		c.serial = serial
	}
}

func WithPackage(pkg string) Option {
	return func(c *config) {
		c.packageName = pkg
	}
}
