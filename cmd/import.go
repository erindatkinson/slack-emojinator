/*
Copyright © 2025 Erin Atkinson
*/
package cmd

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/erindatkinson/slack-emojinator/internal/slack"
	"github.com/erindatkinson/slack-emojinator/internal/utilities"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Add a collection of emoji to a given slack team",
	Run: func(cmd *cobra.Command, args []string) {
		team := viper.GetString("team")
		inputDir := cmd.Flag("directory").Value.String()
		dryRun := utilities.PflagToBool(cmd.Flag("dry-run").Value)
		logger := utilities.NewLogger(
			cmd.Flag("log-level").Value.String(),
			"team", team, "dir", inputDir)
		client := slack.NewSlackClient(
			team,
			viper.GetString("token"),
			viper.GetString("cookie"))

		files, err := os.ReadDir(inputDir)
		if err != nil {
			slog.Error("error reading files")
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

		if !dryRun {
			for _, file := range filteredFiles {
				splits := strings.Split(file.Name(), ".")
				if err := client.ImportEmoji(splits[0], filepath.Join(inputDir, file.Name())); err != nil {
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
	importCmd.Flags().StringP("directory", "d", "./import/", "the directory to import from")
	importCmd.Flags().String("log-level", "info", "enable debug logging")
	importCmd.Flags().Bool("dry-run", false, "do a dry run")
}
