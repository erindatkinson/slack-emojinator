package utilities

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func GetDownloadedEmojiList(directory string) ([]string, error) {
	emojis := []string{}
	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		splits := strings.Split(d.Name(), ".")
		emojis = append(emojis, splits[0])
		return nil
	})
	if err != nil {
		return []string{}, err
	}

	return emojis, nil
}
