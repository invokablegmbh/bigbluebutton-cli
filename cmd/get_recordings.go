package cmd

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"github.com/invokablegmbh/bigbluebutton-cli/pkg/api"
	"github.com/spf13/cobra"
)

func newGetRecordingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-recordings",
		Short: "List recordings",
		Long: `List BigBlueButton recordings with optional filters.

You can filter by meeting ID, recording ID, or state.

Examples:
  bigbluebutton-cli get-recordings
  bigbluebutton-cli get-recordings --meeting-id team-2024
  bigbluebutton-cli get-recordings --state published
  bigbluebutton-cli get-recordings --limit 10 --offset 0
  bigbluebutton-cli get-recordings --format json`,
		RunE: runGetRecordings,
	}

	cmd.Flags().String("meeting-id", "", "Filter by meeting ID (comma-separated for multiple)")
	cmd.Flags().String("record-id", "", "Filter by recording ID (comma-separated for multiple)")
	cmd.Flags().String("state", "", "Filter by state: processing, processed, published, unpublished, deleted, any")
	cmd.Flags().Int("offset", 0, "Pagination offset")
	cmd.Flags().Int("limit", 0, "Results per page (1-100)")
	cmd.Flags().StringSlice("meta", nil, "Metadata filters (key=value)")
	addFormatFlag(cmd)

	return cmd
}

type recordingsOutput struct {
	XMLName    xml.Name          `json:"-" xml:"recordings"`
	Recordings []recordingOutput `json:"recordings" xml:"recording"`
}

type recordingOutput struct {
	Name         string `json:"name" xml:"name"`
	RecordID     string `json:"recordID" xml:"recordID"`
	MeetingID    string `json:"meetingID" xml:"meetingID"`
	State        string `json:"state" xml:"state"`
	Published    bool   `json:"published" xml:"published"`
	StartTime    string `json:"startTime" xml:"startTime"`
	Participants int    `json:"participants" xml:"participants"`
	URL          string `json:"url" xml:"url"`
}

func runGetRecordings(cmd *cobra.Command, args []string) error {
	format, err := getFormat(cmd)
	if err != nil {
		return err
	}

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	req := &api.GetRecordingsRequest{}
	req.MeetingID = optionalString(cmd, "meeting-id")
	req.RecordID = optionalString(cmd, "record-id")
	req.State = optionalString(cmd, "state")
	req.Offset = optionalInt(cmd, "offset")
	req.Limit = optionalInt(cmd, "limit")

	metaPairs, _ := cmd.Flags().GetStringSlice("meta")
	req.Meta = parseMeta(metaPairs)

	client := newClient(cfg)
	resp, err := client.GetRecordings(req)
	if err != nil {
		return err
	}

	if len(resp.Recordings) == 0 && format == formatHuman {
		fmt.Fprintln(os.Stdout, "No recordings found.")
		return nil
	}

	output := recordingsOutput{
		Recordings: make([]recordingOutput, len(resp.Recordings)),
	}
	for i, r := range resp.Recordings {
		startTime := ""
		if r.StartTime > 0 {
			startTime = time.UnixMilli(r.StartTime).Format("2006-01-02")
		}
		url := ""
		if len(r.Playback) > 0 {
			url = r.Playback[0].URL
		}
		output.Recordings[i] = recordingOutput{
			Name:         r.Name,
			RecordID:     r.RecordID,
			MeetingID:    r.MeetingID,
			State:        r.State,
			Published:    r.Published,
			StartTime:    startTime,
			Participants: r.Participants,
			URL:          url,
		}
	}

	return printFormatted(format, output, func() {
		fmt.Fprintf(os.Stdout, "%-20s %-25s %-12s %-10s %-12s %-8s  %s\n",
			"Name", "Record ID", "State", "Published", "Start Time", "Participants", "URL")
		fmt.Fprintf(os.Stdout, "%-20s %-25s %-12s %-10s %-12s %-8s  %s\n",
			"----", "---------", "-----", "---------", "----------", "------------", "---")

		for _, r := range resp.Recordings {
			published := "no"
			if r.Published {
				published = "yes"
			}
			startTime := ""
			if r.StartTime > 0 {
				startTime = time.UnixMilli(r.StartTime).Format("2006-01-02")
			}
			url := ""
			if len(r.Playback) > 0 {
				url = r.Playback[0].URL
			}
			fmt.Fprintf(os.Stdout, "%-20s %-25s %-12s %-10s %-12s %-8d  %s\n",
				truncate(r.Name, 18),
				truncate(r.RecordID, 23),
				r.State,
				published,
				startTime,
				r.Participants,
				url,
			)
		}

		fmt.Fprintf(os.Stdout, "\nTotal: %d recording(s)\n", len(resp.Recordings))
	})
}
