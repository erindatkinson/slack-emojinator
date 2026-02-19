package templates

import "strconv"

type Docs struct {
	Namespace string
	Pages     []*EmojiPage
}

type RanksData struct {
	Ranks string
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
	Name  string
	Count int
}

func (r Rank) toArray() []string {
	return []string{r.Name, strconv.Itoa(r.Count)}
}
