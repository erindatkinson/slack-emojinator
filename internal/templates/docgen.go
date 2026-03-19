package templates

import (
	"os"
	"path"
	"text/template"

	"github.com/erindatkinson/slack-emojinator/internal/cache"
)

func WriteIndex(emojiDir, docsDir string, pages []*cache.EmojiPage) error {
	tpl, err := template.New("index").Parse(MustAssetString("templates/doc_index.md.gotmpl"))
	if err != nil {
		return err
	}
	doc := Docs{Namespace: path.Base(emojiDir), Pages: pages}

	os.RemoveAll(docsDir)
	os.MkdirAll(docsDir, 0700)
	fp, err := os.Create(path.Join(docsDir, "index.md"))
	if err != nil {
		return err
	}
	defer fp.Close()

	if err = tpl.Execute(fp, &doc); err != nil {
		return err
	}

	return nil
}

func WritePages(docsDir string, pages []*cache.EmojiPage) error {
	tpl, err := template.New("docs").Parse(MustAssetString("templates/doc_page.md.gotmpl"))
	if err != nil {
		return err
	}
	for _, page := range pages {
		fp, err := os.Create(path.Join(docsDir, page.Name+".md"))
		if err != nil {
			return err
		}
		defer fp.Close()

		if err = tpl.Execute(fp, *page); err != nil {
			return err
		}
	}
	return nil
}
