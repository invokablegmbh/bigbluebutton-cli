package api

import "encoding/xml"

// BaseResponse is the common response structure for all BBB API calls.
type BaseResponse struct {
	XMLName    xml.Name `xml:"response"`
	ReturnCode string   `xml:"returncode"`
	MessageKey string   `xml:"messageKey"`
	Message    string   `xml:"message"`
}

func (r *BaseResponse) IsSuccess() bool {
	return r.ReturnCode == "SUCCESS"
}

func (r *BaseResponse) ToError() error {
	if r.IsSuccess() {
		return nil
	}
	return &APIError{
		ReturnCode: r.ReturnCode,
		MessageKey: r.MessageKey,
		Message:    r.Message,
	}
}

// --- Administration Responses ---

type CreateResponse struct {
	BaseResponse
	MeetingID         string `xml:"meetingID"`
	InternalMeetingID string `xml:"internalMeetingID"`
	ParentMeetingID   string `xml:"parentMeetingID"`
	AttendeePW        string `xml:"attendeePW"`
	ModeratorPW       string `xml:"moderatorPW"`
	CreateTime        int64  `xml:"createTime"`
	VoiceBridge       string `xml:"voiceBridge"`
	DialNumber        string `xml:"dialNumber"`
	CreateDate        string `xml:"createDate"`
	HasUserJoined     bool   `xml:"hasUserJoined"`
	Duration          int    `xml:"duration"`
	HasBeenForciblyEnded bool `xml:"hasBeenForciblyEnded"`
}

type JoinResponse struct {
	BaseResponse
	MeetingID    string `xml:"meeting_id"`
	UserID       string `xml:"user_id"`
	AuthToken    string `xml:"auth_token"`
	SessionToken string `xml:"session_token"`
	URL          string `xml:"url"`
}

type EndResponse struct {
	BaseResponse
}

type InsertDocumentResponse struct {
	BaseResponse
}

// --- Monitoring Responses ---

type IsMeetingRunningResponse struct {
	BaseResponse
	Running bool `xml:"running"`
}

type Attendee struct {
	UserID          string `xml:"userID" json:"userID"`
	FullName        string `xml:"fullName" json:"fullName"`
	Role            string `xml:"role" json:"role"`
	IsPresenter     bool   `xml:"isPresenter" json:"isPresenter"`
	IsListeningOnly bool   `xml:"isListeningOnly" json:"isListeningOnly"`
	HasJoinedVoice  bool   `xml:"hasJoinedVoice" json:"hasJoinedVoice"`
	HasVideo        bool   `xml:"hasVideo" json:"hasVideo"`
	ClientType      string `xml:"clientType" json:"clientType"`
}

type Meeting struct {
	MeetingName           string     `xml:"meetingName" json:"meetingName"`
	MeetingID             string     `xml:"meetingID" json:"meetingID"`
	InternalMeetingID     string     `xml:"internalMeetingID" json:"internalMeetingID"`
	CreateTime            int64      `xml:"createTime" json:"createTime"`
	CreateDate            string     `xml:"createDate" json:"createDate"`
	VoiceBridge           string     `xml:"voiceBridge" json:"voiceBridge"`
	DialNumber            string     `xml:"dialNumber" json:"dialNumber"`
	AttendeePW            string     `xml:"attendeePW" json:"attendeePW"`
	ModeratorPW           string     `xml:"moderatorPW" json:"moderatorPW"`
	Running               bool       `xml:"running" json:"running"`
	Duration              int        `xml:"duration" json:"duration"`
	HasUserJoined         bool       `xml:"hasUserJoined" json:"hasUserJoined"`
	Recording             bool       `xml:"recording" json:"recording"`
	HasBeenForciblyEnded  bool       `xml:"hasBeenForciblyEnded" json:"hasBeenForciblyEnded"`
	StartTime             int64      `xml:"startTime" json:"startTime"`
	EndTime               int64      `xml:"endTime" json:"endTime"`
	ParticipantCount      int        `xml:"participantCount" json:"participantCount"`
	ListenerCount         int        `xml:"listenerCount" json:"listenerCount"`
	VoiceParticipantCount int        `xml:"voiceParticipantCount" json:"voiceParticipantCount"`
	VideoCount            int        `xml:"videoCount" json:"videoCount"`
	MaxUsers              int        `xml:"maxUsers" json:"maxUsers"`
	ModeratorCount        int        `xml:"moderatorCount" json:"moderatorCount"`
	Attendees             []Attendee `xml:"attendees>attendee" json:"attendees"`
	IsBreakout            bool       `xml:"isBreakout" json:"isBreakout"`
}

type GetMeetingInfoResponse struct {
	BaseResponse
	Meeting
}

type GetMeetingsResponse struct {
	BaseResponse
	Meetings []Meeting `xml:"meetings>meeting"`
}

// --- Recording Responses ---

type PlaybackFormat struct {
	Type           string `xml:"type" json:"type"`
	URL            string `xml:"url" json:"url"`
	ProcessingTime int    `xml:"processingTime" json:"processingTime"`
	Length         int    `xml:"length" json:"length"`
}

type Recording struct {
	RecordID          string           `xml:"recordID" json:"recordID"`
	MeetingID         string           `xml:"meetingID" json:"meetingID"`
	InternalMeetingID string           `xml:"internalMeetingID" json:"internalMeetingID"`
	Name              string           `xml:"name" json:"name"`
	IsBreakout        bool             `xml:"isBreakout" json:"isBreakout"`
	Published         bool             `xml:"published" json:"published"`
	State             string           `xml:"state" json:"state"`
	StartTime         int64            `xml:"startTime" json:"startTime"`
	EndTime           int64            `xml:"endTime" json:"endTime"`
	Participants      int              `xml:"participants" json:"participants"`
	Playback          []PlaybackFormat `xml:"playback>format" json:"playback"`
}

type GetRecordingsResponse struct {
	BaseResponse
	Recordings []Recording `xml:"recordings>recording"`
}

type PublishRecordingsResponse struct {
	BaseResponse
	Published bool `xml:"published"`
}

type DeleteRecordingsResponse struct {
	BaseResponse
	Deleted bool `xml:"deleted"`
}

type UpdateRecordingsResponse struct {
	BaseResponse
	Updated bool `xml:"updated"`
}

// --- Text Track Responses (JSON) ---

type TextTrack struct {
	Href   string `json:"href" xml:"href"`
	Kind   string `json:"kind" xml:"kind"`
	Label  string `json:"label" xml:"label"`
	Lang   string `json:"lang" xml:"lang"`
	Source string `json:"source" xml:"source"`
}

type GetRecordingTextTracksResponse struct {
	Response struct {
		ReturnCode string      `json:"returncode"`
		MessageKey string      `json:"messageKey"`
		Message    string      `json:"message"`
		Tracks     []TextTrack `json:"tracks"`
	} `json:"response"`
}

type PutRecordingTextTrackResponse struct {
	Response struct {
		ReturnCode string `json:"returncode"`
		MessageKey string `json:"messageKey"`
		Message    string `json:"message"`
	} `json:"response"`
}

// --- Engagement Responses ---

type SendChatMessageResponse struct {
	BaseResponse
}

type GetJoinUrlResponse struct {
	BaseResponse
	URL string `xml:"url"`
}
