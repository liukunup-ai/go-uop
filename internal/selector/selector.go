package selector

import (
	"regexp"
	"strings"
)

type SelectorType int

const (
	SelectorTypeText SelectorType = iota
	SelectorTypeID
	SelectorTypeXPath
	SelectorTypeClassName
	SelectorTypePredicate
	SelectorTypeClassChain
)

type Selector struct {
	Type  SelectorType
	Value string
	Index int
	regex *regexp.Regexp
}

func isRegex(value string) bool {
	return len(value) >= 2 &&
		value[0] == '/' &&
		value[len(value)-1] == '/'
}

func ParseRegex(value string) (*regexp.Regexp, bool) {
	if !isRegex(value) {
		return nil, false
	}
	pattern := value[1 : len(value)-1]
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, false
	}
	return re, true
}

func NewSelector(value string) *Selector {
	if re, ok := ParseRegex(value); ok {
		return &Selector{
			Type:  SelectorTypeText,
			Value: value,
			regex: re,
		}
	}

	return &Selector{
		Type:  SelectorTypeText,
		Value: value,
		Index: -1,
	}
}

func ByText(text string) *Selector {
	return NewSelector(text)
}

func ByID(id string) *Selector {
	if re, ok := ParseRegex(id); ok {
		return &Selector{Type: SelectorTypeID, Value: id, regex: re}
	}
	return &Selector{Type: SelectorTypeID, Value: id, Index: -1}
}

func ByXPath(xpath string) *Selector {
	return &Selector{Type: SelectorTypeXPath, Value: xpath, Index: -1}
}

func ByClassName(class string) *Selector {
	return &Selector{Type: SelectorTypeClassName, Value: class, Index: -1}
}

func ByPredicate(predicate string) *Selector {
	return &Selector{Type: SelectorTypePredicate, Value: predicate, Index: -1}
}

func ByClassChain(chain string) *Selector {
	return &Selector{Type: SelectorTypeClassChain, Value: chain, Index: -1}
}

func (l *Selector) SetIndex(idx int) *Selector {
	l.Index = idx
	return l
}

func (l *Selector) Match(text string) bool {
	if l.regex != nil {
		return l.regex.MatchString(text)
	}
	return strings.ToLower(text) == strings.ToLower(l.Value)
}

func (l *Selector) String() string {
	return l.Value
}
