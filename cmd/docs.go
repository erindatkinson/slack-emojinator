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
	Use:   "docs [-d ./emojis] namespace",
	Short: "Generate the docs for a namespace of emojis",
	Long: `This command assumes an archive structure like so:

	./emojis/
	├── namespace1/
	├── namespace2/
	├── namespace3/
	├── namespace4/
	└── namespace5/

	Running 'slack-emojinator docs namespace1' should build a docs directory like so:
	./docs/
	└── namespace1/

	If you need to generate for something outside the local path, running 'slack-emojinator docs -d ../archive/' namespace1
	you would need the file structure to look like:
	../
	├── cwd/
	│	└── .
	└── archive/
		└── namespace1/

	and it should build the docs like so:
	../
	├── cwd/
	│	└── .
	├── archive/
	│	└── namespace1/
	└── docs/
		└── namespace1/	

	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Println(cmd.UsageString())
			return
		}

		namespace := args[0]
		logger := utilities.NewLogger("info", "namespace", namespace)
		// outputRoot := cmd.Flag("output-root").Value.String()

		emojis, err := cache.ListDownloadedEmojis(cmd.Flag("dir").Value.String())
		if err != nil {
			logger.Error("unable to get emojis", "error", err)
			return
		}

		pages := cache.PaginateEmojiList(emojis, namespace)
		if err := templates.WriteIndex(namespace, pages); err != nil {
			logger.Error("error writing index", "error", err)
			return
		}

		if err := templates.WritePages(namespace, pages); err != nil {
			logger.Error("error writing pages", "error", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
	docsCmd.Flags().StringP("dir", "d", "./emojis/", "the root directory to look for emojis in")
}
