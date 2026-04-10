package input

var KeyMap = map[string]int{
	"home":        3,
	"back":        4,
	"left":        21,
	"right":       22,
	"up":          19,
	"down":        20,
	"center":      23,
	"menu":        82,
	"search":      84,
	"enter":       66,
	"delete":      67,
	"del":         67,
	"recent":      187,
	"volume_up":   24,
	"volume_down": 25,
	"volume_mute": 164,
	"camera":      27,
	"power":       26,
}

type Keys struct {
	pressKey     func(key string) error
	pressKeyCode func(code, meta int) error
}

func NewKeys(pressKey func(key string) error, pressKeyCode func(code, meta int) error) *Keys {
	return &Keys{
		pressKey:     pressKey,
		pressKeyCode: pressKeyCode,
	}
}

func (k *Keys) Press(key string) error {
	return k.pressKey(key)
}

func (k *Keys) PressKeyCode(code, meta int) error {
	return k.pressKeyCode(code, meta)
}

func (k *Keys) Home() error {
	return k.Press("home")
}

func (k *Keys) Back() error {
	return k.Press("back")
}

func (k *Keys) Enter() error {
	return k.Press("enter")
}

func (k *Keys) Delete() error {
	return k.Press("delete")
}

func (k *Keys) VolumeUp() error {
	return k.Press("volume_up")
}

func (k *Keys) VolumeDown() error {
	return k.Press("volume_down")
}

func (k *Keys) Power() error {
	return k.Press("power")
}
