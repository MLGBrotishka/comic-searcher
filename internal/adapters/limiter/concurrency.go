package limiter

import (
	"context"

	"golang.org/x/time/rate"
)

type ConcurLimiter struct {
	limiter     *rate.Limiter
	concurLimit int
}

func NewConcurLimiter(concurLimit int) *ConcurLimiter {
	return &ConcurLimiter{
		limiter:     rate.NewLimiter(rate.Limit(concurLimit), concurLimit),
		concurLimit: concurLimit,
	}
}

func (cl *ConcurLimiter) Take(ctx context.Context, _ int) error {
	cl.limiter.Wait(ctx)
	return nil
}
