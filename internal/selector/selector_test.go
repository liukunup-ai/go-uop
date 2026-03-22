package selector

import (
	"testing"
)

func TestByText_ExactMatch(t *testing.T) {
	l := ByText("登录")

	if !l.Match("登录") {
		t.Error("expected exact match")
	}
	if l.Match("登录按钮") {
		t.Error("should not match partial")
	}
}

func TestByText_RegexMatch(t *testing.T) {
	l := ByText("/登.*/")

	if !l.Match("登录") {
		t.Error("expected regex match '登录'")
	}
	if !l.Match("登录页") {
		t.Error("expected regex match '登录页'")
	}
	if l.Match("取消") {
		t.Error("should not match '取消'")
	}
}

func TestByText_RegexUserId(t *testing.T) {
	l := ByText("/^用户\\d+$/")

	if !l.Match("用户1") {
		t.Error("expected match '用户1'")
	}
	if !l.Match("用户123") {
		t.Error("expected match '用户123'")
	}
	if l.Match("用户名") {
		t.Error("should not match '用户名'")
	}
}

func TestByID_Regex(t *testing.T) {
	l := ByID("/btn_.*_confirm/")

	if !l.Match("btn_submit_confirm") {
		t.Error("expected match")
	}
	if !l.Match("btn_cancel_confirm") {
		t.Error("expected match")
	}
}

func TestSelector_SetIndex(t *testing.T) {
	l := ByText("确定").SetIndex(2)

	if l.Index != 2 {
		t.Errorf("expected index 2, got %d", l.Index)
	}
}

func TestSelector_DefaultIndex(t *testing.T) {
	l := ByText("确定")

	if l.Index != -1 {
		t.Errorf("expected default index -1, got %d", l.Index)
	}
}

func TestByXPath(t *testing.T) {
	l := ByXPath("//Button[@text='OK']")

	if l.Type != SelectorTypeXPath {
		t.Errorf("expected XPath type, got %d", l.Type)
	}
	if l.Value != "//Button[@text='OK']" {
		t.Errorf("unexpected value: %s", l.Value)
	}
}

func TestByClassName(t *testing.T) {
	l := ByClassName("android.widget.Button")

	if l.Type != SelectorTypeClassName {
		t.Errorf("expected ClassName type, got %d", l.Type)
	}
}

func TestByPredicate(t *testing.T) {
	l := ByPredicate(`name == "test"`)

	if l.Type != SelectorTypePredicate {
		t.Errorf("expected Predicate type, got %d", l.Type)
	}
}

func TestByClassChain(t *testing.T) {
	l := ByClassChain(`**/Button[*]`)

	if l.Type != SelectorTypeClassChain {
		t.Errorf("expected ClassChain type, got %d", l.Type)
	}
}
