package selector

import (
	"fmt"
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
	SelectorTypeIndex
	SelectorTypePoint
	SelectorTypeCSS
	SelectorTypeTraits
	SelectorTypeEnabled
	SelectorTypeChecked
	SelectorTypeFocused
	SelectorTypeSelected
	SelectorTypeWidth
	SelectorTypeHeight
	SelectorTypeTolerance
	SelectorTypeImage
	SelectorTypeAlgorithm
	SelectorTypeThreshold
	SelectorTypeAbove
	SelectorTypeBelow
	SelectorTypeLeftOf
	SelectorTypeRightOf
	SelectorTypeContainsChild
	SelectorTypeChildOf
	SelectorTypeContainsDescendants
)

func (s SelectorType) String() string {
	switch s {
	case SelectorTypeText:
		return "text"
	case SelectorTypeID:
		return "id"
	case SelectorTypeXPath:
		return "xpath"
	case SelectorTypeClassName:
		return "className"
	case SelectorTypePredicate:
		return "predicate"
	case SelectorTypeClassChain:
		return "classChain"
	case SelectorTypeIndex:
		return "index"
	case SelectorTypePoint:
		return "point"
	case SelectorTypeCSS:
		return "css"
	case SelectorTypeTraits:
		return "traits"
	case SelectorTypeEnabled:
		return "enabled"
	case SelectorTypeChecked:
		return "checked"
	case SelectorTypeFocused:
		return "focused"
	case SelectorTypeSelected:
		return "selected"
	case SelectorTypeWidth:
		return "width"
	case SelectorTypeHeight:
		return "height"
	case SelectorTypeTolerance:
		return "tolerance"
	case SelectorTypeImage:
		return "image"
	case SelectorTypeAlgorithm:
		return "algorithm"
	case SelectorTypeThreshold:
		return "threshold"
	case SelectorTypeAbove:
		return "above"
	case SelectorTypeBelow:
		return "below"
	case SelectorTypeLeftOf:
		return "leftOf"
	case SelectorTypeRightOf:
		return "rightOf"
	case SelectorTypeContainsChild:
		return "containsChild"
	case SelectorTypeChildOf:
		return "childOf"
	case SelectorTypeContainsDescendants:
		return "containsDescendants"
	default:
		return "unknown"
	}
}

type Selector struct {
	Type SelectorType
	Key  string

	Text   string
	ID     string
	XPath  string
	Index  int
	Point  string
	CSS    string
	Traits []string

	Enabled  *bool
	Checked  *bool
	Focused  *bool
	Selected *bool

	Width     int
	Height    int
	Tolerance int

	Image     string
	Algorithm string
	Threshold float64

	Regex *regexp.Regexp

	Nested *Selector
}

func (s *Selector) String() string {
	var parts []string
	if s.Text != "" {
		parts = append(parts, fmt.Sprintf("text=%q", s.Text))
	}
	if s.ID != "" {
		parts = append(parts, fmt.Sprintf("id=%q", s.ID))
	}
	if s.XPath != "" {
		parts = append(parts, fmt.Sprintf("xpath=%q", s.XPath))
	}
	if s.Index > 0 {
		parts = append(parts, fmt.Sprintf("index=%d", s.Index))
	}
	if s.Point != "" {
		parts = append(parts, fmt.Sprintf("point=%q", s.Point))
	}
	if s.CSS != "" {
		parts = append(parts, fmt.Sprintf("css=%q", s.CSS))
	}
	if len(s.Traits) > 0 {
		parts = append(parts, fmt.Sprintf("traits=%v", s.Traits))
	}
	if s.Enabled != nil {
		parts = append(parts, fmt.Sprintf("enabled=%v", *s.Enabled))
	}
	if s.Checked != nil {
		parts = append(parts, fmt.Sprintf("checked=%v", *s.Checked))
	}
	if s.Focused != nil {
		parts = append(parts, fmt.Sprintf("focused=%v", *s.Focused))
	}
	if s.Selected != nil {
		parts = append(parts, fmt.Sprintf("selected=%v", *s.Selected))
	}
	if s.Width > 0 {
		parts = append(parts, fmt.Sprintf("width=%d", s.Width))
	}
	if s.Height > 0 {
		parts = append(parts, fmt.Sprintf("height=%d", s.Height))
	}
	if s.Tolerance > 0 {
		parts = append(parts, fmt.Sprintf("tolerance=%d", s.Tolerance))
	}
	if s.Image != "" {
		parts = append(parts, fmt.Sprintf("image=%q", s.Image))
	}
	if s.Algorithm != "" {
		parts = append(parts, fmt.Sprintf("algorithm=%q", s.Algorithm))
	}
	if s.Threshold > 0 {
		parts = append(parts, fmt.Sprintf("threshold=%v", s.Threshold))
	}
	if s.Nested != nil {
		parts = append(parts, fmt.Sprintf("%s={%s}", s.Type.String(), s.Nested.String()))
	}
	return strings.Join(parts, ", ")
}

