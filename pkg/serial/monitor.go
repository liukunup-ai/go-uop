package serial

import (
	"bytes"
	"fmt"
	"regexp"
	"sync"
	"time"
)

// MatchType 匹配策略
type MatchType int

const (
	MatchOnce MatchType = iota
	MatchContinuous
	MatchRateLimited
)

// Rule 规则
type Rule struct {
	Keyword string
	MatchType
	RateInterval time.Duration
	IsRegex      bool
	enabled      bool
	lastTrigger  time.Time
	mu           sync.Mutex
	regex        *regexp.Regexp
}

// EventHandler 事件处理函数
type EventHandler func(Event)

// Monitor 监视器（实现 Observer 接口）
type Monitor struct {
	rules    []*Rule
	handlers map[string]EventHandler
	mu       sync.RWMutex
}

// NewMonitor 创建监视器
func NewMonitor() *Monitor {
	return &Monitor{
		handlers: make(map[string]EventHandler),
	}
}

// AddRule 添加规则（MatchOnce 或 MatchContinuous）
func (m *Monitor) AddRule(keyword string, matchType MatchType, handler EventHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()

	rule := &Rule{
		Keyword:   keyword,
		MatchType: matchType,
		enabled:   true,
	}

	m.rules = append(m.rules, rule)
	m.handlers[keyword] = handler
}

// AddRateLimitedRule 添加限频规则
func (m *Monitor) AddRateLimitedRule(keyword string, interval time.Duration, handler EventHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()

	rule := &Rule{
		Keyword:      keyword,
		MatchType:    MatchRateLimited,
		RateInterval: interval,
		enabled:      true,
	}

	m.rules = append(m.rules, rule)
	m.handlers[keyword] = handler
}

// EnableRule 启用规则
func (m *Monitor) EnableRule(keyword string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, r := range m.rules {
		if r.Keyword == keyword {
			r.mu.Lock()
			r.enabled = true
			r.mu.Unlock()
			return nil
		}
	}
	return fmt.Errorf("rule not found: %s", keyword)
}

// DisableRule 禁用规则
func (m *Monitor) DisableRule(keyword string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, r := range m.rules {
		if r.Keyword == keyword {
			r.mu.Lock()
			r.enabled = false
			r.mu.Unlock()
			return nil
		}
	}
	return fmt.Errorf("rule not found: %s", keyword)
}

// OnData 处理事件
func (m *Monitor) OnData(e Event) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, rule := range m.rules {
		if !rule.shouldTrigger(e.Data) {
			continue
		}

		matchedEvent := Event{
			Data:      e.Data,
			Timestamp: e.Timestamp,
			Rule:      rule,
		}

		if handler, ok := m.handlers[rule.Keyword]; ok {
			handler(matchedEvent)
		}
	}
}

// shouldTrigger 检查规则是否应触发
func (r *Rule) shouldTrigger(data []byte) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.enabled && r.MatchType == MatchOnce {
		return false
	}

	matched := false
	if r.IsRegex {
		if r.regex == nil {
			r.regex = regexp.MustCompile(r.Keyword)
		}
		matched = r.regex.Match(data)
	} else {
		matched = bytes.Contains(data, []byte(r.Keyword))
	}

	if !matched {
		return false
	}

	switch r.MatchType {
	case MatchOnce:
		r.enabled = false
	case MatchRateLimited:
		if time.Since(r.lastTrigger) < r.RateInterval {
			return false
		}
		r.lastTrigger = time.Now()
	}

	return true
}

// OnError 错误处理（空实现）
func (m *Monitor) OnError(err error) {
	// 日志记录，不影响其他规则
}

// OnClose 关闭处理（空实现）
func (m *Monitor) OnClose() {
}
