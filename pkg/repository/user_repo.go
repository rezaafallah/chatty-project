package repository

import (
	"context"
	"database/sql"
	"my-project/types"
)

type UserRepository interface {
	Create(ctx context.Context, user types.User) error
	GetByUsername(ctx context.Context, username string) (*types.User, error)
	GetByMnemonicHash(ctx context.Context, hash string) (*types.User, error)
}

//Postgres
type UserRepositoryPostgres struct {
	Conn *sql.DB
}

func NewUserRepository(conn *sql.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{Conn: conn}
}

func (r *UserRepositoryPostgres) Create(ctx context.Context, user types.User) error {
	_, err := r.Conn.ExecContext(ctx,
		"INSERT INTO users (id, username, password_hash, mnemonic_hash, public_key, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		user.ID, user.Username, user.PasswordHash, user.MnemonicHash, user.PublicKey, user.CreatedAt,
	)
	return err
}

func (r *UserRepositoryPostgres) GetByUsername(ctx context.Context, username string) (*types.User, error) {
	var user types.User
	row := r.Conn.QueryRowContext(ctx, "SELECT id, password_hash FROM users WHERE username = $1", username)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // user not found
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryPostgres) GetByMnemonicHash(ctx context.Context, hash string) (*types.User, error) {
	var user types.User
	row := r.Conn.QueryRowContext(ctx, "SELECT id FROM users WHERE mnemonic_hash = $1", hash)
	err := row.Scan(&user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}