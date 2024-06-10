package authorizer

import (
	"context"
	"testing"
	"time"

	"my_app/internal/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthorizer(t *testing.T) {
	tokenMaxTime := 1 * time.Hour
	secret := "your_secret"

	authorizer := NewAuthorizer(tokenMaxTime, secret)

	assert.Equal(t, tokenMaxTime, authorizer.tokenMaxTime)
	assert.Equal(t, []byte(secret), authorizer.secret)
	assert.Equal(t, jwt.SigningMethodHS256, authorizer.signMethod)
}

func TestCreateToken(t *testing.T) {
	ctx := context.Background()
	user := &entity.User{ID: 123}
	authorizer := NewAuthorizer(time.Hour, "your_secret")

	token, err := authorizer.CreateToken(ctx, *user)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestVerifyToken(t *testing.T) {
	ctx := context.Background()
	user := &entity.User{ID: 123}
	authorizer := NewAuthorizer(time.Hour, "your_secret")

	token, _ := authorizer.CreateToken(ctx, *user)

	userID, err := authorizer.VerifyToken(ctx, token)

	assert.NoError(t, err)
	assert.Equal(t, user.ID, userID)
}

func TestVerifyTokenError(t *testing.T) {
	ctx := context.Background()
	invalidToken := "invalid_token_string"
	authorizer := NewAuthorizer(time.Hour, "your_secret")

	_, err := authorizer.VerifyToken(ctx, invalidToken)

	assert.Error(t, err)
}
