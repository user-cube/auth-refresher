package ui

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// PrintError prints a formatted error message and exits if exitOnError is true
// If err is nil, only the message is displayed
func PrintError(msg string, err error, exitOnError bool) {
	colors := NewColors()
	if err != nil {
		fmt.Printf("%s %s: %v\n", colors.Red("✗"), msg, err)
	} else {
		fmt.Printf("%s %s\n", colors.Red("✗"), msg)
	}
	if exitOnError {
		os.Exit(1)
	}
}

// PrintSuccess prints a formatted success message
func PrintSuccess(msg string, details ...string) {
	colors := NewColors()
	fmt.Printf("%s %s", colors.Green("✓"), msg)

	for _, detail := range details {
		fmt.Printf(" %s", colors.Cyan(detail))
	}
	fmt.Println()
}

// PrintWarning prints a formatted warning message
func PrintWarning(msg string, details ...string) {
	colors := NewColors()
	fmt.Printf("%s %s", colors.Yellow("!"), msg)

	for _, detail := range details {
		fmt.Printf(" %s", colors.Cyan(detail))
	}
	fmt.Println()
}

// PrintInfo prints a formatted information label and value
func PrintInfo(label string, value string) {
	colors := NewColors()
	fmt.Printf("%s: %s\n", colors.Bold(label), value)
}

// PrintNote prints a formatted note message with an info icon
func PrintNote(msg string, details ...string) {
	colors := NewColors()
	// Using blue color with info icon for notes
	blue := color.New(color.FgBlue, color.Bold).SprintFunc()
	fmt.Printf("%s %s", blue("ℹ"), blue("Note:"))

	fmt.Printf(" %s", msg)

	for _, detail := range details {
		fmt.Printf(" %s", colors.Cyan(detail))
	}
	fmt.Println()
}
