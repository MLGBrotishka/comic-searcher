package xkcd

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	sourceURL string
}

func NewClient(sourceURL string) *Client {
	return &Client{sourceURL: sourceURL}
}

// Получает все комиксы начиная с from до to, которые еще не загружены
func (c *Client) FetchComics(from int, to int, existingComics map[int]bool) ([]Comic, error) {
	var comics []Comic
	i := from
	for {
		if _, exist := existingComics[i]; exist || i == 404 { //Прикол сайта
			i++
			continue
		}
		// Если достигнут лимит, прерываем цикл
		if i > to {
			break
		}
		url := fmt.Sprintf("%s/%d/info.0.json", c.sourceURL, i)
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			// Если возникла ошибка, прерываем цикл
			break
		}
		defer resp.Body.Close()

		var comic Comic
		err = json.NewDecoder(resp.Body).Decode(&comic)
		if err != nil {
			return nil, err
		}

		comics = append(comics, comic)
		i++

	}
	return comics, nil
}
