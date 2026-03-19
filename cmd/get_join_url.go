package cmd

import (
	"fmt"
	"os"

	"github.com/invokablegmbh/bigbluebutton-cli/pkg/api"
	"github.com/spf13/cobra"
)

func newGetJoinUrlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-join-url",
		Short: "Get a join URL from a session token",
		Long: `Generate a join URL using an existing session token.

This is useful for creating additional sessions or replacing existing ones.

Examples:
  bigbluebutton-cli get-join-url --session-token sess-xyz
  bigbluebutton-cli get-join-url --session-token sess-xyz --replace-session --session-name "New Session"`,
		RunE: withWizardHint(runGetJoinUrl),
	}

	cmd.Flags().String("session-token", "", "Session token (required)")
	cmd.Flags().Bool("replace-session", false, "Invalidate the original session")
	cmd.Flags().String("session-name", "", "Name for the new session")

	return cmd
}

func runGetJoinUrl(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	sessionToken, err := requireString(cmd, "session-token", "Session Token")
	if err != nil {
		return err
	}

	req := &api.GetJoinUrlRequest{
		SessionToken: sessionToken,
	}
	req.ReplaceSession = optionalBool(cmd, "replace-session")
	req.SessionName = optionalString(cmd, "session-name")

	client := newClient(cfg)
	resp, err := client.GetJoinUrl(req)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, resp.URL)
	return nil
}
