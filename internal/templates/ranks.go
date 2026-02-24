package templates

import (
	"sort"
	"strings"
	"text/template"

	"github.com/erindatkinson/slack-emojinator/internal/slack"
)

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
