package ui

import (
	"fmt"
	"time"
)

// Spinner represents a simple text spinner for indicating progress
type Spinner struct {
	message   string
	frames    []string
	frameRate time.Duration
	active    bool
	done      chan bool
}

// NewSpinner creates a new spinner with a message
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message:   message,
		frames:    []string{"|", "/", "-", "\\"},
		frameRate: 100 * time.Millisecond,
		active:    false,
		done:      make(chan bool),
	}
}

// Start starts the spinner
func (s *Spinner) Start() {
	if s.active {
		return
	}
	s.active = true

	go func() {
		i := 0
		for {
			select {
			case <-s.done:
				return
			default:
				// Use the Info function to color the spinner
				frame := s.frames[i]
				fmt.Printf("\r%s %s", frame, s.message)
				i = (i + 1) % len(s.frames)
				time.Sleep(s.frameRate)
			}
		}
	}()
}

// Stop stops the spinner and clears the line
func (s *Spinner) Stop() {
	if !s.active {
		return
	}
	s.active = false
	s.done <- true
	fmt.Print("\r")
	// Clear the line
	fmt.Print("\033[K")
}

// WithSpinner runs a function with a spinner
func WithSpinner(message string, fn func() error) error {
	spinner := NewSpinner(message)
	spinner.Start()
	err := fn()
	spinner.Stop()

	if err != nil {
		Error("Error: %v\n", err)
	} else {
		Success("Done: %s\n", message)
	}

	return err
}
