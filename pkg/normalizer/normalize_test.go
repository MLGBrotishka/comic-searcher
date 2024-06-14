package normalizer

import (
	"context"
	"testing"

	"my_app/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected entity.KeywordMap
	}{
		{
			name:  "Test simple normalization",
			input: "I'm going to the park.",
			expected: map[string]bool{
				"go":   true,
				"park": true,
				"i":    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := New()
			result, err := n.Normalize(context.Background(), tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
