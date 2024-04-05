package main

import (
	"flag"
	"fmt"
	db "my_app/pkg/database"
	"my_app/pkg/words"
	"my_app/pkg/xkcd"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SourceURL string `yaml:"source_url"`
	DBFile    string `yaml:"db_file"`
}

func main() {
	// Добавляем флаги
	var outputFlag bool
	var overwriteFlag bool
	var fromId int
	var toId int
	var chunkSize int
	var limitOutput int
	var configPath string
	flag.BoolVar(&outputFlag, "o", false, "Output JSON structure")
	flag.BoolVar(&overwriteFlag, "r", false, "Overwrite existing comics")
	flag.IntVar(&fromId, "f", 0, "Load from id")
	flag.IntVar(&toId, "t", 0, "Load to id")
	flag.IntVar(&chunkSize, "s", 100, "Chunk size load/save")
	flag.IntVar(&limitOutput, "n", 0, "Limit the number of comics")
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

	existingComics := make(map[int]bool)
	if !overwriteFlag {
		existingComics, err = db.LoadExistingComics(config.DBFile)
		if err != nil {
			fmt.Println("Error loading existing comics:", err)
			os.Exit(1)
		}
	}

	// Создаем клиент
	client := xkcd.NewClient(config.SourceURL)
	loadedComics := 0
	printLimit := limitOutput
	for {
		if limitOutput == 0 {
			printLimit = chunkSize
		}

		if toId > 0 {
			if fromId+chunkSize >= toId {
				chunkSize = toId - fromId
			}
		}
		fmt.Println(chunkSize)
		// Получаем комиксы
		comics, err := client.FetchComics(fromId, chunkSize, existingComics)
		if err != nil {
			fmt.Println("Error fetching comics:", err)
			os.Exit(1)
		}

		// Нормализуем
		normalizedComics := words.NormalizeComics(comics, false)
		// Переводим в мапу
		comicsMap := comicsToMap(normalizedComics)

		err = db.SaveComics(config.DBFile, comicsMap, overwriteFlag)
		overwriteFlag = false
		if err != nil {
			fmt.Println("Error saving comics:", err)
			os.Exit(1)
		}
		if outputFlag {
			printComicsInfo(comicsMap, printLimit)
		}
		printLimit -= len(comics)
		loadedComics += len(comics)

		fmt.Println("Loaded new comics: ", loadedComics)
		if len(comics) < chunkSize {
			break
		}
		fromId += len(comics)
		if fromId >= toId {
			break
		}
	}

}

func comicsToMap(comics []words.Comic) db.ComicsMap {
	comicsMap := make(db.ComicsMap)
	for _, comic := range comics {
		data := db.Comic{
			URL:      comic.URL,
			Keywords: comic.Keywords,
		}
		comicsMap[comic.ID] = data
	}
	return comicsMap
}

func printComicsInfo(comicsMap db.ComicsMap, limit int) {
	for id, comic := range comicsMap {
		if limit <= 0 {
			break
		}
		fmt.Println(id)
		fmt.Println(comic.URL)
		for _, value := range comic.Keywords {
			fmt.Printf("%s ", value)
		}
		fmt.Println()
		limit--
	}
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
