package utilities

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vektra/neko"
)

func TestGetDownloadedEmojiList(t *testing.T) {
	tests := neko.Modern(t)

	tests.It("gets all files in path correctly", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp(".", "tmp_*")
		require.Nil(t, err)
		for i := 0; i < 1000; i++ {
			fp, err := os.CreateTemp(tmpDir, "emoji_*.gif")
			require.Nil(t, err)
			fp.Write([]byte("Pretend this is an emoji: " + fmt.Sprint(i)))
			fp.Close()
		}

		emojis, err := GetDownloadedEmojiList(tmpDir)
		assert.Nil(t, err)
		assert.Len(t, emojis, 1000)
		for _, emoji := range emojis {
			assert.Contains(t, emoji, "emoji")
			os.Remove(path.Join(tmpDir, emoji+".gif"))
		}
		os.Remove(tmpDir)

	})

	tests.Run()
}
