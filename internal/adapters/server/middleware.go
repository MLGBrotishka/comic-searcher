package server

import (
	"net/http"

	"my_app/internal/usecase"
)

type AdminMiddleware struct {
	Next        http.Handler
	AuthUseCase usecase.Auth
}

func NewAdminMiddleware(next http.Handler, authUseCase usecase.Auth) *AdminMiddleware {
	return &AdminMiddleware{
		Next:        next,
		AuthUseCase: authUseCase,
	}
}

func (am *AdminMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := am.AuthUseCase.GetUserFromToken(ctx, r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !user.IsAdmin() {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	am.Next.ServeHTTP(w, r)
}
