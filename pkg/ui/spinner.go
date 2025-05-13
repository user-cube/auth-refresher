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

// Enhanced `Stop` method to flush output and ensure the spinner's message is fully cleared
func (s *Spinner) Stop() {
	if !s.active {
		return
	}
	s.active = false
	s.done <- true
	fmt.Print("\r\033[K") // Clear the line
	fmt.Println()         // Move to a new line to ensure no overlap
}

// Added ClearSpinner function to clear spinner output
// Added \r so the whole line is cleared
func ClearSpinner() {
	fmt.Print("\r\033[K") // move to column 0, then erase line
}

// Updated `WithSpinner` to allow suppressing the default completion message
func WithSpinner(message string, fn func() error, suppressCompletionMessage bool) error {
	spinner := NewSpinner(message)
	spinner.Start()
	err := fn()
	spinner.Stop()

	return err
}
