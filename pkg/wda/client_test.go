package wda

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_InvalidURL(t *testing.T) {
	_, err := NewClient("://invalid")
	if err == nil {
		t.Error("expected error for invalid URL")
	}
}

func TestNewClient_ValidURL(t *testing.T) {
	client, err := NewClient("http://localhost:8100")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if client.BaseURL.String() != "http://localhost:8100" {
		t.Errorf("unexpected base URL: %s", client.BaseURL)
	}
}

func TestClient_IsHealthy_True(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	if !client.IsHealthy() {
		t.Error("expected healthy")
	}
}

func TestClient_IsHealthy_False(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	if client.IsHealthy() {
		t.Error("expected unhealthy")
	}
}

func TestClient_StartSession(t *testing.T) {
	var receivedBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/session" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		json.NewDecoder(r.Body).Decode(&receivedBody)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"value": map[string]any{
				"sessionId": "test-session-123",
			},
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	err := client.StartSession("com.example.app")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if client.SessionID != "test-session-123" {
		t.Errorf("unexpected session ID: %s", client.SessionID)
	}

	capabilities := receivedBody["capabilities"].(map[string]any)
	alwaysMatch := capabilities["alwaysMatch"].(map[string]any)
	if alwaysMatch["bundleId"] != "com.example.app" {
		t.Errorf("unexpected bundle ID: %v", alwaysMatch["bundleId"])
	}
}

func TestClient_StopSession(t *testing.T) {
	deleteCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/session/test-session" && r.Method == "DELETE" {
			deleteCalled = true
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"value": nil})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	client.SessionID = "test-session"
	err := client.StopSession()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !deleteCalled {
		t.Error("DELETE session not called")
	}
	if client.SessionID != "" {
		t.Errorf("session ID should be cleared: %s", client.SessionID)
	}
}

func TestClient_Screenshot(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/session/s1/screenshot" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"value": base64.StdEncoding.EncodeToString([]byte("fake-png")),
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	client.SessionID = "s1"
	data, err := client.Screenshot()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(data) != "fake-png" {
		t.Errorf("unexpected data: %s", string(data))
	}
}

func TestClient_Tap(t *testing.T) {
	actionsPath := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actionsPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"value": nil})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	client.SessionID = "s1"
	err := client.Tap(100, 200)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if actionsPath != "/session/s1/actions" {
		t.Errorf("unexpected actions path: %s", actionsPath)
	}
}

func TestClient_FindElement(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/session/s1/element" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var req FindElementRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Using != "accessibility id" {
			t.Errorf("unexpected strategy: %s", req.Using)
		}
		if req.Value != "test-button" {
			t.Errorf("unexpected selector: %s", req.Value)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"value": map[string]any{
				"ELEMENT": "element-123",
			},
			"sessionId": "s1",
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	client.SessionID = "s1"
	elemID, err := client.FindElement("accessibility id", "test-button")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if elemID != "element-123" {
		t.Errorf("unexpected element ID: %s", elemID)
	}
}

func TestClient_FindElements(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/session/s1/elements" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"value": []map[string]any{
				{"ELEMENT": "element-1"},
				{"ELEMENT": "element-2"},
			},
			"sessionId": "s1",
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	client.SessionID = "s1"
	elems, err := client.FindElements("class name", "UIButton")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(elems) != 2 {
		t.Errorf("expected 2 elements, got %d", len(elems))
	}
}

func TestClient_GetSource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/session/s1/source" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"value": "<html/>",
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	client.SessionID = "s1"
	source, err := client.GetSource()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if source != "<html/>" {
		t.Errorf("unexpected source: %s", source)
	}
}

func TestClient_AcceptAlert(t *testing.T) {
	acceptCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/session/s1/alert/accept" {
			acceptCalled = true
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"value": nil})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	client.SessionID = "s1"
	err := client.AcceptAlert()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !acceptCalled {
		t.Error("alert/accept not called")
	}
}

func TestClient_GetAlertText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/session/s1/alert/text" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"value": "Confirm delete?",
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	client.SessionID = "s1"
	text, err := client.GetAlertText()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if text != "Confirm delete?" {
		t.Errorf("unexpected alert text: %s", text)
	}
}

func TestClient_LaunchApp(t *testing.T) {
	launchCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/session/s1/app/launch" {
			launchCalled = true
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"value": nil})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	client.SessionID = "s1"
	err := client.LaunchApp("com.example.app")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !launchCalled {
		t.Error("app/launch not called")
	}
}

func TestClient_Close(t *testing.T) {
	deleteCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			deleteCalled = true
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"value": nil})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	client.SessionID = "s1"
	client.Close()
	if !deleteCalled {
		t.Error("DELETE not called on Close")
	}
}

func TestW3CError(t *testing.T) {
	err := ErrNoSuchElement
	if err.Error() != "no such element" {
		t.Errorf("unexpected error string: %s", err.Error())
	}

	w3cErr := &W3CError{
		Code:    ErrElementClickIntercepted,
		Message: "Element is covered",
	}
	if w3cErr.Error() != "element click intercepted: Element is covered" {
		t.Errorf("unexpected W3CError string: %s", w3cErr.Error())
	}
}

func TestIsW3CError(t *testing.T) {
	errorResp := `{"value":{"error":"no such element","message":"Element not found"}}`
	if !IsW3CError([]byte(errorResp)) {
		t.Error("expected true for error response")
	}

	validResp := `{"value":{"ELEMENT": "elem-1"}}`
	if IsW3CError([]byte(validResp)) {
		t.Error("expected false for valid response")
	}
}

func TestCapabilities(t *testing.T) {
	cap := Capabilities{
		AlwaysMatch: map[string]any{
			"bundleId": "com.example.app",
		},
		FirstMatch: []map[string]any{
			{"bundleId": "com.example.app"},
		},
	}

	data, err := json.Marshal(cap)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	var decoded Capabilities
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if decoded.AlwaysMatch["bundleId"] != "com.example.app" {
		t.Error("alwaysMatch not preserved")
	}
}

func TestClient_XSessionIdHeader(t *testing.T) {
	var sessionHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionHeader = r.Header.Get("X-Session-Id")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"value": nil})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	client.SessionID = "my-session-123"
	client.Screenshot()

	if sessionHeader != "my-session-123" {
		t.Errorf("unexpected X-Session-Id header: %s", sessionHeader)
	}
}
