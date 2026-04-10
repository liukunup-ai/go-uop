package xpath

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
)

type Matcher struct {
	expr string
}

type XMLNode struct {
	XMLName  xml.Name
	Attrs    map[string]string
	Text     string
	Children []*XMLNode
	Parent   *XMLNode
	Depth    int
}

func Compile(expr string) (*Matcher, error) {
	return &Matcher{expr: expr}, nil
}

func (m *Matcher) Match(root *XMLNode) []*XMLNode {
	return m.matchRecursive(root, m.expr)
}

func (m *Matcher) matchNode(node *XMLNode, expr string) []*XMLNode {
	if strings.HasPrefix(expr, "//") {
		return m.matchRecursive(node, expr[2:])
	}
	if strings.HasPrefix(expr, "/") {
		return m.matchDirect(node, expr[1:])
	}
	if strings.HasPrefix(expr, "@") {
		return m.matchAttribute(node, expr[1:])
	}
	if strings.HasPrefix(expr, "*") {
		return m.matchAny(node)
	}
	return nil
}

func (m *Matcher) matchRecursive(node *XMLNode, expr string) []*XMLNode {
	var results []*XMLNode

	if strings.HasPrefix(expr, "*") {
		results = append(results, m.matchAny(node)...)
	} else {
		results = append(results, m.matchDirect(node, expr)...)
	}

	for _, child := range node.Children {
		results = append(results, m.matchRecursive(child, expr)...)
		results = append(results, m.matchRecursive(child, "//"+expr)...)
	}

	return results
}

func (m *Matcher) matchDirect(node *XMLNode, expr string) []*XMLNode {
	if node.XMLName.Local == expr {
		if !m.hasPredicates(expr) {
			return []*XMLNode{node}
		}
		if m.matchPredicates(node, expr) {
			return []*XMLNode{node}
		}
	}
	return nil
}

func (m *Matcher) matchAny(node *XMLNode) []*XMLNode {
	return []*XMLNode{node}
}

func (m *Matcher) hasPredicates(expr string) bool {
	return strings.Contains(expr, "[")
}

func (m *Matcher) matchPredicates(node *XMLNode, expr string) bool {
	predicateRe := regexp.MustCompile(`\[([^\]]+)\]`)
	matches := predicateRe.FindAllStringSubmatch(expr, -1)

	for _, match := range matches {
		predicate := match[1]
		if strings.HasPrefix(predicate, "@") {
			if !m.matchAttributePredicate(node, predicate[1:]) {
				return false
			}
		} else if regexp.MustCompile(`^\d+$`).MatchString(predicate) {
			idx := 0
			fmt.Sscanf(predicate, "%d", &idx)
			parent := node.Parent
			if parent == nil {
				return false
			}
			found := false
			for i, child := range parent.Children {
				if child.XMLName.Local == node.XMLName.Local && i == idx {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		} else if strings.Contains(predicate, "text()") {
			if !m.matchTextPredicate(node, predicate) {
				return false
			}
		}
	}
	return true
}

func (m *Matcher) matchAttributePredicate(node *XMLNode, attrExpr string) bool {
	eqRe := regexp.MustCompile(`(\w+)\s*=\s*['"]([^'"]+)['"]`)
	matches := eqRe.FindStringSubmatch(attrExpr)
	if len(matches) < 3 {
		return false
	}
	attrName := matches[1]
	attrValue := matches[2]
	return node.Attrs[attrName] == attrValue
}

func (m *Matcher) matchTextPredicate(node *XMLNode, predicate string) bool {
	containsRe := regexp.MustCompile(`contains\s*\(\s*text\(\s*\)\s*,\s*['"]([^'"]+)['"]\s*\)`)
	matches := containsRe.FindStringSubmatch(predicate)
	if len(matches) >= 2 {
		return strings.Contains(node.Text, matches[1])
	}
	return false
}

func (m *Matcher) matchAttribute(node *XMLNode, attrExpr string) []*XMLNode {
	eqRe := regexp.MustCompile(`@?(\w+)\s*=\s*['"]([^'"]+)['"]`)
	matches := eqRe.FindStringSubmatch(attrExpr)
	if len(matches) < 3 {
		return nil
	}
	attrName := matches[1]
	attrValue := matches[2]

	if node.Attrs[attrName] == attrValue {
		return []*XMLNode{node}
	}
	return nil
}

func ParseXML(xmlStr string) (*XMLNode, error) {
	decoder := xml.NewDecoder(strings.NewReader(xmlStr))
	return parseNode(decoder, nil, 0)
}

func parseNode(decoder *xml.Decoder, parent *XMLNode, depth int) (*XMLNode, error) {
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		switch se := token.(type) {
		case xml.StartElement:
			node := &XMLNode{
				XMLName: se.Name,
				Attrs:   make(map[string]string),
				Parent:  parent,
				Depth:   depth,
			}
			for _, attr := range se.Attr {
				node.Attrs[attr.Name.Local] = attr.Value
			}
			child, err := parseNode(decoder, node, depth+1)
			if err == nil && child != nil {
				node.Children = append(node.Children, child)
				for child.Children != nil && len(child.Children) > 0 {
					node.Children = append(node.Children, child.Children...)
					child = child.Children[0]
				}
			}
			if node.Parent == nil {
				return node, nil
			}
			return node, nil
		case xml.CharData:
			if parent != nil && strings.TrimSpace(string(se)) != "" {
				parent.Text += string(se)
			}
		case xml.EndElement:
			if parent != nil {
				return parent, nil
			}
		}
	}
	return nil, nil
}

func FindNodesByText(root *XMLNode, text string) []*XMLNode {
	var results []*XMLNode
	if strings.Contains(root.Text, text) {
		results = append(results, root)
	}
	for _, child := range root.Children {
		results = append(results, FindNodesByText(child, text)...)
	}
	return results
}

func FindNodesByAttr(root *XMLNode, attrName, attrValue string) []*XMLNode {
	var results []*XMLNode
	if root.Attrs[attrName] == attrValue {
		results = append(results, root)
	}
	for _, child := range root.Children {
		results = append(results, FindNodesByAttr(child, attrName, attrValue)...)
	}
	return results
}

func GetNodeBounds(node *XMLNode) (x, y, x2, y2 int) {
	if bounds, ok := node.Attrs["bounds"]; ok {
		coordsRe := regexp.MustCompile(`\[(\d+),(\d+)\]\[(\d+),(\d+)\]`)
		matches := coordsRe.FindStringSubmatch(bounds)
		if len(matches) >= 5 {
			fmt.Sscanf(matches[1], "%d", &x)
			fmt.Sscanf(matches[2], "%d", &y)
			fmt.Sscanf(matches[3], "%d", &x2)
			fmt.Sscanf(matches[4], "%d", &y2)
		}
	}
	return
}
