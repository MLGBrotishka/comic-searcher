package limiter

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	mu        *sync.RWMutex
	userLimit map[int]*rate.Limiter
	rateLimit int
}

func NewRateLimiter(rateLimit int) *RateLimiter {
	return &RateLimiter{
		mu:        &sync.RWMutex{},
		userLimit: make(map[int]*rate.Limiter),
		rateLimit: rateLimit,
	}
}

func (rl *RateLimiter) Take(ctx context.Context, userID int) error {
	limiter := rl.getLimiter(userID)
	limiter.Wait(ctx)
	return nil
}

func (rl *RateLimiter) getLimiter(userID int) *rate.Limiter {
	rl.mu.RLock()
	userLimiter, ok := rl.userLimit[userID]
	rl.mu.RUnlock()
	if !ok {
		rl.mu.Lock()
		userLimiter = rate.NewLimiter(rate.Limit(rl.rateLimit), rl.rateLimit)
		rl.userLimit[userID] = userLimiter
		rl.mu.Unlock()
	}
	return userLimiter
}
