package templates

type Docs struct {
	Namespace string
	Pages     []*EmojiPage
}

type ReleaseData struct {
	Start string
	End   string
}

type EmojiPage struct {
	Name     string
	Count    int
	NextPage string
	PrevPage string
	Emojis   []EmojiItem
}

type EmojiItem struct {
	Name    string
	DocPath string
}

type RanksData struct {
	Keys  []string
	Ranks map[string]int
}
