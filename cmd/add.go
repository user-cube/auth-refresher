package cmd

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/user-cube/auth-refresher/pkg/auth"
	"github.com/user-cube/auth-refresher/pkg/ui"
	"gopkg.in/yaml.v3"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new registry",
	Run: func(cmd *cobra.Command, args []string) {
		// Setup signal handling for graceful exit
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			<-c
			ui.PrintInfo("Operation cancelled by user", "")
			cancel()
			os.Exit(0)
		}()

		configPath := filepath.Join(os.Getenv("HOME"), ".auth-refresher", "config.yaml")
		file, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			ui.PrintError("Failed to open config file", err, true)
			return
		}
		defer file.Close()

		var config auth.Config
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil && err.Error() != "EOF" {
			ui.PrintError("Failed to parse config file", err, true)
			return
		}

		if config.Registries == nil {
			config.Registries = make(map[string]auth.Registry)
		}

		name, err := ui.PromptInputWithContext(ctx, "Registry Name", "", nil)
		if err != nil {
			return
		}

		// Updated registry type input to use a selection instead of free typing
		typeOptions := []string{"aws", "helm"}
		typeInput, err := ui.SelectFromList(ctx, "Registry Type", typeOptions)
		if err != nil {
			return
		}

		url, err := ui.PromptInputWithContext(ctx, "Registry URL", "", nil)
		if err != nil {
			return
		}

		region, err := ui.PromptInputWithContext(ctx, "Registry Region", "", nil)
		if err != nil {
			return
		}

		config.Registries[name] = auth.Registry{
			Name:   name,
			Type:   typeInput,
			URL:    url,
			Region: region,
		}

		file.Truncate(0)
		file.Seek(0, 0)
		encoder := yaml.NewEncoder(file)
		if err := encoder.Encode(&config); err != nil {
			ui.PrintError("Failed to save config file", err, true)
			return
		}

		ui.PrintSuccess("Registry added successfully!", name)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
