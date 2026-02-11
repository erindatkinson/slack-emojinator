package templates

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/erindatkinson/slack-emojinator/internal/slack"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/olekukonko/tablewriter"
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

func Render(tpl exec.Template) (string, error) {
	return tpl.ExecuteToString(exec.EmptyContext())
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

type Rank struct {
	Name  string
	Count int
}

func (r Rank) toArray() []string {
	return []string{r.Name, strconv.Itoa(r.Count)}
}

func RenderRanks(emojis []slack.Emoji) (string, error) {

	tpl, err := LoadTemplate("templates/ranks.md.jinja2")
	if err != nil {
		return "", err
	}

	rankMap := make(map[string]*Rank)
	ranks := make([]*Rank, 0)
	for _, emoji := range emojis {
		if _, ok := rankMap[emoji.UserDisplayName]; ok {
			rankMap[emoji.UserDisplayName].Count = rankMap[emoji.UserDisplayName].Count + 1
		} else {
			newRank := &Rank{
				Name:  emoji.UserDisplayName,
				Count: 1,
			}
			rankMap[emoji.UserDisplayName] = newRank
			ranks = append(ranks, newRank)
		}
	}

	sort.Slice(ranks, func(i, j int) bool {
		first := ranks[i].Count
		second := ranks[j].Count
		return first > second
	})

	var builder strings.Builder
	tab := tablewriter.NewWriter(&builder)
	tab.Header([]string{"User", "Count"})
	for _, user := range ranks {
		tab.Append(user.toArray())
	}
	tab.Render()

	return RenderWithData(*tpl, map[string]any{"ranks": builder.String()})

}
