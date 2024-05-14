package hasher

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

type IHasher interface {
	GenerateSalt(ctx context.Context, size int) ([]byte, error)
	HashPassword(ctx context.Context, password string, salt []byte) (string, error)
	VerifyPassword(ctx context.Context, password string, salt []byte, hashedPassword string) (bool, error)
}

type Hasher struct {
}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (h *Hasher) GenerateSalt(ctx context.Context, size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt[:])
	if err != nil {
		return nil, fmt.Errorf("Hasher - GenerateSalt - rand.Read: %w", err)
	}
	return salt, nil
}

func (h *Hasher) HashPassword(ctx context.Context, password string, salt []byte) (string, error) {
	sha := sha512.New()
	_, err := sha.Write(append([]byte(password), salt...))
	if err != nil {
		return "", fmt.Errorf("Hasher - HashPassword - sha512.Write: %w", err)
	}
	hash := sha.Sum(nil)
	return hex.EncodeToString(hash), nil
}

func (h *Hasher) VerifyPassword(ctx context.Context, password string, salt []byte, hashedPassword string) (bool, error) {
	curPassHash, err := h.HashPassword(ctx, password, salt)
	if err != nil {
		return false, fmt.Errorf("Hasher - VerifyPassword - h.HashPassword: %w", err)
	}
	return hashedPassword == curPassHash, nil
}
