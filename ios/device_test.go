package ios

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/liukunup/go-uop/core"
)

func TestNewDevice_InvalidWDAAddress(t *testing.T) {
	_, err := NewDevice("com.example.app", WithAddress("invalid://url"))
	if err == nil {
		t.Error("expected error for invalid WDA address")
	}
}

func TestNewDevice_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/status":
			w.WriteHeader(http.StatusOK)
		case "/session":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	_, err := NewDevice("com.example.app", WithAddress(server.URL))
	if err == nil {
		t.Error("expected error when WDA session fails")
	}
}

func TestDevice_Platform(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/status":
			w.WriteHeader(http.StatusOK)
		case "/session":
			json.NewEncoder(w).Encode(map[string]string{"sessionId": "test-session"})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	device, err := NewDevice("com.example.app", WithAddress(server.URL))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if device.Platform() != core.IOS {
		t.Errorf("expected iOS platform, got %v", device.Platform())
	}
}

func TestDevice_Info(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/status":
			w.WriteHeader(http.StatusOK)
		case "/session":
			json.NewEncoder(w).Encode(map[string]string{"sessionId": "test-session"})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	device, err := NewDevice("com.example.app", WithAddress(server.URL))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := device.Info()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info["platform"] != "ios" {
		t.Errorf("expected platform 'ios', got %v", info["platform"])
	}

	if info["bundleId"] != "com.example.app" {
		t.Errorf("expected bundleId 'com.example.app', got %v", info["bundleId"])
	}
}

func TestDevice_WithUDID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/status":
			w.WriteHeader(http.StatusOK)
		case "/session":
			json.NewEncoder(w).Encode(map[string]string{"sessionId": "test-session"})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	_, err := NewDevice("com.example.app", WithAddress(server.URL), WithUDID("device-uuid"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDevice_ImplementsCoreDevice(t *testing.T) {
	var _ core.Device = (*Device)(nil)
}
