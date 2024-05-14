package authorizer

import (
	"context"
	"fmt"
	"my_app/internal/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Authorizer struct {
	tokenMaxTime time.Duration
	secret       []byte
	signMethod   jwt.SigningMethod
}

func NewAuthorizer(tokenMaxTime time.Duration, secret string) *Authorizer {
	return &Authorizer{
		tokenMaxTime: tokenMaxTime,
		secret:       []byte(secret),
		signMethod:   jwt.SigningMethodHS256, //default
	}
}

func (a *Authorizer) CreateToken(ctx context.Context, user entity.User) (string, error) {
	payload := jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(a.tokenMaxTime).Unix(),
	}
	token := jwt.NewWithClaims(a.signMethod, payload)
	signedtoken, err := token.SignedString(a.secret)
	if err != nil {
		return "", fmt.Errorf("Authorizer - CreateToken - token.SignedString: %w", err)
	}
	return signedtoken, nil
}

func (a *Authorizer) VerifyToken(ctx context.Context, token string) (int, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
		return a.secret, nil
	})
	if err != nil {
		return 0, err
	}
	userID, ok := claims["userId"].(int)
	if !ok {
		return 0, entity.ErrInvalidToken
	}
	return int(userID), nil
}
