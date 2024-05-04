package usecase

import (
	"context"
	"my_app/internal/entity"
)

type (
	Normalizer interface {
		Normalize(context.Context, string) (entity.KeywordMap, error)
	}

	Keyword interface {
		Update(context.Context, map[int]entity.KeywordMap) error
		Search(context.Context, entity.KeywordMap) ([]int, error)
	}

	KeywordRepo interface {
		Replace(context.Context, map[string]entity.IdMap) error
		GetKeywordIds(context.Context, string) (entity.IdMap, error)
	}

	Comic interface {
		Update(context.Context) error
		GetPictures(context.Context, string) ([]string, error)
	}

	ComicRepo interface {
		Store(context.Context, entity.Comic) error
		GetById(context.Context, int) (entity.Comic, error)
		GetAll(context.Context) ([]entity.Comic, error)
		GetAllIds(context.Context) (entity.IdMap, error)
	}

	ComicFetcher interface {
		GetNew(context.Context, entity.IdMap) ([]entity.ComicRaw, error)
	}
)
