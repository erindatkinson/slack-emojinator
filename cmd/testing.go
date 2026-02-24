/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"

	"github.com/erindatkinson/slack-emojinator/internal/cache"
	"github.com/spf13/cobra"
)

// testingCmd represents the testing command
var testingCmd = &cobra.Command{
	Use:   "testing",
	Short: "A brief description of your command",

	RunE: func(cmd *cobra.Command, args []string) error {
		emojis, err := cache.ListDownloadedEmojis(args[0])
		if err != nil {
			return err
		}
		for i, emoji := range emojis {
			if i > 100 {
				return nil
			}
			slog.Info("emoji", "name", emoji.Name, "dir", emoji.Dir, "file", emoji.Filename, "docs", emoji.DocDir)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(testingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
