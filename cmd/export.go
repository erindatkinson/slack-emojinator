/*
Copyright © 2025 Erin Atkinson
*/
package cmd

import (
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
		outputDir := cmd.Flag("directory").Value.String()
		team := viper.GetString("team")
		logger := utilities.NewLogger(
			cmd.Flag("log-level").Value.String(),
			"team", team, "dir", outputDir)

		if err := utilities.CheckEnvs(); err != nil {
			logger.Error(err.Error())
			return
		}
		logger.Info("creating export directory")
		os.MkdirAll(outputDir, 0755)
		client := slack.NewSlackClient(
			team,
			viper.GetString("token"),
			viper.GetString("cookie"),
		)
		logger.Debug("client setup complete")

		logger.Info("retrieving list of current emoji")
		currentEmoji, err := client.ListEmoji()
		if err != nil {
			logger.Error("error retrieving current emoji list", "error", err)
			return
		}

		wp := workerpool.New(2)

		for _, emoji := range currentEmoji {
			loopLog := logger.With("name", emoji.Name)

			request := emoji
			wp.Submit(func() {
				loopLog.Debug("exporting emoji")
				if err := client.ExportEmoji(request, outputDir); err != nil {
					loopLog.Error("error exporting", "error", err)
				}
			})
		}

		wp.StopWait()
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringP("directory", "d", "./export/", "the directory to use to export")
	exportCmd.Flags().String("log-level", "info", "enable debug logging")
}
