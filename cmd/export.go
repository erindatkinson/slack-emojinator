/*
Copyright © 2025 Erin Atkinson
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/erindatkinson/slack-emojinator/internal/slack"
	"github.com/erindatkinson/slack-emojinator/internal/utilities"
	"github.com/gammazero/workerpool"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Pull all emoji from a given slack team",
	Run: func(cmd *cobra.Command, args []string) {
		if err := utilities.CheckEnvs(); err != nil {
			slog.Error(err.Error())
			return
		}

		outputDir := cmd.Flag("directory").Value.String()
		os.MkdirAll(outputDir, 0755)
		client := slack.NewSlackClient(
			viper.GetString("team"),
			viper.GetString("token"),
			viper.GetString("cookie"),
		)
		slog.Debug("Client setup", "team", viper.GetString("team"), "output", outputDir)

		currentEmoji, err := client.ListEmoji()
		if err != nil {
			slog.Error("error listing emoji", "error", err)
			return
		}

		if err = client.ExportEmoji(currentEmoji[0], outputDir); err != nil {
			slog.Error("error exporting emoji", "name", currentEmoji[0].Name, "error", err)
			return
		}

		wp := workerpool.New(2)

		for i, emoji := range currentEmoji {
			slog.Debug("submitting", "name", emoji.Name, "i", i)
			request := emoji
			wp.Submit(func() {
				if err := client.ExportEmoji(request, outputDir); err != nil {
					slog.Error("error exporting", "error", err)
				}
			})
		}

		wp.StopWait()
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringP("directory", "d", "./export/", "the directory to use to export")
}
