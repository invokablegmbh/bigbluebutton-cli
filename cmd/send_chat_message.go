package cmd

import (
	"fmt"
	"os"

	"github.com/invokablegmbh/bigbluebutton-cli/pkg/api"
	"github.com/spf13/cobra"
)

func newSendChatMessageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-chat-message",
		Short: "Send a chat message to a running meeting",
		Long: `Send a chat message to a currently running BigBlueButton meeting.

The message will appear in the public chat. The sender name defaults to "System".

Examples:
  bigbluebutton-cli send-chat-message --meeting-id team-2024 --message "Meeting will end in 5 minutes"
  bigbluebutton-cli send-chat-message --meeting-id team-2024 --message "Welcome!" --user-name "Admin"`,
		RunE: withWizardHint(runSendChatMessage),
	}

	cmd.Flags().String("meeting-id", "", "Meeting ID (required)")
	cmd.Flags().String("message", "", "Chat message to send (required, max 500 chars)")
	cmd.Flags().String("user-name", "", "Sender name (default: System)")

	return cmd
}

func runSendChatMessage(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	meetingID, err := requireString(cmd, "meeting-id", "Meeting ID")
	if err != nil {
		return err
	}
	message, err := requireString(cmd, "message", "Chat message")
	if err != nil {
		return err
	}

	req := &api.SendChatMessageRequest{
		MeetingID: meetingID,
		Message:   message,
	}
	req.UserName = optionalString(cmd, "user-name")

	client := newClient(cfg)
	_, err = client.SendChatMessage(req)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Chat message sent to meeting '%s'.\n", meetingID)
	return nil
}
