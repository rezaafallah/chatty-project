package types

import "github.com/google/uuid"

type Message struct {
	ID         uuid.UUID `json:"id"`
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Content    string    `json:"content"` // Encrypted Base64
	CreatedAt  int64     `json:"created_at"`
}