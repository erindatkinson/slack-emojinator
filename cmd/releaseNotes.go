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
)

var (
	releaseNotesWindowStart time.Time
	releaseNotesWindowEnd   time.Time
	releaseNotesDryRun      bool
)

// releaseNotesCmd represents the releaseNotes command
var releaseNotesCmd = &cobra.Command{
	Use:   "release-notes",
	Short: "Generate and publish release notes",

	Run: func(cmd *cobra.Command, args []string) {
		logger := utilities.ContextLogger(cmd.Context())
		if browser == "" || profile == "" || subdomain == "" {
			logger.Error("error reading configs from env, config, or flags")
			return
		}

		client, err := slack.NewSlackClient(cmd.Context(), browser, profile, subdomain)
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

		emojiMessages := templates.BuildEmojiLists(durationEmojis)

		// message for start of thread
		header, err := templates.RenderHeader(releaseNotesWindowStart, releaseNotesWindowEnd)
		if err != nil {
			logger.Error("unable to render header", "error", err)
			return
		}
		if !releaseNotesDryRun {
			logger.Info("sending chanel header message")
			resp, err := client.PostMessage(channel, header, nil)
			if err != nil {
				logger.Error("unable to post message", "error", err)
				return
			}
			if errStr, ok := resp["error"]; ok {
				logger.Error("error posting message", "error", errStr)
				return
			}

			var thread string = ""
			if ts, ok := resp["ts"]; ok {
				thread = ts.(string)
			}

			logger.Info("sending ranks")
			_, err = client.PostMessage(channel, ranks, &thread)
			if err != nil {
				logger.Error("unable to post ranks to thread", "error", err)
				return
			}

			started := false
			for i, message := range emojiMessages {
				logger.Info("sending page of new emojis", "page", i)
				var markdown string
				if !started {
					markdown = "### New Emojis\n" + message
					started = true
				} else {
					markdown = message
				}

				_, err = client.PostMessage(channel, markdown, &thread)
				if err != nil {
					logger.Error("unable to post followup message", "error", err)
					return
				}
			}
		} else {
			fmt.Println(header)
			fmt.Println(ranks)
			for i, message := range emojiMessages {
				if i == 0 {
					fmt.Printf("### New Emojis\n\n")
				}
				fmt.Println(message)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(releaseNotesCmd)
	now := time.Now()
	releaseNotesCmd.Flags().TimeVar(&releaseNotesWindowStart, "start", now.Add(-14*24*time.Hour), []string{time.RFC822}, "start time")
	releaseNotesCmd.Flags().TimeVar(&releaseNotesWindowEnd, "end", now, []string{time.RFC822}, "end time")
	releaseNotesCmd.Flags().BoolVar(&releaseNotesDryRun, "dry-run", false, "don't post if set")

	// channel flag is set in /cmd/root.go so that it can have the initConfig() call, don't re-add it here.
	// releaseNotesCmd.Flags().StringVarP(&channel, "channel", "c", utilities.ConfigOrEnv("slack", "channel"), "channel to post to")

}
