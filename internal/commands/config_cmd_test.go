package commands

import (
	"strings"
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

func TestValidateRegistryURL(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
		errMsg  string
	}{
		{"valid https", "https://github.com/myorg/myrepo", false, ""},
		{"valid https with path", "https://github.com/ar4mirez/samuel", false, ""},
		{"valid https with port", "https://github.com:443/myorg/myrepo", false, ""},
		{"http rejected", "http://github.com/myorg/myrepo", true, "HTTPS scheme"},
		{"empty scheme", "github.com/myorg/myrepo", true, "HTTPS scheme"},
		{"ftp rejected", "ftp://example.com/repo", true, "HTTPS scheme"},
		{"empty string", "", true, "HTTPS scheme"},
		{"plain text", "not-a-url", true, "HTTPS scheme"},
		{"scheme only", "https://", true, "must have a host"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRegistryURL(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRegistryURL(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateRegistryURL(%q) error = %v, want error containing %q", tt.value, err, tt.errMsg)
				}
			}
		})
	}
}
