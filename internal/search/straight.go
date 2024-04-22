package search

func GetStraightIndex(comics Comics) (map[int]map[string]bool, error) {
	straightIndex := make(map[int]map[string]bool)
	for key, comic := range comics {
		keywordsMap := make(map[string]bool)
		for _, keyword := range comic.Keywords {
			keywordsMap[keyword] = true
		}
		straightIndex[key] = keywordsMap
	}
	return straightIndex, nil
}

func FindStraight(straightIndex map[int]map[string]bool, normalizedInput []string, maxLen int) ([]int, error) {
	ids := make([]int, 0)
	for id := range straightIndex {
		notFound := false
		for _, keyword := range normalizedInput {
			if _, ok := straightIndex[id][keyword]; !ok {
				notFound = true
			}
		}
		if notFound {
			continue
		}
		ids = append(ids, id)
		if len(ids) == maxLen {
			break
		}
	}
	return ids, nil
}
