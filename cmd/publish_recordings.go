package cmd

import (
	"fmt"
	"os"

	"github.com/invokablegmbh/bigbluebutton-cli/pkg/api"
	"github.com/invokablegmbh/bigbluebutton-cli/pkg/interactive"
	"github.com/spf13/cobra"
)

func newPublishRecordingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish-recordings",
		Short: "Publish or unpublish recordings",
		Long: `Publish or unpublish one or more BigBlueButton recordings.

Examples:
  bigbluebutton-cli publish-recordings --record-id rec-123 --publish true
  bigbluebutton-cli publish-recordings --record-id rec-123,rec-456 --publish false`,
		RunE: withWizardHint(runPublishRecordings),
	}

	cmd.Flags().String("record-id", "", "Recording ID(s) to publish/unpublish (required, comma-separated)")
	cmd.Flags().String("publish", "", "Set to 'true' to publish, 'false' to unpublish (required)")

	return cmd
}

func runPublishRecordings(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	recordID, err := requireString(cmd, "record-id", "Recording ID")
	if err != nil {
		return err
	}

	publishStr, _ := cmd.Flags().GetString("publish")
	if publishStr == "" && interactive.IsInteractive() {
		publishStr, err = interactive.AskSelect("Action", []string{"true", "false"})
		if err != nil {
			return err
		}
		cmd.Flags().Set("publish", publishStr)
	}
	if publishStr == "" {
		return fmt.Errorf("--publish is required (true or false)")
	}

	publish := publishStr == "true"

	client := newClient(cfg)
	resp, err := client.PublishRecordings(&api.PublishRecordingsRequest{
		RecordID: recordID,
		Publish:  publish,
	})
	if err != nil {
		return err
	}

	if resp.Published {
		fmt.Fprintf(os.Stdout, "Recording(s) '%s' published successfully.\n", recordID)
	} else {
		fmt.Fprintf(os.Stdout, "Recording(s) '%s' unpublished successfully.\n", recordID)
	}
	return nil
}
