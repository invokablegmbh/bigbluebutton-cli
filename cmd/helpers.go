package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/invokablegmbh/bigbluebutton-cli/pkg/api"
	"github.com/invokablegmbh/bigbluebutton-cli/pkg/config"
	"github.com/invokablegmbh/bigbluebutton-cli/pkg/interactive"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	// Register the interactive config prompt function
	config.SetPromptFunc(func(cfg *config.Config) error {
		if cfg.URL == "" {
			url, err := interactive.AskString("BigBlueButton Server URL", "")
			if err != nil {
				return err
			}
			cfg.URL = url
		}
		if cfg.Secret == "" {
			secret, err := interactive.AskPassword("BigBlueButton Shared Secret")
			if err != nil {
				return err
			}
			cfg.Secret = secret
		}
		return nil
	})
}

// loadConfig loads the BigBlueButton configuration using the cascade.
func loadConfig() (*config.Config, error) {
	return config.Load(flagURL, flagSecret, flagConfig, interactive.IsInteractive())
}

// newClient creates a new API client from the loaded configuration.
func newClient(cfg *config.Config) *api.Client {
	return api.NewClient(cfg.URL, cfg.Secret, flagVerbose)
}

// requireString returns the flag value, or prompts interactively if empty.
// When a value is collected interactively, it is stored back on the cobra flag
// so that the wizard hint can reconstruct the full command.
func requireString(cmd *cobra.Command, flagName, prompt string) (string, error) {
	val, _ := cmd.Flags().GetString(flagName)
	if val != "" {
		return val, nil
	}
	if interactive.IsInteractive() {
		val, err := interactive.AskString(prompt, "")
		if err != nil {
			return "", err
		}
		cmd.Flags().Set(flagName, val)
		return val, nil
	}
	return "", fmt.Errorf("--%s is required", flagName)
}

// optionalString returns the flag value or empty string.
func optionalString(cmd *cobra.Command, flagName string) *string {
	val, _ := cmd.Flags().GetString(flagName)
	if val == "" {
		return nil
	}
	return &val
}

// optionalBool returns the flag value or nil if not explicitly set.
func optionalBool(cmd *cobra.Command, flagName string) *bool {
	if !cmd.Flags().Changed(flagName) {
		return nil
	}
	val, _ := cmd.Flags().GetBool(flagName)
	return &val
}

// optionalInt returns the flag value or nil if not explicitly set.
func optionalInt(cmd *cobra.Command, flagName string) *int {
	if !cmd.Flags().Changed(flagName) {
		return nil
	}
	val, _ := cmd.Flags().GetInt(flagName)
	return &val
}

// parseMeta parses meta key=value pairs from a string slice.
func parseMeta(pairs []string) map[string]string {
	if len(pairs) == 0 {
		return nil
	}
	meta := make(map[string]string)
	for _, pair := range pairs {
		key, value, found := strings.Cut(pair, "=")
		if found {
			meta[key] = value
		}
	}
	return meta
}

// printField prints a labeled field to stdout.
func printField(label, value string) {
	if value != "" {
		fmt.Fprintf(os.Stdout, "%-25s %s\n", label+":", value)
	}
}

// printFieldBool prints a labeled boolean field.
func printFieldBool(label string, value bool) {
	fmt.Fprintf(os.Stdout, "%-25s %v\n", label+":", value)
}

// printFieldInt prints a labeled integer field.
func printFieldInt(label string, value int) {
	fmt.Fprintf(os.Stdout, "%-25s %d\n", label+":", value)
}

// printFieldInt64 prints a labeled int64 field.
func printFieldInt64(label string, value int64) {
	fmt.Fprintf(os.Stdout, "%-25s %d\n", label+":", value)
}

// collectChangedFlags snapshots which flags are currently marked as changed (i.e., set via CLI).
func collectChangedFlags(cmd *cobra.Command) map[string]bool {
	changed := make(map[string]bool)
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			changed[f.Name] = true
		}
	})
	return changed
}

// globalFlags lists persistent flags from the root command that should be excluded from hints.
var globalFlags = map[string]bool{"url": true, "secret": true, "verbose": true, "config": true}

// printWizardHint prints a hint showing how to invoke the command directly without the wizard.
func printWizardHint(cmd *cobra.Command, cliFlags map[string]bool) {
	if !interactive.IsInteractive() {
		return
	}

	// Check if any flags were set interactively (not present in the CLI snapshot)
	anyInteractive := false
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Changed && !cliFlags[f.Name] && !globalFlags[f.Name] {
			anyInteractive = true
		}
	})
	if !anyInteractive {
		return
	}

	var parts []string
	parts = append(parts, "bigbluebutton-cli", cmd.Name())
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed || globalFlags[f.Name] {
			return
		}
		switch f.Value.Type() {
		case "bool":
			if f.Value.String() == "true" {
				parts = append(parts, "--"+f.Name)
			} else if f.DefValue == "true" {
				parts = append(parts, "--"+f.Name+"=false")
			}
		case "stringSlice":
			val := strings.Trim(f.Value.String(), "[]")
			if val != "" {
				for _, v := range strings.Split(val, ",") {
					parts = append(parts, fmt.Sprintf("--%s %q", f.Name, v))
				}
			}
		default:
			if f.Value.String() != "" {
				parts = append(parts, fmt.Sprintf("--%s %q", f.Name, f.Value.String()))
			}
		}
	})

	fmt.Fprintf(os.Stderr, "\nHint: You could also have invoked this command directly instead of using the wizard:\n  %s\n\n", strings.Join(parts, " "))
}

// withWizardHint wraps a RunE function to print a wizard bypass hint after successful execution.
func withWizardHint(fn func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cliFlags := collectChangedFlags(cmd)
		err := fn(cmd, args)
		if err == nil {
			printWizardHint(cmd, cliFlags)
		}
		return err
	}
}

// setFlagFromWizard stores an interactively-collected value back on the cobra flag.
func setFlagFromWizard(cmd *cobra.Command, flagName string, val interface{}) {
	switch v := val.(type) {
	case string:
		cmd.Flags().Set(flagName, v)
	case bool:
		cmd.Flags().Set(flagName, fmt.Sprintf("%v", v))
	case int:
		cmd.Flags().Set(flagName, fmt.Sprintf("%d", v))
	}
}
