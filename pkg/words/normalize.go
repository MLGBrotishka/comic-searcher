package words

import (
	"regexp"
	"slices"

	"github.com/kljensen/snowball"
	"golang.org/x/exp/maps"
)

func NormalizeString(inputString string, dontStripDigits bool) []string {
	//Избавляемся от лишних _ -
	re := regexp.MustCompile(`[-_]+`)
	inputString = re.ReplaceAllString(inputString, "")
	//Избавляемся от апострофов
	re1 := regexp.MustCompile(`'[a-zA-Z]{1,2}\s`)
	inputString = re1.ReplaceAllString(inputString, " ")

	re2 := regexp.MustCompile(`\s[a-zA-Z]{1,2}'`)
	inputString = re2.ReplaceAllString(inputString, " ")

	re3 := regexp.MustCompile(`'`)
	inputString = re3.ReplaceAllString(inputString, " ")

	re4 := regexp.MustCompile(`[\pL\p{Mc}\p{Mn}-_']+`)
	words := re4.FindAllString(inputString, -1)

	// Нормализация слов
	var normalizedWords = make(map[string]bool)
	for _, word := range words {
		if _, ok := customStopWords[word]; ok {
			continue
		}
		normalized, _ := snowball.Stem(word, "english", true)
		normalizedWords[normalized] = true
	}

	// Получение результата
	keys := maps.Keys(normalizedWords)
	slices.Sort(keys)
	return keys
}
