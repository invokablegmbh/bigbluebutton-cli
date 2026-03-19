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

func newCreateJoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-join",
		Short: "Create a meeting and immediately generate a join URL",
		Long: `Create a new BigBlueButton meeting and immediately generate a join URL
in a single step. This prevents the meeting from expiring before anyone joins.

The join URL can be opened in a browser to enter the meeting directly.

Examples:
  bigbluebutton-cli create-join --name "Team Standup" --meeting-id standup-001 --user-name "Alice" --role MODERATOR
  bigbluebutton-cli create-join --name "Webinar" --meeting-id webinar-042 --user-name "Host" --role MODERATOR --record
  bigbluebutton-cli create-join --name "Quick Call" --meeting-id quick-001 --user-name "Bob" --role MODERATOR --url-only`,
		RunE: withWizardHint(runCreateJoin),
	}

	// Create flags
	cmd.Flags().String("name", "", "Meeting name (required)")
	cmd.Flags().String("meeting-id", "", "Unique meeting identifier (required)")
	cmd.Flags().String("welcome", "", "Welcome message")
	cmd.Flags().Int("max-participants", 0, "Maximum number of participants")
	cmd.Flags().String("logout-url", "", "URL to redirect after logout")
	cmd.Flags().Bool("record", false, "Enable recording")
	cmd.Flags().Int("duration", 0, "Meeting duration in minutes (0 = unlimited)")
	cmd.Flags().String("moderator-only-message", "", "Message visible only to moderators")
	cmd.Flags().Bool("auto-start-recording", false, "Automatically start recording")
	cmd.Flags().Bool("allow-start-stop-recording", true, "Allow moderators to start/stop recording")
	cmd.Flags().Bool("webcams-only-for-moderator", false, "Only show webcams to moderators")
	cmd.Flags().Bool("mute-on-start", false, "Mute all users on meeting start")
	cmd.Flags().Bool("allow-mods-to-unmute-users", false, "Allow moderators to unmute users")
	cmd.Flags().String("guest-policy", "", "Guest policy: ALWAYS_ACCEPT, ALWAYS_DENY, ASK_MODERATOR")
	cmd.Flags().String("meeting-layout", "", "Default layout: CUSTOM_LAYOUT, SMART_LAYOUT, PRESENTATION_FOCUS, VIDEO_FOCUS")
	cmd.Flags().Bool("end-when-no-moderator", false, "End meeting when no moderator is present")
	cmd.Flags().Int("end-when-no-moderator-delay", 0, "Minutes before ending when no moderator")
	cmd.Flags().String("banner-text", "", "Banner text displayed in the client")
	cmd.Flags().String("banner-color", "", "Banner background color (hex, e.g. #FF0000)")
	cmd.Flags().String("logo", "", "Branding logo URL")
	cmd.Flags().String("disabled-features", "", "Comma-separated list of disabled features")
	cmd.Flags().Bool("lock-disable-cam", false, "Prevent camera sharing")
	cmd.Flags().Bool("lock-disable-mic", false, "Restrict to listen-only mode")
	cmd.Flags().Bool("lock-disable-private-chat", false, "Disable private messaging")
	cmd.Flags().Bool("lock-disable-public-chat", false, "Disable public chat")
	cmd.Flags().Bool("lock-disable-notes", false, "Disable shared notes")
	cmd.Flags().Bool("lock-hide-user-list", false, "Hide viewer list from attendees")
	cmd.Flags().Bool("lock-on-join", false, "Apply lock settings to new joiners")
	cmd.Flags().Bool("lock-hide-viewers-cursor", false, "Hide whiteboard cursors from viewers")
	cmd.Flags().StringSlice("meta", nil, "Metadata key=value pairs (can specify multiple)")

	// Join flags
	cmd.Flags().String("user-name", "", "Full name of the joining user (required)")
	cmd.Flags().String("role", "", "User role: MODERATOR or VIEWER (required)")
	cmd.Flags().String("user-id", "", "Application-specific user ID")
	cmd.Flags().String("avatar-url", "", "User avatar image URL")
	cmd.Flags().Bool("url-only", false, "Output only the join URL (no extra text)")

	return cmd
}

