package api

import (
	"encoding/json"
	"encoding/xml"
	"testing"
)

func TestBaseResponse_IsSuccess(t *testing.T) {
	tests := []struct {
		returnCode string
		want       bool
	}{
		{"SUCCESS", true},
		{"FAILED", false},
		{"", false},
	}
	for _, tt := range tests {
		r := &BaseResponse{ReturnCode: tt.returnCode}
		if got := r.IsSuccess(); got != tt.want {
			t.Errorf("IsSuccess() for %q = %v, want %v", tt.returnCode, got, tt.want)
		}
	}
}

func TestBaseResponse_ToError(t *testing.T) {
	r := &BaseResponse{ReturnCode: "SUCCESS"}
	if err := r.ToError(); err != nil {
		t.Errorf("SUCCESS should not return error, got %v", err)
	}

	r = &BaseResponse{ReturnCode: "FAILED", MessageKey: "notFound", Message: "Meeting not found"}
	err := r.ToError()
	if err == nil {
		t.Fatal("FAILED should return error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.MessageKey != "notFound" {
		t.Errorf("MessageKey = %s, want notFound", apiErr.MessageKey)
	}
}

func TestCreateResponseParsing(t *testing.T) {
	xmlData := `<response>
		<returncode>SUCCESS</returncode>
		<meetingID>test-123</meetingID>
		<internalMeetingID>internal-456</internalMeetingID>
		<attendeePW>ap</attendeePW>
		<moderatorPW>mp</moderatorPW>
		<createTime>1699900000000</createTime>
		<voiceBridge>70001</voiceBridge>
		<dialNumber>555-1234</dialNumber>
		<createDate>Mon Nov 13 12:00:00 UTC 2023</createDate>
		<hasUserJoined>false</hasUserJoined>
		<duration>60</duration>
		<hasBeenForciblyEnded>false</hasBeenForciblyEnded>
	</response>`

	var resp CreateResponse
	if err := xml.Unmarshal([]byte(xmlData), &resp); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if !resp.IsSuccess() {
		t.Error("expected SUCCESS")
	}
	if resp.MeetingID != "test-123" {
		t.Errorf("MeetingID = %s, want test-123", resp.MeetingID)
	}
	if resp.InternalMeetingID != "internal-456" {
		t.Errorf("InternalMeetingID = %s, want internal-456", resp.InternalMeetingID)
	}
	if resp.Duration != 60 {
		t.Errorf("Duration = %d, want 60", resp.Duration)
	}
	if resp.VoiceBridge != "70001" {
		t.Errorf("VoiceBridge = %s, want 70001", resp.VoiceBridge)
	}
}

func TestGetMeetingsResponseParsing(t *testing.T) {
	xmlData := `<response>
		<returncode>SUCCESS</returncode>
		<meetings>
			<meeting>
				<meetingName>Meeting 1</meetingName>
				<meetingID>m1</meetingID>
				<running>true</running>
				<participantCount>5</participantCount>
				<recording>true</recording>
			</meeting>
			<meeting>
				<meetingName>Meeting 2</meetingName>
				<meetingID>m2</meetingID>
				<running>false</running>
				<participantCount>0</participantCount>
				<recording>false</recording>
			</meeting>
		</meetings>
	</response>`

	var resp GetMeetingsResponse
	if err := xml.Unmarshal([]byte(xmlData), &resp); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if len(resp.Meetings) != 2 {
		t.Fatalf("expected 2 meetings, got %d", len(resp.Meetings))
	}
	if resp.Meetings[0].MeetingName != "Meeting 1" {
		t.Errorf("first meeting name = %s, want Meeting 1", resp.Meetings[0].MeetingName)
	}
	if !resp.Meetings[0].Running {
		t.Error("first meeting should be running")
	}
	if resp.Meetings[0].ParticipantCount != 5 {
		t.Errorf("participant count = %d, want 5", resp.Meetings[0].ParticipantCount)
	}
}

func TestGetMeetingInfoResponseParsing(t *testing.T) {
	xmlData := `<response>
		<returncode>SUCCESS</returncode>
		<meetingName>Test</meetingName>
		<meetingID>test-1</meetingID>
		<running>true</running>
		<participantCount>3</participantCount>
		<attendees>
			<attendee>
				<userID>u1</userID>
				<fullName>John</fullName>
				<role>MODERATOR</role>
				<isPresenter>true</isPresenter>
				<isListeningOnly>false</isListeningOnly>
				<hasJoinedVoice>true</hasJoinedVoice>
				<hasVideo>true</hasVideo>
				<clientType>HTML5</clientType>
			</attendee>
		</attendees>
	</response>`

	var resp GetMeetingInfoResponse
	if err := xml.Unmarshal([]byte(xmlData), &resp); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if len(resp.Attendees) != 1 {
		t.Fatalf("expected 1 attendee, got %d", len(resp.Attendees))
	}
	if resp.Attendees[0].FullName != "John" {
		t.Errorf("attendee name = %s, want John", resp.Attendees[0].FullName)
	}
	if resp.Attendees[0].Role != "MODERATOR" {
		t.Errorf("attendee role = %s, want MODERATOR", resp.Attendees[0].Role)
	}
}

func TestGetRecordingsResponseParsing(t *testing.T) {
	xmlData := `<response>
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
				<playback>
					<format>
						<type>presentation</type>
						<url>https://example.com/playback/rec-1</url>
						<length>60</length>
					</format>
				</playback>
			</recording>
		</recordings>
	</response>`

	var resp GetRecordingsResponse
	if err := xml.Unmarshal([]byte(xmlData), &resp); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if len(resp.Recordings) != 1 {
		t.Fatalf("expected 1 recording, got %d", len(resp.Recordings))
	}
	rec := resp.Recordings[0]
	if rec.RecordID != "rec-1" {
		t.Errorf("RecordID = %s, want rec-1", rec.RecordID)
	}
	if !rec.Published {
		t.Error("recording should be published")
	}
	if len(rec.Playback) != 1 {
		t.Fatalf("expected 1 playback format, got %d", len(rec.Playback))
	}
	if rec.Playback[0].Type != "presentation" {
		t.Errorf("playback type = %s, want presentation", rec.Playback[0].Type)
	}
}

func TestGetRecordingTextTracksResponseParsing(t *testing.T) {
	jsonData := `{"response":{"returncode":"SUCCESS","tracks":[{"href":"https://example.com/track.vtt","kind":"subtitles","label":"English","lang":"en","source":"upload"}]}}`

	var resp GetRecordingTextTracksResponse
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if resp.Response.ReturnCode != "SUCCESS" {
		t.Errorf("ReturnCode = %s, want SUCCESS", resp.Response.ReturnCode)
	}
	if len(resp.Response.Tracks) != 1 {
		t.Fatalf("expected 1 track, got %d", len(resp.Response.Tracks))
	}
	track := resp.Response.Tracks[0]
	if track.Lang != "en" {
		t.Errorf("Lang = %s, want en", track.Lang)
	}
	if track.Kind != "subtitles" {
		t.Errorf("Kind = %s, want subtitles", track.Kind)
	}
}

func TestErrorResponseParsing(t *testing.T) {
	xmlData := `<response>
		<returncode>FAILED</returncode>
		<messageKey>notFound</messageKey>
		<message>We could not find a meeting with that meeting ID</message>
	</response>`

	var resp BaseResponse
	if err := xml.Unmarshal([]byte(xmlData), &resp); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if resp.IsSuccess() {
		t.Error("should not be success")
	}
	if resp.MessageKey != "notFound" {
		t.Errorf("MessageKey = %s, want notFound", resp.MessageKey)
	}
}
