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

	Auth interface {
		SignIn(ctx context.Context, login string, password string) (token string, err error)
		SignUp(ctx context.Context, login string, password string) (token string, err error)
		GetUserFromToken(ctx context.Context, token string) (user entity.User, err error)
	}

	UserRepo interface {
		GetById(ctx context.Context, id int) (entity.User, error)
		GetByLogin(ctx context.Context, login string) (entity.User, error)
		Store(context.Context, entity.User) (UserId int, err error)
	}

	Hasher interface {
		GenerateSalt(ctx context.Context, size int) ([]byte, error)
		HashPassword(ctx context.Context, password string, salt []byte) (string, error)
		VerifyPassword(ctx context.Context, password string, salt []byte, hashedPassword string) (bool, error)
	}

	Authorizer interface {
		CreateToken(ctx context.Context, user entity.User) (token string, err error)
		VerifyToken(ctx context.Context, token string) (UserId int, err error)
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
