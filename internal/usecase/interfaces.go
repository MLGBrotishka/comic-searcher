package usecase

import (
	"context"
	"my_app/internal/entity"
)

type (
	Authorizer interface {
		CreateToken(ctx context.Context, user entity.User) (token string, err error)
		VerifyToken(ctx context.Context, token string) (userID int, err error)
	}

	Hasher interface {
		GenerateSalt(ctx context.Context, size int) ([]byte, error)
		HashPassword(ctx context.Context, password string, salt []byte) (string, error)
		VerifyPassword(ctx context.Context, password string, salt []byte, hashedPassword string) (bool, error)
	}
)

type (
	Normalizer interface {
		Normalize(context.Context, string) (entity.KeywordMap, error)
	}

	Keyword interface {
		Update(context.Context, map[int]entity.KeywordMap) error
		Search(context.Context, entity.KeywordMap) ([]int, error)
	}
)

type ComicFetcher interface {
	GetNew(context.Context, entity.IdMap) ([]entity.ComicRaw, error)
}

type (
	UserRepo interface {
		GetById(ctx context.Context, id int) (entity.User, error)
		GetByLogin(ctx context.Context, login string) (entity.User, error)
		Store(context.Context, entity.User) (userID int, err error)
	}

	ComicRepo interface {
		Store(context.Context, entity.Comic) error
		GetById(context.Context, int) (entity.Comic, error)
		GetAll(context.Context) ([]entity.Comic, error)
		GetAllIds(context.Context) (entity.IdMap, error)
	}
)
