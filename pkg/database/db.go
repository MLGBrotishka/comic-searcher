package db

import (
	"encoding/json"
	"os"
)

func SaveComics(filePath string, comics ComicsMap) error {
	// Преобразование структуры в JSON
	jsonData, err := json.Marshal(comics)
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
