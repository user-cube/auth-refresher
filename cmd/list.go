package cmd

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/user-cube/auth-refresher/pkg/auth"
	"github.com/user-cube/auth-refresher/pkg/ui"
	"gopkg.in/yaml.v3"
)

// Update references to Config
var config auth.Config

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registries in a table format",
	Run: func(cmd *cobra.Command, args []string) {

		configPath := filepath.Join(os.Getenv("HOME"), ".auth-refresher", "config.yaml")
		file, err := os.Open(configPath)
		if err != nil {
			ui.PrintError("Failed to open config file", err, true)
			return
		}
		defer func() {
			if err := file.Close(); err != nil {
				ui.PrintError("Failed to close file", err, true)
			}
		}()

		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			ui.PrintError("Failed to parse config file", err, true)
			return
		}

		// Sort registries by name and then by type
		sortedKeys := make([]string, 0, len(config.Registries))
		for key := range config.Registries {
			sortedKeys = append(sortedKeys, key)
		}
		sort.Slice(sortedKeys, func(i, j int) bool {
			if config.Registries[sortedKeys[i]].Name == config.Registries[sortedKeys[j]].Name {
				return config.Registries[sortedKeys[i]].Type < config.Registries[sortedKeys[j]].Type
			}
			return config.Registries[sortedKeys[i]].Name < config.Registries[sortedKeys[j]].Name
		})

		// Use sortedKeys to iterate and display registries
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Name", "Type", "URL", "Region", "Last Login", "Last Logout"})
		for _, key := range sortedKeys {
			registry := config.Registries[key]
			t.AppendRow(table.Row{registry.Name, registry.Type, registry.URL, registry.Region, registry.LastLogin, registry.LastLogout})
		}

		t.Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
