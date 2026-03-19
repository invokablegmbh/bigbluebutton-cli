package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestServer(t *testing.T, expectedCall string, response string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, expectedCall) {
			t.Errorf("expected path to contain %s, got %s", expectedCall, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
}

func newTestClient(serverURL, secret string) *Client {
	return &Client{
		BaseURL:    serverURL,
		Secret:     secret,
		HTTPClient: http.DefaultClient,
	}
}

func TestClient_Create(t *testing.T) {
	srv := newTestServer(t, "create", `<response>
		<returncode>SUCCESS</returncode>
		<meetingID>test-meeting</meetingID>
		<internalMeetingID>internal-123</internalMeetingID>
		<voiceBridge>70001</voiceBridge>
		<createTime>1699900000000</createTime>
		<duration>60</duration>
		<hasUserJoined>false</hasUserJoined>
	</response>`)
	defer srv.Close()

	client := newTestClient(srv.URL, "test-secret")
	resp, err := client.Create(&CreateRequest{
		Name:      "Test Meeting",
		MeetingID: "test-meeting",
		Record:    BoolPtr(true),
		Duration:  IntPtr(60),
	})
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}

	if resp.MeetingID != "test-meeting" {
		t.Errorf("MeetingID = %s, want test-meeting", resp.MeetingID)
	}
	if resp.VoiceBridge != "70001" {
		t.Errorf("VoiceBridge = %s, want 70001", resp.VoiceBridge)
	}
	if resp.Duration != 60 {
		t.Errorf("Duration = %d, want 60", resp.Duration)
	}
}

func TestClient_Create_Error(t *testing.T) {
	srv := newTestServer(t, "create", `<response>
		<returncode>FAILED</returncode>
		<messageKey>idNotUnique</messageKey>
		<message>A meeting already exists with that meeting ID.</message>
	</response>`)
	defer srv.Close()

	client := newTestClient(srv.URL, "test-secret")
	_, err := client.Create(&CreateRequest{
		Name:      "Test",
		MeetingID: "existing",
	})
	if err == nil {
		t.Fatal("expected error for FAILED response")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.MessageKey != "idNotUnique" {
		t.Errorf("MessageKey = %s, want idNotUnique", apiErr.MessageKey)
	}
}

func TestClient_Join(t *testing.T) {
	srv := newTestServer(t, "join", `<response>
		<returncode>SUCCESS</returncode>
		<meeting_id>test-meeting</meeting_id>
		<user_id>user-123</user_id>
		<auth_token>auth-abc</auth_token>
		<session_token>sess-xyz</session_token>
		<url>https://bbb.example.com/client?sessionToken=sess-xyz</url>
	</response>`)
	defer srv.Close()

	client := newTestClient(srv.URL, "test-secret")
	resp, err := client.Join(&JoinRequest{
		FullName:  "John Doe",
		MeetingID: "test-meeting",
		Role:      "MODERATOR",
	})
	if err != nil {
		t.Fatalf("Join error: %v", err)
	}

	if resp.MeetingID != "test-meeting" {
		t.Errorf("MeetingID = %s, want test-meeting", resp.MeetingID)
	}
	if resp.UserID != "user-123" {
		t.Errorf("UserID = %s, want user-123", resp.UserID)
	}
	if resp.SessionToken != "sess-xyz" {
		t.Errorf("SessionToken = %s, want sess-xyz", resp.SessionToken)
	}
}

func TestClient_JoinURL(t *testing.T) {
	client := &Client{
		BaseURL: "https://bbb.example.com/bigbluebutton/api",
		Secret:  "test-secret",
	}

	url := client.JoinURL(&JoinRequest{
		FullName:  "John",
		MeetingID: "m1",
		Role:      "VIEWER",
	})

	if !strings.Contains(url, "bigbluebutton/api/join?") {
		t.Errorf("URL should contain api/join, got %s", url)
	}
	if !strings.Contains(url, "fullName=John") {
		t.Errorf("URL should contain fullName, got %s", url)
	}
	if !strings.Contains(url, "checksum=") {
		t.Errorf("URL should contain checksum, got %s", url)
	}
}

func TestClient_End(t *testing.T) {
	srv := newTestServer(t, "end", `<response>
		<returncode>SUCCESS</returncode>
		<messageKey>sentEndMeetingRequest</messageKey>
		<message>A request to end the meeting was sent.</message>
	</response>`)
	defer srv.Close()

	client := newTestClient(srv.URL, "test-secret")
	_, err := client.End(&EndRequest{MeetingID: "test-meeting"})
	if err != nil {
		t.Fatalf("End error: %v", err)
	}
}

func TestClient_InsertDocument(t *testing.T) {
	var receivedBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		receivedBody = string(body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<response><returncode>SUCCESS</returncode></response>`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL, "test-secret")
	xmlBody := `<modules><module name="presentation"><document url="https://example.com/slides.pdf"/></module></modules>`
	_, err := client.InsertDocument(&InsertDocumentRequest{
		MeetingID: "test-meeting",
		Body:      xmlBody,
	})
	if err != nil {
		t.Fatalf("InsertDocument error: %v", err)
	}
	if receivedBody != xmlBody {
		t.Errorf("body = %s, want %s", receivedBody, xmlBody)
	}
}
