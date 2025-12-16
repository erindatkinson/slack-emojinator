package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

		resp.Body.Close()

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

func (c *Client) ExportEmoji(emoji Emoji, dir string) error {
	name, err := parseFile(emoji.URL)
	if err != nil {
		return err
	}

	fp, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		return err
	}
	defer fp.Close()

	resp, err := http.Get(emoji.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("bad request (%d)", resp.StatusCode)
	}

	fp.ReadFrom(resp.Body)
	return nil
}

func (c *Client) ImportEmoji(name, fPath string) error {
	slog.Info("importing emoji", "name", name)

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	addField(writer, "mode", "data")
	addField(writer, "name", name)
	addField(writer, "token", c.token)
	imgWriter, _ := writer.CreateFormFile("image", filepath.Base(fPath))
	fp, err := os.Open(fPath)
	if err != nil {
		return err
	}
	defer fp.Close()
	io.Copy(imgWriter, fp)
	writer.Close()
	contentType := writer.FormDataContentType()

	uri := fmt.Sprintf("https://%s.slack.com/api/emoji.add", c.team)

	req, err := http.NewRequest(http.MethodPost, uri, buf)
	if err != nil {
		return err
	}
	c.setHeaders(req)
	req.Header.Set("Content-Type", contentType)

	for attempts := 0; attempts < 3; attempts++ {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			retry := resp.Header.Get("Retry-After")
			seconds, err := strconv.Atoi(retry)
			if err != nil {
				return err
			}
			time.Sleep(time.Duration(seconds) * time.Second)
			continue
		}

		defer resp.Body.Close()

		var data map[string]any
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return err
		}
		slog.Info("response", "code", resp.StatusCode, "data", data)
		return nil
	}

	return fmt.Errorf("attempted 3 times and failed")
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Cookie", c.cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
}

func addField(wrapper *multipart.Writer, name, data string) error {
	writer, err := wrapper.CreateFormField(name)
	if err != nil {
		return err
	}
	writer.Write([]byte(data))
	return nil
}

func parseFile(uri string) (string, error) {
	obj, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	splits := strings.Split(obj.Path, "/")
	name, err := url.PathUnescape(splits[2])
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(obj.Path)
	return name + ext, nil
}
