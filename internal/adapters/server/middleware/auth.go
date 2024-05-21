package middleware

import (
	"context"
	"my_app/internal/entity"
	"net/http"
)

type AuthMiddleware struct {
	uc entity.AuthUseCase
}

func NewAuthMiddleware(authUseCase entity.AuthUseCase) *AuthMiddleware {
	return &AuthMiddleware{
		uc: authUseCase,
	}
}

func (am *AuthMiddleware) Auth(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user, err := am.uc.GetUserFromToken(ctx, req.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx = context.WithValue(ctx, "UserID", user.ID)
		ctx = context.WithValue(ctx, "Role", user.IsAdmin())
		f(w, req.WithContext(ctx))
	})
}
