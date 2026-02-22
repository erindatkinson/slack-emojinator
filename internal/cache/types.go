package cache

type EmojiItem struct {
	Name    string
	DocPath string
}

type EmojiPage struct {
	Name     string
	Count    int
	NextPage string
	PrevPage string
	Emojis   []EmojiItem
}
