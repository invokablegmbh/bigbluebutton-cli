package cmd

import (
	"fmt"
	"os"

	"github.com/invokablegmbh/bigbluebutton-cli/pkg/api"
	"github.com/spf13/cobra"
)

func newUpdateRecordingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-recordings",
		Short: "Update recording metadata",
		Long: `Update metadata for one or more BigBlueButton recordings.

Examples:
  bigbluebutton-cli update-recordings --record-id rec-123 --meta name="New Title"
  bigbluebutton-cli update-recordings --record-id rec-123 --meta name="Title",description="A description"`,
		RunE: withWizardHint(runUpdateRecordings),
	}

	cmd.Flags().String("record-id", "", "Recording ID(s) to update (required, comma-separated)")
	cmd.Flags().StringSlice("meta", nil, "Metadata key=value pairs (required, can specify multiple)")

	return cmd
}

func runUpdateRecordings(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	recordID, err := requireString(cmd, "record-id", "Recording ID")
	if err != nil {
		return err
	}

	metaPairs, _ := cmd.Flags().GetStringSlice("meta")
	meta := parseMeta(metaPairs)

	client := newClient(cfg)
	resp, err := client.UpdateRecordings(&api.UpdateRecordingsRequest{
		RecordID: recordID,
		Meta:     meta,
	})
	if err != nil {
		return err
	}

	if resp.Updated {
		fmt.Fprintf(os.Stdout, "Recording(s) '%s' updated successfully.\n", recordID)
	} else {
		fmt.Fprintf(os.Stdout, "Recording(s) '%s' could not be updated.\n", recordID)
	}
	return nil
}
