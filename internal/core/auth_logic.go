package core

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"my-project/internal/adapter/postgres"
	"my-project/internal/crypto"
	"my-project/types"
	"github.com/google/uuid"
)

type AuthLogic struct {
	DB        *postgres.DB
	JWTSecret string 
}

func NewAuthLogic(db *postgres.DB, secret string) *AuthLogic {
	return &AuthLogic{
		DB:        db,
		JWTSecret: secret,
	}
}

//Register
func (a *AuthLogic) Register(req types.RegisterReq) (string, error) {
    mnemonic, _ := crypto.GenerateMnemonic()
    pubKey, _ := crypto.GenerateKeyPair(mnemonic)
    
    user := types.User{
        ID:           uuid.New(),
        Username:     req.Username,
        PasswordHash: crypto.HashString(req.Password),
        MnemonicHash: crypto.HashString(mnemonic),
        PublicKey:    pubKey,
        CreatedAt:    time.Now().Unix(),
    }

    _, err := a.DB.Conn.Exec(
        "INSERT INTO users (id, username, password_hash, mnemonic_hash, public_key, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
        user.ID, user.Username, user.PasswordHash, user.MnemonicHash, user.PublicKey, user.CreatedAt,
    )
    return mnemonic, err
}

//Login
func (a *AuthLogic) Login(username, password string) (string, error) {
	var user types.User
	row := a.DB.Conn.QueryRow("SELECT id, password_hash FROM users WHERE username = $1", username)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}

	if crypto.HashString(password) != user.PasswordHash {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(a.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *AuthLogic) RecoverAccount(mnemonic string) (string, error) {
	mnemonicHash := crypto.HashString(mnemonic)
	var user types.User

	row := a.DB.Conn.QueryRow("SELECT id FROM users WHERE mnemonic_hash = $1", mnemonicHash)
	err := row.Scan(&user.ID)
	if err != nil {
		return "", errors.New("invalid mnemonic")
	}

	return a.generateToken(user.ID)
}

func (a *AuthLogic) generateToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(a.JWTSecret))
}