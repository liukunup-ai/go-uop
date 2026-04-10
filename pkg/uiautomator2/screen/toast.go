package screen

type Toast struct {
	lastToast  func() (string, error)
	clearToast func() error
}

func NewToast(lastToast func() (string, error), clearToast func() error) *Toast {
	return &Toast{
		lastToast:  lastToast,
		clearToast: clearToast,
	}
}

func (t *Toast) GetLast() (string, error) {
	return t.lastToast()
}

func (t *Toast) Clear() error {
	return t.clearToast()
}
