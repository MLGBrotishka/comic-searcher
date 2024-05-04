package usecase

import (
	"context"
	"fmt"
	"my_app/internal/entity"
)

type ComicUseCase struct {
	repo       ComicRepo
	fetcher    ComicFetcher
	normalizer Normalizer
	keyword    Keyword
}

func NewComic(r ComicRepo, f ComicFetcher, n Normalizer, k Keyword) *ComicUseCase {
	return &ComicUseCase{
		repo:       r,
		fetcher:    f,
		normalizer: n,
		keyword:    k,
	}
}

func (uc *ComicUseCase) Update(ctx context.Context) error {
	idsExistingComics, err := uc.repo.GetAllIds(ctx)
	if err != nil {
		return fmt.Errorf("ComicUseCase - Update - uc.repo.GetAllIds: %w", err)
	}
	newRawComics, err := uc.fetcher.GetNew(ctx, idsExistingComics)
	if err != nil {
		return fmt.Errorf("ComicUseCase - Update - uc.fetcher.GetNew: %w", err)
	}
	newComics, err := normilizeRawComics(ctx, uc.normalizer, newRawComics)
	if err != nil {
		return fmt.Errorf("ComicUseCase - Update - normilizeRawComics: %w", err)
	}
	idKeywords := make(map[int]entity.KeywordMap)
	for _, comic := range newComics {
		err = uc.repo.Store(ctx, comic)
		if err != nil {
			return fmt.Errorf("ComicUseCase - Update - uc.repo.Store: %w", err)
		}
		idKeywords[comic.ID] = comic.Keywords
	}
	err = uc.keyword.Update(ctx, idKeywords)
	if err != nil {
		return fmt.Errorf("ComicUseCase - Update - uc.searcher.Update: %w", err)
	}
	return nil
}

func (uc *ComicUseCase) GetPictures(ctx context.Context, query string) ([]string, error) {
	searchKeywords, err := uc.normalizer.Normalize(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ComicUseCase - GetPictures - uc.normalizer.Normalize: %w", err)
	}
	if len(searchKeywords) == 0 {
		return nil, fmt.Errorf("ComicUseCase - GetPictures - uc.normalizer.Normalize: %w", entity.ErrBadRequest)
	}
	idFound, err := uc.keyword.Search(ctx, searchKeywords)
	if len(idFound) == 0 {
		return nil, fmt.Errorf("ComicUseCase - GetPictures - uc.searcher.Search: %w", entity.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("ComicUseCase - GetPictures - uc.searcher.Search: %w", err)
	}
	var urlFound []string
	for _, id := range idFound {
		comic, err := uc.repo.GetById(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("ComicUseCase - GetPictures - uc.repo.GetById: %w", err)
		}
		urlFound = append(urlFound, comic.URL)
	}
	return urlFound, nil
}

func normilizeRawComics(ctx context.Context, normalizer Normalizer, rawComics []entity.ComicRaw) ([]entity.Comic, error) {
	var comics []entity.Comic
	for _, rawComic := range rawComics {
		var err error
		comic := entity.Comic{
			ID:  rawComic.ID,
			URL: rawComic.URL,
		}
		comicData := rawComic.GetData()
		comic.Keywords, err = normalizer.Normalize(ctx, comicData)
		if err != nil {
			return nil, fmt.Errorf("ComicUseCase - normilizeRawComics - uc.normalizer.Normalize: %w", err)
		}
		comics = append(comics, comic)
	}
	return comics, nil
}
