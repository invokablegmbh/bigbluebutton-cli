package testutil

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// NewFakeBBBServer creates a test HTTP server that responds to BBB API calls
// with canned responses from the testdata directory.
func NewFakeBBBServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Extract the API call name from the path (last segment)
		parts := strings.Split(strings.TrimRight(path, "/"), "/")
		apiCall := parts[len(parts)-1]

		// Verify checksum is present
		if r.URL.Query().Get("checksum") == "" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<response><returncode>FAILED</returncode><messageKey>checksumError</messageKey><message>You did not pass the checksum security check</message></response>`))
			return
		}

		response, ok := fixtures[apiCall]
		if !ok {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<response><returncode>FAILED</returncode><messageKey>unknownCall</messageKey><message>Unknown API call: ` + apiCall + `</message></response>`))
			return
		}

		w.Header().Set("Content-Type", "text/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
}

// NewFakeBBBServerWithHandler creates a test server with a custom handler.
func NewFakeBBBServerWithHandler(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(handler)
}
