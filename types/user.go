// types/user.go
package types

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	MnemonicHash string    `json:"-"`
	PublicKey    string    `json:"public_key"`
	CreatedAt    int64     `json:"created_at"` // Unix Time
}

type Message struct {
	ID         uuid.UUID `json:"id"`
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Payload    string    `json:"payload"` // Base64 Encrypted Content
	CreatedAt  int64     `json:"created_at"`
}