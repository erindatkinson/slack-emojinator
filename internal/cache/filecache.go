package cache

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

func ListDownloadedEmojis(directory string) (emojis []EmojiItem, err error) {
	emojis = make([]EmojiItem, 0)
	err = filepath.WalkDir(directory, func(fPath string, d fs.DirEntry, err error) error {
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
		checkedDir := path.Dir(fPath)
		docDir := strings.ReplaceAll(checkedDir, "emojis", "docs")
		if docDir == checkedDir {
			docDir = "./docs"
		}
		emoji := EmojiItem{
			Name:     strings.Split(d.Name(), ".")[0],
			Filename: d.Name(),
			Dir:      checkedDir,
			DocDir:   docDir,
		}

		emojis = append(emojis, emoji)
		return nil
	})

	slices.SortFunc(emojis, func(a EmojiItem, b EmojiItem) int {
		return strings.Compare(a.Name, b.Name)
	})

	return
}

func PaginateEmojiList(list []EmojiItem) []*EmojiPage {
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
			page.PrevPage = path.Join("/docs/", pages[i-1].Name+".md")
		}

		if i < len(pages)-1 {
			page.NextPage = path.Join("/docs/", pages[i+1].Name+".md")
		}
	}
	return pages
}
