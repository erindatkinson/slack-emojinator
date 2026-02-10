/*
Copyright © 2025 Erin Atkinson
*/
package cmd

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/erindatkinson/slack-emojinator/internal/slack"
	"github.com/erindatkinson/slack-emojinator/internal/utilities"
	// "github.com/markkurossi/tabulate"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Rank struct {
	Name string
	Count int
}

var releaseNotesWindowStart time.Time
var releaseNotesWindowEnd time.Time

// releaseNotesCmd represents the releaseNotes command
var releaseNotesCmd = &cobra.Command{
	Use:   "release-notes",
	Short: "Generate and publish release notes",

	Run: func(cmd *cobra.Command, args []string) {
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

		durationEmojis := lo.Filter(emojis, func(item slack.Emoji, index int) {
			if item.Created > releaseNotesWindowStart
		}
		// ranks := []Rank{}
		// tab := tabulate.New(tabulate.ASCII)
		// tab.Header("Name").SetAlign(tabulate.ML)
		// tab.Header("Count")
		// err := tabulate.Reflect(tab, 0, nil, &ranks)
		// tab.Print(os.Stdout)

		logger.Info("first", "emoji", emojis[0])
		ranks := make(map[string]int)
		for _, emoji := range emojis {
			
			if _, ok := ranks[emoji.UserDisplayName] {

			}
		}
		// emoji.UserDisplayName

		headerTpl, err := gonja.FromString(utilities.MustAssetString("templates/header.md.jinja2"))
		if err != nil {
			logger.Error("unable to read template", "error", err)
			return
		}

		tpl, err := gonja.FromString(utilities.MustAssetString("templates/release_notes.md.jinja2"))
		if err != nil {
			logger.Error("unable to read template", "error", err)
			return
		}

		data := exec.EmptyContext()
		data.Set("start", releaseNotesWindowStart.Format(time.DateOnly))
		data.Set("end", releaseNotesWindowEnd.Format(time.DateOnly))

		render, err := headerTpl.ExecuteToString(data)
		if err != nil {
			logger.Error("error rendering template", "error", err)
			return
		}

		fmt.Println(render)

		data = exec.EmptyContext()
		data.Set("emojis", []string{})
		renderBody, err := tpl.ExecuteToString(data)
		if err != nil {
			slog.Error("error rendering template", "error", err)
			return
		}

		fmt.Println(renderBody)
	},
}

func init() {
	rootCmd.AddCommand(releaseNotesCmd)
	now := time.Now()
	releaseNotesCmd.Flags().TimeVar(&releaseNotesWindowStart, "start", now.Add(-14*24*time.Hour), []string{time.RFC822}, "start time")
	releaseNotesCmd.Flags().TimeVar(&releaseNotesWindowEnd, "end", now, []string{time.RFC822}, "end time")
}
