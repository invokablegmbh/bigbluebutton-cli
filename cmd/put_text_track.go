package cmd

import (
	"fmt"
	"os"

	"github.com/invokablegmbh/bigbluebutton-cli/pkg/api"
	"github.com/spf13/cobra"
)

func newPutTextTrackCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "put-text-track",
		Short: "Upload a subtitle/caption file for a recording",
		Long: `Upload a subtitle or caption file (SRT, SSA/ASS, or WebVTT format)
for a BigBlueButton recording.

Examples:
  bigbluebutton-cli put-text-track --record-id rec-123 --kind subtitles --lang en --label "English" --file subtitles.vtt
  bigbluebutton-cli put-text-track --record-id rec-123 --kind captions --lang es --label "Spanish" --file captions.srt`,
		RunE: withWizardHint(runPutTextTrack),
	}

	cmd.Flags().String("record-id", "", "Recording ID (required)")
	cmd.Flags().String("kind", "", "Track type: subtitles or captions (required)")
	cmd.Flags().String("lang", "", "Language tag, e.g. en, es, fr (required)")
	cmd.Flags().String("label", "", "Human-readable label, e.g. 'English' (required)")
	cmd.Flags().String("file", "", "Path to subtitle file (required)")

	return cmd
}

func runPutTextTrack(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	recordID, err := requireString(cmd, "record-id", "Recording ID")
	if err != nil {
		return err
	}
	kind, err := requireString(cmd, "kind", "Track type (subtitles/captions)")
	if err != nil {
		return err
	}
	lang, err := requireString(cmd, "lang", "Language tag (e.g. en)")
	if err != nil {
		return err
	}
	label, err := requireString(cmd, "label", "Track label (e.g. English)")
	if err != nil {
		return err
	}
	filePath, err := requireString(cmd, "file", "Subtitle file path")
	if err != nil {
		return err
	}

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filePath)
	}

	client := newClient(cfg)
	_, err = client.PutRecordingTextTrack(&api.PutRecordingTextTrackRequest{
		RecordID: recordID,
		Kind:     kind,
		Lang:     lang,
		Label:    label,
		FilePath: filePath,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Text track uploaded successfully for recording '%s'.\n", recordID)
	return nil
}
