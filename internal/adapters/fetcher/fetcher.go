package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"my_app/internal/entity"
	"net/http"
	"sync"
)

type ComicFetcher struct {
	sourceURL string
	parallel  int
	logger    Logger
	muErr     sync.Mutex
	muRes     sync.Mutex
	wg        sync.WaitGroup
}

func NewFetcher(sourceURL string, parallel int, logger Logger) *ComicFetcher {
	return &ComicFetcher{
		sourceURL: sourceURL,
		parallel:  parallel,
		logger:    logger,
	}
}

func (c *ComicFetcher) GetNew(ctx context.Context, existingComics entity.IdMap) ([]entity.ComicRaw, error) {
	var comics []entity.ComicRaw
	jobs := make(chan int)

	errorCount := 0

	for w := 1; w <= c.parallel; w++ {
		c.wg.Add(1)
		go c.worker(w, jobs, &comics, &errorCount)
	}

	currId := 0
	for errorCount < 5 {
		currId++
		if existingComics[currId] {
			continue
		}
		jobs <- currId
	}
	close(jobs)
	c.wg.Wait()
	return comics, nil
}

func (c *ComicFetcher) worker(id int, jobs <-chan int, comics *[]entity.ComicRaw, errorCount *int) {
	defer c.wg.Done()
	for j := range jobs {
		comic, err := c.GetComic(context.Background(), j)
		if err != nil {
			c.muErr.Lock()
			(*errorCount)++
			c.muErr.Unlock()
			c.logger.Debug("Worker %d failed to get comic %d: %v", id, j, err)
			continue
		}
		c.muRes.Lock()
		(*comics) = append((*comics), comic)
		c.muRes.Unlock()
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
