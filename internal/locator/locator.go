package locator

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

type Locator struct {
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

func NewLocator(value string) *Locator {
	if re, ok := ParseRegex(value); ok {
		return &Locator{
			Type:  SelectorTypeText,
			Value: value,
			regex: re,
		}
	}

	return &Locator{
		Type:  SelectorTypeText,
		Value: value,
		Index: -1,
	}
}

func ByText(text string) *Locator {
	return NewLocator(text)
}

func ByID(id string) *Locator {
	if re, ok := ParseRegex(id); ok {
		return &Locator{Type: SelectorTypeID, Value: id, regex: re}
	}
	return &Locator{Type: SelectorTypeID, Value: id, Index: -1}
}

func ByXPath(xpath string) *Locator {
	return &Locator{Type: SelectorTypeXPath, Value: xpath, Index: -1}
}

func ByClassName(class string) *Locator {
	return &Locator{Type: SelectorTypeClassName, Value: class, Index: -1}
}

func ByPredicate(predicate string) *Locator {
	return &Locator{Type: SelectorTypePredicate, Value: predicate, Index: -1}
}

func ByClassChain(chain string) *Locator {
	return &Locator{Type: SelectorTypeClassChain, Value: chain, Index: -1}
}

func (l *Locator) SetIndex(idx int) *Locator {
	l.Index = idx
	return l
}

func (l *Locator) Match(text string) bool {
	if l.regex != nil {
		return l.regex.MatchString(text)
	}
	return strings.ToLower(text) == strings.ToLower(l.Value)
}

func (l *Locator) String() string {
	return l.Value
}
