package ui

import (
	"context"
	"fmt"

	"github.com/manifoldco/promptui"
)

// SelectFromList prompts the user to select an item from a list with context support
func SelectFromList(ctx context.Context, label string, items []string) (string, error) {
	resultChan := make(chan string, 1)
	errorChan := make(chan error, 1)

	go func() {
		prompt := promptui.Select{
			Label: label,
			Items: items,
			Templates: &promptui.SelectTemplates{
				Active:   "▶ {{ . | cyan }}", // Highlight the active selection in cyan
				Inactive: "  {{ . }}",
				Selected: "✔ {{ . | green }}", // Show the selected item in green
			},
		}

		_, result, err := prompt.Run()
		if err != nil {
			if err == promptui.ErrInterrupt {
				errorChan <- fmt.Errorf("operation cancelled by user")
				return
			}
			errorChan <- err
			return
		}
		resultChan <- result
	}()

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("operation cancelled by user")
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return "", err
	}
}

// ConfirmWithContext prompts the user for a yes/no confirmation with context support
func ConfirmWithContext(ctx context.Context, label string) (bool, error) {
	resultChan := make(chan bool, 1)
	errorChan := make(chan error, 1)

	go func() {
		prompt := promptui.Prompt{
			Label:     label,
			IsConfirm: true,
		}

		_, err := prompt.Run()
		if err != nil {
			if err == promptui.ErrAbort {
				resultChan <- false
				return
			}
			errorChan <- err
			return
		}
		resultChan <- true
	}()

	select {
	case <-ctx.Done():
		return false, fmt.Errorf("operation cancelled by user")
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return false, err
	}
}

// PromptInputWithContext prompts the user for text input with context support
// Add an optional `mask` parameter to enable masking sensitive input
func PromptInputWithContext(ctx context.Context, label string, defaultValue string, validate promptui.ValidateFunc, mask bool) (string, error) {
	resultChan := make(chan string, 1)
	errorChan := make(chan error, 1)

	go func() {
		prompt := promptui.Prompt{
			Label:    label,
			Default:  defaultValue,
			Validate: validate,
		}
		if mask {
			prompt.Mask = '*'
		}
		result, err := prompt.Run()
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("operation cancelled by user")
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return "", err
	}
}

// Updated `PromptInput` wrapper to include masking support
func PromptInput(ctx context.Context, label string, mask bool) (string, error) {
	return PromptInputWithContext(ctx, label, "", nil, mask)
}
