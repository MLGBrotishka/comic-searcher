package normalizer

import (
	"context"
	"fmt"
	"regexp"

	"github.com/kljensen/snowball"
)

type Normalizer struct {
	//Maybe conf
}

func New() *Normalizer {
	return &Normalizer{}
}

func (n *Normalizer) Normalize(ctx context.Context, inputString string) (map[string]bool, error) {
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
		if stopWords[word] {
			continue
		}
		normalized, err := snowball.Stem(word, "english", true)
		if err != nil {
			return nil, fmt.Errorf("Normalizer - Normalize - snowball.Stem: %w", err)
		}
		normalizedWords[normalized] = true
	}

	return normalizedWords, nil
}
