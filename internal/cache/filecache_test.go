package cache

import (
	"os"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vektra/neko"
)

func helpBuildTestDir(t *testing.T) (string, func(*testing.T)) {
	dir, err := os.MkdirTemp(".", "tmp_*")
	require.Nil(t, err)
	var files []string
	for i := 0; i < 1000; i++ {
		fp, err := os.CreateTemp(dir, faker.FirstName()+"_*.png")
		require.Nil(t, err)
		files = append(files, fp.Name())
		fp.WriteString(faker.LastName())
		fp.Close()
	}
	return dir, func(st *testing.T) {
		for _, file := range files {
			err := os.Remove(file)
			require.Nil(st, err)
		}
		err := os.Remove(dir)
		require.Nil(st, err)
	}

}
func TestListDownloadedEmojis(t *testing.T) {
	tests := neko.Modern(t)

	tests.It("lists data correctly", func(t *testing.T) {
		dir, closer := helpBuildTestDir(t)
		defer closer(t)

		emojis, err := ListDownloadedEmojis(dir)
		require.Nil(t, err)
		assert.Len(t, emojis, 1000)
		for _, emoji := range emojis {
			assert.Contains(t, emoji.Filename, ".png")
			assert.Contains(t, emoji.DocDir, "docs")
		}
	})

	tests.Run()
}
