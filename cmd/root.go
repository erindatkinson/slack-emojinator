/*
Copyright © 2025 Erin Atkinson
*/
package cmd

import (
	"os"

	"github.com/erindatkinson/slack-emojinator/internal/utilities"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "slack-emojinator",
	Short: "A tool to bulk import and export slack emojis",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(utilities.InitConfig)
}
