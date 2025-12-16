/*
Copyright © 2025 Erin Atkinson
*/
package cmd

import (
	"os"

	"github.com/erindatkinson/slack-emojinator/internal/utilities"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "slack-emojinator",
	Short: "A tool to bulk import and export slack emojis",
	// Run: func(cmd *cobra.Command, args []string) {

	// 	client := slack.NewSlackClient(
	// 		viper.GetString("team"),
	// 		viper.GetString("token"),
	// 		viper.GetString("cookie"),
	// 	)

	// 	emojis, err := client.ListEmoji()
	// 	if err != nil {
	// 		slog.Error("error getting emojis", "error", err)
	// 	}

	// 	slog.Info("emojis", "count", len(emojis))
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("slack")

	for _, env := range utilities.Envs {
		viper.BindEnv(env)
	}

	viper.SetDefault("concurrency", "1")
	viper.AutomaticEnv() // read in environment variables that match

}
