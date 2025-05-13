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
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to a selected registry",
	RunE: func(cmd *cobra.Command, args []string) error {
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
			return err
		}
		defer file.Close()

		// Call LoginToRegistry without spinner
		err = auth.LoginToRegistry(ctx, configPath)
		if err != nil {
			ui.PrintError("Failed to login to registry", err, true)
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
