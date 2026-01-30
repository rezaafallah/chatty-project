package types

import "my-project/pkg/uid"

type User struct {
	ID           uid.ID    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	MnemonicHash string    `json:"-"`
	PublicKey    string    `json:"public_key"`
	CreatedAt    int64     `json:"created_at"`
}

type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=8"`
}