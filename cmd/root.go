package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	// Default to discarding log output; verbose mode enables it.
	log.SetOutput(io.Discard)
}

var (
	flagURL     string
	flagSecret  string
	flagVerbose bool
	flagConfig  string
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "bigbluebutton-cli",
		Short: "BigBlueButton CLI - command line interface for the BigBlueButton API",
		Long: `bigbluebutton-cli is a command line tool that provides convenient access
to all BigBlueButton API endpoints.

It supports interactive mode (wizard-style prompts for missing parameters)
and non-interactive mode (all parameters via command line flags).

Configuration is resolved in the following order:
  1. Command line flags (--url, --secret)
  2. Environment variables (BBB_URL, BBB_SECRET)
  3. Config file (/etc/bigbluebutton-cli/config.yaml or ~/.bigbluebutton-cli/config.yaml)
  4. BigBlueButton autodetection (if installed on the same host)
  5. Interactive wizard (prompts for URL and secret)`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if flagVerbose {
				log.SetOutput(os.Stderr)
				log.SetFlags(log.Ltime | log.Lmicroseconds)
			} else {
				log.SetOutput(io.Discard)
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&flagURL, "url", "u", "", "BigBlueButton server URL (env: BBB_URL)")
	rootCmd.PersistentFlags().StringVarP(&flagSecret, "secret", "s", "", "BigBlueButton shared secret (env: BBB_SECRET)")
	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "Enable verbose/debug output")
	rootCmd.PersistentFlags().StringVarP(&flagConfig, "config", "c", "", "Path to config file")

	registerCommands(rootCmd)

	return rootCmd
}

func registerCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newCreateCmd())
	rootCmd.AddCommand(newCreateJoinCmd())
	rootCmd.AddCommand(newJoinCmd())
	rootCmd.AddCommand(newEndCmd())
	rootCmd.AddCommand(newInsertDocumentCmd())
	rootCmd.AddCommand(newIsMeetingRunningCmd())
	rootCmd.AddCommand(newGetMeetingInfoCmd())
	rootCmd.AddCommand(newGetMeetingsCmd())
	rootCmd.AddCommand(newGetRecordingsCmd())
	rootCmd.AddCommand(newPublishRecordingsCmd())
	rootCmd.AddCommand(newDeleteRecordingsCmd())
	rootCmd.AddCommand(newUpdateRecordingsCmd())
	rootCmd.AddCommand(newGetTextTracksCmd())
	rootCmd.AddCommand(newPutTextTrackCmd())
	rootCmd.AddCommand(newSendChatMessageCmd())
	rootCmd.AddCommand(newGetJoinUrlCmd())
}

func Execute() error {
	rootCmd := NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}
	return nil
}
