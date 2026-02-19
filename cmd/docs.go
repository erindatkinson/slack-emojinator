package cmd

import (
	"os"
	"path"
	"text/template"

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
		emojis, err := templates.ListEmojiFiles(inputRoot, namespace)
		if err != nil {
			logger.Error("unable to get emojis", "error", err)
			return
		}

		indexTpl, err := template.New("index").Parse(string(templates.MustAsset("templates/doc_index.md.gotmpl")))
		if err != nil {
			logger.Error("error loading template", "error", err)
			return
		}

		tpl, err := template.New("docs").Parse(string(templates.MustAsset("templates/doc_page.md.gotmpl")))
		if err != nil {
			logger.Error("error loading template", "error", err)
			return
		}

		logger.Info("paginating emojis")
		pages := templates.PaginateEmojiList(emojis, namespace)
		doc := templates.Docs{Namespace: namespace, Pages: pages}
		os.MkdirAll(path.Join("docs/", namespace), 0700)
		indexFp, err := os.Create(path.Join("docs/", namespace, "index.md"))
		if err = indexTpl.Execute(indexFp, &doc); err != nil {
			logger.Error("error writing index", "error", err)
			indexFp.Close()
			return
		}
		indexFp.Close()

		for _, page := range pages {
			fp, err := os.Create(path.Join("docs/", namespace, page.Name+".md"))
			if err != nil {
				logger.Error("couldn't make file", "error", err)
				return
			}

			if err = tpl.Execute(fp, *page); err != nil {
				logger.Error("error writing file", "error", err)
				fp.Close()
				return
			}
			fp.Close()
		}

	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
	docsCmd.Flags().StringP("output-root", "o", "./docs/", "the root directory to output a namespace to")
	docsCmd.Flags().StringP("input-root", "i", "./emojis/", "the root directory to look for a namespace in")
}
