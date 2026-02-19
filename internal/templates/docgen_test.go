package templates

import (
	"fmt"
	"os"
	"path"
	"testing"

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
		emojis, err := ListEmojiFiles(tmpDir, "tmpDir")
		require.Nil(t, err)

		removeTestFiles(tmpDir, emojis)
	})

	tests.It("doesn't duplicate emojis across pages", func(t *testing.T) {
		tmpDir := genTestFiles(t)
		emojis, err := ListEmojiFiles(tmpDir, "tmpDir")
		require.Nil(t, err)

		removeTestFiles(tmpDir, emojis)
	})

	tests.It("paginates files correctly", func(t *testing.T) {
		tmpDir := genTestFiles(t)
		emojis, err := ListEmojiFiles(tmpDir, "tmpDir")
		require.Nil(t, err)

		removeTestFiles(tmpDir, emojis)
	})

	tests.Run()
}
