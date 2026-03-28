package uop

import (
	"errors"
	"fmt"
	"strings"

	"github.com/liukunup/go-uop/core"
	fluentpkg "github.com/liukunup/go-uop/internal/fluent"
	"github.com/liukunup/go-uop/internal/selector"
)

var ErrAssertionFailed = errors.New("assertion failed")

type AssertError struct {
	Message string
	Cause   error
}

func (e *AssertError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *AssertError) Unwrap() error {
	return e.Cause
}

func AssertVisible(device core.Device, loc *selector.Selector) error {
	finder, ok := device.(fluentpkg.ElementFinder)
	if !ok {
		return &AssertError{Message: "device does not support element finding"}
	}

	elem, err := finder.FindElement(loc)
	if err != nil {
		return &AssertError{Message: "element not found", Cause: err}
	}

	if !elem.Visible {
		return &AssertError{Message: fmt.Sprintf("element %q is not visible", loc.Value)}
	}

	return nil
}

func AssertNotVisible(device core.Device, loc *selector.Selector) error {
	finder, ok := device.(fluentpkg.ElementFinder)
	if !ok {
		return &AssertError{Message: "device does not support element finding"}
	}

	elem, err := finder.FindElement(loc)
	if err != nil {
		return nil
	}

	if elem.Visible {
		return &AssertError{Message: fmt.Sprintf("element %q should not be visible", loc.Value)}
	}

	return nil
}

func AssertTrue(condition bool, msg string) error {
	if !condition {
		if msg == "" {
			return &AssertError{Message: "expected true but got false"}
		}
		return &AssertError{Message: msg}
	}
	return nil
}

func AssertFalse(condition bool, msg string) error {
	if condition {
		if msg == "" {
			return &AssertError{Message: "expected false but got true"}
		}
		return &AssertError{Message: msg}
	}
	return nil
}

func AssertEqual(expected, actual any, msg string) error {
	if expected != actual {
		if msg == "" {
			return &AssertError{Message: fmt.Sprintf("expected %v but got %v", expected, actual)}
		}
		return &AssertError{Message: msg}
	}
	return nil
}

func AssertNotEqual(notExpected, actual any, msg string) error {
	if notExpected == actual {
		if msg == "" {
			return &AssertError{Message: fmt.Sprintf("expected not %v but got %v", notExpected, actual)}
		}
		return &AssertError{Message: msg}
	}
	return nil
}

func ExpectError(err error, msg string) error {
	if err == nil {
		if msg == "" {
			return &AssertError{Message: "expected error but got nil"}
		}
		return &AssertError{Message: msg}
	}
	return nil
}

func AssertNoError(err error, msg string) error {
	if err != nil {
		if msg == "" {
			return &AssertError{Message: fmt.Sprintf("unexpected error: %v", err)}
		}
		return &AssertError{Message: msg, Cause: err}
	}
	return nil
}

func AssertText(device core.Device, loc *selector.Selector, expectedText string) error {
	finder, ok := device.(fluentpkg.ElementFinder)
	if !ok {
		return &AssertError{Message: "device does not support element finding"}
	}

	elem, err := finder.FindElement(loc)
	if err != nil {
		return &AssertError{Message: "element not found", Cause: err}
	}

	if elem.Text != expectedText {
		return &AssertError{Message: fmt.Sprintf("expected text %q but got %q", expectedText, elem.Text)}
	}

	return nil
}

func AssertContains(s, substr string) error {
	if !strings.Contains(s, substr) {
		return &AssertError{Message: fmt.Sprintf("expected %q to contain %q", s, substr)}
	}
	return nil
}

func AssertNotContains(s, substr string) error {
	if strings.Contains(s, substr) {
		return &AssertError{Message: fmt.Sprintf("expected %q to not contain %q", s, substr)}
	}
	return nil
}
