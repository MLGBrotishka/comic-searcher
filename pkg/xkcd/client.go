package xkcd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	sourceURL string
}

func NewClient(sourceURL string) *Client {
	return &Client{sourceURL: sourceURL}
}

func (c *Client) FetchComics(limit int) ([]Comic, error) {
	var comics []Comic
	i := 1
	for {
		url := fmt.Sprintf("%s/%d/info.0.json", c.sourceURL, i)
		resp, err := http.Get(url)
		if err != nil {
			// Если возникла ошибка, прерываем цикл
			break
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var comic Comic
		err = json.Unmarshal(body, &comic)
		if err != nil {
			return nil, err
		}

		comics = append(comics, comic)
		i++

		// Если limit не равен 0 и достигнут лимит, прерываем цикл
		if limit > 0 && len(comics) >= limit {
			break
		}
	}
	return comics, nil
}
