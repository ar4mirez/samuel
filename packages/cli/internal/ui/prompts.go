package ui

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

// SelectOption represents an option in a select prompt
type SelectOption struct {
	Name        string
	Description string
	Value       string
}

// Select prompts the user to select one option from a list
func Select(label string, options []SelectOption) (SelectOption, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▸ {{ .Name | cyan }} - {{ .Description | faint }}",
		Inactive: "  {{ .Name }} - {{ .Description | faint }}",
		Selected: "{{ \"✓\" | green }} {{ .Name | cyan }}",
	}

	prompt := promptui.Select{
		Label:     label,
		Items:     options,
		Templates: templates,
		Size:      10,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return SelectOption{}, err
	}

	return options[idx], nil
}

// MultiSelect prompts the user to select multiple options
// Returns the selected options
func MultiSelect(label string, options []SelectOption, defaults []string) ([]SelectOption, error) {
	// Track selected state
	selected := make(map[int]bool)
	for i, opt := range options {
		for _, d := range defaults {
			if opt.Value == d {
				selected[i] = true
				break
			}
		}
	}

	for {
		// Build display items
		type displayItem struct {
			Index       int
			Name        string
			Description string
			Selected    bool
		}

		items := make([]displayItem, len(options)+1)
		for i, opt := range options {
			items[i] = displayItem{
				Index:       i,
				Name:        opt.Name,
				Description: opt.Description,
				Selected:    selected[i],
			}
		}
		// Add "Done" option at the end
		items[len(options)] = displayItem{
			Index:       -1,
			Name:        "Done",
			Description: "Finish selection",
			Selected:    false,
		}

		templates := &promptui.SelectTemplates{
			Label: "{{ . }}",
			Active: `▸ {{ if eq .Index -1 }}{{ .Name | green }}{{ else }}{{ if .Selected }}[✓]{{ else }}[ ]{{ end }} {{ .Name | cyan }}{{ end }} - {{ .Description | faint }}`,
			Inactive: `  {{ if eq .Index -1 }}{{ .Name }}{{ else }}{{ if .Selected }}[✓]{{ else }}[ ]{{ end }} {{ .Name }}{{ end }} - {{ .Description | faint }}`,
			Selected: "",
		}

		prompt := promptui.Select{
			Label:     label + " (space to toggle, enter on Done to finish)",
			Items:     items,
			Templates: templates,
			Size:      12,
		}

		idx, _, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		// If "Done" selected, return results
		if idx == len(options) {
			break
		}

		// Toggle selection
		selected[idx] = !selected[idx]
	}

	// Collect selected options
	var result []SelectOption
	for i, opt := range options {
		if selected[i] {
			result = append(result, opt)
		}
	}

	return result, nil
}

// Confirm prompts for yes/no confirmation
func Confirm(label string, defaultYes bool) (bool, error) {
	suffix := " [y/N]"
	if defaultYes {
		suffix = " [Y/n]"
	}

	prompt := promptui.Prompt{
		Label:     label + suffix,
		IsConfirm: true,
		Default:   "",
	}

	result, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			return false, nil
		}
		// Empty response uses default
		if result == "" {
			return defaultYes, nil
		}
		return false, err
	}

	result = strings.ToLower(strings.TrimSpace(result))
	return result == "y" || result == "yes", nil
}

// Input prompts for text input
func Input(label string, defaultValue string, validate func(string) error) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Default:  defaultValue,
		Validate: validate,
	}

	return prompt.Run()
}

// InputWithPlaceholder prompts for text input with a placeholder hint
func InputWithPlaceholder(label string, placeholder string) (string, error) {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}: ",
		Valid:   "{{ . | green }}: ",
		Invalid: "{{ . | red }}: ",
		Success: "{{ . | bold }}: ",
	}

	prompt := promptui.Prompt{
		Label:     label,
		Templates: templates,
		Default:   "",
	}

	if placeholder != "" {
		fmt.Printf("  (e.g., %s)\n", placeholder)
	}

	return prompt.Run()
}
