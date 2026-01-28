package core

import (
	"context"
	"encoding/json"
	"my-project/internal/adapter/redis"
	"my-project/internal/repository"
	"my-project/types"
	"github.com/google/uuid"
)

type ChatLogic struct {
	Repo   repository.MessageRepository
	Queue  *redis.Client
}

func NewChatLogic(repo repository.MessageRepository, q *redis.Client) *ChatLogic {
	return &ChatLogic{Repo: repo, Queue: q}
}

func (l *ChatLogic) ProcessIncomingMessage(ctx context.Context, rawMsg []byte) error {
	var msg types.Message
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return err
	}

	err := l.Repo.Save(ctx, msg)
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