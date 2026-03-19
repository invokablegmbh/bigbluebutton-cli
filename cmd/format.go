package cmd

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	formatHuman = "human"
	formatJSON  = "json"
	formatXML   = "xml"
)

// addFormatFlag adds the --format flag to a command.
func addFormatFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("format", "f", formatHuman, "Output format: human, json, or xml")
}

// getFormat returns the validated format flag value.
func getFormat(cmd *cobra.Command) (string, error) {
	format, _ := cmd.Flags().GetString("format")
	switch format {
	case formatHuman, formatJSON, formatXML:
		return format, nil
	default:
		return "", fmt.Errorf("invalid format %q: must be human, json, or xml", format)
	}
}

// xmlWrapper wraps a list of items for XML marshaling with a root element.
type xmlWrapper struct {
	XMLName xml.Name
	Items   any `xml:"item"`
}

// printFormatted outputs data in the requested format. For "human" format it
// calls the provided humanFn. For json/xml it marshals the data.
func printFormatted(format string, data any, humanFn func()) error {
	switch format {
	case formatJSON:
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(data)
	case formatXML:
		enc := xml.NewEncoder(os.Stdout)
		enc.Indent("", "  ")
		if err := enc.Encode(data); err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout)
		return nil
	default:
		humanFn()
		return nil
	}
}
