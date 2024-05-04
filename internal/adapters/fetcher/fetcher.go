package xkcd

import (
	"context"
	"encoding/json"
	"fmt"
	"my_app/internal/entity"
	"my_app/pkg/logger"
	"net/http"
	"sync"
)

type ComicFetcher struct {
	sourceURL string
	parallel  int
	logger    logger.Logger
	mu        sync.Mutex
	wg        sync.WaitGroup
}

func NewFetcher(sourceURL string, parallel int, logger logger.Logger) *ComicFetcher {
	return &ComicFetcher{
		sourceURL: sourceURL,
		parallel:  parallel,
		logger:    logger,
	}
}

func (c *ComicFetcher) GetNew(ctx context.Context, existingComics map[int]bool) ([]entity.ComicRaw, error) {
	var comics []entity.ComicRaw
	jobs := make(chan int)
	results := make(chan entity.ComicRaw)

	errorCount := 0

	for w := 1; w <= c.parallel; w++ {
		c.wg.Add(1)
		go c.worker(w, jobs, results, &errorCount)
	}

	currId := 0
	for errorCount < 10 {
		currId++
		if existingComics[currId] {
			continue
		}
		jobs <- currId
	}
	close(jobs)
	c.wg.Wait()
	close(results)
	for comic := range results {
		comics = append(comics, comic)
	}
	return comics, nil
}

func (c *ComicFetcher) worker(id int, jobs <-chan int, results chan<- entity.ComicRaw, errorCount *int) {
	defer c.wg.Done()
	for j := range jobs {
		comic, err := c.GetComic(context.Background(), j)
		if err != nil {
			c.mu.Lock()
			(*errorCount)++
			c.mu.Unlock()
			continue
		}
		results <- comic
		c.logger.Info("Worker :%d got comic : %d", id, j)
	}
}

// Получает комикс
func (c *ComicFetcher) GetComic(ctx context.Context, id int) (entity.ComicRaw, error) {
	var comic entity.ComicRaw
	url := fmt.Sprintf("%s/%d/info.0.json", c.sourceURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return entity.ComicRaw{}, fmt.Errorf("ComicFetcher - FetchComic - http.Get: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return entity.ComicRaw{}, fmt.Errorf("ComicFetcher - FetchComic - http.Get: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&comic)
	if err != nil {
		return entity.ComicRaw{}, fmt.Errorf("ComicFetcher - FetchComic - json.Decode: %w", err)
	}
	return comic, nil
}
