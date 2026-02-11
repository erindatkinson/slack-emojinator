/*
Copyright © 2025 Erin Atkinson
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/erindatkinson/slack-emojinator/internal/slack"
	"github.com/erindatkinson/slack-emojinator/internal/templates"
	"github.com/erindatkinson/slack-emojinator/internal/utilities"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Rank struct {
	Name  string
	Count int
}

var releaseNotesWindowStart time.Time
var releaseNotesWindowEnd time.Time

// releaseNotesCmd represents the releaseNotes command
var releaseNotesCmd = &cobra.Command{
	Use:   "release-notes",
	Short: "Generate and publish release notes",

	Run: func(cmd *cobra.Command, args []string) {
		channel := viper.GetString("release_channel")
		team := viper.GetString("team")
		logger := utilities.NewLogger("info", "team", team)

		client := slack.NewSlackClient(
			team,
			viper.GetString("token"),
			viper.GetString("cookie"))

		emojis, err := client.ListEmoji()
		if err != nil {
			logger.Error("unable to retrieve emoji list")
			return
		}

		durationEmojis := lo.Filter(emojis, func(item slack.Emoji, index int) bool {
			if item.Created > releaseNotesWindowStart.Unix() {
				if item.Created < releaseNotesWindowEnd.Unix() {
					return true
				}
			}
			return false
		})

		ranks, err := templates.RenderRanks(durationEmojis)
		if err != nil {
			logger.Error("unable to render the rank list", "error", err)
			return
		}

		emojiMessaes := templates.BuildEmojiLists(durationEmojis)

		// message for start of thread
		headerTpl, err := templates.LoadTemplate("templates/header.md.jinja2")
		if err != nil {
			logger.Error("unable to load header template", "error", err)
			return
		}
		data := map[string]any{
			"start": releaseNotesWindowStart.Format(time.RFC822),
			"end":   releaseNotesWindowEnd.Format(time.RFC822),
		}
		header, err := templates.RenderWithData(*headerTpl, data)
		if err != nil {
			logger.Error("unable to render header template", "error", err)
			return
		}

		logger.Info("sending chanel header message")
		resp, err := client.PostMessage(header, channel, "", false)
		if err != nil {
			logger.Error("unable to post message", "error", err)
			return
		}

		logger.Info("sending ranks")
		thread := resp["ts"].(string)
		if _, err := client.PostMessage(ranks, channel, thread, false); err != nil {
			logger.Error("unable to post ranks to thread", "error", err)
			return
		}

		started := false
		for i, message := range emojiMessaes {
			logger.Info("sending page of new emojis", "page", i)
			var markdown string
			if !started {
				markdown = "## New Emojis\n" + message
				started = true
			} else {
				markdown = message
			}

			resp, err := client.PostMessage(markdown, channel, thread, false)
			if err != nil {
				logger.Error("unable to post followup message", "error", err)
				return
			}

			if !resp["ok"].(bool) {
				logger.Info("debug", "page", i, "resp", resp, "len", len(markdown))
				fmt.Println(markdown)
			}

		}

	},
}

func init() {
	rootCmd.AddCommand(releaseNotesCmd)
	now := time.Now()
	releaseNotesCmd.Flags().TimeVar(&releaseNotesWindowStart, "start", now.Add(-14*24*time.Hour), []string{time.RFC822}, "start time")
	releaseNotesCmd.Flags().TimeVar(&releaseNotesWindowEnd, "end", now, []string{time.RFC822}, "end time")
}
