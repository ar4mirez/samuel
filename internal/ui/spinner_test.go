package ui

import (
	"sync"
	"testing"
	"time"
)

func TestSpinner_StartStop(t *testing.T) {
	s := NewSpinner("testing")
	s.Start()
	time.Sleep(250 * time.Millisecond)
	s.Stop()
}

func TestSpinner_StopIdempotent(t *testing.T) {
	s := NewSpinner("testing")
	s.Start()
	time.Sleep(150 * time.Millisecond)

	// Calling Stop multiple times must not panic (double close).
	s.Stop()
	s.Stop()
	s.Stop()
}

func TestSpinner_StopWithoutStart(t *testing.T) {
	s := NewSpinner("testing")
	// Stop without Start must not panic.
	s.Stop()
}

func TestSpinner_ConcurrentStop(t *testing.T) {
	s := NewSpinner("testing")
	s.Start()
	time.Sleep(150 * time.Millisecond)

	// Race multiple Stop calls from different goroutines.
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Stop()
		}()
	}
	wg.Wait()
}

func TestSpinner_RapidStartStop(t *testing.T) {
	// Start and stop immediately — goroutine must exit cleanly.
	s := NewSpinner("testing")
	s.Start()
	s.Stop()
	// Allow goroutine to observe the done channel and exit.
	time.Sleep(50 * time.Millisecond)
}

func TestSpinner_Success(t *testing.T) {
	s := NewSpinner("testing")
	s.Start()
	time.Sleep(150 * time.Millisecond)
	s.Success("done")
}

func TestSpinner_Error(t *testing.T) {
	s := NewSpinner("testing")
	s.Start()
	time.Sleep(150 * time.Millisecond)
	s.Error("failed")
}
