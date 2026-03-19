package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name       string
		baseURL    string
		wantSuffix string
	}{
		{"bare domain", "https://bbb.example.com", "https://bbb.example.com/bigbluebutton/api"},
		{"with bigbluebutton", "https://bbb.example.com/bigbluebutton", "https://bbb.example.com/bigbluebutton/api"},
		{"with full api path", "https://bbb.example.com/bigbluebutton/api", "https://bbb.example.com/bigbluebutton/api"},
		{"trailing slash", "https://bbb.example.com/", "https://bbb.example.com/bigbluebutton/api"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.baseURL, "secret", false)
			if c.BaseURL != tt.wantSuffix {
				t.Errorf("BaseURL = %s, want %s", c.BaseURL, tt.wantSuffix)
			}
		})
	}
}

func TestDoGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.Contains(r.URL.RawQuery, "checksum=") {
			t.Error("request missing checksum parameter")
		}
		if !strings.Contains(r.URL.RawQuery, "meetingID=test123") {
			t.Error("request missing meetingID parameter")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<response><returncode>SUCCESS</returncode></response>`))
	}))
	defer srv.Close()

	c := &Client{
		BaseURL:    srv.URL,
		Secret:     "test-secret",
		HTTPClient: http.DefaultClient,
	}

	body, err := c.doGet("isMeetingRunning", map[string][]string{
		"meetingID": {"test123"},
	})
	if err != nil {
		t.Fatalf("doGet error: %v", err)
	}
	if !strings.Contains(string(body), "SUCCESS") {
		t.Errorf("response body = %s, want SUCCESS", string(body))
	}
}

func TestDoGetHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := &Client{
		BaseURL:    srv.URL,
		Secret:     "test-secret",
		HTTPClient: http.DefaultClient,
	}

	_, err := c.doGet("getMeetings", nil)
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error = %v, want HTTP 500 mention", err)
	}
}

func TestDoPost(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/xml" {
			t.Errorf("expected Content-Type application/xml, got %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<response><returncode>SUCCESS</returncode></response>`))
	}))
	defer srv.Close()

	c := &Client{
		BaseURL:    srv.URL,
		Secret:     "test-secret",
		HTTPClient: http.DefaultClient,
	}

	body, err := c.doPost("insertDocument", map[string][]string{
		"meetingID": {"test123"},
	}, "application/xml", strings.NewReader("<modules/>"))
	if err != nil {
		t.Fatalf("doPost error: %v", err)
	}
	if !strings.Contains(string(body), "SUCCESS") {
		t.Errorf("response body = %s, want SUCCESS", string(body))
	}
}
