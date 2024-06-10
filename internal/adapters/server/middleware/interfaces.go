package middleware

import (
	"context"
	"my_app/internal/entity"
)

type (
	AuthUseCase interface {
		GetUserFromToken(ctx context.Context, token string) (user entity.User, err error)
	}

	Limiter interface {
		Take(ctx context.Context, userID int) error
	}
)

type userIDKey string
type roleKey string
