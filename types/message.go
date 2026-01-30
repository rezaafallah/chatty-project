package types

import "my-project/pkg/uid"

type Message struct {
	ID         uid.ID `json:"id"`
	SenderID   uid.ID `json:"sender_id"`
	ReceiverID uid.ID `json:"receiver_id"`
	Content    string `json:"content"`
	CreatedAt  int64  `json:"created_at"`
}