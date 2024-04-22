package search

import (
	"fmt"
	db "my_app/pkg/database"
	"my_app/pkg/words"
	"testing"
	"time"
)

var indexMap IndexMap
var straightIndex map[int]map[string]bool
var normalizedInput []string
var maxLen int

// Инициализация данных перед тестами
func init() {
	var err error
	// Создание индекса
	indexMap, err = LoadIndex("../../index.json")
	if err != nil {
		panic(err)
	}

	comics, err := db.LoadComicsJsonl("../../database.jsonl")
	if err != nil {
		panic(err)
	}
	// Создание прямого индекса
	straightIndex, err = GetStraightIndex(convertComics(comics))
	if err != nil {
		panic(err)
	}
	normalizedInput = words.NormalizeString("Number of communicating civilizations in our galaxy", false)
	maxLen = 10
	// Здесь вы можете сохранить indexMap и straightIndex в глобальные переменные,
	// чтобы использовать их в бенчмарках
}

func convertComics(dbComics map[int]db.Comic) Comics {
	searchComics := make(Comics)
	for id, comic := range dbComics {
		searchComics[id] = Comic{
			URL:      comic.URL,
			Keywords: comic.Keywords,
		}
	}
	return searchComics
}

// Бенчмарк для FindByIndex
func BenchmarkFindByIndex(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindByIndex(indexMap, normalizedInput, maxLen)
	}
	fmt.Println("number of iterations: ", b.N)
	fmt.Println("elapsed:", b.Elapsed()/time.Duration(b.N))
}

// Бенчмарк для FindStraight
func BenchmarkFindStraight(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindStraight(straightIndex, normalizedInput, maxLen)
	}
	fmt.Println("number of iterations: ", b.N)
	fmt.Println("elapsed:", b.Elapsed()/time.Duration(b.N))
}
