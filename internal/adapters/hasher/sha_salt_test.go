package hasher

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSalt(t *testing.T) {
	ctx := context.Background()
	size := 16 // Example size, adjust based on your requirements
	salt, err := NewHasher().GenerateSalt(ctx, size)
	assert.NoError(t, err)
	assert.Equal(t, size, len(salt))
}

func TestHashPassword(t *testing.T) {
	ctx := context.Background()
	password := "testPassword123"
	salt := []byte{0x01, 0x02, 0x03} // Example salt, adjust as needed
	hashedPassword, err := NewHasher().HashPassword(ctx, password, salt)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
}

func TestVerifyPassword(t *testing.T) {
	ctx := context.Background()
	password := "testPassword123"
	salt := []byte{0x01, 0x02, 0x03}
	hashedPassword, _ := NewHasher().HashPassword(ctx, password, salt)
	isValid, err := NewHasher().VerifyPassword(ctx, password, salt, hashedPassword)
	assert.NoError(t, err)
	assert.True(t, isValid)
}
