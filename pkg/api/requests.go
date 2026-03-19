package api

import (
	"net/url"
	"strconv"
)

// Helper to set a string pointer value into url.Values if non-nil.
func setOptionalString(v url.Values, key string, val *string) {
	if val != nil {
		v.Set(key, *val)
	}
}

func setOptionalBool(v url.Values, key string, val *bool) {
	if val != nil {
		v.Set(key, strconv.FormatBool(*val))
	}
}

func setOptionalInt(v url.Values, key string, val *int) {
	if val != nil {
		v.Set(key, strconv.Itoa(*val))
	}
}

// Ptr helpers for building requests.
func StringPtr(s string) *string { return &s }
func BoolPtr(b bool) *bool       { return &b }
func IntPtr(i int) *int          { return &i }

// --- Administration Requests ---

type CreateRequest struct {
	Name      string
	MeetingID string

	// Optional
	AttendeePW                     *string
	ModeratorPW                    *string
	Welcome                        *string
	DialNumber                     *string
	VoiceBridge                    *string
	MaxParticipants                *int
	LogoutURL                      *string
	Record                         *bool
	Duration                       *int
	IsBreakout                     *bool
	ParentMeetingID                *string
	Sequence                       *int
	FreeJoin                       *bool
	BreakoutRoomsPrivateChatEnabled *bool
	BreakoutRoomsRecord            *bool
	ModeratorOnlyMessage           *string
	AutoStartRecording             *bool
	AllowStartStopRecording        *bool
	WebcamsOnlyForModerator        *bool
	BannerText                     *string
	BannerColor                    *string
	MuteOnStart                    *bool
	AllowModsToUnmuteUsers         *bool
	LockSettingsDisableCam         *bool
	LockSettingsDisableMic         *bool
	LockSettingsDisablePrivateChat *bool
	LockSettingsDisablePublicChat  *bool
	LockSettingsDisableNotes       *bool
	LockSettingsHideUserList       *bool
	LockSettingsLockOnJoin         *bool
	LockSettingsLockOnJoinConfigurable *bool
	LockSettingsHideViewersCursor  *bool
	GuestPolicy                    *string
	MeetingKeepEvents              *bool
	EndWhenNoModerator             *bool
	EndWhenNoModeratorDelayInMinutes *int
	MeetingLayout                  *string
	LearningDashboardCleanupDelayInMinutes *int
	AllowModsToEjectCameras        *bool
	AllowRequestsWithoutSession    *bool
	UserCameraCap                  *int
	MeetingCameraCap               *int
	MeetingExpireIfNoUserJoinedInMinutes *int
	MeetingExpireWhenLastUserLeftInMinutes *int
	Logo                           *string
	DisabledFeatures               *string
	NotifyRecordingIsOn            *bool
	PresentationUploadExternalUrl  *string
	PresentationUploadExternalDescription *string
	RecordFullDurationMedia        *bool
	PreUploadedPresentation        *string
	PreUploadedPresentationName    *string
	Meta                           map[string]string
}

