package cmd

import (
	"fmt"
	"os"

	"github.com/invokablegmbh/bigbluebutton-cli/pkg/api"
	"github.com/invokablegmbh/bigbluebutton-cli/pkg/interactive"
	"github.com/spf13/cobra"
)

func newJoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "join",
		Short: "Generate a join URL for a BigBlueButton meeting",
		Long: `Generate a join URL that can be used to join a BigBlueButton meeting.

The URL can be opened in a browser to join the meeting directly.

Examples:
  bigbluebutton-cli join --meeting-id team-2024 --name "John Doe" --role MODERATOR
  bigbluebutton-cli join --meeting-id team-2024 --name "Jane Smith" --role VIEWER
  bigbluebutton-cli join --meeting-id team-2024 --name "Bot" --role VIEWER --user-id bot-1`,
		RunE: withWizardHint(runJoin),
	}

	cmd.Flags().String("meeting-id", "", "Meeting ID to join (required)")
	cmd.Flags().String("name", "", "Full name of the user (required)")
	cmd.Flags().String("role", "", "User role: MODERATOR or VIEWER (required)")
	cmd.Flags().String("user-id", "", "Application-specific user ID")
	cmd.Flags().String("avatar-url", "", "User avatar image URL")
	cmd.Flags().String("redirect", "", "Set to 'false' to get XML response instead of redirect")
	cmd.Flags().String("logout-url", "", "URL to redirect after meeting")
	cmd.Flags().String("guest", "", "Set to 'true' for guest access")
	cmd.Flags().String("enforce-layout", "", "Force a specific layout")
	cmd.Flags().Bool("url-only", false, "Output only the join URL (no extra text)")

	return cmd
}

func runJoin(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	meetingID, err := requireString(cmd, "meeting-id", "Meeting ID")
	if err != nil {
		return err
	}
	fullName, err := requireString(cmd, "name", "Full Name")
	if err != nil {
		return err
	}

	role, _ := cmd.Flags().GetString("role")
	if role == "" && interactive.IsInteractive() {
		role, err = interactive.AskSelect("Role", []string{"MODERATOR", "VIEWER"})
		if err != nil {
			return err
		}
		cmd.Flags().Set("role", role)
	}
	if role == "" {
		return fmt.Errorf("--role is required (MODERATOR or VIEWER)")
	}

	req := &api.JoinRequest{
		FullName:  fullName,
		MeetingID: meetingID,
		Role:      role,
	}
	req.UserID = optionalString(cmd, "user-id")
	req.AvatarURL = optionalString(cmd, "avatar-url")
	req.LogoutURL = optionalString(cmd, "logout-url")
	req.Guest = optionalString(cmd, "guest")
	req.EnforceLayout = optionalString(cmd, "enforce-layout")

	client := newClient(cfg)

	urlOnly, _ := cmd.Flags().GetBool("url-only")
	if urlOnly {
		// Just output the URL without making an API call
		joinURL := client.JoinURL(req)
		fmt.Fprintln(os.Stdout, joinURL)
		return nil
	}

	// Get join details via API
	resp, err := client.Join(req)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Join URL generated successfully!\n\n")
	printField("Meeting ID", resp.MeetingID)
	printField("User ID", resp.UserID)
	printField("Session Token", resp.SessionToken)
	printField("URL", resp.URL)

	return nil
}
