package api

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestChecksum(t *testing.T) {
	tests := []struct {
		name        string
		apiCall     string
		queryString string
		secret      string
	}{
		{
			name:        "create with params",
			apiCall:     "create",
			queryString: "name=Test+Meeting&meetingID=abc123",
			secret:      "my-secret",
		},
		{
			name:        "getMeetings no params",
			apiCall:     "getMeetings",
			queryString: "",
			secret:      "my-secret",
		},
		{
			name:        "end with meetingID",
			apiCall:     "end",
			queryString: "meetingID=abc123",
			secret:      "supersecret",
		},
		{
			name:        "empty secret",
			apiCall:     "create",
			queryString: "name=Test",
			secret:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Checksum(tt.apiCall, tt.queryString, tt.secret)

			// Verify independently
			h := sha256.New()
			h.Write([]byte(tt.apiCall + tt.queryString + tt.secret))
			want := hex.EncodeToString(h.Sum(nil))

			if got != want {
				t.Errorf("Checksum() = %s, want %s", got, want)
			}

			// Verify it's a valid hex-encoded SHA-256 (64 chars)
			if len(got) != 64 {
				t.Errorf("Checksum length = %d, want 64", len(got))
			}
		})
	}
}

func TestChecksumDeterministic(t *testing.T) {
	a := Checksum("create", "name=Test", "secret")
	b := Checksum("create", "name=Test", "secret")
	if a != b {
		t.Error("Checksum is not deterministic")
	}
}

func TestChecksumDifferentInputs(t *testing.T) {
	a := Checksum("create", "name=Test", "secret1")
	b := Checksum("create", "name=Test", "secret2")
	if a == b {
		t.Error("Different secrets should produce different checksums")
	}
}
