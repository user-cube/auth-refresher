package auth

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/user-cube/auth-refresher/pkg/ui"
	"gopkg.in/yaml.v3"
)

type Config struct {
	CurrentRegistry string              `yaml:"last_used_registry"`
	Registries      map[string]Registry `yaml:"registries"`
}

type Registry struct {
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
	URL        string `yaml:"url"`
	Region     string `yaml:"region"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	LastLogin  string `yaml:"last_login"`  // Field to store the last login date
	LastLogout string `yaml:"last_logout"` // Field to store the last logout date
}

// LoadConfig loads the configuration from the given file path
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

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

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return err
	}

	// Ensure the current registry appears on top with extra info
	currentRegistry := config.CurrentRegistry
	var sortedOptions []string
	for key, registry := range config.Registries {
		if key == currentRegistry {
			sortedOptions = append([]string{key + " (last one used)|" + registry.Type}, sortedOptions...)
		} else {
			sortedOptions = append(sortedOptions, key+"|"+registry.Type)
		}
	}
	sort.Strings(sortedOptions[1:]) // Sort the rest of the options, excluding the first (current registry)

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

	// Strip the label '(last one used)' from the selected registry name
	cleanedSelected := strings.TrimSuffix(selected, " (last one used)")
	registry, exists := config.Registries[cleanedSelected]
	if !exists {
		return fmt.Errorf("registry '%s' not found in the configuration", cleanedSelected)
	}

	// Ensure the registry type is not empty
	if registry.Type == "" {
		return fmt.Errorf("registry '%s' has no type defined in the configuration", cleanedSelected)
	}

	// Validate the registry type before starting the spinner
	if registry.Type != "aws" && registry.Type != "helm" && registry.Type != "docker" {
		return fmt.Errorf("unsupported registry type: %s", registry.Type)
	}

	if registry.Type == "docker" {
		password := registry.Password
		if password == "" {
			var err error
			password, err = ui.PromptInput(ctx, "Enter your Docker password", true) // Enable masking for password input
			if err != nil {
				fmt.Println("Error reading password:", err)
				return err
			}
		}
		registry.Password = password // Update the registry object with the password
	}

	// Start the spinner after gathering necessary inputs
	err = ui.WithSpinner("Logging in to the selected registry", func() error {
		switch registry.Type {
		case "docker":
			loginCmd := exec.CommandContext(ctx, "docker", "login", "--username", registry.Username, "--password", registry.Password, registry.URL)
			if err := loginCmd.Run(); err != nil {
				return err
			}
		case "aws":
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
		case "helm":
			cmd := exec.CommandContext(ctx, "aws", "ecr", "get-login-password", "--region", registry.Region)
			output, err := cmd.Output()
			if err != nil {
				return err
			}
			loginCmd := exec.CommandContext(ctx, "helm", "registry", "login", registry.URL, "--username", "AWS", "--password", string(output))
			if err := loginCmd.Run(); err != nil {
				return err
			}
		}
		return nil
	}, false)
	if err != nil {
		return err
	}

	// Update the `last_used_registry` field in the configuration
	config.CurrentRegistry = cleanedSelected
	registry.Password = ""                                        // Clear the password field for security reasons
	registry.LastLogin = time.Now().Format("2006-01-02 15:04:05") // Update the `LastLogin` field with the current date
	config.Registries[cleanedSelected] = registry                 // Update the registry entry in the configuration

	file, err = os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to open config file for writing: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	encoder := yaml.NewEncoder(file)
	defer func() {
		if err := encoder.Close(); err != nil {
			fmt.Println("Error closing encoder:", err)
		}
	}()
	if err := encoder.Encode(&config); err != nil {
		return fmt.Errorf("failed to write updated config: %w", err)
	}

	return nil
}
