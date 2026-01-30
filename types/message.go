package types

import "my-project/pkg/uid"

type Message struct {
	ID         uid.ID `json:"id"`
	SenderID   uid.ID `json:"sender_id"`
	ReceiverID uid.ID `json:"receiver_id"`
	Content    string `json:"content"`
	CreatedAt  int64  `json:"created_at"`
	// New fields for Edit/Delete logic
	EditedAt  int64 `json:"edited_at,omitempty"`  // 0 means not edited
	DeletedAt int64 `json:"deleted_at,omitempty"` // 0 means not deleted
}

// Job represents a background task for Redis queue
type Job struct {
	Type    string `json:"type"`
	Payload []byte `json:"payload"`
}

// EditMessageReq is the payload for editing a message
type EditMessageReq struct {
	MessageID uid.ID `json:"message_id"`
	Content   string `json:"content"`
	UserID    uid.ID `json:"-"` // Filled from context/session
}

// DeleteMessageReq is the payload for deleting a message
type DeleteMessageReq struct {
	MessageID uid.ID `json:"message_id"`
	UserID    uid.ID `json:"-"` // Filled from context/session
}