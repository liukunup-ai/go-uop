package wda

import (
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

func TestClient_Screenshot(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/screenshot" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]string{
			"value": "iVBORw0KGgoAAAANSUhEUg==",
		})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	data, err := client.Screenshot()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected screenshot data")
	}
}

func TestClient_Tap(t *testing.T) {
	var tapPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tapPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"value": nil})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	err := client.Tap(100, 200)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if tapPath != "/wda/tap/0/100/200" {
		t.Errorf("unexpected tap path: %s", tapPath)
	}
}

func TestClient_SendKeys(t *testing.T) {
	var method, body string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)
		body = req["value"].([]interface{})[0].(string)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"value": nil})
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	err := client.SendKeys("hello")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if method != "POST" {
		t.Errorf("expected POST, got %s", method)
	}
	if body != "hello" {
		t.Errorf("expected 'hello', got '%s'", body)
	}
}
