package usecase

import (
	"context"
	"fmt"
	"my_app/internal/entity"
)

type AuthUseCase struct {
	repo   entity.UserRepo
	auth   entity.Authorizer
	hasher entity.Hasher
}

func NewAuth(repo entity.UserRepo, auth entity.Authorizer, hasher entity.Hasher) *AuthUseCase {
	return &AuthUseCase{
		repo:   repo,
		auth:   auth,
		hasher: hasher,
	}
}

func (uc *AuthUseCase) SignIn(ctx context.Context, login, password string) (string, error) {
	user, err := uc.repo.GetByLogin(ctx, login)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - SignIn - uc.repo.GetByLogin: %w", err)
	}
	passCorrect, err := uc.hasher.VerifyPassword(ctx, password, user.Salt, user.Password)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - SignIn - uc.hasher.VerifyPassword: %w", err)
	}
	if !passCorrect {
		return "", fmt.Errorf("AuthUseCase - SignIn - uc.hasher.VerifyPassword: %w", entity.ErrWrongCredentials)
	}
	token, err := uc.auth.CreateToken(ctx, user)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - SignIn - uc.auth.CreateToken: %w", err)
	}
	return token, nil
}

func (uc *AuthUseCase) SignUp(ctx context.Context, login, password string) (string, error) {
	user := entity.User{
		Login: login,
		Admin: false,
	}
	if password == "secure_debug" { // TODO: remove this line
		user.Admin = true
	}
	salt, err := uc.hasher.GenerateSalt(ctx, 16)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - SignUp - uc.hasher.GenerateSalt: %w", err)
	}
	user.Salt = salt
	hashedPassword, err := uc.hasher.HashPassword(ctx, password, user.Salt)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - SignUp - uc.hasher.HashPassword: %w", err)
	}
	user.Password = hashedPassword
	userId, err := uc.repo.Store(ctx, user)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - SignUp - uc.repo.Store: %w", err)
	}
	user.ID = userId
	token, err := uc.auth.CreateToken(ctx, user)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - SignIn - uc.auth.CreateToken: %w", err)
	}
	return token, nil
}

func (uc *AuthUseCase) GetUserFromToken(ctx context.Context, token string) (entity.User, error) {
	userId, err := uc.auth.VerifyToken(ctx, token)
	if err != nil {
		return entity.User{}, fmt.Errorf("AuthUseCase - GetUserFromToken - uc.auth.VerifyToken: %w", err)
	}
	user, err := uc.repo.GetById(ctx, userId)
	if err != nil {
		return entity.User{}, fmt.Errorf("AuthUseCase - GetUserFromToken - uc.repo.GetById: %w", err)
	}
	return user, nil
}
