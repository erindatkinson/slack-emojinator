package slack

type EmojiJsonFile struct {
	Ok    bool
	Emoji map[string]string
}

type Pagination struct {
	Count int64 `json:"count"`
	Page  int64 `json:"page"`
	Pages int64 `json:"pages"`
	Total int64 `json:"total"`
}

type Emoji struct {
	Name            string   `json:"name"`
	Created         int64    `json:"created"`
	IsAlias         int64    `json:"is_alias"`
	AliasFor        string   `json:"alias_for"`
	Synonyms        []string `json:"synonyms"`
	CanDelete       bool     `json:"can_delete"`
	IsBad           bool     `json:"is_bad"`
	AvatarHash      string   `json:"avatar_hash"`
	TeamId          string   `json:"team_id"`
	URL             string   `json:"url"`
	UserDisplayName string   `json:"user_display_name"`
	UserID          string   `json:"user_id"`
}

type EmojiList struct {
	Ok                    bool       `json:"ok"`
	Emoji                 []Emoji    `json:"emoji"`
	DisabledEmoji         []Emoji    `json:"disabled_emoji"`
	CustomEmojiTotalCount int64      `json:"custom_emoji_total_count"`
	Paging                Pagination `json:"paging"`
}
