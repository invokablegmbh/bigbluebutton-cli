package api

import "fmt"

// APIError represents an error returned by the BigBlueButton API.
type APIError struct {
	ReturnCode string
	MessageKey string
	Message    string
}

func (e *APIError) Error() string {
	if e.MessageKey != "" {
		return fmt.Sprintf("BBB API error: %s (key: %s)", e.Message, e.MessageKey)
	}
	return fmt.Sprintf("BBB API error: %s", e.Message)
}
