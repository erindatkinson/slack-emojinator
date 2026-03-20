/*
Copyright © 2025 Erin Atkinson
*/
package cmd

import (
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/erindatkinson/emoji-archiver/internal/slack"
	"github.com/erindatkinson/emoji-archiver/internal/utilities"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var importDryRun bool

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Add a collection of emoji to a given slack team",
	Run: func(cmd *cobra.Command, args []string) {
		logger := utilities.ContextLogger(cmd.Context())
		if browser == "" || profile == "" || subdomain == "" {
			slog.Error("error reading configs from env, config, or flags")
			return
		}

		importDir := path.Join(directory, subdomain)
		client, err := slack.NewSlackClient(cmd.Context(), browser, profile, subdomain)
		if err != nil {
			logger.Error("error creating slack client", "error", err)
			return
		}

		files, err := os.ReadDir(importDir)
		if err != nil {
			logger.Error("error reading files", "error", err)
			return
		}
		logger.Info("found emojis to import", "count", len(files))

		emojis, err := client.ListEmoji()
		if err != nil {
			logger.Error("error listing emojis", "err", err)
			return
		}
		logger.Info("found existing emojis", "count", len(emojis))

		filteredFiles := lo.Filter(files, func(item os.DirEntry, index int) bool {
			if item.Name() == ".DS_Store" {
				return false
			}

			splits := strings.Split(item.Name(), ".")
			_, ok := lo.Find(emojis, func(emoji slack.Emoji) bool {
				return splits[0] == emoji.Name
			})
			if ok {
				logger.Debug("filtering out file", "emoji", splits[0])
			}
			return !ok
		})
		logger.Info("emojis to upload", "count", len(filteredFiles))

		if !importDryRun {
			for _, file := range filteredFiles {
				splits := strings.Split(file.Name(), ".")
				if err := client.ImportEmoji(splits[0], filepath.Join(importDir, file.Name())); err != nil {
					logger.Error("error importing", "error", err)
					return
				}
			}
		} else {
			logger.Info("skipping import due to dry-run")
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().BoolVar(&importDryRun, "dry-run", false, "do a dry run")
}
