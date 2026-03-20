/*
Copyright © 2025 Erin Atkinson
*/
package cmd

import (
	"os"

	"github.com/erindatkinson/emoji-archiver/internal/utilities"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var browser, profile, subdomain, channel, directory, logLevel string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "emoji-archiver",
	Short: "A tool to bulk import and export slack emojis",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger := utilities.NewLogger(logLevel,
			"subdomain", subdomain,
			"root-directory", directory,
		)
		cmd.SetContext(utilities.ToContext(cmd.Context(), logger))
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	initConfig()
	rootCmd.PersistentFlags().StringVarP(&subdomain, "subdomain", "s", utilities.ConfigOrEnv("slack", "subdomain"), "what subdomain to pull a slack token for")
	rootCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "./emojis/", "base directory to use")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "log-level to use")
	rootCmd.PersistentFlags().StringVarP(&browser, "browser", "b", utilities.ConfigOrEnv("slack", "browser"), "browser to look for token")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", utilities.ConfigOrEnv("slack", "profile"), "profile to look for token")
	// releaseNotes channel is entered here since it has to be post initConfig for ConfigOrEnv to work, but calling
	// initConfig multiple times causes a panic
	releaseNotesCmd.Flags().StringVarP(&channel, "channel", "c", utilities.ConfigOrEnv("slack", "channel"), "channel to post to")

}

// initConfig reads in config file
func initConfig() {
	viper.SetConfigName(".config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/emoji-archiver/")
	viper.AddConfigPath("$HOME/.emojinator")
	viper.AddConfigPath(".")
	viper.ReadInConfig()
}