func (r *CreateRequest) ToValues() url.Values {
	v := url.Values{}
	v.Set("name", r.Name)
	v.Set("meetingID", r.MeetingID)

	setOptionalString(v, "attendeePW", r.AttendeePW)
	setOptionalString(v, "moderatorPW", r.ModeratorPW)
	setOptionalString(v, "welcome", r.Welcome)
	setOptionalString(v, "dialNumber", r.DialNumber)
	setOptionalString(v, "voiceBridge", r.VoiceBridge)
	setOptionalInt(v, "maxParticipants", r.MaxParticipants)
	setOptionalString(v, "logoutURL", r.LogoutURL)
	setOptionalBool(v, "record", r.Record)
	setOptionalInt(v, "duration", r.Duration)
	setOptionalBool(v, "isBreakout", r.IsBreakout)
	setOptionalString(v, "parentMeetingID", r.ParentMeetingID)
	setOptionalInt(v, "sequence", r.Sequence)
	setOptionalBool(v, "freeJoin", r.FreeJoin)
	setOptionalBool(v, "breakoutRoomsPrivateChatEnabled", r.BreakoutRoomsPrivateChatEnabled)
	setOptionalBool(v, "breakoutRoomsRecord", r.BreakoutRoomsRecord)
	setOptionalString(v, "moderatorOnlyMessage", r.ModeratorOnlyMessage)
	setOptionalBool(v, "autoStartRecording", r.AutoStartRecording)
	setOptionalBool(v, "allowStartStopRecording", r.AllowStartStopRecording)
	setOptionalBool(v, "webcamsOnlyForModerator", r.WebcamsOnlyForModerator)
	setOptionalString(v, "bannerText", r.BannerText)
	setOptionalString(v, "bannerColor", r.BannerColor)
	setOptionalBool(v, "muteOnStart", r.MuteOnStart)
	setOptionalBool(v, "allowModsToUnmuteUsers", r.AllowModsToUnmuteUsers)
	setOptionalBool(v, "lockSettingsDisableCam", r.LockSettingsDisableCam)
	setOptionalBool(v, "lockSettingsDisableMic", r.LockSettingsDisableMic)
	setOptionalBool(v, "lockSettingsDisablePrivateChat", r.LockSettingsDisablePrivateChat)
	setOptionalBool(v, "lockSettingsDisablePublicChat", r.LockSettingsDisablePublicChat)
	setOptionalBool(v, "lockSettingsDisableNotes", r.LockSettingsDisableNotes)
	setOptionalBool(v, "lockSettingsHideUserList", r.LockSettingsHideUserList)
	setOptionalBool(v, "lockSettingsLockOnJoin", r.LockSettingsLockOnJoin)
	setOptionalBool(v, "lockSettingsLockOnJoinConfigurable", r.LockSettingsLockOnJoinConfigurable)
	setOptionalBool(v, "lockSettingsHideViewersCursor", r.LockSettingsHideViewersCursor)
	setOptionalString(v, "guestPolicy", r.GuestPolicy)
	setOptionalBool(v, "meetingKeepEvents", r.MeetingKeepEvents)
	setOptionalBool(v, "endWhenNoModerator", r.EndWhenNoModerator)
	setOptionalInt(v, "endWhenNoModeratorDelayInMinutes", r.EndWhenNoModeratorDelayInMinutes)
	setOptionalString(v, "meetingLayout", r.MeetingLayout)
	setOptionalInt(v, "learningDashboardCleanupDelayInMinutes", r.LearningDashboardCleanupDelayInMinutes)
	setOptionalBool(v, "allowModsToEjectCameras", r.AllowModsToEjectCameras)
	setOptionalBool(v, "allowRequestsWithoutSession", r.AllowRequestsWithoutSession)
	setOptionalInt(v, "userCameraCap", r.UserCameraCap)
	setOptionalInt(v, "meetingCameraCap", r.MeetingCameraCap)
	setOptionalInt(v, "meetingExpireIfNoUserJoinedInMinutes", r.MeetingExpireIfNoUserJoinedInMinutes)
	setOptionalInt(v, "meetingExpireWhenLastUserLeftInMinutes", r.MeetingExpireWhenLastUserLeftInMinutes)
	setOptionalString(v, "logo", r.Logo)
	setOptionalString(v, "disabledFeatures", r.DisabledFeatures)
	setOptionalBool(v, "notifyRecordingIsOn", r.NotifyRecordingIsOn)
	setOptionalString(v, "presentationUploadExternalUrl", r.PresentationUploadExternalUrl)
	setOptionalString(v, "presentationUploadExternalDescription", r.PresentationUploadExternalDescription)
	setOptionalBool(v, "recordFullDurationMedia", r.RecordFullDurationMedia)
	setOptionalString(v, "preUploadedPresentation", r.PreUploadedPresentation)
	setOptionalString(v, "preUploadedPresentationName", r.PreUploadedPresentationName)

	for key, val := range r.Meta {
		v.Set("meta_"+key, val)
	}

	return v
}

type JoinRequest struct {
	FullName  string
	MeetingID string
	Role      string

	// Optional
	Password            *string
	CreateTime          *string
	UserID              *string
	WebVoiceConf        *string
	AvatarURL           *string
	Redirect            *string
	ErrorRedirectUrl    *string
	LogoutURL           *string
	Guest               *string
	ExcludeFromDashboard *string
	EnforceLayout       *string
}

func (r *JoinRequest) ToValues() url.Values {
	v := url.Values{}
	v.Set("fullName", r.FullName)
	v.Set("meetingID", r.MeetingID)
	if r.Role != "" {
		v.Set("role", r.Role)
	}

	setOptionalString(v, "password", r.Password)
	setOptionalString(v, "createTime", r.CreateTime)
	setOptionalString(v, "userID", r.UserID)
	setOptionalString(v, "webVoiceConf", r.WebVoiceConf)
	setOptionalString(v, "avatarURL", r.AvatarURL)
	setOptionalString(v, "redirect", r.Redirect)
	setOptionalString(v, "errorRedirectUrl", r.ErrorRedirectUrl)
	setOptionalString(v, "logoutURL", r.LogoutURL)
	setOptionalString(v, "guest", r.Guest)
	setOptionalString(v, "excludeFromDashboard", r.ExcludeFromDashboard)
	setOptionalString(v, "enforceLayout", r.EnforceLayout)

	return v
}

