package search

import (
	"encoding/json"
	"os"
)

type Comics map[int]Comic

type Comic struct {
	URL      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

type IndexMap map[string]map[int]bool

func BuildIndex(IndexPath string, comics Comics) error {
	indexedComics := make(IndexMap)

	// _, err := os.Stat(IndexPath)
	// if !os.IsNotExist(err) {
	// 	jsonData, err := json.Unmarshal(comics)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return
	// }
	// if err != nil {
	// 	return err
	// }

	for id, comic := range comics {
		for _, keyword := range comic.Keywords {
			if indexedComics[keyword] == nil {
				indexedComics[keyword] = make(map[int]bool)
			}
			indexedComics[keyword][id] = true
		}
	}
	jsonData, err := json.Marshal(indexedComics)
	if err != nil {
		return err
	}
	err = os.WriteFile(IndexPath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadIndex(IndexPath string) (IndexMap, error) {
	data, err := os.ReadFile(IndexPath)
	if err != nil {
		return nil, err
	}
	var indexMap IndexMap
	err = json.Unmarshal(data, &indexMap)
	if err != nil {
		return nil, err
	}
	return indexMap, nil
}

func FindByIndex(indexMap IndexMap, normalizedInput []string, maxLen int) ([]int, error) {
	var intersection map[int]bool
	for _, keyword := range normalizedInput {
		if intersection == nil {
			intersection = indexMap[keyword]
			continue
		}

		for k := range intersection {
			if _, ok := indexMap[keyword][k]; !ok {
				delete(intersection, k)
			}
		}
	}
	ids := make([]int, 0)
	for k := range intersection {
		ids = append(ids, k)
		if len(ids) == maxLen {
			break
		}
	}
	return ids, nil
}
