package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/invokablegmbh/bigbluebutton-cli/pkg/api"
	"github.com/invokablegmbh/bigbluebutton-cli/pkg/interactive"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new BigBlueButton meeting",
		Long: `Create a new BigBlueButton meeting with the specified parameters.

In interactive mode, you will be prompted for required parameters and
optionally for advanced settings. In non-interactive mode, provide
all required parameters via flags.

Examples:
  bigbluebutton-cli create --name "Team Meeting" --meeting-id team-2024
  bigbluebutton-cli create --name "Recorded Session" --meeting-id rec-001 --record
  bigbluebutton-cli create --name "Locked Meeting" --meeting-id locked-001 --mute-on-start --guest-policy ASK_MODERATOR`,
		RunE: withWizardHint(runCreate),
	}

	// Required
	cmd.Flags().String("name", "", "Meeting name (required)")
	cmd.Flags().String("meeting-id", "", "Unique meeting identifier (required)")

	// Common optional
	cmd.Flags().String("welcome", "", "Welcome message")
	cmd.Flags().String("dial-number", "", "Phone dial-in number")
	cmd.Flags().String("voice-bridge", "", "FreeSWITCH voice bridge number")
	cmd.Flags().Int("max-participants", 0, "Maximum number of participants")
	cmd.Flags().String("logout-url", "", "URL to redirect after logout")
	cmd.Flags().Bool("record", false, "Enable recording")
	cmd.Flags().Int("duration", 0, "Meeting duration in minutes (0 = unlimited)")
	cmd.Flags().String("moderator-only-message", "", "Message visible only to moderators")
	cmd.Flags().Bool("auto-start-recording", false, "Automatically start recording")
	cmd.Flags().Bool("allow-start-stop-recording", true, "Allow moderators to start/stop recording")
	cmd.Flags().Bool("webcams-only-for-moderator", false, "Only show webcams to moderators")
	cmd.Flags().String("banner-text", "", "Banner text displayed in the client")
	cmd.Flags().String("banner-color", "", "Banner background color (hex, e.g. #FF0000)")
	cmd.Flags().Bool("mute-on-start", false, "Mute all users on meeting start")
	cmd.Flags().Bool("allow-mods-to-unmute-users", false, "Allow moderators to unmute users")
	cmd.Flags().String("guest-policy", "", "Guest policy: ALWAYS_ACCEPT, ALWAYS_DENY, ASK_MODERATOR")
	cmd.Flags().String("meeting-layout", "", "Default layout: CUSTOM_LAYOUT, SMART_LAYOUT, PRESENTATION_FOCUS, VIDEO_FOCUS")
	cmd.Flags().Bool("end-when-no-moderator", false, "End meeting when no moderator is present")
	cmd.Flags().Int("end-when-no-moderator-delay", 0, "Minutes before ending when no moderator")
	cmd.Flags().String("logo", "", "Branding logo URL")
	cmd.Flags().String("disabled-features", "", "Comma-separated list of disabled features")
	cmd.Flags().Bool("notify-recording-is-on", false, "Show recording consent notification")

	// Lock settings
	cmd.Flags().Bool("lock-disable-cam", false, "Prevent camera sharing")
	cmd.Flags().Bool("lock-disable-mic", false, "Restrict to listen-only mode")
	cmd.Flags().Bool("lock-disable-private-chat", false, "Disable private messaging")
	cmd.Flags().Bool("lock-disable-public-chat", false, "Disable public chat")
	cmd.Flags().Bool("lock-disable-notes", false, "Disable shared notes")
	cmd.Flags().Bool("lock-hide-user-list", false, "Hide viewer list from attendees")
	cmd.Flags().Bool("lock-on-join", false, "Apply lock settings to new joiners")
	cmd.Flags().Bool("lock-hide-viewers-cursor", false, "Hide whiteboard cursors from viewers")

	// Breakout
	cmd.Flags().Bool("is-breakout", false, "Create as breakout room")
	cmd.Flags().String("parent-meeting-id", "", "Parent meeting ID (for breakout rooms)")
	cmd.Flags().Int("sequence", 0, "Breakout room sequence number")
	cmd.Flags().Bool("free-join", false, "Allow user to choose breakout room")

	// Metadata
	cmd.Flags().StringSlice("meta", nil, "Metadata key=value pairs (can specify multiple)")

	return cmd
}