type EndRequest struct {
	MeetingID string
}

func (r *EndRequest) ToValues() url.Values {
	v := url.Values{}
	v.Set("meetingID", r.MeetingID)
	return v
}

type InsertDocumentRequest struct {
	MeetingID string
	Body      string // XML body content
}

func (r *InsertDocumentRequest) ToValues() url.Values {
	v := url.Values{}
	v.Set("meetingID", r.MeetingID)
	return v
}

// --- Monitoring Requests ---

type IsMeetingRunningRequest struct {
	MeetingID string
}

func (r *IsMeetingRunningRequest) ToValues() url.Values {
	v := url.Values{}
	v.Set("meetingID", r.MeetingID)
	return v
}

type GetMeetingInfoRequest struct {
	MeetingID string
}

func (r *GetMeetingInfoRequest) ToValues() url.Values {
	v := url.Values{}
	v.Set("meetingID", r.MeetingID)
	return v
}

// GetMeetingsRequest has no parameters beyond checksum.
type GetMeetingsRequest struct{}

func (r *GetMeetingsRequest) ToValues() url.Values {
	return url.Values{}
}

// --- Recording Requests ---

type GetRecordingsRequest struct {
	MeetingID *string
	RecordID  *string
	State     *string
	Offset    *int
	Limit     *int
	Meta      map[string]string
}

func (r *GetRecordingsRequest) ToValues() url.Values {
	v := url.Values{}
	setOptionalString(v, "meetingID", r.MeetingID)
	setOptionalString(v, "recordID", r.RecordID)
	setOptionalString(v, "state", r.State)
	setOptionalInt(v, "offset", r.Offset)
	setOptionalInt(v, "limit", r.Limit)
	for key, val := range r.Meta {
		v.Set("meta_"+key, val)
	}
	return v
}

type PublishRecordingsRequest struct {
	RecordID string
	Publish  bool
}

func (r *PublishRecordingsRequest) ToValues() url.Values {
	v := url.Values{}
	v.Set("recordID", r.RecordID)
	v.Set("publish", strconv.FormatBool(r.Publish))
	return v
}

type DeleteRecordingsRequest struct {
	RecordID string
}

func (r *DeleteRecordingsRequest) ToValues() url.Values {
	v := url.Values{}
	v.Set("recordID", r.RecordID)
	return v
}

type UpdateRecordingsRequest struct {
	RecordID string
	Meta     map[string]string
}

func (r *UpdateRecordingsRequest) ToValues() url.Values {
	v := url.Values{}
	v.Set("recordID", r.RecordID)
	for key, val := range r.Meta {
		v.Set("meta_"+key, val)
	}
	return v
}

type GetRecordingTextTracksRequest struct {
	RecordID string
}

func (r *GetRecordingTextTracksRequest) ToValues() url.Values {
	v := url.Values{}
	v.Set("recordID", r.RecordID)
	return v
}

type PutRecordingTextTrackRequest struct {
	RecordID string
	Kind     string
	Lang     string
	Label    string
	FilePath string // local path to subtitle file
}

func (r *PutRecordingTextTrackRequest) ToValues() url.Values {
	v := url.Values{}
	v.Set("recordID", r.RecordID)
	v.Set("kind", r.Kind)
	v.Set("lang", r.Lang)
	v.Set("label", r.Label)
	return v
}

// --- Engagement Requests ---

type SendChatMessageRequest struct {
	MeetingID string
	Message   string
	UserName  *string
}

func (r *SendChatMessageRequest) ToValues() url.Values {
	v := url.Values{}
	v.Set("meetingID", r.MeetingID)
	v.Set("message", r.Message)
	setOptionalString(v, "userName", r.UserName)
	return v
}

type GetJoinUrlRequest struct {
	SessionToken   string
	ReplaceSession *bool
	SessionName    *string
}

func (r *GetJoinUrlRequest) ToValues() url.Values {
	v := url.Values{}
	v.Set("sessionToken", r.SessionToken)
	setOptionalBool(v, "replaceSession", r.ReplaceSession)
	setOptionalString(v, "sessionName", r.SessionName)
	return v
}
