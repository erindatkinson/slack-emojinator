package templates

type Docs struct {
	Namespace string
	Pages     []*EmojiPage
}

type RanksData struct {
	Ranks []Rank
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

type Rank struct {
	Name    string
	Count   int
	Padding string
}
