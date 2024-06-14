package limiter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRateLimiter(t *testing.T) {
	l := NewRateLimiter(10)
	assert.NotNil(t, l)
	assert.Equal(t, 10, l.rateLimit)
}

func TestRateTake(t *testing.T) {
	l := NewRateLimiter(10)
	ctx := context.Background()

	err := l.Take(ctx, 1)
	assert.NoError(t, err)
}
