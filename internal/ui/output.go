package ui

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	// Colors
	successColor = color.New(color.FgGreen)
	errorColor   = color.New(color.FgRed)
	warnColor    = color.New(color.FgYellow)
	infoColor    = color.New(color.FgCyan)
	boldColor    = color.New(color.Bold)
	dimColor     = color.New(color.Faint)

	// Symbols
	SuccessSymbol = "✓"
	ErrorSymbol   = "✗"
	WarnSymbol    = "⚠"
	InfoSymbol    = "→"
	PendingSymbol = "○"
	ActiveSymbol  = "●"
)

// DisableColors turns off colored output
func DisableColors() {
	color.NoColor = true
}

// Success prints a success message with green checkmark
func Success(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	successColor.Fprintf(os.Stdout, "%s %s\n", SuccessSymbol, msg)
}

// Error prints an error message with red X
func Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	errorColor.Fprintf(os.Stderr, "%s %s\n", ErrorSymbol, msg)
}

// Warn prints a warning message with yellow symbol
func Warn(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	warnColor.Fprintf(os.Stdout, "%s %s\n", WarnSymbol, msg)
}

// Info prints an info message with cyan arrow
func Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	infoColor.Fprintf(os.Stdout, "%s %s\n", InfoSymbol, msg)
}

// Print prints a plain message
func Print(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format+"\n", args...)
}

// Bold prints bold text
func Bold(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	boldColor.Fprintln(os.Stdout, msg)
}

// Dim prints dimmed/faint text
func Dim(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	dimColor.Fprintln(os.Stdout, msg)
}

// Header prints a section header
func Header(title string) {
	fmt.Println()
	boldColor.Println(title)
	fmt.Println()
}

// Section prints a subsection header
func Section(title string) {
	fmt.Printf("\n%s:\n", title)
}

// ListItem prints a list item with proper indentation
func ListItem(indent int, format string, args ...interface{}) {
	padding := ""
	for i := 0; i < indent; i++ {
		padding += "  "
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s%s\n", padding, msg)
}

// SuccessItem prints a success list item
func SuccessItem(indent int, format string, args ...interface{}) {
	padding := ""
	for i := 0; i < indent; i++ {
		padding += "  "
	}
	msg := fmt.Sprintf(format, args...)
	successColor.Fprintf(os.Stdout, "%s%s %s\n", padding, SuccessSymbol, msg)
}

// WarnItem prints a warning list item
func WarnItem(indent int, format string, args ...interface{}) {
	padding := ""
	for i := 0; i < indent; i++ {
		padding += "  "
	}
	msg := fmt.Sprintf(format, args...)
	warnColor.Fprintf(os.Stdout, "%s%s %s\n", padding, WarnSymbol, msg)
}

// ErrorItem prints an error list item
func ErrorItem(indent int, format string, args ...interface{}) {
	padding := ""
	for i := 0; i < indent; i++ {
		padding += "  "
	}
	msg := fmt.Sprintf(format, args...)
	errorColor.Fprintf(os.Stdout, "%s%s %s\n", padding, ErrorSymbol, msg)
}

// Table helpers for aligned output

// TableRow prints a row with key-value alignment
func TableRow(key, value string) {
	fmt.Printf("  %-20s %s\n", key+":", value)
}

// ColoredTableRow prints a row with colored value
func ColoredTableRow(key, value string, c *color.Color) {
	fmt.Printf("  %-20s ", key+":")
	c.Println(value)
}
