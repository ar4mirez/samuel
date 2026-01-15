package commands

import (
	"testing"
)

func TestIsValidConfigKey(t *testing.T) {
	tests := []struct {
		key  string
		want bool
	}{
		{"version", true},
		{"registry", true},
		{"installed.languages", true},
		{"installed.frameworks", true},
		{"installed.workflows", true},
		{"invalid", false},
		{"", false},
		{"VERSION", false}, // Case sensitive
		{"installed", false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := isValidConfigKey(tt.key)
			if got != tt.want {
				t.Errorf("isValidConfigKey(%q) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestFormatConfigValue(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  string
	}{
		{"string", "hello", "hello"},
		{"empty string", "", "(empty)"},
		{"string slice", []string{"a", "b", "c"}, "a, b, c"},
		{"empty slice", []string{}, "(none)"},
		{"nil", nil, "<nil>"},
		{"int", 42, "42"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatConfigValue(tt.value)
			if got != tt.want {
				t.Errorf("formatConfigValue(%v) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}

func TestFormatConfigValue_SingleItem(t *testing.T) {
	got := formatConfigValue([]string{"single"})
	if got != "single" {
		t.Errorf("formatConfigValue([single]) = %q, want %q", got, "single")
	}
}
