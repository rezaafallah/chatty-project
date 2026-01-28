package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenManager interface {
	Generate(userID uuid.UUID) (string, error)
}

type JWTManager struct {
	SecretKey string
	Expiry    time.Duration
}

func NewJWTManager(secret string, expiry time.Duration) *JWTManager {
	return &JWTManager{
		SecretKey: secret,
		Expiry:    expiry,
	}
}

func (m *JWTManager) Generate(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(m.Expiry).Unix(),
	})
	return token.SignedString([]byte(m.SecretKey))
}