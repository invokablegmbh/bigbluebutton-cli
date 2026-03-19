package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// newFakeServer creates a test BBB server returning canned responses based on API call.
func newFakeServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		parts := strings.Split(strings.TrimRight(path, "/"), "/")
		apiCall := parts[len(parts)-1]

		responses := map[string]string{
			"create": `<response>
				<returncode>SUCCESS</returncode>
				<meetingID>test-meeting</meetingID>
				<internalMeetingID>internal-123</internalMeetingID>
				<voiceBridge>70001</voiceBridge>
				<createTime>1699900000000</createTime>
				<createDate>Mon Nov 13 12:00:00 UTC 2023</createDate>
				<duration>0</duration>
				<hasUserJoined>false</hasUserJoined>
			</response>`,
			"join": `<response>
				<returncode>SUCCESS</returncode>
				<meeting_id>test-meeting</meeting_id>
				<user_id>user-123</user_id>
				<auth_token>auth-abc</auth_token>
				<session_token>sess-xyz</session_token>
				<url>https://bbb.example.com/client?sessionToken=sess-xyz</url>
			</response>`,
			"end": `<response>
				<returncode>SUCCESS</returncode>
				<messageKey>sentEndMeetingRequest</messageKey>
				<message>Meeting ended.</message>
			</response>`,
			"insertDocument": `<response><returncode>SUCCESS</returncode></response>`,
			"isMeetingRunning": `<response>
				<returncode>SUCCESS</returncode>
				<running>true</running>
			</response>`,
			"getMeetingInfo": `<response>
				<returncode>SUCCESS</returncode>
				<meetingName>Test Meeting</meetingName>
				<meetingID>test-meeting</meetingID>
				<internalMeetingID>internal-123</internalMeetingID>
				<createDate>Mon Nov 13 12:00:00 UTC 2023</createDate>
				<running>true</running>
				<participantCount>2</participantCount>
				<moderatorCount>1</moderatorCount>
				<attendees>
					<attendee>
						<userID>u1</userID>
						<fullName>John Doe</fullName>
						<role>MODERATOR</role>
						<isPresenter>true</isPresenter>
						<isListeningOnly>false</isListeningOnly>
						<hasJoinedVoice>true</hasJoinedVoice>
						<hasVideo>true</hasVideo>
					</attendee>
				</attendees>
			</response>`,
			"getMeetings": `<response>
				<returncode>SUCCESS</returncode>
				<meetings>
					<meeting>
						<meetingName>Test Meeting 1</meetingName>
						<meetingID>meeting-1</meetingID>
						<running>true</running>
						<participantCount>5</participantCount>
						<recording>false</recording>
					</meeting>
				</meetings>
			</response>`,
			"getRecordings": `<response>
				<returncode>SUCCESS</returncode>
				<recordings>
					<recording>
						<recordID>rec-1</recordID>
						<meetingID>m1</meetingID>
						<name>Recording 1</name>
						<published>true</published>
						<state>published</state>
						<startTime>1699900000000</startTime>
						<endTime>1699903600000</endTime>
						<participants>5</participants>
					</recording>
				</recordings>
			</response>`,
			"publishRecordings": `<response>
				<returncode>SUCCESS</returncode>
				<published>true</published>
			</response>`,
			"deleteRecordings": `<response>
				<returncode>SUCCESS</returncode>
				<deleted>true</deleted>
			</response>`,
			"updateRecordings": `<response>
				<returncode>SUCCESS</returncode>
				<updated>true</updated>
			</response>`,
			"getRecordingTextTracks": `{"response":{"returncode":"SUCCESS","tracks":[{"href":"https://example.com/track.vtt","kind":"subtitles","label":"English","lang":"en","source":"upload"}]}}`,
			"sendChatMessage": `<response>
				<returncode>SUCCESS</returncode>
				<messageKey>sentChatMessage</messageKey>
			</response>`,
			"getJoinUrl": `<response>
				<returncode>SUCCESS</returncode>
				<url>https://bbb.example.com/join?token=abc</url>
			</response>`,
		}

		resp, ok := responses[apiCall]
		if !ok {
			w.Write([]byte(`<response><returncode>FAILED</returncode><messageKey>unknownCall</messageKey></response>`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resp))
	}))
}

