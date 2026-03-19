package cmd

import (
	"fmt"
	"os"

	"github.com/invokablegmbh/bigbluebutton-cli/pkg/api"
	"github.com/spf13/cobra"
)

func newIsMeetingRunningCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "is-meeting-running",
		Short: "Check if a meeting is currently running",
		Long: `Check if a BigBlueButton meeting is currently running.

Returns "true" or "false".

Examples:
  bigbluebutton-cli is-meeting-running --meeting-id team-2024`,
		RunE: withWizardHint(runIsMeetingRunning),
	}

	cmd.Flags().String("meeting-id", "", "Meeting ID to check (required)")

	return cmd
}

func runIsMeetingRunning(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	meetingID, err := requireString(cmd, "meeting-id", "Meeting ID")
	if err != nil {
		return err
	}

	client := newClient(cfg)
	resp, err := client.IsMeetingRunning(&api.IsMeetingRunningRequest{MeetingID: meetingID})
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%v\n", resp.Running)
	return nil
}
