package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"github.com/user-cube/auth-refresher/pkg/auth"
	"github.com/user-cube/auth-refresher/pkg/ui"
	"gopkg.in/yaml.v3"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from a selected registry",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := filepath.Join(os.Getenv("HOME"), ".auth-refresher", "config.yaml")
		file, err := os.Open(configPath)
		if err != nil {
			ui.PrintError("Failed to open config file", err, true)
			return
		}

		var config auth.Config
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			ui.PrintError("Failed to parse config file", err, true)
			return
		}

		// Sort registries by name and then by type
		keys := make([]string, 0, len(config.Registries))
		for key := range config.Registries {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool {
			if config.Registries[keys[i]].Name == config.Registries[keys[j]].Name {
				return config.Registries[keys[i]].Type < config.Registries[keys[j]].Type
			}
			return config.Registries[keys[i]].Name < config.Registries[keys[j]].Name
		})

		selected, err := ui.SelectFromList(cmd.Context(), "Select a registry to logout", keys)
		if err != nil {
			if err.Error() == "operation cancelled by user" {
				return // Gracefully handle user cancellation
			}
			ui.PrintError("Failed to select a registry", err, true)
			return
		}

		// Clear credentials for the selected registry
		registry := config.Registries[selected]

		// Perform Docker logout for both AWS and Docker registries
		if registry.Type == "docker" || registry.Type == "aws" {
			logoutCmd := exec.Command("docker", "logout", registry.URL)
			if err := logoutCmd.Run(); err != nil {
				ui.PrintError("Failed to perform Docker logout", err, true)
				return
			}
		}

		// Perform Helm logout for Helm registries
		if registry.Type == "helm" {
			logoutCmd := exec.Command("helm", "registry", "logout", registry.URL)
			if err := logoutCmd.Run(); err != nil {
				ui.PrintError("Failed to perform Helm logout", err, true)
				return
			}
		}

		// Only update the `LastLogout` field with the current date
		registry.LastLogout = time.Now().Format("2006-01-02 15:04:05") // Set the last logout date
		config.Registries[selected] = registry                         // Update the registry entry in the configuration

		// Save the updated configuration
		file, err = os.Create(configPath)
		if err != nil {
			ui.PrintError("Failed to open config file for writing", err, true)
			return
		}
		defer func() {
			if err := file.Close(); err != nil {
				ui.PrintError("Failed to close file", err, true)
			}
		}()

		encoder := yaml.NewEncoder(file)
		defer func() {
			if err := encoder.Close(); err != nil {
				ui.PrintError("Failed to close encoder", err, true)
			}
		}()
		if err := encoder.Encode(&config); err != nil {
			ui.PrintError("Failed to write updated config", err, true)
			return
		}

		ui.PrintSuccess("Successfully logged out from registry:", selected)
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}