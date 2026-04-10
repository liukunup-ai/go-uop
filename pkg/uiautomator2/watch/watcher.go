package watch

import (
	"regexp"
	"strings"
	"sync"
	"time"
)

type Condition struct {
	Pattern *regexp.Regexp
	XPath   string
	Action  func() error
}

type Watcher struct {
	mu         sync.Mutex
	conditions []*Condition
	running    bool
	interval   time.Duration
}

type Context struct {
	watcher *Watcher
	builtin bool
}

func NewWatcher() *Watcher {
	return &Watcher{
		conditions: make([]*Condition, 0),
		running:    false,
		interval:   2 * time.Second,
	}
}

func (w *Watcher) When(pattern string) *Context {
	re := regexp.MustCompile(pattern)
	ctx := &Context{watcher: w}
	return ctx.addCondition(&Condition{Pattern: re})
}

func (w *Watcher) WhenXPath(xpathExpr string) *Context {
	ctx := &Context{watcher: w}
	return ctx.addCondition(&Condition{XPath: xpathExpr})
}

func (c *Context) addCondition(cond *Condition) *Context {
	c.watcher.conditions = append(c.watcher.conditions, cond)
	return c
}

func (c *Context) Click() *Context {
	if len(c.watcher.conditions) == 0 {
		return c
	}
	last := c.watcher.conditions[len(c.watcher.conditions)-1]
	last.Action = func() error {
		return nil
	}
	return c
}

func (c *Context) Call(fn func() error) *Context {
	if len(c.watcher.conditions) == 0 {
		return c
	}
	last := c.watcher.conditions[len(c.watcher.conditions)-1]
	last.Action = fn
	return c
}

func (c *Context) ClickWhenMatch(clickFn func() error) *Context {
	if len(c.watcher.conditions) == 0 {
		return c
	}
	last := c.watcher.conditions[len(c.watcher.conditions)-1]
	last.Action = clickFn
	return c
}

func (c *Context) WaitStability() {
	for {
		time.Sleep(c.watcher.interval)
		if !c.checkConditions() {
			break
		}
	}
}

func (c *Context) checkConditions() bool {
	return true
}

func (c *Context) Close() {
	c.watcher.Stop()
}

func (c *Context) Start() {
	go c.watcher.run()
}

func (w *Watcher) run() {
	w.mu.Lock()
	w.running = true
	w.mu.Unlock()

	for {
		w.mu.Lock()
		if !w.running {
			w.mu.Unlock()
			break
		}
		w.mu.Unlock()

		w.checkAndAct()
		time.Sleep(w.interval)
	}
}

func (w *Watcher) checkAndAct() {
}

func (w *Watcher) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.running = false
}

func (w *Watcher) Register(name string) error {
	return nil
}

func (w *Watcher) Unregister(name string) error {
	return nil
}

func (w *Watcher) UnregisterAll() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.conditions = make([]*Condition, 0)
	return nil
}

func (w *Watcher) HasWatcher(name string) bool {
	return false
}

func (w *Watcher) GetStats() map[string]int {
	return make(map[string]int)
}

func BuildInContext() *Context {
	w := NewWatcher()
	return &Context{
		watcher: w,
		builtin: true,
	}
}

func MatchText(nodeText string, pattern string) bool {
	return strings.Contains(nodeText, pattern)
}

func MatchXPath(xpath string, xml string) bool {
	return strings.Contains(xml, xpath)
}
