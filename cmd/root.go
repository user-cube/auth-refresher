package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "auth-refresher",
	Short: "A CLI tool for managing Docker/ECR registry logins",
	Long: `Auth Refresher is a command-line tool designed to simplify the process of managing
Docker and ECR registry logins. It provides an intuitive interface for selecting registries
from a configuration file and handles login operations with support for AWS and Helm registries.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
