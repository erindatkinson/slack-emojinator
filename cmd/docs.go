package cmd

import (
	"path"

	"github.com/erindatkinson/slack-emojinator/internal/cache"
	"github.com/erindatkinson/slack-emojinator/internal/templates"
	"github.com/erindatkinson/slack-emojinator/internal/utilities"
	"github.com/spf13/cobra"
)

type emojiFile struct {
	Name string
	Path string
}

var docsRootDir string

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate the docs for a namespace of emojis",
	Run: func(cmd *cobra.Command, args []string) {
		logger := utilities.ContextLogger(cmd.Context())
		emojiDir := path.Join(directory, subdomain)
		docsDir := path.Join(docsRootDir, subdomain)
		emojis, err := cache.ListDownloadedEmojis(emojiDir)
		if err != nil {
			logger.Error("unable to get emojis", "error", err)
			return
		}

		pages := cache.PaginateEmojiList(emojis, docsDir)
		if err := templates.WriteIndex(emojiDir, docsDir, pages); err != nil {
			logger.Error("error writing index", "error", err)
			return
		}

		if err := templates.WritePages(docsDir, pages); err != nil {
			logger.Error("error writing pages", "error", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
	docsCmd.Flags().StringVar(&docsRootDir, "docs-dir", "docs/", "the root directory to write docs into")
}
