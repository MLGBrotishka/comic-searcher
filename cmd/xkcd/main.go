package main

import (
	"flag"
	"fmt"
	"my_app/internal/search"
	"my_app/internal/xkcd"
	db "my_app/pkg/database"
	"my_app/pkg/words"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SourceURL string `yaml:"source_url"`
	DBPath    string `yaml:"db_file"`
	Parallel  int    `yaml:"parallel"`
	IndexPath string `yaml:"index_file"`
}

func main() {
	// Добавляем флаги
	var configPath string
	var inputString string
	var debugFrom int
	var useIndexSearch bool
	flag.StringVar(&configPath, "c", "", "Config path")
	flag.StringVar(&inputString, "s", "", "Input string to find")
	flag.BoolVar(&useIndexSearch, "i", false, "Use index search")
	flag.IntVar(&debugFrom, "f", 0, "From id load, need to debug")
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

	// Проверяем строку
	if inputString == "" {
		fmt.Println("Please provide an input string with -s")
		return
	}

	config, err := loadConfig(configPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	// Получение новых комиксов и сохранение их в базу данных
	err = xkcd.GetNewComics(config.DBPath, config.SourceURL, config.Parallel, debugFrom)
	if err != nil {
		fmt.Println("Error getting new comics:", err)
		panic(err)
	}

	// Загрузка комиксов из базы данных
	comics, err := db.LoadComicsJsonl(config.DBPath)
	if err != nil {
		fmt.Println("Error loading comics:", err)
		panic(err)
	}

	// Преобразование комиксов из формата базы данных в формат поиска
	searchComics := convertComics(comics)

	// Построение индекса для поиска по ключевым словам
	err = search.BuildIndex(config.IndexPath, searchComics)
	if err != nil {
		fmt.Println("Error building indexes:", err)
		panic(err)
	}

	// Загрузка индекса для поиска
	indexMap, err := search.LoadIndex(config.IndexPath)
	if err != nil {
		fmt.Println("Error loading indexes:", err)
		panic(err)
	}

	// Нормализация входной строки для поиска
	normalizedInput := words.NormalizeString(inputString, false)
	if len(normalizedInput) == 0 {
		fmt.Println("Please provide more information")
		return
	}

	// Поиск комиксов по индексу или напрямую
	var foundIds []int
	if useIndexSearch {
		foundIds, _ = search.FindByIndex(indexMap, normalizedInput, 10)
	} else {
		straightIndex, _ := search.GetStraightIndex(searchComics)
		foundIds, _ = search.FindStraight(straightIndex, normalizedInput, 10)
	}

	// Вывод результатов поиска
	if len(foundIds) == 0 {
		fmt.Println("Not found")
		return
	}
	fmt.Println("Found:", len(foundIds), "comics")
	for _, id := range foundIds {
		fmt.Println(id, comics[id].URL)
	}
}

// Преобразует комиксы из формата базы данных в формат поиска
func convertComics(dbComics map[int]db.Comic) search.Comics {
	searchComics := make(search.Comics)
	for id, comic := range dbComics {
		searchComics[id] = search.Comic{
			URL:      comic.URL,
			Keywords: comic.Keywords,
		}
	}
	return searchComics
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
