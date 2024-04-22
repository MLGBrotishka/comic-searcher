package db

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type ComicJsonl struct {
	ID    int   `json:"id"`
	Comic Comic `json:"comic"`
}

func LoadExistingComicsJsonl(filePath string) (map[int]bool, error) {
	// Проверяем, существует ли файл
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Если файл не существует, возвращаем пустую карту
		return make(map[int]bool), nil
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	existingComicsMap := make(map[int]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var comic ComicJsonl
		err := json.Unmarshal(scanner.Bytes(), &comic)
		if err != nil {
			return existingComicsMap, err
		}
		existingComicsMap[comic.ID] = true
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return existingComicsMap, nil
}

func LoadComicsJsonl(filePath string) (map[int]Comic, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	comics := make(map[int]Comic)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var comic ComicJsonl
		err := json.Unmarshal(scanner.Bytes(), &comic)
		if err != nil {
			return comics, err
		}
		comics[comic.ID] = comic.Comic
	}
	return comics, nil
}

// Быстрая запись, без загрузки в память
func SaveComicsJsonl(filePath string, comics ComicsMap) error {
	// Дозаписываем в файл формата jsonl комикс
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	for id, comic := range comics {
		comicToWrite := ComicJsonl{
			ID:    id,
			Comic: comic,
		}
		j, err := json.Marshal(comicToWrite)
		if err != nil {
			return fmt.Errorf("could not json marshal data: %w", err)
		}
		f.Write(j)
		f.Write([]byte("\n"))
	}
	return nil
}
