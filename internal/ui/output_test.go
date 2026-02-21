package ui

import (
	"os"
	"strings"
	"testing"

	"github.com/fatih/color"
)

// captureStdout runs fn while capturing os.Stdout and color.Output,
// returning the captured text.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}

	origStdout := os.Stdout
	origColorOut := color.Output
	os.Stdout = w
	color.Output = w

	fn()

	w.Close()
	os.Stdout = origStdout
	color.Output = origColorOut

	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	r.Close()

	return string(buf[:n])
}

// captureStderr runs fn while capturing os.Stderr, returning the captured text.
func captureStderr(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}

	orig := os.Stderr
	os.Stderr = w

	fn()

	w.Close()
	os.Stderr = orig

	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	r.Close()

	return string(buf[:n])
}

func TestDisableColors(t *testing.T) {
	orig := color.NoColor
	defer func() { color.NoColor = orig }()

	color.NoColor = false
	DisableColors()

	if !color.NoColor {
		t.Error("DisableColors did not set color.NoColor to true")
	}
}

func TestSuccess(t *testing.T) {
	orig := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = orig }()

	tests := []struct {
		name     string
		format   string
		args     []interface{}
		contains string
	}{
		{"simple", "done", nil, SuccessSymbol + " done"},
		{"formatted", "step %d of %d", []interface{}{1, 3}, SuccessSymbol + " step 1 of 3"},
		{"string_arg", "installed %s", []interface{}{"go-guide"}, SuccessSymbol + " installed go-guide"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureStdout(t, func() {
				Success(tt.format, tt.args...)
			})
			if !strings.Contains(got, tt.contains) {
				t.Errorf("got %q, want contains %q", got, tt.contains)
			}
			if !strings.HasSuffix(strings.TrimRight(got, "\r"), "\n") {
				t.Error("output should end with newline")
			}
		})
	}
}

func TestError(t *testing.T) {
	orig := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = orig }()

	tests := []struct {
		name     string
		format   string
		args     []interface{}
		contains string
	}{
		{"simple", "failed", nil, ErrorSymbol + " failed"},
		{"formatted", "exit code %d", []interface{}{1}, ErrorSymbol + " exit code 1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureStderr(t, func() {
				Error(tt.format, tt.args...)
			})
			if !strings.Contains(got, tt.contains) {
				t.Errorf("got %q, want contains %q", got, tt.contains)
			}
		})
	}
}

func TestWarn(t *testing.T) {
	orig := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = orig }()

	got := captureStdout(t, func() {
		Warn("check %s", "permissions")
	})
	want := WarnSymbol + " check permissions"
	if !strings.Contains(got, want) {
		t.Errorf("got %q, want contains %q", got, want)
	}
}

func TestInfo(t *testing.T) {
	orig := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = orig }()

	got := captureStdout(t, func() {
		Info("downloading %s", "v1.0.0")
	})
	want := InfoSymbol + " downloading v1.0.0"
	if !strings.Contains(got, want) {
		t.Errorf("got %q, want contains %q", got, want)
	}
}

func TestPrint(t *testing.T) {
	got := captureStdout(t, func() {
		Print("hello %s", "world")
	})
	if !strings.Contains(got, "hello world") {
		t.Errorf("got %q, want contains %q", got, "hello world")
	}
	if !strings.HasSuffix(strings.TrimRight(got, "\r"), "\n") {
		t.Error("Print should append newline")
	}
}

func TestBold(t *testing.T) {
	orig := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = orig }()

	got := captureStdout(t, func() {
		Bold("important %s", "message")
	})
	if !strings.Contains(got, "important message") {
		t.Errorf("got %q, want contains %q", got, "important message")
	}
}

func TestDim(t *testing.T) {
	orig := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = orig }()

	got := captureStdout(t, func() {
		Dim("subtle %s", "hint")
	})
	if !strings.Contains(got, "subtle hint") {
		t.Errorf("got %q, want contains %q", got, "subtle hint")
	}
}

