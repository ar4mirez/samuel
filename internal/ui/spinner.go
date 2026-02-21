package ui

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

// Spinner provides a simple loading spinner
type Spinner struct {
	bar      *progressbar.ProgressBar
	message  string
	done     chan struct{}
	stopOnce sync.Once
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
		done:    make(chan struct{}),
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-s.done:
				return
			case <-ticker.C:
				_ = s.bar.Add(1)
			}
		}
	}()
}

// Stop halts the spinner. Safe to call multiple times and concurrently.
func (s *Spinner) Stop() {
	s.stopOnce.Do(func() {
		close(s.done)
		_ = s.bar.Finish()
	})
}

// Success stops the spinner and prints a success message
func (s *Spinner) Success(message string) {
	s.Stop()
	Success("%s", message)
}

// Error stops the spinner and prints an error message
func (s *Spinner) Error(message string) {
	s.Stop()
	Error("%s", message)
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
	_ = p.bar.Add(n)
}

// Finish completes the progress bar
func (p *ProgressBar) Finish() {
	_ = p.bar.Finish()
	fmt.Println()
}
