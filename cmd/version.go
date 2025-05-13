package cmd

import (
	"runtime"

	"github.com/spf13/cobra"
	"github.com/user-cube/gocli-template/pkg/ui"
)

// Version information
var (
	BuildDate = "unknown"
	GitCommit = "unknown"
	Version   = "v1.0.0"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information of $PROJECT_NAME",
	Long: `Display the version, build date, and git commit of your $PROJECT_NAME installation.

Examples:
  # Show version information
  $PROJECT_NAME version`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintInfo("$PROJECT_NAME", Version)
		ui.PrintInfo("Git Commit", GitCommit)
		ui.PrintInfo("Built", BuildDate)
		ui.PrintInfo("Platform", runtime.GOOS+"/"+runtime.GOARCH)
		ui.PrintInfo("Go Version", runtime.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
