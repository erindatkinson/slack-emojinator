package templates

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
)

func ParseName(emoji string) {

}

func ListEmojiFiles(directory string) (emojis []EmojiItem, err error) {
	emojis = make([]EmojiItem, 0)
	err = filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if path == directory {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		emojis = append(emojis, EmojiItem{
			Name:    d.Name(),
			DocPath: path,
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

func PaginateEmojiList(list []EmojiItem) map[string][]EmojiItem {
	pages := make(map[string][]EmojiItem)
	count := 0
	for i := 0; i < len(list); i = i + 100 {
		char := string(list[i].Name[0])
		page := list[i : i+100]
		pages[fmt.Sprintf("page-%s-%06d", char, count)] = page
		count++
	}

	return pages
}
