package auth

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/user-cube/auth-refresher/pkg/ui"
	"gopkg.in/yaml.v3"
)

type Config struct {
	CurrentRegistry string              `yaml:"current_registry"`
	Name            string              `yaml:"name"`
	Registries      map[string]Registry `yaml:"registries"`
}

type Registry struct {
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	URL    string `yaml:"url"`
	Region string `yaml:"region"`
}

// LoadConfig loads the configuration from the given file path
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// Login handles the login process for the selected registry
func Login(registry Registry) error {
	switch registry.Type {
	case "aws":
		ui.PrintInfo("Logging into AWS ECR", registry.URL)
		cmd := fmt.Sprintf("aws ecr get-login-password --region %s | docker login --username AWS --password-stdin %s", registry.Region, registry.URL)
		if err := executeCommand(cmd); err != nil {
			return fmt.Errorf("failed to login to AWS ECR: %w", err)
		}
	case "helm":
		ui.PrintInfo("Logging into Helm registry", registry.URL)
		cmd := fmt.Sprintf("helm registry login %s", registry.URL)
		if err := executeCommand(cmd); err != nil {
			return fmt.Errorf("failed to login to Helm registry: %w", err)
		}
	default:
		return fmt.Errorf("unsupported registry type: %s", registry.Type)
	}

	ui.PrintSuccess("Successfully logged into", registry.Name)
	return nil
}

// executeCommand runs a shell command
func executeCommand(cmd string) error {
	fmt.Println("Executing:", cmd)
	return nil // Replace with actual command execution logic
}

// Updated `LoginToRegistry` function to ensure the spinner starts only after the user selects a registry
func LoginToRegistry(ctx context.Context, configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return err
	}

	var sortedOptions []string
	for key, registry := range config.Registries {
		sortedOptions = append(sortedOptions, key+"|"+registry.Type)
	}
	sort.Strings(sortedOptions)

	var options []string
	for _, sortedOption := range sortedOptions {
		parts := strings.Split(sortedOption, "|")
		options = append(options, parts[0])
	}

	selected, err := ui.SelectFromList(ctx, "Select a registry to login", options)
	if err != nil {
		if err.Error() == "operation cancelled by user" {
			return nil // Gracefully handle user cancellation
		}
		return err
	}

	registry := config.Registries[selected]

	// Updated spinner handling to use the new ClearSpinner function
	return ui.WithSpinner("Logging in to the selected registry", func() error {

		if registry.Type == "aws" {
			cmd := exec.CommandContext(ctx, "aws", "ecr", "get-login-password", "--region", registry.Region)
			output, err := cmd.Output()
			if err != nil {
				return err
			}

			loginCmd := exec.CommandContext(ctx, "docker", "login", "--username", "AWS", "--password-stdin", registry.URL)
			loginCmd.Stdin = bytes.NewReader(output)
			if err := loginCmd.Run(); err != nil {
				return err
			}

			// Ensure spinner is cleared before printing success message
			ui.ClearSpinner()
			ui.PrintSuccess("Successfully logged in to registry:", registry.Name)
		} else if registry.Type == "helm" {
			ui.ClearSpinner()
			ui.PrintInfo("Helm registry login is not implemented yet.", "")
		} else {
			return fmt.Errorf("unsupported registry type: %s", registry.Type)
		}
		return nil
	}, false)
}
