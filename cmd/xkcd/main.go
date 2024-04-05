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
	var limit int
	flag.BoolVar(&outputFlag, "o", false, "Output JSON structure")
	flag.IntVar(&limit, "n", 0, "Limit the number of comics")
	flag.Parse()

	// Получаем конфиг
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	configPath := filepath.Join(exPath, "config.yaml")
	config, err := loadConfig(configPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	// Создаем клиент
	client := xkcd.NewClient(config.SourceURL)
	// Получаем комиксы
	comics, err := client.FetchComics(limit)
	if err != nil {
		fmt.Println("Error fetching comics:", err)
		os.Exit(1)
	}

	// Нормализуем
	normalizedComics := words.NormalizeComics(comics, false)
	// Переводим в мапу
	comicsMap := comicsToMap(normalizedComics)

	if outputFlag {
		fmt.Println(comicsMap)
	} else {
		err = db.SaveComics(config.DBFile, comicsMap)
		if err != nil {
			fmt.Println("Error saving comics:", err)
			os.Exit(1)
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