// settingDef describes an additional setting for the interactive wizard.
type settingDef struct {
	label    string   // display name in multi-select
	flagName string   // corresponding CLI flag name, empty if wizard-only
	kind     string   // "bool", "string", "int", "select"
	options  []string // for "select" kind
	apply    func(req *api.CreateRequest, val interface{})
}

// additionalSettings returns all available settings for the create wizard.
func additionalSettings() []settingDef {
	return []settingDef{
		// General
		{label: "Welcome message", flagName: "welcome", kind: "string", apply: func(r *api.CreateRequest, v interface{}) { s := v.(string); r.Welcome = &s }},
		{label: "Max participants", flagName: "max-participants", kind: "int", apply: func(r *api.CreateRequest, v interface{}) { i := v.(int); r.MaxParticipants = &i }},
		{label: "Duration (minutes, 0 = unlimited)", flagName: "duration", kind: "int", apply: func(r *api.CreateRequest, v interface{}) { i := v.(int); r.Duration = &i }},
		{label: "Logout URL", flagName: "logout-url", kind: "string", apply: func(r *api.CreateRequest, v interface{}) { s := v.(string); r.LogoutURL = &s }},
		{label: "Dial number", flagName: "dial-number", kind: "string", apply: func(r *api.CreateRequest, v interface{}) { s := v.(string); r.DialNumber = &s }},
		{label: "Voice bridge", flagName: "voice-bridge", kind: "string", apply: func(r *api.CreateRequest, v interface{}) { s := v.(string); r.VoiceBridge = &s }},
		{label: "Banner text", flagName: "banner-text", kind: "string", apply: func(r *api.CreateRequest, v interface{}) { s := v.(string); r.BannerText = &s }},
		{label: "Banner color (hex)", flagName: "banner-color", kind: "string", apply: func(r *api.CreateRequest, v interface{}) { s := v.(string); r.BannerColor = &s }},
		{label: "Logo URL", flagName: "logo", kind: "string", apply: func(r *api.CreateRequest, v interface{}) { s := v.(string); r.Logo = &s }},
		{label: "Moderator-only message", flagName: "moderator-only-message", kind: "string", apply: func(r *api.CreateRequest, v interface{}) { s := v.(string); r.ModeratorOnlyMessage = &s }},
		{label: "Disabled features (comma-separated)", flagName: "disabled-features", kind: "string", apply: func(r *api.CreateRequest, v interface{}) { s := v.(string); r.DisabledFeatures = &s }},

		// Recording
		{label: "Enable recording", flagName: "record", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.Record = &b }},
		{label: "Auto-start recording", flagName: "auto-start-recording", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.AutoStartRecording = &b }},
		{label: "Allow start/stop recording", flagName: "allow-start-stop-recording", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.AllowStartStopRecording = &b }},
		{label: "Record full duration media", flagName: "", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.RecordFullDurationMedia = &b }},
		{label: "Notify recording is on", flagName: "notify-recording-is-on", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.NotifyRecordingIsOn = &b }},

		// Participants & Moderation
		{label: "Mute on start", flagName: "mute-on-start", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.MuteOnStart = &b }},
		{label: "Allow mods to unmute users", flagName: "allow-mods-to-unmute-users", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.AllowModsToUnmuteUsers = &b }},
		{label: "Allow mods to eject cameras", flagName: "", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.AllowModsToEjectCameras = &b }},
		{label: "Webcams only for moderator", flagName: "webcams-only-for-moderator", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.WebcamsOnlyForModerator = &b }},
		{label: "Guest policy", flagName: "guest-policy", kind: "select", options: []string{"ALWAYS_ACCEPT", "ALWAYS_DENY", "ASK_MODERATOR"}, apply: func(r *api.CreateRequest, v interface{}) { s := v.(string); r.GuestPolicy = &s }},
		{label: "Meeting layout", flagName: "meeting-layout", kind: "select", options: []string{"CUSTOM_LAYOUT", "SMART_LAYOUT", "PRESENTATION_FOCUS", "VIDEO_FOCUS"}, apply: func(r *api.CreateRequest, v interface{}) { s := v.(string); r.MeetingLayout = &s }},
		{label: "End when no moderator", flagName: "end-when-no-moderator", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.EndWhenNoModerator = &b }},
		{label: "End when no moderator delay (minutes)", flagName: "end-when-no-moderator-delay", kind: "int", apply: func(r *api.CreateRequest, v interface{}) { i := v.(int); r.EndWhenNoModeratorDelayInMinutes = &i }},
		{label: "User camera cap", flagName: "", kind: "int", apply: func(r *api.CreateRequest, v interface{}) { i := v.(int); r.UserCameraCap = &i }},
		{label: "Meeting camera cap", flagName: "", kind: "int", apply: func(r *api.CreateRequest, v interface{}) { i := v.(int); r.MeetingCameraCap = &i }},

		// Lock settings
		{label: "Lock: disable camera", flagName: "lock-disable-cam", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.LockSettingsDisableCam = &b }},
		{label: "Lock: disable microphone", flagName: "lock-disable-mic", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.LockSettingsDisableMic = &b }},
		{label: "Lock: disable private chat", flagName: "lock-disable-private-chat", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.LockSettingsDisablePrivateChat = &b }},
		{label: "Lock: disable public chat", flagName: "lock-disable-public-chat", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.LockSettingsDisablePublicChat = &b }},
		{label: "Lock: disable shared notes", flagName: "lock-disable-notes", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.LockSettingsDisableNotes = &b }},
		{label: "Lock: hide user list", flagName: "lock-hide-user-list", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.LockSettingsHideUserList = &b }},
		{label: "Lock: lock on join", flagName: "lock-on-join", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.LockSettingsLockOnJoin = &b }},
		{label: "Lock: hide viewers cursor", flagName: "lock-hide-viewers-cursor", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.LockSettingsHideViewersCursor = &b }},

		// Breakout
		{label: "Is breakout room", flagName: "is-breakout", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.IsBreakout = &b }},
		{label: "Parent meeting ID", flagName: "parent-meeting-id", kind: "string", apply: func(r *api.CreateRequest, v interface{}) { s := v.(string); r.ParentMeetingID = &s }},
		{label: "Breakout sequence number", flagName: "sequence", kind: "int", apply: func(r *api.CreateRequest, v interface{}) { i := v.(int); r.Sequence = &i }},
		{label: "Free join (breakout)", flagName: "free-join", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.FreeJoin = &b }},
		{label: "Breakout rooms private chat", flagName: "", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.BreakoutRoomsPrivateChatEnabled = &b }},
		{label: "Breakout rooms record", flagName: "", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.BreakoutRoomsRecord = &b }},

		// Expiry & Misc
		{label: "Meeting expire if no user joined (minutes)", flagName: "", kind: "int", apply: func(r *api.CreateRequest, v interface{}) { i := v.(int); r.MeetingExpireIfNoUserJoinedInMinutes = &i }},
		{label: "Meeting expire when last user left (minutes)", flagName: "", kind: "int", apply: func(r *api.CreateRequest, v interface{}) { i := v.(int); r.MeetingExpireWhenLastUserLeftInMinutes = &i }},
		{label: "Learning dashboard cleanup delay (minutes)", flagName: "", kind: "int", apply: func(r *api.CreateRequest, v interface{}) { i := v.(int); r.LearningDashboardCleanupDelayInMinutes = &i }},
		{label: "Keep meeting events", flagName: "", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.MeetingKeepEvents = &b }},
		{label: "Allow requests without session", flagName: "", kind: "bool", apply: func(r *api.CreateRequest, v interface{}) { b := v.(bool); r.AllowRequestsWithoutSession = &b }},
		{label: "Pre-uploaded presentation URL", flagName: "", kind: "string", apply: func(r *api.CreateRequest, v interface{}) { s := v.(string); r.PreUploadedPresentation = &s }},
		{label: "Pre-uploaded presentation name", flagName: "", kind: "string", apply: func(r *api.CreateRequest, v interface{}) { s := v.(string); r.PreUploadedPresentationName = &s }},
		{label: "Presentation upload external URL", flagName: "", kind: "string", apply: func(r *api.CreateRequest, v interface{}) { s := v.(string); r.PresentationUploadExternalUrl = &s }},
		{label: "Presentation upload external description", flagName: "", kind: "string", apply: func(r *api.CreateRequest, v interface{}) { s := v.(string); r.PresentationUploadExternalDescription = &s }},
	}
}

