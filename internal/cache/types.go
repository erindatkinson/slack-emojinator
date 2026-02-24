package cache

type EmojiItem struct {
	Name     string
	Filename string
	Dir      string
	DocDir   string
}

type EmojiPage struct {
	Name     string
	Count    int
	NextPage string
	PrevPage string
	Emojis   []EmojiItem
}
