package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "dev"

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of bigbluebutton-cli",
		Long:  "Display the current version of the bigbluebutton-cli tool.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "bigbluebutton-cli version %s\n", Version)
		},
	}
}
