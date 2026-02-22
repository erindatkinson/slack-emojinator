package templates

import (
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/erindatkinson/slack-emojinator/internal/slack"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
)

func LoadTemplate(path string) (*exec.Template, error) {
	return gonja.FromString(MustAssetString(path))
}

func RenderWithData(tpl exec.Template, data map[string]any) (string, error) {
	ctx := exec.EmptyContext()
	for k, v := range data {
		ctx.Set(k, v)
	}
	return tpl.ExecuteToString(ctx)
}

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

// RenderRanks iterates through
func RenderRanks(emojis []slack.Emoji) (string, error) {
	tpl := template.New("ranks")

	// Using a 'dict' as it makes this an easier loop
	// 1 loop thru all emojis, 1 loop thru names to sort, 1 loop thru names to print
	rankMap := make(map[string]int)
	keys := make([]string, 0)
	maxLen := 0
	for _, emoji := range emojis {
		if len(emoji.UserDisplayName) > maxLen {
			maxLen = len(emoji.UserDisplayName)
		}

		if _, ok := rankMap[emoji.UserDisplayName]; ok {
			rankMap[emoji.UserDisplayName] += 1
		} else {
			rankMap[emoji.UserDisplayName] = 1
			keys = append(keys, emoji.UserDisplayName)
		}
	}

	sort.Slice(keys, func(i, j int) bool {
		return rankMap[keys[i]] > rankMap[keys[j]]
	})

	renderData := RanksData{
		Keys:  keys,
		Ranks: rankMap,
	}

	// add padding function for spacing counts based on longest name
	tpl = tpl.Funcs(template.FuncMap{
		"padding": func(value string) string {
			var out string
			for i := 0; i < maxLen-len(value)+1; i++ {
				out += " "
			}
			return out
		}})

	tpl, err := tpl.Parse(
		MustAssetString("templates/ranks.md.gotmpl"))
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	if err = tpl.Execute(&builder, renderData); err != nil {
		return "", err
	}

	return builder.String(), nil
}
