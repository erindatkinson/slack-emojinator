package templates

import "strconv"

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
