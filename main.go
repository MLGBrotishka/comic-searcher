package main

import (
	"flag"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/bbalet/stopwords"
	"github.com/kljensen/snowball"
	"golang.org/x/exp/maps"
)

func main() {
	//Создаем флаги
	var inputString string
	var dontStripDigits bool
	flag.StringVar(&inputString, "s", "", "Input string to normalize")
	flag.BoolVar(&dontStripDigits, "n", false, "Don't strip digits")

	// Парсинг аргументов командной строки
	flag.Parse()

	// Если строка не была передана, выводим сообщение об ошибке и завершаем работу
	if inputString == "" {
		fmt.Println("Please provide an input string with -s")
		return
	}

	//Избавляемся от апострофов
	re1 := regexp.MustCompile(`'[a-zA-Z]{1,2}\s`)
	inputString = re1.ReplaceAllString(inputString, "")

	re2 := regexp.MustCompile(`\s[a-zA-Z]{1,2}'`)
	inputString = re2.ReplaceAllString(inputString, "")

	re3 := regexp.MustCompile(`'`)
	inputString = re3.ReplaceAllString(inputString, " ")

	// Очистка строки от стоп-слов
	if dontStripDigits {
		stopwords.DontStripDigits()
	}
	cleanedStrings := stopwords.CleanString(inputString, "en", false)
	words := strings.Fields(cleanedStrings)

	// Нормализация слов
	var normalizedWords = make(map[string]bool)
	for _, word := range words {
		if _, ok := customStopWords[word]; ok {
			continue
		}
		normalized, _ := snowball.Stem(word, "english", true)
		normalizedWords[normalized] = true
	}

	// Вывод результата
	keys := maps.Keys(normalizedWords)
	slices.Sort(keys)
	fmt.Println(strings.Join(keys, " "))
}
