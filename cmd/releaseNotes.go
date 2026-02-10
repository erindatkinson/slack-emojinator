/*
Copyright © 2025 Erin Atkinson
*/
package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/bndr/gotabulate"
	"github.com/erindatkinson/slack-emojinator/internal/slack"
	"github.com/erindatkinson/slack-emojinator/internal/utilities"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
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

		header, body, err := renderTemplates(durationEmojis)
		if err != nil {
			logger.Error("unable to render templates", "error", err)
			return
		}

		// Slack has 12k message limit
		if len(body) > 12_000 {
			fmt.Println("Message over slack limits, posting to standard out instead")
			fmt.Println("----------------------------------------------------------")
			fmt.Println(header)
			fmt.Println(body)

		} else {
			resp1, err := client.PostMessage(header, channel, "", false)
			if err != nil {
				logger.Error("error posting message header", "error", err)
				return
			}

			_, err = client.PostMessage(body, channel, resp1["ts"].(string), false)
			if err != nil {
				logger.Error("error posting message body", "error", err)
				return
			}
		}
	},
}

func renderTemplates(emojis []slack.Emoji) (string, string, error) {
	ranks := buildRanks(emojis)
	tab := gotabulate.Create(ranks)
	tab.SetHeaders([]string{"Name", "Count"})
	tab.SetAlign("center")
	rankString := tab.Render("simple")

	headerTpl, err := gonja.FromString(utilities.MustAssetString("templates/header.md.jinja2"))
	if err != nil {
		return "", "", err
	}

	tpl, err := gonja.FromString(utilities.MustAssetString("templates/release_notes.md.jinja2"))
	if err != nil {
		return "", "", err
	}

	data := exec.EmptyContext()
	data.Set("start", releaseNotesWindowStart.Format(time.DateOnly))
	data.Set("end", releaseNotesWindowEnd.Format(time.DateOnly))

	renderHeader, err := headerTpl.ExecuteToString(data)
	if err != nil {
		return "", "", err
	}

	data = exec.EmptyContext()
	data.Set("emojis", buildEmojiString(emojis))
	data.Set("ranks", rankString)
	renderBody, err := tpl.ExecuteToString(data)
	if err != nil {
		return "", "", err
	}

	return renderHeader, renderBody, nil
}

func buildEmojiString(emojis []slack.Emoji) []string {
	out := make([]string, 0)
	for _, emoji := range emojis {
		out = append(out, fmt.Sprintf("* :%s: | `:%s:`", emoji.Name, emoji.Name))
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i] < out[j]
	})
	return out
}

func buildRanks(emojis []slack.Emoji) [][]interface{} {
	ranks := make(map[string]*Rank)
	for _, emoji := range emojis {
		if _, ok := ranks[emoji.UserDisplayName]; ok {
			ranks[emoji.UserDisplayName].Count = ranks[emoji.UserDisplayName].Count + 1
		} else {
			ranks[emoji.UserDisplayName] = &Rank{
				Name:  emoji.UserDisplayName,
				Count: 1,
			}
		}
	}

	var rankArray [][]interface{}
	for _, rank := range ranks {
		rankArray = append(rankArray, []interface{}{rank.Name, rank.Count})
	}

	sort.Slice(rankArray, func(i, j int) bool {
		first := rankArray[i][1].(int)
		second := rankArray[j][1].(int)
		return first > second
	})
	return rankArray
}

func init() {
	rootCmd.AddCommand(releaseNotesCmd)
	now := time.Now()
	releaseNotesCmd.Flags().TimeVar(&releaseNotesWindowStart, "start", now.Add(-14*24*time.Hour), []string{time.RFC822}, "start time")
	releaseNotesCmd.Flags().TimeVar(&releaseNotesWindowEnd, "end", now, []string{time.RFC822}, "end time")
}
