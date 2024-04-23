package xkcd

import (
	"fmt"
	db "my_app/pkg/database"
	"my_app/pkg/words"
	"sync"
)

func GetNewComics(DBPath string, SourceURL string, Parallel int, debugFrom int) error {
	existingComics, err := db.LoadExistingComicsJsonl(DBPath)

	if err != nil {
		fmt.Println("Error loading existing comics:", err)
		return err
	}

	// Создаем клиент
	client := NewClient(SourceURL)

	jobs := make(chan int)
	results := make(chan error)
	var goWg sync.WaitGroup
	var dbMx sync.Mutex

	for w := 1; w <= Parallel; w++ {
		goWg.Add(1)
		go worker(w, &goWg, existingComics, client, DBPath, &dbMx, jobs, results)
	}

	// Горутина для подсчета ошибок
	errorCount := 0

	go func(errorCount *int) {
		count := 0
		for err := range results {
			if err != nil {
				*errorCount++
				if count >= 10 {
					break
				}
			}
		}
	}(&errorCount)

	// Основной цикл
	// fromId := 0
	// debug
	fromId := debugFrom
	for {
		if errorCount < 10 {
			fromId++
			jobs <- fromId
		} else {
			break
		}
	}
	close(jobs)
	goWg.Wait()
	return nil
}

func worker(id int, goWg *sync.WaitGroup, existingComics map[int]bool, client *Client, filePath string, dbMx *sync.Mutex, jobs <-chan int, results chan<- error) {
	defer goWg.Done()
	for j := range jobs {
		if j == -1 {
			return
		}
		if _, ok := existingComics[j]; ok {
			continue
		}
		comics, err := client.FetchComics(j, j, existingComics)
		if err != nil {
			results <- err
			continue
		}
		comicsMap := comicsToNormalizedMap(comics)
		dbMx.Lock()
		err = db.SaveComicsJsonl(filePath, comicsMap)
		dbMx.Unlock()
		if err != nil {
			results <- err
			continue
		}
		fmt.Println("Worker :", id, " saved comic ", j)
		results <- nil

	}
}

func comicsToNormalizedMap(comics []Comic) db.ComicsMap {
	comicsMap := make(db.ComicsMap)
	for _, comic := range comics {
		inputString := comic.Transcript + " " + comic.Alt + " " + comic.Title
		keywords := words.NormalizeString(inputString, false)
		data := db.Comic{
			URL:      comic.URL,
			Keywords: keywords,
		}
		comicsMap[comic.ID] = data
	}
	return comicsMap
}
