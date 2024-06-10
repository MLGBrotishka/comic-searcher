package middleware

import (
	"context"
	"net/http"
)

type AuthMiddleware struct {
	uc AuthUseCase
}

func NewAuthMiddleware(authUseCase AuthUseCase) *AuthMiddleware {
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
		ctx = context.WithValue(ctx, userIDKey("UserID"), user.ID)
		ctx = context.WithValue(ctx, roleKey("Role"), user.IsAdmin())
		f(w, req.WithContext(ctx))
	})
}
