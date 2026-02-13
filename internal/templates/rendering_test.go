package templates

import (
	"fmt"
	"path"
	"testing"

	"github.com/erindatkinson/slack-emojinator/internal/slack"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vektra/neko"
)

func TestLoadTemplate(t *testing.T) {
	tests := neko.Modern(t)

	tests.It("can load the expected templates", func(t *testing.T) {
		dir := "templates/"
		expectedTemplates := []string{
			"header.md.jinja2",
			"ranks.md.jinja2",
			"stats.md.jinja2",
		}

		for _, tplFile := range expectedTemplates {
			_, err := LoadTemplate(path.Join(dir, tplFile))
			assert.Nil(t, err)
		}
	})

	tests.Run()
}

func TestRenderWithData(t *testing.T) {
	tests := neko.Modern(t)

	tests.It("renders the header correctly", func(t *testing.T) {
		tpl, err := LoadTemplate("templates/header.md.jinja2")
		require.Nil(t, err)

		data := map[string]any{
			"start": "Start",
			"end":   "End",
		}

		rendered, err := RenderWithData(*tpl, data)
		assert.Nil(t, err)
		assert.Equal(t, "# :sby-a-new-emoji: Emoji Release Notes Start - End", rendered)

	})

	tests.Run()
}

func TestBuildEmojiLists(t *testing.T) {
	tests := neko.Modern(t)

	tests.It("renders simple lists", func(t *testing.T) {

		emojis := []slack.Emoji{
			{Name: "a-test"},
			{Name: "b-test"},
			{Name: "c-test"},
			{Name: "d-test"},
			{Name: "e-test"},
		}
		rendered := BuildEmojiLists(emojis)
		expected := []string{"* :a-test: | `:a-test:`\n" +
			"* :b-test: | `:b-test:`\n" +
			"* :c-test: | `:c-test:`\n" +
			"* :d-test: | `:d-test:`\n" +
			"* :e-test: | `:e-test:`\n"}

		assert.Equal(t, expected, rendered)
	})

	tests.It("renders large lists", func(t *testing.T) {
		var emojis []slack.Emoji = make([]slack.Emoji, 0)
		for i := 0; i < 1000; i++ {
			emojis = append(emojis, slack.Emoji{Name: fmt.Sprintf("test-%d", i)})
		}
		require.Len(t, emojis, 1000)

		rendered := BuildEmojiLists(emojis)
		// 0-9:     10  @ len("* :test-_: | `:test-_:`\n")     => 24 == 240
		// 10-99:   90  @ len("* :test-__: | `:test-__:`\n")   => 26 == 2_340
		// 100-999: 900 @ len("* :test-___: | `:test-___:`\n") => 28 == 25_200
		// Sum: 27_780
		// 10_000 limit batches: 3

		assert.Len(t, rendered, 3)
		assert.LessOrEqual(t, len(rendered[0]), 10_000)
		assert.LessOrEqual(t, len(rendered[1]), 10_000)
		assert.Less(t, len(rendered[2]), 10_000)
	})

	tests.Run()
}

func TestRenderRanks(t *testing.T) {
	tests := neko.Modern(t)

	tests.It("handles the ranks correctly", func(t *testing.T) {
		// default len is 25 so this should be okay with formatting expectation
		users := []string{faker.Username(), faker.Username(), faker.Username()}
		alternate := true
		emojis := make([]slack.Emoji, 0)
		for i := 0; i < 1000; i++ {
			// 50% user 0, 25% users 1&2
			var user string
			if i%2 == 0 {
				user = users[0]
			} else {
				if alternate {
					user = users[1]
				} else {
					user = users[2]
				}
				alternate = !alternate
			}
			emojis = append(emojis, slack.Emoji{
				UserDisplayName: user,
			})
		}

		rendered, err := RenderRanks(emojis)
		require.Nil(t, err)
		expected := fmt.Sprintf(
			"# :sby-a-new-emoji: Emoji Release Notes\n\n## Uploaders\n\n```\n┌─────────┬───────┐\n│  USER   │ COUNT │\n├─────────┼───────┤\n│ %s │ 500   │\n│ %s │ 250   │\n│ %s │ 250   │\n└─────────┴───────┘\n\n```\n",
			users[0], users[1], users[2])
		assert.Equal(t, rendered, expected)

	})

	tests.Run()
}
