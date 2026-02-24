package templates

import (
	"fmt"
	"sort"
	"strings"

	"github.com/erindatkinson/slack-emojinator/internal/slack"
)

/*
BuildEmojiLists returns an array of strings that are less than 12k long,
if the whole list is less than 12k there will be only one string for later
looping for posting to the thread
*/
func BuildEmojiLists(emojis []slack.Emoji) []string {
	tpl := "* :%s: | `:%s:`\n"
	var sorted []slack.Emoji
	copy(emojis, sorted)
	sort.Slice(sorted, func(i, j int) bool {
		first := sorted[i].Name
		second := sorted[j].Name
		return first > second
	})

	batches := make([]string, 0)
	batch := ""
	for _, emoji := range emojis {
		rendered := fmt.Sprintf(tpl, emoji.Name, emoji.Name)
		if len(batch)+len(rendered) > 10_000 {
			// clone existing batch string and save to batches
			batches = append(batches, strings.Clone(batch))

			// start new batch
			batch = ""
		}
		batch = batch + rendered
	}

	batches = append(batches, strings.Clone(batch))
	return batches

}
