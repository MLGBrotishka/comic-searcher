package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"

	"github.com/bbalet/stopwords"
	"github.com/kljensen/snowball"
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
	re := regexp.MustCompile(`'[^ ]*`)
	inputString = re.ReplaceAllString(inputString, "")

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
	var keys []string
	for word := range normalizedWords {
		keys = append(keys, word)
	}
	fmt.Println(strings.Join(keys, " "))
}