func TestHeader(t *testing.T) {
	orig := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = orig }()

	got := captureStdout(t, func() {
		Header("My Section")
	})
	if !strings.Contains(got, "My Section") {
		t.Errorf("got %q, want contains %q", got, "My Section")
	}
	// Header adds blank lines above and below the title
	lines := strings.Split(got, "\n")
	if len(lines) < 3 {
		t.Errorf("expected at least 3 lines (blank, title, blank), got %d", len(lines))
	}
}

func TestSection(t *testing.T) {
	got := captureStdout(t, func() {
		Section("Details")
	})
	if !strings.Contains(got, "Details:") {
		t.Errorf("got %q, want contains %q", got, "Details:")
	}
}

func TestListItem(t *testing.T) {
	tests := []struct {
		name   string
		indent int
		format string
		args   []interface{}
		want   string
	}{
		{"no_indent", 0, "item %d", []interface{}{1}, "item 1\n"},
		{"indent_1", 1, "nested", nil, "  nested\n"},
		{"indent_2", 2, "deep", nil, "    deep\n"},
		{"indent_3", 3, "deeper %s", []interface{}{"value"}, "      deeper value\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureStdout(t, func() {
				ListItem(tt.indent, tt.format, tt.args...)
			})
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSuccessItem(t *testing.T) {
	orig := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = orig }()

	tests := []struct {
		name     string
		indent   int
		format   string
		args     []interface{}
		contains string
	}{
		{"no_indent", 0, "ok", nil, SuccessSymbol + " ok"},
		{"indent_1", 1, "good", nil, "  " + SuccessSymbol + " good"},
		{"formatted", 0, "installed %s", []interface{}{"go"}, SuccessSymbol + " installed go"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureStdout(t, func() {
				SuccessItem(tt.indent, tt.format, tt.args...)
			})
			if !strings.Contains(got, tt.contains) {
				t.Errorf("got %q, want contains %q", got, tt.contains)
			}
		})
	}
}

func TestWarnItem(t *testing.T) {
	orig := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = orig }()

	tests := []struct {
		name     string
		indent   int
		format   string
		args     []interface{}
		contains string
	}{
		{"no_indent", 0, "caution", nil, WarnSymbol + " caution"},
		{"indent_2", 2, "watch out", nil, "    " + WarnSymbol + " watch out"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureStdout(t, func() {
				WarnItem(tt.indent, tt.format, tt.args...)
			})
			if !strings.Contains(got, tt.contains) {
				t.Errorf("got %q, want contains %q", got, tt.contains)
			}
		})
	}
}

func TestErrorItem(t *testing.T) {
	orig := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = orig }()

	tests := []struct {
		name     string
		indent   int
		format   string
		args     []interface{}
		contains string
	}{
		{"no_indent", 0, "bad", nil, ErrorSymbol + " bad"},
		{"indent_1", 1, "missing %s", []interface{}{"file"}, "  " + ErrorSymbol + " missing file"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureStdout(t, func() {
				ErrorItem(tt.indent, tt.format, tt.args...)
			})
			if !strings.Contains(got, tt.contains) {
				t.Errorf("got %q, want contains %q", got, tt.contains)
			}
		})
	}
}

func TestTableRow(t *testing.T) {
	got := captureStdout(t, func() {
		TableRow("Version", "1.2.3")
	})
	if !strings.Contains(got, "Version:") {
		t.Errorf("got %q, want contains 'Version:'", got)
	}
	if !strings.Contains(got, "1.2.3") {
		t.Errorf("got %q, want contains '1.2.3'", got)
	}
	// Verify indentation (2 leading spaces)
	if !strings.HasPrefix(got, "  ") {
		t.Errorf("got %q, want 2-space indent prefix", got)
	}
}

func TestColoredTableRow(t *testing.T) {
	orig := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = orig }()

	got := captureStdout(t, func() {
		ColoredTableRow("Status", "healthy", color.New(color.FgGreen))
	})
	if !strings.Contains(got, "Status:") {
		t.Errorf("got %q, want contains 'Status:'", got)
	}
	if !strings.Contains(got, "healthy") {
		t.Errorf("got %q, want contains 'healthy'", got)
	}
	if !strings.HasPrefix(got, "  ") {
		t.Errorf("got %q, want 2-space indent prefix", got)
	}
}
