package core

import (
	"my-project/internal/adapter/postgres"
	"my-project/internal/crypto"
	"my-project/types"
	"github.com/google/uuid"
	"time"
)

type AuthLogic struct {
	DB *postgres.DB
}

func NewAuthLogic(db *postgres.DB) *AuthLogic {
	return &AuthLogic{DB: db}
}

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