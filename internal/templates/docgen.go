package templates

import (
	"fmt"
	"io/fs"
	"log/slog"
	"path"
	"path/filepath"
	"slices"
)

func ParseName(emoji string) {

}

func ListEmojiFiles(directory string, namespace string) (emojis []EmojiItem, err error) {
	emojis = make([]EmojiItem, 0)
	err = filepath.WalkDir(path.Join(directory, namespace), func(fPath string, d fs.DirEntry, err error) error {
		if fPath == directory {
			return nil
		}
		if d.Name() == ".DS_Store" {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		if d.Name() == "" {
			slog.Info("wtf", "path", fPath, "error", err)
		}
		emojis = append(emojis, EmojiItem{
			Name:    d.Name(),
			DocPath: path.Join("/emojis/", namespace, path.Base(fPath)),
		})
		return nil
	})

	slices.SortFunc(emojis, func(a EmojiItem, b EmojiItem) int {
		if a.Name < b.Name {
			return -1
		} else if a.Name > b.Name {
			return 1
		} else {
			return 0
		}
	})
	return
}

func PaginateEmojiList(list []EmojiItem, namespace string) []*EmojiPage {
	pages := []*EmojiPage{}
	count := 0
	for i := 0; i < len(list); i = i + 100 {
		char := string(list[i].Name[0])
		var emojis []EmojiItem
		if i+100 > len(list) {
			emojis = list[i : len(list)-1]
		} else {
			emojis = list[i : i+100]
		}
		name := fmt.Sprintf("page-%s-%06d", char, count)
		page := EmojiPage{
			Name:     name,
			Count:    count,
			Emojis:   emojis,
			PrevPage: "",
			NextPage: "",
		}
		pages = append(pages, &page)
		count++
	}

	for i, page := range pages {
		if i > 0 {
			page.PrevPage = path.Join("/docs/", namespace, pages[i-1].Name+".md")
		}

		if i < len(pages)-1 {
			page.NextPage = path.Join("/docs/", namespace, pages[i+1].Name+".md")
		}
	}
	return pages
}
