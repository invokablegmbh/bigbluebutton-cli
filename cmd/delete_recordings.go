package cmd

import (
	"fmt"
	"os"

	"github.com/invokablegmbh/bigbluebutton-cli/pkg/api"
	"github.com/spf13/cobra"
)

func newDeleteRecordingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-recordings",
		Short: "Delete recordings",
		Long: `Delete one or more BigBlueButton recordings.

This action is irreversible.

Examples:
  bigbluebutton-cli delete-recordings --record-id rec-123
  bigbluebutton-cli delete-recordings --record-id rec-123,rec-456`,
		RunE: withWizardHint(runDeleteRecordings),
	}

	cmd.Flags().String("record-id", "", "Recording ID(s) to delete (required, comma-separated)")

	return cmd
}

func runDeleteRecordings(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	recordID, err := requireString(cmd, "record-id", "Recording ID")
	if err != nil {
		return err
	}

	client := newClient(cfg)
	resp, err := client.DeleteRecordings(&api.DeleteRecordingsRequest{RecordID: recordID})
	if err != nil {
		return err
	}

	if resp.Deleted {
		fmt.Fprintf(os.Stdout, "Recording(s) '%s' deleted successfully.\n", recordID)
	} else {
		fmt.Fprintf(os.Stdout, "Recording(s) '%s' could not be deleted.\n", recordID)
	}
	return nil
}
