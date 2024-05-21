package middleware

import (
	"my_app/internal/entity"
	"net/http"
)

type LimiterMiddleware struct {
	limiter entity.Limiter
}

func NewLimiterMiddleware(limiter entity.Limiter) *LimiterMiddleware {
	return &LimiterMiddleware{
		limiter: limiter,
	}
}

func (lm *LimiterMiddleware) WaitLimit(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		userID := ctx.Value("UserID").(int)
		if userID == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		lm.limiter.Take(ctx, userID)

		f(w, req)
	})
}
