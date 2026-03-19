package cmd

import (
	"fmt"
	"os"

	"github.com/invokablegmbh/bigbluebutton-cli/pkg/api"
	"github.com/spf13/cobra"
)

func newInsertDocumentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insert-document",
		Short: "Insert a presentation document into a running meeting",
		Long: `Insert a presentation document into a currently running BigBlueButton meeting.

The document URL should point to a supported presentation file (PDF, PPT, etc.).

Examples:
  bigbluebutton-cli insert-document --meeting-id team-2024 --document-url https://example.com/slides.pdf
  bigbluebutton-cli insert-document --meeting-id team-2024 --document-url https://example.com/deck.pptx --filename "Q4 Review"`,
		RunE: withWizardHint(runInsertDocument),
	}

	cmd.Flags().String("meeting-id", "", "Meeting ID (required)")
	cmd.Flags().String("document-url", "", "URL of the document to insert (required)")
	cmd.Flags().String("filename", "", "Display name for the document")
	cmd.Flags().Bool("downloadable", false, "Allow participants to download the document")
	cmd.Flags().Bool("removable", false, "Allow the document to be removed")

	return cmd
}

func runInsertDocument(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	meetingID, err := requireString(cmd, "meeting-id", "Meeting ID")
	if err != nil {
		return err
	}
	docURL, err := requireString(cmd, "document-url", "Document URL")
	if err != nil {
		return err
	}

	filename, _ := cmd.Flags().GetString("filename")
	downloadable, _ := cmd.Flags().GetBool("downloadable")
	removable, _ := cmd.Flags().GetBool("removable")

	// Build XML body for the document
	downloadAttr := "false"
	if downloadable {
		downloadAttr = "true"
	}
	removeAttr := "false"
	if removable {
		removeAttr = "true"
	}

	xmlBody := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<modules>
  <module name="presentation">
    <document url="%s" downloadable="%s" removable="%s"`, docURL, downloadAttr, removeAttr)
	if filename != "" {
		xmlBody += fmt.Sprintf(` filename="%s"`, filename)
	}
	xmlBody += ` />
  </module>
</modules>`

	client := newClient(cfg)
	_, err = client.InsertDocument(&api.InsertDocumentRequest{
		MeetingID: meetingID,
		Body:      xmlBody,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Document inserted into meeting '%s'.\n", meetingID)
	return nil
}
