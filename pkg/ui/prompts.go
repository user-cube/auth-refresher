package ui

import (
	"github.com/manifoldco/promptui"
)

// SelectFromList prompts the user to select an item from a list
func SelectFromList(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return items[idx], nil
}

// Confirm prompts the user for a yes/no confirmation
func Confirm(label string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	_, err := prompt.Run()
	if err != nil {
		// promptui returns an error when the user selects "n"
		if err == promptui.ErrAbort {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// PromptInput prompts the user for text input
func PromptInput(label string, defaultValue string, validate promptui.ValidateFunc) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Default:  defaultValue,
		Validate: validate,
	}

	return prompt.Run()
}
