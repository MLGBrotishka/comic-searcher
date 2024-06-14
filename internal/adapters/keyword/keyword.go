package keyword

import (
	"context"
	"fmt"
	"my_app/internal/entity"
)

type Keyword struct {
	repo KeywordRepo
}

func New(r KeywordRepo) *Keyword {
	return &Keyword{
		repo: r,
	}
}

func NewLoad(r KeywordRepo, comicrepo ComicRepo) (*Keyword, error) {
	k := &Keyword{
		repo: r,
	}
	comics, err := comicrepo.GetAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Keyword - NewLoad - comicrepo.GetAll: %w", err)
	}
	idKeywords := make(map[int]entity.KeywordMap, len(comics))
	for _, comic := range comics {
		idKeywords[comic.ID] = comic.Keywords
	}
	err = k.Update(context.Background(), idKeywords)
	if err != nil {
		return nil, fmt.Errorf("Keyword - NewLoad - Update: %w", err)
	}
	return k, nil
}

func (k *Keyword) Update(ctx context.Context, idKeywords map[int]entity.KeywordMap) error {
	invertedIndex := make(map[string]entity.IdMap)
	for id := range idKeywords {
		for keyword := range idKeywords[id] {
			if invertedIndex[keyword] == nil {
				invertedIndex[keyword] = make(entity.IdMap)
			}
			invertedIndex[keyword][id] = true
		}
	}
	err := k.repo.Replace(ctx, invertedIndex)
	if err != nil {
		return fmt.Errorf("Keyword - Update - k.repo.Replace: %w", err)
	}
	return nil
}

func (k *Keyword) Search(ctx context.Context, keywords entity.KeywordMap) ([]int, error) {
	if len(keywords) == 0 {
		return nil, fmt.Errorf("Keyword - Search: %s", "keywords are empty")
	}
	var idIntersect entity.IdMap
	var err error
	for keyword := range keywords {
		if idIntersect == nil {
			idIntersect, err = k.repo.GetKeywordIds(ctx, keyword)
			if err != nil {
				return nil, fmt.Errorf("Keyword - Search - k.repo.GetKeywordIds: %w", err)
			}
			continue
		}

		keywordIds, err := k.repo.GetKeywordIds(ctx, keyword)
		if err != nil {
			return nil, fmt.Errorf("Keyword - Search - k.repo.GetKeywordIds: %w", err)
		}

		//Intersect
		for id := range idIntersect {
			if !keywordIds[id] {
				delete(idIntersect, id)
			}
		}
	}
	ids := make([]int, 0, len(idIntersect))
	for id := range idIntersect {
		ids = append(ids, id)
	}
	return ids, nil
}
