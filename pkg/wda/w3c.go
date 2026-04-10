package wda

import (
	"encoding/json"
	"fmt"
)

// W3C WebDriver error codes
type ErrorCode string

const (
	ErrElementNotInteractable  ErrorCode = "element not interactable"
	ErrInvalidSessionID        ErrorCode = "invalid session id"
	ErrNoSuchElement           ErrorCode = "no such element"
	ErrStaleElementReference   ErrorCode = "stale element reference"
	ErrElementClickIntercepted ErrorCode = "element click intercepted"
	ErrUnknownError            ErrorCode = "unknown error"
	ErrNoSuchAlert             ErrorCode = "no such alert"
)

// W3C Response envelope - all responses wrap in this structure
type W3CResponse struct {
	Value     any    `json:"value,omitempty"`
	SessionID string `json:"sessionId,omitempty"`
}

// W3C Error response
type W3CErrorResponse struct {
	Value struct {
		Error      ErrorCode `json:"error"`
		Message    string    `json:"message"`
		Stacktrace string    `json:"stacktrace,omitempty"`
	} `json:"value"`
}

// W3C Capabilities (alwaysMatch/firstMatch structure)
type Capabilities struct {
	AlwaysMatch map[string]any   `json:"alwaysMatch,omitempty"`
	FirstMatch  []map[string]any `json:"firstMatch,omitempty"`
}

// NewSessionRequest for POST /session
type NewSessionRequest struct {
	Capabilities Capabilities `json:"capabilities"`
}

// NewSessionResponse from POST /session
type NewSessionResponse struct {
	SessionID    string         `json:"sessionId"`
	Capabilities map[string]any `json:"capabilities"`
}

// Element reference (W3C style: element-UUID-1)
type ElementID string

// FindElementRequest for POST /session/{id}/element
type FindElementRequest struct {
	Using string `json:"using"`
	Value string `json:"value"`
}

// Element rect response
type ElementRect struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// W3C Actions for touch/pointer
type ActionsRequest struct {
	Actions []interface{} `json:"actions"`
}

// Source response
type SourceResponse struct {
	Value string `json:"value"`
}

// Screenshot response
type ScreenshotResponse struct {
	Value string `json:"value"`
}

// Alert response
type AlertResponse struct {
	Value string `json:"value"`
}

// W3C Error implements error interface
func (e ErrorCode) Error() string {
	return string(e)
}

// W3CError wraps error code with message
type W3CError struct {
	Code    ErrorCode
	Message string
}

func (e *W3CError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewW3CError creates a W3CError from error response
func NewW3CError(resp *W3CErrorResponse) *W3CError {
	return &W3CError{
		Code:    resp.Value.Error,
		Message: resp.Value.Message,
	}
}

// IsW3CError checks if response is a W3C error
func IsW3CError(respBody []byte) bool {
	var errResp W3CErrorResponse
	if err := json.Unmarshal(respBody, &errResp); err != nil {
		return false
	}
	return errResp.Value.Error != ""
}
