package db

import (
	"encoding/json"
	"os"
)

// Загружает существующие комиксы из файла и возвращает сет индексов
func LoadExistingComics(filePath string) (map[int]bool, error) {
	// Проверяем, существует ли файл
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Если файл не существует, возвращаем пустую карту
		return make(map[int]bool), nil
	}

	// Читаем данные из файла
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Преобразуем данные в структуру ComicsMap
	var existingComics ComicsMap
	err = json.Unmarshal(data, &existingComics)
	if err != nil {
		return nil, err
	}

	// Создаем сет индексов существующих комиксов
	existingComicsMap := make(map[int]bool)
	for id := range existingComics {
		existingComicsMap[id] = true
	}

	return existingComicsMap, nil
}

//todo: быстрая запись, без загрузки в память

func SaveComics(filePath string, comics ComicsMap, overwrite bool) error {
	// Проверяем, существует ли файл
	if _, err := os.Stat(filePath); os.IsNotExist(err) || overwrite {
		// Если файл не существует или включена перезапись, создаем его с новыми комиксами
		jsonData, err := json.Marshal(comics)
		if err != nil {
			return err
		}
		return os.WriteFile(filePath, jsonData, 0644)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var existingComics ComicsMap
	err = json.Unmarshal(data, &existingComics)
	if err != nil {
		return err
	}

	for id, comic := range comics {
		if _, exists := existingComics[id]; !exists {
			existingComics[id] = comic
		}
	}

	// Преобразование структуры в JSON
	jsonData, err := json.Marshal(existingComics)
	if err != nil {
		return err
	}

	// Запись JSON в файл
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
