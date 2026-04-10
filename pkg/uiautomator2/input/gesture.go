package input

type Gesture struct {
	tap         func(x, y int) error
	doubleClick func(x, y int) error
	longClick   func(x, y, duration int) error
	swipe       func(sx, sy, ex, ey, duration int) error
	drag        func(sx, sy, ex, ey, duration int) error
	touch       func(action string, x, y int) error
}

func NewGesture(
	tap func(x, y int) error,
	doubleClick func(x, y int) error,
	longClick func(x, y, duration int) error,
	swipe func(sx, sy, ex, ey, duration int) error,
	drag func(sx, sy, ex, ey, duration int) error,
	touch func(action string, x, y int) error,
) *Gesture {
	return &Gesture{
		tap:         tap,
		doubleClick: doubleClick,
		longClick:   longClick,
		swipe:       swipe,
		drag:        drag,
		touch:       touch,
	}
}

func (g *Gesture) Tap(x, y int) error {
	return g.tap(x, y)
}

func (g *Gesture) DoubleClick(x, y int) error {
	return g.doubleClick(x, y)
}

func (g *Gesture) LongClick(x, y int, duration int) error {
	return g.longClick(x, y, duration)
}

func (g *Gesture) Swipe(sx, sy, ex, ey, duration int) error {
	return g.swipe(sx, sy, ex, ey, duration)
}

func (g *Gesture) Drag(sx, sy, ex, ey, duration int) error {
	return g.drag(sx, sy, ex, ey, duration)
}

func (g *Gesture) TouchDown(x, y int) error {
	return g.touch("down", x, y)
}

func (g *Gesture) TouchMove(x, y int) error {
	return g.touch("move", x, y)
}

func (g *Gesture) TouchUp(x, y int) error {
	return g.touch("up", x, y)
}

func (g *Gesture) SwipeUp(duration int) error {
	return g.swipe(540, 1800, 540, 200, duration)
}

func (g *Gesture) SwipeDown(duration int) error {
	return g.swipe(540, 200, 540, 1800, duration)
}

func (g *Gesture) SwipeLeft(duration int) error {
	return g.swipe(540, 960, 100, 960, duration)
}

func (g *Gesture) SwipeRight(duration int) error {
	return g.swipe(100, 960, 540, 960, duration)
}
