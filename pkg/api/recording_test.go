package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetRecordings(t *testing.T) {
	srv := newTestServer(t, "getRecordings", `<response>
		<returncode>SUCCESS</returncode>
		<recordings>
			<recording>
				<recordID>rec-1</recordID>
				<meetingID>m1</meetingID>
				<name>Recording 1</name>
				<published>true</published>
				<state>published</state>
				<participants>5</participants>
				<playback>
					<format>
						<type>presentation</type>
						<url>https://example.com/playback/rec-1</url>
						<length>60</length>
					</format>
				</playback>
			</recording>
		</recordings>
	</response>`)
	defer srv.Close()

	client := newTestClient(srv.URL, "test-secret")
	resp, err := client.GetRecordings(&GetRecordingsRequest{
		State: StringPtr("published"),
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if len(resp.Recordings) != 1 {
		t.Fatalf("expected 1 recording, got %d", len(resp.Recordings))
	}
	if resp.Recordings[0].RecordID != "rec-1" {
		t.Errorf("RecordID = %s, want rec-1", resp.Recordings[0].RecordID)
	}
	if !resp.Recordings[0].Published {
		t.Error("recording should be published")
	}
}

func TestClient_GetRecordings_WithFilters(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("meetingID") != "m1" {
			t.Errorf("meetingID = %s, want m1", q.Get("meetingID"))
		}
		if q.Get("state") != "published" {
			t.Errorf("state = %s, want published", q.Get("state"))
		}
		if q.Get("limit") != "10" {
			t.Errorf("limit = %s, want 10", q.Get("limit"))
		}
		w.Write([]byte(`<response><returncode>SUCCESS</returncode><recordings></recordings></response>`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL, "test-secret")
	_, err := client.GetRecordings(&GetRecordingsRequest{
		MeetingID: StringPtr("m1"),
		State:     StringPtr("published"),
		Limit:     IntPtr(10),
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestClient_PublishRecordings(t *testing.T) {
	srv := newTestServer(t, "publishRecordings", `<response>
		<returncode>SUCCESS</returncode>
		<published>true</published>
	</response>`)
	defer srv.Close()

	client := newTestClient(srv.URL, "test-secret")
	resp, err := client.PublishRecordings(&PublishRecordingsRequest{
		RecordID: "rec-1",
		Publish:  true,
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !resp.Published {
		t.Error("expected published=true")
	}
}

func TestClient_DeleteRecordings(t *testing.T) {
	srv := newTestServer(t, "deleteRecordings", `<response>
		<returncode>SUCCESS</returncode>
		<deleted>true</deleted>
	</response>`)
	defer srv.Close()

	client := newTestClient(srv.URL, "test-secret")
	resp, err := client.DeleteRecordings(&DeleteRecordingsRequest{RecordID: "rec-1"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !resp.Deleted {
		t.Error("expected deleted=true")
	}
}

func TestClient_UpdateRecordings(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("recordID") != "rec-1" {
			t.Errorf("recordID = %s, want rec-1", q.Get("recordID"))
		}
		if q.Get("meta_name") != "New Title" {
			t.Errorf("meta_name = %s, want New Title", q.Get("meta_name"))
		}
		w.Write([]byte(`<response><returncode>SUCCESS</returncode><updated>true</updated></response>`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL, "test-secret")
	resp, err := client.UpdateRecordings(&UpdateRecordingsRequest{
		RecordID: "rec-1",
		Meta:     map[string]string{"name": "New Title"},
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !resp.Updated {
		t.Error("expected updated=true")
	}
}

func TestClient_GetRecordingTextTracks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"response":{"returncode":"SUCCESS","tracks":[{"href":"https://example.com/track.vtt","kind":"subtitles","label":"English","lang":"en","source":"upload"}]}}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL, "test-secret")
	resp, err := client.GetRecordingTextTracks(&GetRecordingTextTracksRequest{RecordID: "rec-1"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if len(resp.Response.Tracks) != 1 {
		t.Fatalf("expected 1 track, got %d", len(resp.Response.Tracks))
	}
	if resp.Response.Tracks[0].Lang != "en" {
		t.Errorf("Lang = %s, want en", resp.Response.Tracks[0].Lang)
	}
}

func TestClient_GetRecordingTextTracks_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"response":{"returncode":"FAILED","messageKey":"notFound","message":"Recording not found"}}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL, "test-secret")
	_, err := client.GetRecordingTextTracks(&GetRecordingTextTracksRequest{RecordID: "nonexistent"})
	if err == nil {
		t.Fatal("expected error for FAILED response")
	}
}
