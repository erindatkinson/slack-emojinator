package cmd

import (
	"github.com/erindatkinson/slack-emojinator/internal/cache"
	"github.com/erindatkinson/slack-emojinator/internal/templates"
	"github.com/erindatkinson/slack-emojinator/internal/utilities"
	"github.com/spf13/cobra"
)

type emojiFile struct {
	Name string
	Path string
}

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:   "docs [-d ./emojis]",
	Short: "Generate the docs for a namespace of emojis",
	Long: `This command assumes an archive structure like so:

	./emojis/namespace/

	Running 'slack-emojinator docs namespace1' should build a docs directory like so:
	./docs/namespace/

	`,
	Run: func(cmd *cobra.Command, args []string) {
		dir := cmd.Flag("dir").Value.String()

		logger := utilities.NewLogger("info")
		// outputRoot := cmd.Flag("output-root").Value.String()

		emojis, err := cache.ListDownloadedEmojis(dir)
		if err != nil {
			logger.Error("unable to get emojis", "error", err)
			return
		}

		pages := cache.PaginateEmojiList(emojis)
		if err := templates.WriteIndex(emojis[0].DocDir, pages); err != nil {
			logger.Error("error writing index", "error", err)
			return
		}

		if err := templates.WritePages(emojis[0].DocDir, pages); err != nil {
			logger.Error("error writing pages", "error", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
	docsCmd.Flags().StringP("dir", "d", "./emojis/", "the root directory to look for emojis in")
}