func runCreateJoin(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	// Resolve required create params
	name, err := requireString(cmd, "name", "Meeting Name")
	if err != nil {
		return err
	}
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

	// Resolve required join params
	userName, err := requireString(cmd, "user-name", "Your Name")
	if err != nil {
		return err
	}

	role, _ := cmd.Flags().GetString("role")
	if role == "" && interactive.IsInteractive() {
		role, err = interactive.AskSelect("Role", []string{"MODERATOR", "VIEWER"})
		if err != nil {
			return err
		}
		cmd.Flags().Set("role", role)
	}
	if role == "" {
		return fmt.Errorf("--role is required (MODERATOR or VIEWER)")
	}

	// Build create request
	createReq := &api.CreateRequest{
		Name:      name,
		MeetingID: meetingID,
	}
	createReq.Welcome = optionalString(cmd, "welcome")
	createReq.MaxParticipants = optionalInt(cmd, "max-participants")
	createReq.LogoutURL = optionalString(cmd, "logout-url")
	createReq.Record = optionalBool(cmd, "record")
	createReq.Duration = optionalInt(cmd, "duration")
	createReq.ModeratorOnlyMessage = optionalString(cmd, "moderator-only-message")
	createReq.AutoStartRecording = optionalBool(cmd, "auto-start-recording")
	createReq.AllowStartStopRecording = optionalBool(cmd, "allow-start-stop-recording")
	createReq.WebcamsOnlyForModerator = optionalBool(cmd, "webcams-only-for-moderator")
	createReq.MuteOnStart = optionalBool(cmd, "mute-on-start")
	createReq.AllowModsToUnmuteUsers = optionalBool(cmd, "allow-mods-to-unmute-users")
	createReq.GuestPolicy = optionalString(cmd, "guest-policy")
	createReq.MeetingLayout = optionalString(cmd, "meeting-layout")
	createReq.EndWhenNoModerator = optionalBool(cmd, "end-when-no-moderator")
	createReq.EndWhenNoModeratorDelayInMinutes = optionalInt(cmd, "end-when-no-moderator-delay")
	createReq.BannerText = optionalString(cmd, "banner-text")
	createReq.BannerColor = optionalString(cmd, "banner-color")
	createReq.Logo = optionalString(cmd, "logo")
	createReq.DisabledFeatures = optionalString(cmd, "disabled-features")
	createReq.LockSettingsDisableCam = optionalBool(cmd, "lock-disable-cam")
	createReq.LockSettingsDisableMic = optionalBool(cmd, "lock-disable-mic")
	createReq.LockSettingsDisablePrivateChat = optionalBool(cmd, "lock-disable-private-chat")
	createReq.LockSettingsDisablePublicChat = optionalBool(cmd, "lock-disable-public-chat")
	createReq.LockSettingsDisableNotes = optionalBool(cmd, "lock-disable-notes")
	createReq.LockSettingsHideUserList = optionalBool(cmd, "lock-hide-user-list")
	createReq.LockSettingsLockOnJoin = optionalBool(cmd, "lock-on-join")
	createReq.LockSettingsHideViewersCursor = optionalBool(cmd, "lock-hide-viewers-cursor")
	metaPairs, _ := cmd.Flags().GetStringSlice("meta")
	createReq.Meta = parseMeta(metaPairs)

	// Interactive wizard: ask whether to configure additional settings
	if interactive.IsInteractive() && !hasAnyExplicitCreateFlags(cmd) {
		configure, err := interactive.AskConfirm("Configure additional meeting settings?", false)
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

			selectedSet := make(map[string]bool, len(selected))
			for _, l := range selected {
				selectedSet[l] = true
			}

			for _, s := range settings {
				if !selectedSet[s.label] {
					continue
				}
				val, err := promptForSetting(s)
				if err != nil {
					return err
				}
				s.apply(createReq, val)
				if s.flagName != "" {
					setFlagFromWizard(cmd, s.flagName, val)
				}
			}
		}
	}

	client := newClient(cfg)

	// Step 1: Create the meeting
	createResp, err := client.Create(createReq)
	if err != nil {
		return fmt.Errorf("creating meeting: %w", err)
	}

	// Step 2: Generate join URL immediately
	joinReq := &api.JoinRequest{
		FullName:  userName,
		MeetingID: meetingID,
		Role:      role,
	}
	joinReq.UserID = optionalString(cmd, "user-id")
	joinReq.AvatarURL = optionalString(cmd, "avatar-url")

	joinURL := client.JoinURL(joinReq)

	urlOnly, _ := cmd.Flags().GetBool("url-only")
	if urlOnly {
		fmt.Fprintln(os.Stdout, joinURL)
		return nil
	}

	fmt.Fprintf(os.Stdout, "Meeting created and join URL generated!\n\n")
	printField("Meeting ID", createResp.MeetingID)
	printField("Internal Meeting ID", createResp.InternalMeetingID)
	printField("Voice Bridge", createResp.VoiceBridge)
	printField("Dial Number", createResp.DialNumber)
	printField("Create Date", createResp.CreateDate)
	printFieldInt("Duration", createResp.Duration)
	fmt.Fprintf(os.Stdout, "\n")
	printField("Join as", userName)
	printField("Role", role)
	printField("Join URL", joinURL)

	return nil
}
