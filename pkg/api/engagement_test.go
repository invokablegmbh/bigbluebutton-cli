package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_SendChatMessage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("meetingID") != "m1" {
			t.Errorf("meetingID = %s, want m1", q.Get("meetingID"))
		}
		if q.Get("message") != "Hello everyone" {
			t.Errorf("message = %s, want Hello everyone", q.Get("message"))
		}
		if q.Get("userName") != "Admin" {
			t.Errorf("userName = %s, want Admin", q.Get("userName"))
		}
		w.Write([]byte(`<response><returncode>SUCCESS</returncode><messageKey>sentChatMessage</messageKey></response>`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL, "test-secret")
	_, err := client.SendChatMessage(&SendChatMessageRequest{
		MeetingID: "m1",
		Message:   "Hello everyone",
		UserName:  StringPtr("Admin"),
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestClient_SendChatMessage_Error(t *testing.T) {
	srv := newTestServer(t, "sendChatMessage", `<response>
		<returncode>FAILED</returncode>
		<messageKey>notFound</messageKey>
		<message>Meeting not found</message>
	</response>`)
	defer srv.Close()

	client := newTestClient(srv.URL, "test-secret")
	_, err := client.SendChatMessage(&SendChatMessageRequest{
		MeetingID: "nonexistent",
		Message:   "Hello",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestClient_GetJoinUrl(t *testing.T) {
	srv := newTestServer(t, "getJoinUrl", `<response>
		<returncode>SUCCESS</returncode>
		<url>https://bbb.example.com/join?token=abc</url>
	</response>`)
	defer srv.Close()

	client := newTestClient(srv.URL, "test-secret")
	resp, err := client.GetJoinUrl(&GetJoinUrlRequest{
		SessionToken: "tok-123",
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if resp.URL != "https://bbb.example.com/join?token=abc" {
		t.Errorf("URL = %s, want https://bbb.example.com/join?token=abc", resp.URL)
	}
}

func TestClient_GetJoinUrl_WithOptions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("sessionToken") != "tok-123" {
			t.Errorf("sessionToken = %s, want tok-123", q.Get("sessionToken"))
		}
		if q.Get("replaceSession") != "true" {
			t.Errorf("replaceSession = %s, want true", q.Get("replaceSession"))
		}
		if q.Get("sessionName") != "New Session" {
			t.Errorf("sessionName = %s, want New Session", q.Get("sessionName"))
		}
		w.Write([]byte(`<response><returncode>SUCCESS</returncode><url>https://example.com/join</url></response>`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL, "test-secret")
	_, err := client.GetJoinUrl(&GetJoinUrlRequest{
		SessionToken:   "tok-123",
		ReplaceSession: BoolPtr(true),
		SessionName:    StringPtr("New Session"),
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}
