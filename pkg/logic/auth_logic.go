package logic

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"my-project/internal/auth"
	"my-project/pkg/repository"
	"my-project/pkg/utils/crypto"
	"my-project/types"
)

type AuthLogic struct {
	Repo      repository.UserRepository
	TokenMgr  *auth.JWTManager
}

func NewAuthLogic(repo repository.UserRepository, tokenMgr *auth.JWTManager) *AuthLogic {
	return &AuthLogic{
		Repo:     repo,
		TokenMgr: tokenMgr,
	}
}

func (a *AuthLogic) Register(ctx context.Context, req types.RegisterReq) (string, error) {
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

	err := a.Repo.Create(ctx, user)
	return mnemonic, err
}

func (a *AuthLogic) Login(ctx context.Context, username, password string) (string, error) {
	user, err := a.Repo.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if user == nil || crypto.HashString(password) != user.PasswordHash {
		return "", errors.New("invalid credentials")
	}

	return a.TokenMgr.Generate(user.ID)
}

func (a *AuthLogic) RecoverAccount(ctx context.Context, mnemonic string) (string, error) {
	mnemonicHash := crypto.HashString(mnemonic)
	
	user, err := a.Repo.GetByMnemonicHash(ctx, mnemonicHash)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid mnemonic")
	}

	return a.TokenMgr.Generate(user.ID)
}