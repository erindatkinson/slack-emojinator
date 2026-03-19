package cache

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

func ListDownloadedEmojis(emojiDir string) (emojis []EmojiItem, err error) {
	emojis = make([]EmojiItem, 0)
	err = filepath.WalkDir(emojiDir, func(fPath string, d fs.DirEntry, err error) error {
		if fPath == emojiDir {
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
		emoji := EmojiItem{
			Name:     strings.Split(d.Name(), ".")[0],
			Filename: d.Name(),
			Dir:      path.Dir(fPath),
		}

		emojis = append(emojis, emoji)
		return nil
	})

	slices.SortFunc(emojis, func(a EmojiItem, b EmojiItem) int {
		return strings.Compare(a.Name, b.Name)
	})

	return
}

func PaginateEmojiList(list []EmojiItem, docsDir string) []*EmojiPage {
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
			page.PrevPage = path.Join("/", docsDir, pages[i-1].Name+".md")
		}

		if i < len(pages)-1 {
			page.NextPage = path.Join("/", docsDir, pages[i+1].Name+".md")
		}
	}
	return pages
}
