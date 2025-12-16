package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
)

type Client struct {
	cookie string
	token  string
	team   string
}

func NewSlackClient(team, token, cookie string) *Client {
	return &Client{
		cookie: cookie,
		token:  token,
		team:   team,
	}
}

func (c *Client) PostMessage() {

}

func (c *Client) ListEmoji() ([]Emoji, error) {
	emojis := make([]Emoji, 0)
	uri := fmt.Sprintf("https://%s.slack.com/api/emoji.adminList", c.team)
	params := url.Values{}
	params.Set("query", "")
	params.Set("page", "1")
	params.Set("count", "1000")
	params.Set("token", c.token)

	for {
		slog.Info("Downloading list", "page", params.Get("page"))
		payload := bytes.NewBufferString(params.Encode())
		req, err := http.NewRequest(http.MethodPost, uri, payload)
		if err != nil {
			return []Emoji{}, err
		}

		c.setHeaders(req)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return []Emoji{}, err
		}
		if resp.StatusCode != http.StatusOK {
			return []Emoji{}, fmt.Errorf("response code error: %d", resp.StatusCode)
		}

		data := EmojiList{}
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return []Emoji{}, err
		}

		if data.Ok {
			emojis = append(emojis, data.Emoji...)
			if data.Paging.Page+1 > data.Paging.Pages {
				return emojis, nil
			} else {
				params.Set("page", fmt.Sprint(data.Paging.Page+1))
			}
		} else {
			return []Emoji{}, fmt.Errorf("bad response: %v", data)
		}
	}
}

func (c *Client) ExportEmoji() {

}

func (c *Client) ImportEmoji() {

}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Cookie", c.cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
}
