package cmd

import (
	"fmt"
	"os"

	"github.com/invokablegmbh/bigbluebutton-cli/pkg/api"
	"github.com/spf13/cobra"
)

func newEndCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "end",
		Short: "End a running BigBlueButton meeting",
		Long: `End a currently running BigBlueButton meeting.

This will immediately terminate the meeting and disconnect all participants.

Examples:
  bigbluebutton-cli end --meeting-id team-2024`,
		RunE: withWizardHint(runEnd),
	}

	cmd.Flags().String("meeting-id", "", "Meeting ID to end (required)")

	return cmd
}

func runEnd(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	meetingID, err := requireString(cmd, "meeting-id", "Meeting ID")
	if err != nil {
		return err
	}

	client := newClient(cfg)
	_, err = client.End(&api.EndRequest{MeetingID: meetingID})
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Meeting '%s' has been ended.\n", meetingID)
	return nil
}
