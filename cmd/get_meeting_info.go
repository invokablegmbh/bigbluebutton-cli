package cmd

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"

	"github.com/invokablegmbh/bigbluebutton-cli/pkg/api"
	"github.com/spf13/cobra"
)

func newGetMeetingInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-meeting-info",
		Short: "Get detailed information about a meeting",
		Long: `Get detailed information about a BigBlueButton meeting including
participant list, recording status, and meeting configuration.

Examples:
  bigbluebutton-cli get-meeting-info --meeting-id team-2024
  bigbluebutton-cli get-meeting-info --meeting-id team-2024 --format json`,
		RunE: withWizardHint(runGetMeetingInfo),
	}

	cmd.Flags().String("meeting-id", "", "Meeting ID (required)")
	addFormatFlag(cmd)

	return cmd
}

type meetingInfoOutput struct {
	XMLName              xml.Name              `json:"-" xml:"meeting"`
	MeetingName          string                `json:"meetingName" xml:"meetingName"`
	MeetingID            string                `json:"meetingID" xml:"meetingID"`
	InternalMeetingID    string                `json:"internalMeetingID" xml:"internalMeetingID"`
	CreateDate           string                `json:"createDate" xml:"createDate"`
	VoiceBridge          string                `json:"voiceBridge" xml:"voiceBridge"`
	DialNumber           string                `json:"dialNumber" xml:"dialNumber"`
	Running              bool                  `json:"running" xml:"running"`
	Recording            bool                  `json:"recording" xml:"recording"`
	Duration             int                   `json:"duration" xml:"duration"`
	HasUserJoined        bool                  `json:"hasUserJoined" xml:"hasUserJoined"`
	HasBeenForciblyEnded bool                  `json:"hasBeenForciblyEnded" xml:"hasBeenForciblyEnded"`
	ParticipantCount     int                   `json:"participantCount" xml:"participantCount"`
	ListenerCount        int                   `json:"listenerCount" xml:"listenerCount"`
	VoiceParticipantCount int                  `json:"voiceParticipantCount" xml:"voiceParticipantCount"`
	VideoCount           int                   `json:"videoCount" xml:"videoCount"`
	ModeratorCount       int                   `json:"moderatorCount" xml:"moderatorCount"`
	IsBreakout           bool                  `json:"isBreakout" xml:"isBreakout"`
	Attendees            []meetingInfoAttendee `json:"attendees" xml:"attendee"`
}

type meetingInfoAttendee struct {
	FullName        string `json:"fullName" xml:"fullName"`
	Role            string `json:"role" xml:"role"`
	IsPresenter     bool   `json:"isPresenter" xml:"isPresenter"`
	HasJoinedVoice  bool   `json:"hasJoinedVoice" xml:"hasJoinedVoice"`
	HasVideo        bool   `json:"hasVideo" xml:"hasVideo"`
	IsListeningOnly bool   `json:"isListeningOnly" xml:"isListeningOnly"`
}

func runGetMeetingInfo(cmd *cobra.Command, args []string) error {
	format, err := getFormat(cmd)
	if err != nil {
		return err
	}

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	meetingID, err := requireString(cmd, "meeting-id", "Meeting ID")
	if err != nil {
		return err
	}

	client := newClient(cfg)
	resp, err := client.GetMeetingInfo(&api.GetMeetingInfoRequest{MeetingID: meetingID})
	if err != nil {
		return err
	}

	m := &resp.Meeting

	output := meetingInfoOutput{
		MeetingName:           m.MeetingName,
		MeetingID:             m.MeetingID,
		InternalMeetingID:     m.InternalMeetingID,
		CreateDate:            m.CreateDate,
		VoiceBridge:           m.VoiceBridge,
		DialNumber:            m.DialNumber,
		Running:               m.Running,
		Recording:             m.Recording,
		Duration:              m.Duration,
		HasUserJoined:         m.HasUserJoined,
		HasBeenForciblyEnded:  m.HasBeenForciblyEnded,
		ParticipantCount:      m.ParticipantCount,
		ListenerCount:         m.ListenerCount,
		VoiceParticipantCount: m.VoiceParticipantCount,
		VideoCount:            m.VideoCount,
		ModeratorCount:        m.ModeratorCount,
		IsBreakout:            m.IsBreakout,
		Attendees:             make([]meetingInfoAttendee, len(m.Attendees)),
	}
	for i, a := range m.Attendees {
		output.Attendees[i] = meetingInfoAttendee{
			FullName:        a.FullName,
			Role:            a.Role,
			IsPresenter:     a.IsPresenter,
			HasJoinedVoice:  a.HasJoinedVoice,
			HasVideo:        a.HasVideo,
			IsListeningOnly: a.IsListeningOnly,
		}
	}

	return printFormatted(format, output, func() {
		printField("Meeting Name", m.MeetingName)
		printField("Meeting ID", m.MeetingID)
		printField("Internal Meeting ID", m.InternalMeetingID)
		printField("Create Date", m.CreateDate)
		printField("Voice Bridge", m.VoiceBridge)
		printField("Dial Number", m.DialNumber)
		printFieldBool("Running", m.Running)
		printFieldBool("Recording", m.Recording)
		printFieldInt("Duration", m.Duration)
		printFieldBool("Has User Joined", m.HasUserJoined)
		printFieldBool("Forcibly Ended", m.HasBeenForciblyEnded)
		printFieldInt("Participants", m.ParticipantCount)
		printFieldInt("Listeners", m.ListenerCount)
		printFieldInt("Voice Participants", m.VoiceParticipantCount)
		printFieldInt("Video Count", m.VideoCount)
		printFieldInt("Moderator Count", m.ModeratorCount)
		printFieldBool("Is Breakout", m.IsBreakout)

		if len(m.Attendees) > 0 {
			fmt.Fprintf(os.Stdout, "\nAttendees:\n")
			fmt.Fprintf(os.Stdout, "  %-20s %-12s %-10s %-8s %-8s %-8s\n",
				"Name", "Role", "Presenter", "Voice", "Video", "Listen")
			fmt.Fprintf(os.Stdout, "  %-20s %-12s %-10s %-8s %-8s %-8s\n",
				"----", "----", "---------", "-----", "-----", "------")
			for _, a := range m.Attendees {
				fmt.Fprintf(os.Stdout, "  %-20s %-12s %-10s %-8s %-8s %-8s\n",
					a.FullName,
					a.Role,
					strconv.FormatBool(a.IsPresenter),
					strconv.FormatBool(a.HasJoinedVoice),
					strconv.FormatBool(a.HasVideo),
					strconv.FormatBool(a.IsListeningOnly),
				)
			}
		}
	})
}
