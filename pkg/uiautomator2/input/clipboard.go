package input

type Clipboard struct {
	get func() (string, error)
	set func(text, label string) error
}

func NewClipboard(get func() (string, error), set func(text, label string) error) *Clipboard {
	return &Clipboard{
		get: get,
		set: set,
	}
}

func (c *Clipboard) Get() (string, error) {
	return c.get()
}

func (c *Clipboard) Set(text string, label string) error {
	return c.set(text, label)
}
