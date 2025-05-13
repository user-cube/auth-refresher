// Package ui provides user interface utilities for the gclone application.
package ui

import (
	"github.com/fatih/color"
)

var (
	// Success prints text in green color
	Success = color.New(color.FgGreen).PrintfFunc()
	// Info prints text in cyan color
	Info = color.New(color.FgCyan).PrintfFunc()
	// Warning prints text in yellow color
	Warning = color.New(color.FgYellow).PrintfFunc()
	// Error prints text in red color
	Error = color.New(color.FgRed).PrintfFunc()
	// Highlight prints text in yellow color, useful for emphasizing values
	Highlight = color.New(color.FgYellow).SprintFunc()
	// Normal prints text in white color
	Normal = color.New(color.FgWhite).PrintfFunc()
)

// Colors creates and returns commonly used colored print functions
type Colors struct {
	Red    func(a ...interface{}) string
	Green  func(a ...interface{}) string
	Yellow func(a ...interface{}) string
	Cyan   func(a ...interface{}) string
	Bold   func(a ...interface{}) string
	Faint  func(a ...interface{}) string
}

// NewColors returns initialized color functions for consistent UI styling
func NewColors() *Colors {
	return &Colors{
		Red:    color.New(color.FgRed, color.Bold).SprintFunc(),
		Green:  color.New(color.FgGreen, color.Bold).SprintFunc(),
		Yellow: color.New(color.FgYellow, color.Bold).SprintFunc(),
		Cyan:   color.New(color.FgCyan, color.Bold).SprintFunc(),
		Bold:   color.New(color.Bold).SprintFunc(),
		Faint:  color.New(color.Faint).SprintFunc(),
	}
}

// PrintKeyValue prints a key-value pair with the key highlighted
func PrintKeyValue(key, value string) {
	Normal("  %s: %s\n", key, value)
}
