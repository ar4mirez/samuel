package ui

import (
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

// Spinner provides a simple loading spinner
type Spinner struct {
	bar     *progressbar.ProgressBar
	message string
}

// NewSpinner creates a new spinner with the given message
func NewSpinner(message string) *Spinner {
	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionSetDescription(message),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionClearOnFinish(),
	)

	return &Spinner{
		bar:     bar,
		message: message,
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	go func() {
		for {
			if s.bar == nil {
				return
			}
			s.bar.Add(1)
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

// Stop halts the spinner
func (s *Spinner) Stop() {
	if s.bar != nil {
		s.bar.Finish()
		s.bar = nil
	}
}

// Success stops the spinner and prints a success message
func (s *Spinner) Success(message string) {
	s.Stop()
	Success(message)
}

// Error stops the spinner and prints an error message
func (s *Spinner) Error(message string) {
	s.Stop()
	Error(message)
}

// ProgressBar creates a determinate progress bar
type ProgressBar struct {
	bar *progressbar.ProgressBar
}

// NewProgressBar creates a new progress bar with the given max value
func NewProgressBar(max int, description string) *ProgressBar {
	bar := progressbar.NewOptions(max,
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionSetDescription(description),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(40),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	return &ProgressBar{bar: bar}
}

// Add increments the progress bar
func (p *ProgressBar) Add(n int) {
	p.bar.Add(n)
}

// Finish completes the progress bar
func (p *ProgressBar) Finish() {
	p.bar.Finish()
	fmt.Println()
}
