package middleware

import (
	"net/http"
)

type LimiterMiddleware struct {
	limiter Limiter
}

func NewLimiterMiddleware(limiter Limiter) *LimiterMiddleware {
	return &LimiterMiddleware{
		limiter: limiter,
	}
}

func (lm *LimiterMiddleware) WaitLimit(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		userID := ctx.Value(userIDKey("UserID")).(int)
		if userID == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		_ = lm.limiter.Take(ctx, userID)

		f(w, req)
	})
}
