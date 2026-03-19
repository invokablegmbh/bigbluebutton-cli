package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the BigBlueButton connection configuration.
type Config struct {
	URL    string
	Secret string
}

// Load resolves configuration using the cascade:
// CLI flags > ENV variables > config file > BBB autodetection > interactive wizard.
// Set allowInteractive to true to enable the wizard fallback.
func Load(flagURL, flagSecret, configPath string, allowInteractive bool) (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("/etc/bigbluebutton-cli")
	v.AddConfigPath("$HOME/.bigbluebutton-cli")
	v.SetEnvPrefix("BBB")
	v.AutomaticEnv()

	if configPath != "" {
		v.SetConfigFile(configPath)
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok && configPath != "" {
			return nil, fmt.Errorf("reading config file: %w", err)
		}
		log.Printf("No config file found, continuing with other sources")
	} else {
		log.Printf("Using config file: %s", v.ConfigFileUsed())
	}

	// Step 1: CLI flags take highest priority
	if flagURL != "" {
		v.Set("url", flagURL)
	}
	if flagSecret != "" {
		v.Set("secret", flagSecret)
	}

	cfg := &Config{
		URL:    v.GetString("url"),
		Secret: v.GetString("secret"),
	}

	// Step 4: Autodetect from local BBB installation
	if cfg.URL == "" || cfg.Secret == "" {
		log.Printf("Attempting BigBlueButton autodetection...")
		if auto, err := Autodetect(); err == nil {
			log.Printf("Autodetection successful: URL=%s", auto.URL)
			if cfg.URL == "" {
				cfg.URL = auto.URL
			}
			if cfg.Secret == "" {
				cfg.Secret = auto.Secret
			}
		} else {
			log.Printf("Autodetection failed: %v", err)
		}
	}

	// Step 5: Interactive wizard
	if allowInteractive && (cfg.URL == "" || cfg.Secret == "") {
		log.Printf("Starting interactive configuration wizard...")
		if err := interactiveConfig(cfg); err != nil {
			return nil, fmt.Errorf("interactive configuration: %w", err)
		}
	}

	// Final validation
	if cfg.URL == "" || cfg.Secret == "" {
		return nil, fmt.Errorf("BigBlueButton URL and secret are required.\n" +
			"Provide them via:\n" +
			"  - Flags: --url and --secret\n" +
			"  - Environment: BBB_URL and BBB_SECRET\n" +
			"  - Config file: /etc/bigbluebutton-cli/config.yaml\n" +
			"  - Or run in a terminal for interactive setup")
	}

	// Normalize URL
	cfg.URL = strings.TrimRight(cfg.URL, "/")

	return cfg, nil
}

// interactiveConfig prompts the user for missing configuration values.
func interactiveConfig(cfg *Config) error {
	// Import cycle prevention: this calls into the interactive package
	// through a function variable that gets set by the cmd package.
	if promptFunc == nil {
		return fmt.Errorf("interactive mode not available")
	}
	return promptFunc(cfg)
}

// PromptFunc is a function type for interactive prompting.
type PromptFunc func(cfg *Config) error

var promptFunc PromptFunc

// SetPromptFunc sets the function used for interactive configuration.
func SetPromptFunc(f PromptFunc) {
	promptFunc = f
}
