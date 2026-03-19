package api

import (
	"testing"
)

func TestClient_IsMeetingRunning(t *testing.T) {
	tests := []struct {
		name     string
		response string
		want     bool
	}{
		{
			name: "meeting running",
			response: `<response><returncode>SUCCESS</returncode><running>true</running></response>`,
			want: true,
		},
		{
			name: "meeting not running",
			response: `<response><returncode>SUCCESS</returncode><running>false</running></response>`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := newTestServer(t, "isMeetingRunning", tt.response)
			defer srv.Close()

			client := newTestClient(srv.URL, "test-secret")
			resp, err := client.IsMeetingRunning(&IsMeetingRunningRequest{MeetingID: "test"})
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			if resp.Running != tt.want {
				t.Errorf("Running = %v, want %v", resp.Running, tt.want)
			}
		})
	}
}

func TestClient_GetMeetingInfo(t *testing.T) {
	srv := newTestServer(t, "getMeetingInfo", `<response>
		<returncode>SUCCESS</returncode>
		<meetingName>Test Meeting</meetingName>
		<meetingID>test-1</meetingID>
		<running>true</running>
		<participantCount>3</participantCount>
		<moderatorCount>1</moderatorCount>
		<attendees>
			<attendee>
				<userID>u1</userID>
				<fullName>John</fullName>
				<role>MODERATOR</role>
				<isPresenter>true</isPresenter>
			</attendee>
			<attendee>
				<userID>u2</userID>
				<fullName>Jane</fullName>
				<role>VIEWER</role>
				<isPresenter>false</isPresenter>
			</attendee>
		</attendees>
	</response>`)
	defer srv.Close()

	client := newTestClient(srv.URL, "test-secret")
	resp, err := client.GetMeetingInfo(&GetMeetingInfoRequest{MeetingID: "test-1"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if resp.MeetingName != "Test Meeting" {
		t.Errorf("MeetingName = %s, want Test Meeting", resp.MeetingName)
	}
	if !resp.Running {
		t.Error("meeting should be running")
	}
	if resp.ParticipantCount != 3 {
		t.Errorf("ParticipantCount = %d, want 3", resp.ParticipantCount)
	}
	if len(resp.Attendees) != 2 {
		t.Fatalf("expected 2 attendees, got %d", len(resp.Attendees))
	}
	if resp.Attendees[0].Role != "MODERATOR" {
		t.Errorf("first attendee role = %s, want MODERATOR", resp.Attendees[0].Role)
	}
}

func TestClient_GetMeetings(t *testing.T) {
	srv := newTestServer(t, "getMeetings", `<response>
		<returncode>SUCCESS</returncode>
		<meetings>
			<meeting>
				<meetingName>Meeting 1</meetingName>
				<meetingID>m1</meetingID>
				<running>true</running>
				<participantCount>5</participantCount>
			</meeting>
			<meeting>
				<meetingName>Meeting 2</meetingName>
				<meetingID>m2</meetingID>
				<running>true</running>
				<participantCount>10</participantCount>
			</meeting>
		</meetings>
	</response>`)
	defer srv.Close()

	client := newTestClient(srv.URL, "test-secret")
	resp, err := client.GetMeetings()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if len(resp.Meetings) != 2 {
		t.Fatalf("expected 2 meetings, got %d", len(resp.Meetings))
	}
	if resp.Meetings[0].MeetingID != "m1" {
		t.Errorf("first meeting ID = %s, want m1", resp.Meetings[0].MeetingID)
	}
	if resp.Meetings[1].ParticipantCount != 10 {
		t.Errorf("second meeting participants = %d, want 10", resp.Meetings[1].ParticipantCount)
	}
}

func TestClient_GetMeetings_Empty(t *testing.T) {
	srv := newTestServer(t, "getMeetings", `<response>
		<returncode>SUCCESS</returncode>
		<meetings></meetings>
	</response>`)
	defer srv.Close()

	client := newTestClient(srv.URL, "test-secret")
	resp, err := client.GetMeetings()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if len(resp.Meetings) != 0 {
		t.Errorf("expected 0 meetings, got %d", len(resp.Meetings))
	}
}
