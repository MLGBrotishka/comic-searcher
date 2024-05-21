package entity

import (
	"context"
)

//TODO: раскидать интерфейсы по мере использования, так не дело

// Интерфейсы, связанные с нормализацией и обработкой ключевых слов
type (
	Normalizer interface {
		Normalize(context.Context, string) (KeywordMap, error)
	}

	Keyword interface {
		Update(context.Context, map[int]KeywordMap) error
		Search(context.Context, KeywordMap) ([]int, error)
	}

	KeywordRepo interface {
		Replace(context.Context, map[string]IdMap) error
		GetKeywordIds(context.Context, string) (IdMap, error)
	}
)

// Интерфейсы, связанные с аутентификацией и авторизацией пользователей
type (
	AuthUseCase interface {
		SignIn(ctx context.Context, login string, password string) (token string, err error)
		SignUp(ctx context.Context, login string, password string) (token string, err error)
		GetUserFromToken(ctx context.Context, token string) (user User, err error)
	}

	UserRepo interface {
		GetById(ctx context.Context, id int) (User, error)
		GetByLogin(ctx context.Context, login string) (User, error)
		Store(context.Context, User) (userID int, err error)
	}

	Hasher interface {
		GenerateSalt(ctx context.Context, size int) ([]byte, error)
		HashPassword(ctx context.Context, password string, salt []byte) (string, error)
		VerifyPassword(ctx context.Context, password string, salt []byte, hashedPassword string) (bool, error)
	}

	Authorizer interface {
		CreateToken(ctx context.Context, user User) (token string, err error)
		VerifyToken(ctx context.Context, token string) (userID int, err error)
	}

	Limiter interface {
		Take(ctx context.Context, userID int) error
	}
)

// Интерфейсы, связанные с управлением комиксами
type (
	ComicUseCase interface {
		Update(context.Context) error
		GetPictures(context.Context, string) ([]string, error)
	}

	ComicRepo interface {
		Store(context.Context, Comic) error
		GetById(context.Context, int) (Comic, error)
		GetAll(context.Context) ([]Comic, error)
		GetAllIds(context.Context) (IdMap, error)
	}

	ComicFetcher interface {
		GetNew(context.Context, IdMap) ([]ComicRaw, error)
	}
)
