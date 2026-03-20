package templates

import "github.com/erindatkinson/emoji-archiver/internal/cache"

type Docs struct {
	Namespace string
	Pages     []*cache.EmojiPage
}

type ReleaseData struct {
	Start string
	End   string
}

type RanksData struct {
	Keys  []string
	Ranks map[string]int
}
