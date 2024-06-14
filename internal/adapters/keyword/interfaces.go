package keyword

import (
	"context"
	"my_app/internal/entity"
)

type KeywordRepo interface {
	Replace(context.Context, map[string]entity.IdMap) error
	GetKeywordIds(context.Context, string) (entity.IdMap, error)
}

type ComicRepo interface {
	GetAll(context.Context) ([]entity.Comic, error)
}