func executeCommand(t *testing.T, args ...string) (string, error) {
	t.Helper()

	// Capture stdout since commands write to os.Stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd := NewRootCmd()
	buf := new(bytes.Buffer)
	rootCmd.SetOut(w)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	var outBuf bytes.Buffer
	outBuf.ReadFrom(r)
	// Combine cobra output and stdout captures
	outBuf.Write(buf.Bytes())
	return outBuf.String(), err
}

func TestVersionCommand(t *testing.T) {
	out, err := executeCommand(t, "version")
	if err != nil {
		t.Fatalf("version error: %v", err)
	}
	if !strings.Contains(out, "bigbluebutton-cli version") {
		t.Errorf("output = %s, want version string", out)
	}
}

func TestCreateCommand(t *testing.T) {
	srv := newFakeServer(t)
	defer srv.Close()

	out, err := executeCommand(t,
		"create",
		"--url", srv.URL,
		"--secret", "test-secret",
		"--name", "Test Meeting",
		"--meeting-id", "test-meeting",
	)
	if err != nil {
		t.Fatalf("create error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "Meeting created successfully") {
		t.Errorf("output = %s, want success message", out)
	}
	if !strings.Contains(out, "test-meeting") {
		t.Errorf("output = %s, want meeting ID", out)
	}
}

func TestCreateJoinCommand(t *testing.T) {
	srv := newFakeServer(t)
	defer srv.Close()

	out, err := executeCommand(t,
		"create-join",
		"--url", srv.URL,
		"--secret", "test-secret",
		"--name", "Test Meeting",
		"--meeting-id", "test-meeting",
		"--user-name", "Alice",
		"--role", "MODERATOR",
	)
	if err != nil {
		t.Fatalf("create-join error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "Meeting created and join URL generated") {
		t.Errorf("output = %s, want success message", out)
	}
	if !strings.Contains(out, "test-meeting") {
		t.Errorf("output = %s, want meeting ID", out)
	}
	if !strings.Contains(out, "Alice") {
		t.Errorf("output = %s, want user name", out)
	}
	if !strings.Contains(out, "Join URL") {
		t.Errorf("output = %s, want join URL field", out)
	}
	if !strings.Contains(out, "checksum=") {
		t.Errorf("output = %s, want checksum in join URL", out)
	}
}

func TestCreateJoinCommand_URLOnly(t *testing.T) {
	srv := newFakeServer(t)
	defer srv.Close()

	out, err := executeCommand(t,
		"create-join",
		"--url", srv.URL,
		"--secret", "test-secret",
		"--name", "Test Meeting",
		"--meeting-id", "test-meeting",
		"--user-name", "Bob",
		"--role", "VIEWER",
		"--url-only",
	)
	if err != nil {
		t.Fatalf("create-join error: %v\noutput: %s", err, out)
	}
	// Should output only the URL, no extra text
	if strings.Contains(out, "Meeting created") {
		t.Errorf("url-only mode should not print extra text")
	}
	if !strings.Contains(out, "join?") {
		t.Errorf("output = %s, want join URL", out)
	}
	if !strings.Contains(out, "checksum=") {
		t.Errorf("output = %s, want checksum in URL", out)
	}
}

func TestCreateJoinCommand_MissingRequiredFlags(t *testing.T) {
	srv := newFakeServer(t)
	defer srv.Close()

	_, err := executeCommand(t,
		"create-join",
		"--url", srv.URL,
		"--secret", "test-secret",
		"--name", "Test",
		"--meeting-id", "test-meeting",
		// missing --user-name and --role
	)
	if err == nil {
		t.Fatal("expected error for missing required flags")
	}
}

func TestCreateCommand_MissingRequiredFlags(t *testing.T) {
	srv := newFakeServer(t)
	defer srv.Close()

	_, err := executeCommand(t,
		"create",
		"--url", srv.URL,
		"--secret", "test-secret",
	)
	if err == nil {
		t.Fatal("expected error for missing required flags")
	}
}

func TestJoinCommand(t *testing.T) {
	srv := newFakeServer(t)
	defer srv.Close()

	out, err := executeCommand(t,
		"join",
		"--url", srv.URL,
		"--secret", "test-secret",
		"--meeting-id", "test-meeting",
		"--name", "John Doe",
		"--role", "MODERATOR",
	)
	if err != nil {
		t.Fatalf("join error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "Join URL generated successfully") {
		t.Errorf("output = %s, want success message", out)
	}
}

func TestJoinCommand_URLOnly(t *testing.T) {
	srv := newFakeServer(t)
	defer srv.Close()

	out, err := executeCommand(t,
		"join",
		"--url", srv.URL,
		"--secret", "test-secret",
		"--meeting-id", "test-meeting",
		"--name", "John Doe",
		"--role", "VIEWER",
		"--url-only",
	)
	if err != nil {
		t.Fatalf("join error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "join?") {
		t.Errorf("output = %s, want join URL", out)
	}
	if !strings.Contains(out, "checksum=") {
		t.Errorf("output = %s, want checksum in URL", out)
	}
}

func TestEndCommand(t *testing.T) {
	srv := newFakeServer(t)
	defer srv.Close()

	out, err := executeCommand(t,
		"end",
		"--url", srv.URL,
		"--secret", "test-secret",
		"--meeting-id", "test-meeting",
	)
	if err != nil {
		t.Fatalf("end error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "has been ended") {
		t.Errorf("output = %s, want ended message", out)
	}
}

func TestIsMeetingRunningCommand(t *testing.T) {
	srv := newFakeServer(t)
	defer srv.Close()

	out, err := executeCommand(t,
		"is-meeting-running",
		"--url", srv.URL,
		"--secret", "test-secret",
		"--meeting-id", "test-meeting",
	)
	if err != nil {
		t.Fatalf("is-meeting-running error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "true") {
		t.Errorf("output = %s, want true", out)
	}
}

func TestGetMeetingInfoCommand(t *testing.T) {
	srv := newFakeServer(t)
	defer srv.Close()

	out, err := executeCommand(t,
		"get-meeting-info",
		"--url", srv.URL,
		"--secret", "test-secret",
		"--meeting-id", "test-meeting",
	)
	if err != nil {
		t.Fatalf("get-meeting-info error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "Test Meeting") {
		t.Errorf("output = %s, want meeting name", out)
	}
	if !strings.Contains(out, "John Doe") {
		t.Errorf("output = %s, want attendee name", out)
	}
}

func TestGetMeetingsCommand(t *testing.T) {
	srv := newFakeServer(t)
	defer srv.Close()

	out, err := executeCommand(t,
		"get-meetings",
		"--url", srv.URL,
		"--secret", "test-secret",
	)
	if err != nil {
		t.Fatalf("get-meetings error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "Test Meeting 1") {
		t.Errorf("output = %s, want meeting name", out)
	}
	if !strings.Contains(out, "Total: 1 meeting") {
		t.Errorf("output = %s, want total count", out)
	}
}

func TestGetRecordingsCommand(t *testing.T) {
	srv := newFakeServer(t)
	defer srv.Close()

	out, err := executeCommand(t,
		"get-recordings",
		"--url", srv.URL,
		"--secret", "test-secret",
	)
	if err != nil {
		t.Fatalf("get-recordings error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "rec-1") {
		t.Errorf("output = %s, want recording ID", out)
	}
	if !strings.Contains(out, "Total: 1 recording") {
		t.Errorf("output = %s, want total count", out)
	}
}

func TestPublishRecordingsCommand(t *testing.T) {
	srv := newFakeServer(t)
	defer srv.Close()

	out, err := executeCommand(t,
		"publish-recordings",
		"--url", srv.URL,
		"--secret", "test-secret",
		"--record-id", "rec-1",
		"--publish", "true",
	)
	if err != nil {
		t.Fatalf("publish-recordings error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "published successfully") {
		t.Errorf("output = %s, want success message", out)
	}
}

func TestDeleteRecordingsCommand(t *testing.T) {
	srv := newFakeServer(t)
	defer srv.Close()

	out, err := executeCommand(t,
		"delete-recordings",
		"--url", srv.URL,
		"--secret", "test-secret",
		"--record-id", "rec-1",
	)
	if err != nil {
		t.Fatalf("delete-recordings error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "deleted successfully") {
		t.Errorf("output = %s, want success message", out)
	}
}

func TestUpdateRecordingsCommand(t *testing.T) {
	srv := newFakeServer(t)
	defer srv.Close()

	out, err := executeCommand(t,
		"update-recordings",
		"--url", srv.URL,
		"--secret", "test-secret",
		"--record-id", "rec-1",
		"--meta", "name=New Title",
	)
	if err != nil {
		t.Fatalf("update-recordings error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "updated successfully") {
		t.Errorf("output = %s, want success message", out)
	}
}

func TestGetTextTracksCommand(t *testing.T) {
	srv := newFakeServer(t)
	defer srv.Close()

	out, err := executeCommand(t,
		"get-text-tracks",
		"--url", srv.URL,
		"--secret", "test-secret",
		"--record-id", "rec-1",
	)
	if err != nil {
		t.Fatalf("get-text-tracks error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "English") {
		t.Errorf("output = %s, want track label", out)
	}
	if !strings.Contains(out, "subtitles") {
		t.Errorf("output = %s, want track kind", out)
	}
}

func TestSendChatMessageCommand(t *testing.T) {
	srv := newFakeServer(t)
	defer srv.Close()

	out, err := executeCommand(t,
		"send-chat-message",
		"--url", srv.URL,
		"--secret", "test-secret",
		"--meeting-id", "test-meeting",
		"--message", "Hello everyone",
	)
	if err != nil {
		t.Fatalf("send-chat-message error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "Chat message sent") {
		t.Errorf("output = %s, want success message", out)
	}
}

func TestGetJoinUrlCommand(t *testing.T) {
	srv := newFakeServer(t)
	defer srv.Close()

	out, err := executeCommand(t,
		"get-join-url",
		"--url", srv.URL,
		"--secret", "test-secret",
		"--session-token", "tok-123",
	)
	if err != nil {
		t.Fatalf("get-join-url error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "https://bbb.example.com/join") {
		t.Errorf("output = %s, want join URL", out)
	}
}

func TestInsertDocumentCommand(t *testing.T) {
	srv := newFakeServer(t)
	defer srv.Close()

	out, err := executeCommand(t,
		"insert-document",
		"--url", srv.URL,
		"--secret", "test-secret",
		"--meeting-id", "test-meeting",
		"--document-url", "https://example.com/slides.pdf",
	)
	if err != nil {
		t.Fatalf("insert-document error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "Document inserted") {
		t.Errorf("output = %s, want success message", out)
	}
}

func TestHelpOutput(t *testing.T) {
	out, err := executeCommand(t, "--help")
	if err != nil {
		t.Fatalf("help error: %v", err)
	}
	// Verify all commands are listed
	commands := []string{
		"create", "create-join", "join", "end", "insert-document",
		"is-meeting-running", "get-meeting-info", "get-meetings",
		"get-recordings", "publish-recordings", "delete-recordings",
		"update-recordings", "get-text-tracks", "put-text-track",
		"send-chat-message", "get-join-url", "version",
	}
	for _, cmd := range commands {
		if !strings.Contains(out, cmd) {
			t.Errorf("help output missing command: %s", cmd)
		}
	}
}

func TestCommandHelp(t *testing.T) {
	commands := []string{
		"create", "create-join", "join", "end", "insert-document",
		"is-meeting-running", "get-meeting-info", "get-meetings",
		"get-recordings", "publish-recordings", "delete-recordings",
		"update-recordings", "get-text-tracks", "put-text-track",
		"send-chat-message", "get-join-url",
	}

	for _, cmd := range commands {
		t.Run(cmd, func(t *testing.T) {
			out, err := executeCommand(t, cmd, "--help")
			if err != nil {
				t.Fatalf("%s --help error: %v", cmd, err)
			}
			if out == "" {
				t.Errorf("%s --help produced no output", cmd)
			}
			// Each command should have examples
			if !strings.Contains(out, "bigbluebutton-cli "+cmd) {
				t.Errorf("%s help missing usage example", cmd)
			}
		})
	}
}
