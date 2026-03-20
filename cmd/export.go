// /*
// Copyright © 2025 Erin Atkinson
// */
package cmd

import (
	"os"
	"path"
	"slices"

	"github.com/erindatkinson/emoji-archiver/internal/cache"
	"github.com/erindatkinson/emoji-archiver/internal/slack"
	"github.com/erindatkinson/emoji-archiver/internal/utilities"
	"github.com/gammazero/workerpool"
	"github.com/spf13/cobra"
)

var concurrency int

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Pull all emoji from a given slack team",
	Run: func(cmd *cobra.Command, args []string) {
		logger := utilities.ContextLogger(cmd.Context())
		if browser == "" || profile == "" || subdomain == "" {
			logger.Error("error reading configs from env, config, or flags")
			return
		}

		logger.Info("creating export directory")
		exportDir := path.Join(directory, subdomain)
		os.MkdirAll(exportDir, 0755)

		client, err := slack.NewSlackClient(cmd.Context(), browser, profile, subdomain)
		if err != nil {
			logger.Error("unable to create slack client", "error", err)
			return
		}
		logger.Debug("client setup complete")
		logger.Info("retrieving list of current emoji")
		currentEmoji, err := client.ListEmoji()
		if err != nil {
			logger.Error("error retrieving current emoji list", "error", err)
			return
		}
		logger.Info("listing downloaded emojis from filesystem")
		cached, err := cache.ListDownloadedEmojis(exportDir)
		if err != nil {
			logger.Error("unable to get cached emojis", "error", err)
			return
		}

		logger.Info("exporting emojis")
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
				if err := client.ExportEmoji(request, exportDir); err != nil {
					loopLog.Error("error exporting", "error", err)
				}

			})
		}

		wp.StopWait()
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().IntVar(&concurrency, "concurrency", 1, "concurrency to use to download")
}
