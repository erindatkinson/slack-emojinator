package cmd

import (
	"github.com/spf13/cobra"
)

var version = "dev"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "prints the version",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
