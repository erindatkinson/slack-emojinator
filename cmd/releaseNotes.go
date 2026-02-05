/*
Copyright © 2025 Erin Atkinson
*/
package cmd

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/erindatkinson/slack-emojinator/internal/utilities"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/spf13/cobra"
)

var releaseNotesWindowStart time.Time
var releaseNotesWindowEnd time.Time

// releaseNotesCmd represents the releaseNotes command
var releaseNotesCmd = &cobra.Command{
	Use:   "release-notes",
	Short: "A brief description of your command",

	Run: func(cmd *cobra.Command, args []string) {
		headerTpl, err := gonja.FromString(utilities.MustAssetString("templates/header.md.jinja2"))
		if err != nil {
			slog.Error("unable to read template", "error", err)
			return
		}

		tpl, err := gonja.FromString(utilities.MustAssetString("templates/release_notes.md.jinja2"))
		if err != nil {
			slog.Error("unable to read template", "error", err)
			return
		}

		data := exec.EmptyContext()
		data.Set("start", releaseNotesWindowStart.Format(time.DateOnly))
		data.Set("end", releaseNotesWindowEnd.Format(time.DateOnly))

		render, err := headerTpl.ExecuteToString(data)
		if err != nil {
			slog.Error("error rendering template", "error", err)
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
