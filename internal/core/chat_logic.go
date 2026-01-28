package core

import (
	"context"
	"encoding/json"
	"my-project/internal/adapter/postgres"
	"my-project/internal/adapter/redis"
	"my-project/types"
	"github.com/google/uuid"
)

type ChatLogic struct {
	Repo   *postgres.DB
	Queue  *redis.Client
}

func NewChatLogic(r *postgres.DB, q *redis.Client) *ChatLogic {
	return &ChatLogic{Repo: r, Queue: q}
}

func (l *ChatLogic) ProcessIncomingMessage(ctx context.Context, rawMsg []byte) error {
	var msg types.Message
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return err
	}

	// 1.(Permanent Storage)
	_, err := l.Repo.Conn.ExecContext(ctx, 
		"INSERT INTO messages (id, sender_id, receiver_id, content, created_at) VALUES ($1, $2, $3, $4, $5)",
		msg.ID, msg.SenderID, msg.ReceiverID, msg.Content, msg.CreatedAt)
	if err != nil {
		return err
	}

	// 
	// key: history:USER_ID
	senderKey := "history:" + msg.SenderID.String()
	receiverKey := "history:" + msg.ReceiverID.String()

	_ = l.Queue.CacheMessage(ctx, senderKey, rawMsg)
	_ = l.Queue.CacheMessage(ctx, receiverKey, rawMsg)

	return l.Queue.Publish(ctx, "chat.broadcast", rawMsg)
}

// GetHistory:
func (l *ChatLogic) GetHistory(userID uuid.UUID) ([]types.Message, error) {
	ctx := context.Background()
	key := "history:" + userID.String()

	rawMsgs, err := l.Queue.GetRecentMessages(ctx, key)
	if err != nil {
		return nil, err
	}

	var messages []types.Message
	for _, raw := range rawMsgs {
		var msg types.Message
		if err := json.Unmarshal([]byte(raw), &msg); err == nil {
			messages = append(messages, msg)
		}
	}
	
	return messages, nil
}