package core

import (
	"context"
	"encoding/json"
	"my-project/internal/adapter/postgres"
	"my-project/internal/adapter/redis"
	"my-project/internal/policy"
	"my-project/types"
	"github.com/google/uuid"
)

type ChatLogic struct {
	Repo   *postgres.DB //Repository Wrapper
	Queue  *redis.Client
	Policy *policy.ChatPolicy
}

func NewChatLogic(r *postgres.DB, q *redis.Client) *ChatLogic {
	return &ChatLogic{
		Repo:   r,
		Queue:  q,
	}
}

// ProcessIncomingMessage
func (l *ChatLogic) ProcessIncomingMessage(ctx context.Context, rawMsg []byte) error {
	var msg types.Message
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return err
	}

	_, err := l.Repo.Conn.ExecContext(ctx, 
		"INSERT INTO messages (id, sender_id, receiver_id, content, created_at) VALUES ($1, $2, $3, $4, $5)",
		msg.ID, msg.SenderID, msg.ReceiverID, msg.Content, msg.CreatedAt)
	
	if err != nil {
		return err
	}

	// 2.Redis Pub/Sub
	return l.Queue.Publish(ctx, "chat.broadcast", rawMsg)
}

func (l *ChatLogic) GetHistory(userID uuid.UUID) {
}