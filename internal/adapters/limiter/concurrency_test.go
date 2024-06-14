package limiter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConcurLimiter(t *testing.T) {
	l := NewConcurLimiter(10)
	assert.NotNil(t, l)
	assert.Equal(t, 10, l.concurLimit)
}

func TestConcurTake(t *testing.T) {
	l := NewConcurLimiter(10)
	ctx := context.Background()

	err := l.Take(ctx, 0)
	assert.NoError(t, err)
}
