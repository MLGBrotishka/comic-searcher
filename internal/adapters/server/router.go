package server

import (
	"net/http"

	"my_app/internal/adapters/limiter"
	"my_app/internal/adapters/server/middleware"
	"my_app/internal/entity"
	"my_app/pkg/logger"
)

func NewRouter(handler *http.ServeMux, uc entity.ComicUseCase, auth entity.AuthUseCase, concurLimit int, rateLimit int, l logger.Interface) {
	rateLimiter := middleware.NewLimiterMiddleware(limiter.NewRateLimiter(rateLimit))
	concurLimiter := middleware.NewLimiterMiddleware(limiter.NewConcurLimiter(concurLimit))
	authMiddleware := middleware.NewAuthMiddleware(auth)
	mid := func(f http.HandlerFunc) http.HandlerFunc {
		return authMiddleware.Auth(concurLimiter.WaitLimit(rateLimiter.WaitLimit(f)))
	}
	newComicRoutes(handler, uc, mid, l)
	mid2 := func(f http.HandlerFunc) http.HandlerFunc {
		return f
	}
	newAuthRoutes(handler, auth, mid2, l)

}
