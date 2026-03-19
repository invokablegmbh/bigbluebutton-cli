package interactive

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	"golang.org/x/term"
)

const doneLabel = "✔ Done - confirm selection"

// IsInteractive returns true if stdin is a terminal.
func IsInteractive() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

// newForm creates a form with a constrained height so titles stay visible.
func newForm(groups ...*huh.Group) *huh.Form {
	return huh.NewForm(groups...).WithHeight(6)
}

// AskString prompts the user for a string value.
func AskString(label, defaultVal string) (string, error) {
	var result string
	input := huh.NewInput().
		Title(label).
		Value(&result)
	if defaultVal != "" {
		input.Placeholder(defaultVal)
	}
	err := newForm(huh.NewGroup(input)).Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %w", err)
	}
	if result == "" {
		result = defaultVal
	}
	return result, nil
}

// AskPassword prompts the user for a password (hidden input).
func AskPassword(label string) (string, error) {
	var result string
	err := newForm(
		huh.NewGroup(
			huh.NewInput().
				Title(label).
				EchoMode(huh.EchoModePassword).
				Value(&result),
		),
	).Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %w", err)
	}
	return result, nil
}

// AskConfirm prompts the user for a yes/no confirmation.
func AskConfirm(label string, defaultVal bool) (bool, error) {
	result := defaultVal
	err := newForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(label).
				Value(&result),
		),
	).Run()
	if err != nil {
		return false, fmt.Errorf("prompt failed: %w", err)
	}
	return result, nil
}

// AskSelect prompts the user to select from a list of options.
func AskSelect(label string, options []string) (string, error) {
	var result string
	opts := make([]huh.Option[string], len(options))
	for i, o := range options {
		opts[i] = huh.NewOption(o, o)
	}
	height := len(options) + 4
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(label).
				Options(opts...).
				Value(&result),
		),
	).WithHeight(height).Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %w", err)
	}
	return result, nil
}

// AskMultiSelect prompts the user to select multiple options from a list.
// Toggle items with Enter, submit with Alt+Enter or by selecting the "Done" entry.
func AskMultiSelect(label string, options []string) ([]string, error) {
	var result []string

	// Add all real options plus a "Done" entry at the bottom.
	opts := make([]huh.Option[string], 0, len(options)+1)
	for _, o := range options {
		opts = append(opts, huh.NewOption(o, o))
	}
	opts = append(opts, huh.NewOption(doneLabel, doneLabel))

	// Custom keymap: Enter toggles, Alt+Enter submits.
	km := huh.NewDefaultKeyMap()
	km.MultiSelect.Toggle = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "toggle"))
	km.MultiSelect.Next = key.NewBinding(key.WithKeys("alt+enter"), key.WithHelp("alt+enter", "confirm"))
	km.MultiSelect.Submit = key.NewBinding(key.WithKeys("alt+enter"), key.WithHelp("alt+enter", "confirm"))

	ms := huh.NewMultiSelect[string]().
		Title(label).
		Options(opts...).
		Value(&result).
		Height(15).
		Filterable(false)

	err := huh.NewForm(huh.NewGroup(ms)).
		WithKeyMap(km).
		WithHeight(17).
		Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	// If "Done" was toggled, treat it as submit and strip it from results.
	filtered := make([]string, 0, len(result))
	for _, r := range result {
		if r != doneLabel {
			filtered = append(filtered, r)
		}
	}
	return filtered, nil
}

// AskInt prompts the user for an integer value.
func AskInt(label string, defaultVal int) (int, error) {
	var result string
	err := newForm(
		huh.NewGroup(
			huh.NewInput().
				Title(label).
				Placeholder(fmt.Sprintf("%d", defaultVal)).
				Value(&result),
		),
	).Run()
	if err != nil {
		return 0, fmt.Errorf("prompt failed: %w", err)
	}
	if result == "" {
		return defaultVal, nil
	}
	var val int
	_, err = fmt.Sscanf(result, "%d", &val)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", result)
	}
	return val, nil
}
