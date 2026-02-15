package templates

import (
	"fmt"
	"os"
	"path"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vektra/neko"
)

func genTestFiles(t *testing.T) (dir string) {
	tmpDir, err := os.MkdirTemp(".", "tmp_*")
	require.Nil(t, err)
	for i := 0; i < 1000; i++ {
		fp, err := os.CreateTemp(tmpDir, "emoji_*.gif")
		require.Nil(t, err)
		fp.Write([]byte("Pretend this is an emoji: " + fmt.Sprint(i)))
		fp.Close()
	}
	return tmpDir

}

func removeTestFiles(tmpDir string, emojis []EmojiItem) {
	for _, emoji := range emojis {
		os.Remove(path.Join(tmpDir, emoji.Name))
	}
	os.Remove(tmpDir)
}
func TestPaginateEmojiList(t *testing.T) {

	tests := neko.Modern(t)

	tests.It("generates the correct keys", func(t *testing.T) {
		tmpDir := genTestFiles(t)
		emojis, err := ListEmojiFiles(tmpDir)
		require.Nil(t, err)
		pages := PaginateEmojiList(emojis)
		keys := make([]string, 0)
		expectedKeys := []string{
			"page-e-000000",
			"page-e-000001",
			"page-e-000002",
			"page-e-000003",
			"page-e-000004",
			"page-e-000005",
			"page-e-000006",
			"page-e-000007",
			"page-e-000008",
			"page-e-000009",
		}

		for key := range pages {
			keys = append(keys, key)
		}

		slices.Sort(keys)
		assert.Equal(t, expectedKeys, keys)
		removeTestFiles(tmpDir, emojis)
	})

	tests.It("doesn't duplicate emojis across pages", func(t *testing.T) {
		tmpDir := genTestFiles(t)
		emojis, err := ListEmojiFiles(tmpDir)
		require.Nil(t, err)

		pages := PaginateEmojiList(emojis)
		var lastEmoji *EmojiItem = nil
		for _, page := range pages {
			if lastEmoji != nil {
				assert.NotContains(t, page, *lastEmoji)
			}
			lastEmoji = &page[len(page)-1]
		}

		removeTestFiles(tmpDir, emojis)
	})

	tests.It("paginates files correctly", func(t *testing.T) {
		tmpDir := genTestFiles(t)
		emojis, err := ListEmojiFiles(tmpDir)
		require.Nil(t, err)

		pages := PaginateEmojiList(emojis)
		assert.Len(t, pages, 10)
		for _, page := range pages {
			assert.Len(t, page, 100)
		}

		for _, emoji := range emojis {
			assert.Contains(t, emoji.Name, "emoji")
			os.Remove(path.Join(tmpDir, emoji.Name))
		}
		os.Remove(tmpDir)

	})

	tests.Run()
}
