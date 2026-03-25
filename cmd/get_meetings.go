package cmd

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newGetMeetingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-meetings",
		Short: "List all running meetings",
		Long: `List all currently running BigBlueButton meetings.

Displays meeting name, ID, participant count, and recording status.

Examples:
  bigbluebutton-cli get-meetings
  bigbluebutton-cli get-meetings --format json
  bigbluebutton-cli get-meetings --format xml`,
		RunE: runGetMeetings,
	}

	addFormatFlag(cmd)

	return cmd
}

type meetingsOutput struct {
	XMLName  xml.Name        `json:"-" xml:"meetings"`
	Meetings []meetingOutput `json:"meetings" xml:"meeting"`
}

type meetingOutput struct {
	Name             string `json:"name" xml:"name"`
	MeetingID        string `json:"meetingID" xml:"meetingID"`
	ParticipantCount int    `json:"participantCount" xml:"participantCount"`
	Recording        bool   `json:"recording" xml:"recording"`
	Running          bool   `json:"running" xml:"running"`
}

func runGetMeetings(cmd *cobra.Command, args []string) error {
	format, err := getFormat(cmd)
	if err != nil {
		return err
	}

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	client := newClient(cfg)
	resp, err := client.GetMeetings()
	if err != nil {
		return err
	}

	if len(resp.Meetings) == 0 && format == formatHuman {
		fmt.Fprintln(os.Stdout, "No meetings are currently running.")
		return nil
	}

	output := meetingsOutput{
		Meetings: make([]meetingOutput, len(resp.Meetings)),
	}
	for i, m := range resp.Meetings {
		output.Meetings[i] = meetingOutput{
			Name:             m.MeetingName,
			MeetingID:        m.MeetingID,
			ParticipantCount: m.ParticipantCount,
			Recording:        m.Recording,
			Running:          m.Running,
		}
	}

	return printFormatted(format, output, func() {
		fmt.Fprintf(os.Stdout, "%-30s %-36s %-12s %-10s %-10s\n",
			"Name", "Meeting ID", "Participants", "Recording", "Running")
		fmt.Fprintf(os.Stdout, "%-30s %-36s %-12s %-10s %-10s\n",
			"----", "------------------------------------", "------------", "---------", "-------")

		for _, m := range resp.Meetings {
			recording := "no"
			if m.Recording {
				recording = "yes"
			}
			running := "no"
			if m.Running {
				running = "yes"
			}
			fmt.Fprintf(os.Stdout, "%-30s %-36s %-12d %-10s %-10s\n",
				truncate(m.MeetingName, 28),
				m.MeetingID,
				m.ParticipantCount,
				recording,
				running,
			)
		}

		fmt.Fprintf(os.Stdout, "\nTotal: %d meeting(s)\n", len(resp.Meetings))
	})
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-2] + ".."
}