func (s *Selector) IsEmpty() bool {
	return s.Text == "" && s.ID == "" && s.XPath == "" && s.Index == 0 &&
		s.Point == "" && s.CSS == "" && len(s.Traits) == 0 &&
		s.Width == 0 && s.Height == 0 && s.Image == "" && s.Nested == nil
}

func isRegex(value string) bool {
	return len(value) >= 2 && value[0] == '/' && value[len(value)-1] == '/'
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

func ByText(text string) *Selector {
	return &Selector{Type: SelectorTypeText, Text: text}
}

func ByID(id string) *Selector {
	if re, ok := ParseRegex(id); ok {
		return &Selector{Type: SelectorTypeID, ID: id, Regex: re}
	}
	return &Selector{Type: SelectorTypeID, ID: id}
}

func ByXPath(xpath string) *Selector {
	return &Selector{Type: SelectorTypeXPath, XPath: xpath}
}

func ByClassName(class string) *Selector {
	return &Selector{Type: SelectorTypeClassName, ID: class}
}

func ByPredicate(predicate string) *Selector {
	return &Selector{Type: SelectorTypePredicate, ID: predicate}
}

func ByClassChain(chain string) *Selector {
	return &Selector{Type: SelectorTypeClassChain, ID: chain}
}

func ByIndex(idx int) *Selector {
	return &Selector{Type: SelectorTypeIndex, Index: idx}
}

func ByPoint(x, y interface{}) *Selector {
	var pointStr string
	switch v := y.(type) {
	case int:
		pointStr = fmt.Sprintf("%v,%d", x, v)
	case float64:
		pointStr = fmt.Sprintf("%v,%.1f", x, v)
	case string:
		pointStr = fmt.Sprintf("%v,%s", x, v)
	default:
		pointStr = fmt.Sprintf("%v", x)
	}
	return &Selector{Type: SelectorTypePoint, Point: pointStr}
}

func ByCSS(css string) *Selector {
	return &Selector{Type: SelectorTypeCSS, CSS: css}
}

func ByTraits(traits ...string) *Selector {
	return &Selector{Type: SelectorTypeTraits, Traits: traits}
}

func ByEnabled(enabled bool) *Selector {
	return &Selector{Type: SelectorTypeEnabled, Enabled: &enabled}
}

func ByChecked(checked bool) *Selector {
	return &Selector{Type: SelectorTypeChecked, Checked: &checked}
}

func ByFocused(focused bool) *Selector {
	return &Selector{Type: SelectorTypeFocused, Focused: &focused}
}

func BySelected(selected bool) *Selector {
	return &Selector{Type: SelectorTypeSelected, Selected: &selected}
}

func ByImage(image string) *Selector {
	return &Selector{Type: SelectorTypeImage, Image: image}
}

func ByAlgorithm(algorithm string) *Selector {
	return &Selector{Type: SelectorTypeAlgorithm, Algorithm: algorithm}
}

func ByThreshold(threshold float64) *Selector {
	return &Selector{Type: SelectorTypeThreshold, Threshold: threshold}
}

func Above(selector *Selector) *Selector {
	return &Selector{Type: SelectorTypeAbove, Nested: selector}
}

func Below(selector *Selector) *Selector {
	return &Selector{Type: SelectorTypeBelow, Nested: selector}
}

func LeftOf(selector *Selector) *Selector {
	return &Selector{Type: SelectorTypeLeftOf, Nested: selector}
}

func RightOf(selector *Selector) *Selector {
	return &Selector{Type: SelectorTypeRightOf, Nested: selector}
}

func ContainsChild(selector *Selector) *Selector {
	return &Selector{Type: SelectorTypeContainsChild, Nested: selector}
}

func ChildOf(selector *Selector) *Selector {
	return &Selector{Type: SelectorTypeChildOf, Nested: selector}
}

func ContainsDescendants(selector *Selector) *Selector {
	return &Selector{Type: SelectorTypeContainsDescendants, Nested: selector}
}

func (l *Selector) SetIndex(idx int) *Selector {
	l.Index = idx
	return l
}

func (l *Selector) Match(text string) bool {
	if l.Regex != nil {
		return l.Regex.MatchString(text)
	}
	return strings.EqualFold(text, l.Text)
}
