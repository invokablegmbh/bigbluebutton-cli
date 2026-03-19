package api

import (
	"testing"
)

func TestCreateRequestToValues(t *testing.T) {
	req := &CreateRequest{
		Name:      "Test Meeting",
		MeetingID: "test-123",
		Record:    BoolPtr(true),
		Duration:  IntPtr(60),
		Welcome:   StringPtr("Welcome!"),
		Meta:      map[string]string{"category": "test"},
	}

	v := req.ToValues()

	if v.Get("name") != "Test Meeting" {
		t.Errorf("name = %s, want Test Meeting", v.Get("name"))
	}
	if v.Get("meetingID") != "test-123" {
		t.Errorf("meetingID = %s, want test-123", v.Get("meetingID"))
	}
	if v.Get("record") != "true" {
		t.Errorf("record = %s, want true", v.Get("record"))
	}
	if v.Get("duration") != "60" {
		t.Errorf("duration = %s, want 60", v.Get("duration"))
	}
	if v.Get("welcome") != "Welcome!" {
		t.Errorf("welcome = %s, want Welcome!", v.Get("welcome"))
	}
	if v.Get("meta_category") != "test" {
		t.Errorf("meta_category = %s, want test", v.Get("meta_category"))
	}
}

func TestCreateRequestToValuesOptionalNil(t *testing.T) {
	req := &CreateRequest{
		Name:      "Test",
		MeetingID: "m1",
	}

	v := req.ToValues()

	// Optional fields should not be present
	if v.Get("record") != "" {
		t.Errorf("record should be empty, got %s", v.Get("record"))
	}
	if v.Get("duration") != "" {
		t.Errorf("duration should be empty, got %s", v.Get("duration"))
	}
	if v.Get("welcome") != "" {
		t.Errorf("welcome should be empty, got %s", v.Get("welcome"))
	}
}

func TestJoinRequestToValues(t *testing.T) {
	req := &JoinRequest{
		FullName:  "John Doe",
		MeetingID: "test-123",
		Role:      "MODERATOR",
		UserID:    StringPtr("user-1"),
	}

	v := req.ToValues()

	if v.Get("fullName") != "John Doe" {
		t.Errorf("fullName = %s, want John Doe", v.Get("fullName"))
	}
	if v.Get("role") != "MODERATOR" {
		t.Errorf("role = %s, want MODERATOR", v.Get("role"))
	}
	if v.Get("userID") != "user-1" {
		t.Errorf("userID = %s, want user-1", v.Get("userID"))
	}
}

func TestEndRequestToValues(t *testing.T) {
	req := &EndRequest{MeetingID: "test-123"}
	v := req.ToValues()
	if v.Get("meetingID") != "test-123" {
		t.Errorf("meetingID = %s, want test-123", v.Get("meetingID"))
	}
}

func TestGetRecordingsRequestToValues(t *testing.T) {
	req := &GetRecordingsRequest{
		MeetingID: StringPtr("m1"),
		State:     StringPtr("published"),
		Limit:     IntPtr(10),
		Offset:    IntPtr(0),
	}

	v := req.ToValues()

	if v.Get("meetingID") != "m1" {
		t.Errorf("meetingID = %s, want m1", v.Get("meetingID"))
	}
	if v.Get("state") != "published" {
		t.Errorf("state = %s, want published", v.Get("state"))
	}
	if v.Get("limit") != "10" {
		t.Errorf("limit = %s, want 10", v.Get("limit"))
	}
}

func TestPublishRecordingsRequestToValues(t *testing.T) {
	req := &PublishRecordingsRequest{RecordID: "rec-1", Publish: true}
	v := req.ToValues()
	if v.Get("recordID") != "rec-1" {
		t.Errorf("recordID = %s, want rec-1", v.Get("recordID"))
	}
	if v.Get("publish") != "true" {
		t.Errorf("publish = %s, want true", v.Get("publish"))
	}
}

func TestSendChatMessageRequestToValues(t *testing.T) {
	req := &SendChatMessageRequest{
		MeetingID: "m1",
		Message:   "Hello",
		UserName:  StringPtr("Admin"),
	}
	v := req.ToValues()
	if v.Get("meetingID") != "m1" {
		t.Errorf("meetingID = %s, want m1", v.Get("meetingID"))
	}
	if v.Get("message") != "Hello" {
		t.Errorf("message = %s, want Hello", v.Get("message"))
	}
	if v.Get("userName") != "Admin" {
		t.Errorf("userName = %s, want Admin", v.Get("userName"))
	}
}

func TestGetJoinUrlRequestToValues(t *testing.T) {
	req := &GetJoinUrlRequest{
		SessionToken:   "tok-123",
		ReplaceSession: BoolPtr(true),
		SessionName:    StringPtr("New Session"),
	}
	v := req.ToValues()
	if v.Get("sessionToken") != "tok-123" {
		t.Errorf("sessionToken = %s, want tok-123", v.Get("sessionToken"))
	}
	if v.Get("replaceSession") != "true" {
		t.Errorf("replaceSession = %s, want true", v.Get("replaceSession"))
	}
	if v.Get("sessionName") != "New Session" {
		t.Errorf("sessionName = %s, want New Session", v.Get("sessionName"))
	}
}

func TestPtrHelpers(t *testing.T) {
	s := StringPtr("hello")
	if *s != "hello" {
		t.Errorf("StringPtr = %s, want hello", *s)
	}

	b := BoolPtr(true)
	if *b != true {
		t.Error("BoolPtr should be true")
	}

	i := IntPtr(42)
	if *i != 42 {
		t.Errorf("IntPtr = %d, want 42", *i)
	}
}