// promptForSetting asks the user for a value based on the setting kind.
func promptForSetting(s settingDef) (interface{}, error) {
	switch s.kind {
	case "bool":
		return interactive.AskConfirm(s.label, false)
	case "string":
		return interactive.AskString(s.label, "")
	case "int":
		return interactive.AskInt(s.label, 0)
	case "select":
		return interactive.AskSelect(s.label, s.options)
	}
	return nil, fmt.Errorf("unknown setting kind: %s", s.kind)
}

func runCreate(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	// Step 1: Ask for meeting name
	name, err := requireString(cmd, "name", "Meeting Name")
	if err != nil {
		return err
	}

	// Step 2: Ask for meeting ID (leave empty for random)
	meetingID, _ := cmd.Flags().GetString("meeting-id")
	if meetingID == "" && interactive.IsInteractive() {
		meetingID, err = interactive.AskString("Meeting ID (leave empty to auto-generate)", "")
		if err != nil {
			return err
		}
	}
	if meetingID == "" {
		b := make([]byte, 8)
		if _, err := rand.Read(b); err != nil {
			return fmt.Errorf("failed to generate meeting ID: %w", err)
		}
		meetingID = hex.EncodeToString(b)
	}
	cmd.Flags().Set("meeting-id", meetingID)

	req := &api.CreateRequest{
		Name:      name,
		MeetingID: meetingID,
	}

	// Apply any flags that were explicitly set on the command line
	req.Welcome = optionalString(cmd, "welcome")
	req.DialNumber = optionalString(cmd, "dial-number")
	req.VoiceBridge = optionalString(cmd, "voice-bridge")
	req.MaxParticipants = optionalInt(cmd, "max-participants")
	req.LogoutURL = optionalString(cmd, "logout-url")
	req.Record = optionalBool(cmd, "record")
	req.Duration = optionalInt(cmd, "duration")
	req.ModeratorOnlyMessage = optionalString(cmd, "moderator-only-message")
	req.AutoStartRecording = optionalBool(cmd, "auto-start-recording")
	req.AllowStartStopRecording = optionalBool(cmd, "allow-start-stop-recording")
	req.WebcamsOnlyForModerator = optionalBool(cmd, "webcams-only-for-moderator")
	req.BannerText = optionalString(cmd, "banner-text")
	req.BannerColor = optionalString(cmd, "banner-color")
	req.MuteOnStart = optionalBool(cmd, "mute-on-start")
	req.AllowModsToUnmuteUsers = optionalBool(cmd, "allow-mods-to-unmute-users")
	req.GuestPolicy = optionalString(cmd, "guest-policy")
	req.MeetingLayout = optionalString(cmd, "meeting-layout")
	req.EndWhenNoModerator = optionalBool(cmd, "end-when-no-moderator")
	req.EndWhenNoModeratorDelayInMinutes = optionalInt(cmd, "end-when-no-moderator-delay")
	req.Logo = optionalString(cmd, "logo")
	req.DisabledFeatures = optionalString(cmd, "disabled-features")
	req.NotifyRecordingIsOn = optionalBool(cmd, "notify-recording-is-on")
	req.LockSettingsDisableCam = optionalBool(cmd, "lock-disable-cam")
	req.LockSettingsDisableMic = optionalBool(cmd, "lock-disable-mic")
	req.LockSettingsDisablePrivateChat = optionalBool(cmd, "lock-disable-private-chat")
	req.LockSettingsDisablePublicChat = optionalBool(cmd, "lock-disable-public-chat")
	req.LockSettingsDisableNotes = optionalBool(cmd, "lock-disable-notes")
	req.LockSettingsHideUserList = optionalBool(cmd, "lock-hide-user-list")
	req.LockSettingsLockOnJoin = optionalBool(cmd, "lock-on-join")
	req.LockSettingsHideViewersCursor = optionalBool(cmd, "lock-hide-viewers-cursor")
	req.IsBreakout = optionalBool(cmd, "is-breakout")
	req.ParentMeetingID = optionalString(cmd, "parent-meeting-id")
	req.Sequence = optionalInt(cmd, "sequence")
	req.FreeJoin = optionalBool(cmd, "free-join")

	metaPairs, _ := cmd.Flags().GetStringSlice("meta")
	req.Meta = parseMeta(metaPairs)

	// Step 3: Interactive wizard - ask whether to configure additional settings
	if interactive.IsInteractive() && !hasAnyExplicitCreateFlags(cmd) {
		configure, err := interactive.AskConfirm("Configure additional settings?", false)
		if err != nil {
			return err
		}

		if configure {
			settings := additionalSettings()
			labels := make([]string, len(settings))
			for i, s := range settings {
				labels[i] = s.label
			}

			selected, err := interactive.AskMultiSelect("Select settings to configure", labels)
			if err != nil {
				return err
			}

			// Build lookup of selected labels
			selectedSet := make(map[string]bool, len(selected))
			for _, l := range selected {
				selectedSet[l] = true
			}

			// Prompt for each selected setting
			for _, s := range settings {
				if !selectedSet[s.label] {
					continue
				}
				val, err := promptForSetting(s)
				if err != nil {
					return err
				}
				s.apply(req, val)
				if s.flagName != "" {
					setFlagFromWizard(cmd, s.flagName, val)
				}
			}
		}
	}

	client := newClient(cfg)
	resp, err := client.Create(req)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Meeting created successfully!\n\n")
	printField("Meeting ID", resp.MeetingID)
	printField("Internal Meeting ID", resp.InternalMeetingID)
	printField("Voice Bridge", resp.VoiceBridge)
	printField("Dial Number", resp.DialNumber)
	printField("Create Date", resp.CreateDate)
	printFieldInt("Duration", resp.Duration)
	printFieldBool("Has User Joined", resp.HasUserJoined)

	return nil
}

// hasAnyExplicitCreateFlags returns true if any optional create flags were explicitly set.
func hasAnyExplicitCreateFlags(cmd *cobra.Command) bool {
	optionalFlags := []string{
		"welcome", "dial-number", "voice-bridge", "max-participants", "logout-url",
		"record", "duration", "moderator-only-message", "auto-start-recording",
		"allow-start-stop-recording", "webcams-only-for-moderator", "banner-text",
		"banner-color", "mute-on-start", "allow-mods-to-unmute-users", "guest-policy",
		"meeting-layout", "end-when-no-moderator", "end-when-no-moderator-delay",
		"logo", "disabled-features", "notify-recording-is-on",
		"lock-disable-cam", "lock-disable-mic", "lock-disable-private-chat",
		"lock-disable-public-chat", "lock-disable-notes", "lock-hide-user-list",
		"lock-on-join", "lock-hide-viewers-cursor",
		"is-breakout", "parent-meeting-id", "sequence", "free-join", "meta",
	}
	for _, f := range optionalFlags {
		if cmd.Flags().Changed(f) {
			return true
		}
	}
	return false
}
