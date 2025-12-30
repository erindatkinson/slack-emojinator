/*
Copyright © 2025 Erin Atkinson
*/
package cmd

import (
	"os"
	"slices"
	"strconv"

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
		concurrency, _ := strconv.Atoi(cmd.Flag("concurrency").Value.String())

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
		cached, err := utilities.GetDownloadedEmojiList(outputDir)
		if err != nil {
			logger.Error("unable to get cached emojis", "error", err)
		}

		wp := workerpool.New(concurrency)

		for _, emoji := range currentEmoji {
			request := emoji
			wp.Submit(func() {
				loopLog := logger.With("name", request.Name)
				if slices.Contains(cached, request.Name) {
					loopLog.Debug("already downloaded, skipping")
					return
				}

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
	exportCmd.Flags().IntP("concurrency", "c", 2, "worker concurrency")
}
