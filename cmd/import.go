/*
Copyright © 2025 Erin Atkinson
*/
package cmd

import (
	"log/slog"
	"path/filepath"

	"github.com/erindatkinson/slack-emojinator/internal/slack"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Add a collection of emoji to a given slack team",
	Run: func(cmd *cobra.Command, args []string) {
		inputDir := cmd.Flag("directory").Value.String()
		client := slack.NewSlackClient(
			viper.GetString("team"),
			viper.GetString("token"),
			viper.GetString("cookie"))

		if err := client.ImportEmoji("1password_test", filepath.Join(inputDir, "1password.png")); err != nil {
			slog.Error("error importing", "error", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringP("directory", "d", "./import/", "the directory to import from")
}
