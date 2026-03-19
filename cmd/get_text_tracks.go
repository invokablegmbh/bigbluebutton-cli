package cmd

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/invokablegmbh/bigbluebutton-cli/pkg/api"
	"github.com/spf13/cobra"
)

func newGetTextTracksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-text-tracks",
		Short: "Get text tracks (subtitles/captions) for a recording",
		Long: `Get the list of text tracks (subtitles and captions) associated with
a BigBlueButton recording.

Examples:
  bigbluebutton-cli get-text-tracks --record-id rec-123
  bigbluebutton-cli get-text-tracks --record-id rec-123 --format json`,
		RunE: withWizardHint(runGetTextTracks),
	}

	cmd.Flags().String("record-id", "", "Recording ID (required)")
	addFormatFlag(cmd)

	return cmd
}

type textTracksOutput struct {
	XMLName xml.Name          `json:"-" xml:"textTracks"`
	Tracks  []textTrackOutput `json:"tracks" xml:"track"`
}

type textTrackOutput struct {
	Label  string `json:"label" xml:"label"`
	Kind   string `json:"kind" xml:"kind"`
	Lang   string `json:"lang" xml:"lang"`
	Source string `json:"source" xml:"source"`
	Href   string `json:"href" xml:"href"`
}

func runGetTextTracks(cmd *cobra.Command, args []string) error {
	format, err := getFormat(cmd)
	if err != nil {
		return err
	}

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	recordID, err := requireString(cmd, "record-id", "Recording ID")
	if err != nil {
		return err
	}

	client := newClient(cfg)
	resp, err := client.GetRecordingTextTracks(&api.GetRecordingTextTracksRequest{RecordID: recordID})
	if err != nil {
		return err
	}

	tracks := resp.Response.Tracks
	if len(tracks) == 0 && format == formatHuman {
		fmt.Fprintln(os.Stdout, "No text tracks found for this recording.")
		return nil
	}

	output := textTracksOutput{
		Tracks: make([]textTrackOutput, len(tracks)),
	}
	for i, t := range tracks {
		output.Tracks[i] = textTrackOutput{
			Label:  t.Label,
			Kind:   t.Kind,
			Lang:   t.Lang,
			Source: t.Source,
			Href:   t.Href,
		}
	}

	return printFormatted(format, output, func() {
		fmt.Fprintf(os.Stdout, "%-15s %-12s %-8s %-10s %s\n",
			"Label", "Kind", "Lang", "Source", "URL")
		fmt.Fprintf(os.Stdout, "%-15s %-12s %-8s %-10s %s\n",
			"-----", "----", "----", "------", "---")

		for _, t := range tracks {
			fmt.Fprintf(os.Stdout, "%-15s %-12s %-8s %-10s %s\n",
				t.Label, t.Kind, t.Lang, t.Source, t.Href)
		}
	})
}
