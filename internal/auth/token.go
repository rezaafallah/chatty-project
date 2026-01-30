package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"my-project/pkg/uid"
)

// JWTManager responsible for generating tokens
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

// Generate creates a signed JWT token
func (m *JWTManager) Generate(userID uid.ID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(m.Expiry).Unix(),
	})
	return token.SignedString([]byte(m.SecretKey))
}