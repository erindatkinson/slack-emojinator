package cmd

import (
	"fmt"
	"slices"
	"strings"

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
	Use:   "docs [namespace]",
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
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Println(cmd.UsageString())
			return
		}

		namespace := args[0]
		logger := utilities.NewLogger("info", "namespace", namespace)
		// outputRoot := cmd.Flag("output-root").Value.String()
		inputRoot := cmd.Flag("input-root").Value.String()

		logger.Info("getting emoji list")
		emojis, err := templates.ListEmojiFiles(inputRoot)
		if err != nil {
			logger.Error("unable to get emojis", "error", err)
			return
		}

		logger.Info("loading template")
		tpl, err := templates.LoadTemplate("templates/doc_page.md.jinja2")
		if err != nil {
			logger.Error("unable to load stored template", "error", err)
			return
		}

		logger.Info("paginating emojis")
		pages := templates.PaginateEmojiList(emojis)
		keys := make([]string, 0, len(pages))
		for k := range pages {
			keys = append(keys, k)
		}

		logger.Info("sorting keys", "pages", len(keys))
		slices.Sort(keys)

		logger.Info("building template data")
		for i := 0; i < len(keys); i++ {
			data := make(map[string]any)
			key := keys[i]
			keySplits := strings.Split(key, "-")
			data["page_count"] = keySplits[2]
			if i > 0 {
				data["prev_page"] = keys[i-1]
			} else if i < len(keys)-1 {
				data["next_page"] = keys[i+1]
			}

			page := pages[key]
			emojiData := make(map[string]string)
			for _, emoji := range page {
				emojiData["name"] = emoji.Name
				emojiData["path"] = emoji.DocPath
			}
			data["emoji"] = emojiData
			rendered, err := templates.RenderWithData(*tpl, data)
			if err != nil {
				logger.Error("unable to render data", "error", err)
				return
			}
			fmt.Println(rendered)
			if i > 0 {
				return
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
	docsCmd.Flags().StringP("output-root", "o", "./docs/", "the root directory to output a namespace to")
	docsCmd.Flags().StringP("input-root", "i", "./emojis/", "the root directory to look for a namespace in")
}
