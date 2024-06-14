package server

import (
	"context"
	"my_app/internal/entity"
)

type (
	AuthUseCase interface {
		SignIn(ctx context.Context, login string, password string) (token string, err error)
		SignUp(ctx context.Context, login string, password string) (token string, err error)
		GetUserFromToken(ctx context.Context, token string) (user entity.User, err error)
	}

	ComicUseCase interface {
		Update(context.Context) error
		GetPictures(context.Context, string) ([]string, error)
	}
)
