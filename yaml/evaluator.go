package yaml

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	reVariable = regexp.MustCompile(`\$\{([^}]+)\}`)
	reENV      = regexp.MustCompile(`\$\{ENV\(([^)]+)\)\}`)
)

type Context struct {
	Variables map[string]interface{}
	Env       map[string]string
}

func NewContext() *Context {
	return &Context{
		Variables: make(map[string]interface{}),
		Env:       make(map[string]string),
	}
}

func (c *Context) SetVariable(name string, value interface{}) {
	c.Variables[name] = value
}

func (c *Context) GetVariable(name string) interface{} {
	return c.Variables[name]
}

func (c *Context) Evaluate(input string) (string, error) {
	result := reVariable.ReplaceAllStringFunc(input, func(match string) string {
		expr := match[2 : len(match)-1]
		return c.evalExpr(expr)
	})

	return result, nil
}

func (c *Context) evalExpr(expr string) string {
	if matches := reENV.FindStringSubmatch(expr); len(matches) > 1 {
		key := matches[1]
		if val, ok := c.Env[key]; ok {
			return val
		}
		return ""
	}

	if strings.HasPrefix(expr, "variables.") {
		name := strings.TrimPrefix(expr, "variables.")
		if val, ok := c.Variables[name]; ok {
			return fmt.Sprintf("%v", val)
		}
		return ""
	}

	if val, ok := c.Variables[expr]; ok {
		return fmt.Sprintf("%v", val)
	}

	return "${" + expr + "}"
}
