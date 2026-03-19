package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const bbbPropertiesPath = "/etc/bigbluebutton/bbb-web.properties"

// Autodetect attempts to detect BigBlueButton configuration from
// the local installation at /etc/bigbluebutton/bbb-web.properties.
func Autodetect() (*Config, error) {
	return autodetectFromFile(bbbPropertiesPath)
}

func autodetectFromFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening properties file: %w", err)
	}
	defer f.Close()

	cfg := &Config{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		switch key {
		case "bigbluebutton.web.serverURL":
			cfg.URL = value
		case "securitySalt":
			cfg.Secret = value
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading properties file: %w", err)
	}

	if cfg.URL == "" || cfg.Secret == "" {
		return nil, fmt.Errorf("incomplete configuration in %s", path)
	}

	return cfg, nil
}
