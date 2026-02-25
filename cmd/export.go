// /*
// Copyright © 2025 Erin Atkinson
// */
package cmd

import (
	"encoding/json"
	"os"
	"slices"
	"strconv"

	"github.com/erindatkinson/slack-emojinator/internal/cache"
	"github.com/erindatkinson/slack-emojinator/internal/slack"
	"github.com/erindatkinson/slack-emojinator/internal/utilities"
	"github.com/gammazero/workerpool"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func EmojisFromFile(jsonPath string) ([]slack.Emoji, error) {
	jsonBytes, err := os.ReadFile(jsonPath)
	if err != nil {
		return []slack.Emoji{}, nil
	}
	var jsonEmojis slack.EmojiJsonFile
	if err = json.Unmarshal(jsonBytes, &jsonEmojis); err != nil {
		return []slack.Emoji{}, nil
	}

	currentEmoji := make([]slack.Emoji, 0)
	for e, url := range jsonEmojis.Emoji {
		currentEmoji = append(currentEmoji, slack.Emoji{
			Name: e,
			URL:  url,
		})
	}
	return currentEmoji, nil
}

func EmojisFromApi(client *slack.Client) ([]slack.Emoji, error) {
	currentEmoji, err := client.ListEmoji()
	if err != nil {
		return []slack.Emoji{}, err
	}
	return currentEmoji, nil
}

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Pull all emoji from a given slack team",
	Run: func(cmd *cobra.Command, args []string) {
		var client *slack.Client
		var currentEmoji []slack.Emoji
		var err error
		outputDir := cmd.Flag("directory").Value.String()
		team := viper.GetString("team")
		concurrency, _ := strconv.Atoi(cmd.Flag("concurrency").Value.String())
		logger := utilities.NewLogger(
			cmd.Flag("log-level").Value.String(),
			"team", team, "dir", outputDir)

		jsonPath := cmd.Flag("json").Value.String()
		if jsonPath != "" {
			if currentEmoji, err = EmojisFromFile(jsonPath); err != nil {
				logger.Error("unable to load emojis from file", "error", err)
				return
			}
		} else {
			if err := utilities.CheckEnvs(); err != nil {
				logger.Error("error getting environment vars", "error", err)
				return
			}
			client = slack.NewSlackClient(
				team,
				viper.GetString("token"),
				viper.GetString("cookie"),
			)
			if currentEmoji, err = EmojisFromApi(client); err != nil {
				logger.Error("unable to load emojis from slack", "error", err)
			}

		}
		os.MkdirAll(outputDir, 0755)
		cached, err := cache.ListDownloadedEmojis(outputDir)
		if err != nil {
			logger.Error("unable to get cached emojis", "error", err)
		}

		wp := workerpool.New(concurrency)

		for _, emoji := range currentEmoji {
			request := emoji
			wp.Submit(func() {
				loopLog := logger.With("name", request.Name)
				if slices.ContainsFunc(cached, func(e cache.EmojiItem) bool {
					return e.Name == request.Name
				}) {
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
	exportCmd.Flags().StringP("json", "j", "", "json file to use instead of the api")
	exportCmd.Flags().String("log-level", "info", "enable debug logging")
	exportCmd.Flags().IntP("concurrency", "c", 2, "worker concurrency")
}
