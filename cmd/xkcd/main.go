package main

import (
	"flag"
	"fmt"
	db "my_app/pkg/database"
	"my_app/pkg/words"
	"my_app/pkg/xkcd"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SourceURL string `yaml:"source_url"`
	DBFile    string `yaml:"db_file"`
	Parallel  int    `yaml:"parallel"`
}

func main() {
	// Добавляем флаги
	var configPath string
	flag.StringVar(&configPath, "c", "", "Config path")
	flag.Parse()

	// Получаем конфиг
	if configPath == "" {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		configPath = filepath.Join(exPath, "config.yaml")
	}

	config, err := loadConfig(configPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	existingComics, err := db.LoadExistingComicsJsonl(config.DBFile)
	if err != nil {
		fmt.Println("Error loading existing comics:", err)
		os.Exit(1)
	}

	// Создаем клиент
	client := xkcd.NewClient(config.SourceURL)

	jobs := make(chan int)
	results := make(chan error)
	var goWg sync.WaitGroup
	var dbMx sync.Mutex

	for w := 1; w <= config.Parallel; w++ {
		goWg.Add(1)
		go worker(w, &goWg, existingComics, client, config.DBFile, &dbMx, jobs, results)
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
	fromId := 0
	for {
		if errorCount < 10 {
			fromId++
			jobs <- fromId
		} else {
			break
		}
	}
	fmt.Println("Остановлено из-за 10 ошибок")
	close(jobs)
	goWg.Wait()

}

func worker(id int, goWg *sync.WaitGroup, existingComics map[int]bool, client *xkcd.Client, filePath string, dbMx *sync.Mutex, jobs <-chan int, results chan<- error) {
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

func comicsToNormalizedMap(comics []xkcd.Comic) db.ComicsMap {
	comicsMap := make(db.ComicsMap)
	for _, comic := range comics {
		inputString := comic.Transcript + " " + comic.Alt
		keywords := words.NormalizeString(inputString, false)
		data := db.Comic{
			URL:      comic.URL,
			Keywords: keywords,
		}
		comicsMap[comic.ID] = data
	}
	return comicsMap
}

func loadConfig(path string) (*Config, error) {
	// Читаем файл конфигурации
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config

	// Парсим YAML в структуру Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
