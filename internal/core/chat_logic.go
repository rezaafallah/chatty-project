package core

import (
	"context"
	"encoding/json"
	"my-project/internal/port"
	"my-project/internal/repository"
	"my-project/types"
	"errors"
	"github.com/google/uuid"
)

type ChatLogic struct {
	Repo   repository.MessageRepository
	Broker port.MessageBroker
}

func NewChatLogic(repo repository.MessageRepository, broker port.MessageBroker) *ChatLogic {
	return &ChatLogic{Repo: repo, Broker: broker}
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

	if msg.Content == "" {
        return errors.New("message content cannot be empty")
    }

	// 
	// key: history:USER_ID
	senderKey := "history:" + msg.SenderID.String()
	receiverKey := "history:" + msg.ReceiverID.String()
	_ = l.Broker.CacheMessage(ctx, senderKey, rawMsg)
	_ = l.Broker.CacheMessage(ctx, receiverKey, rawMsg)

	return l.Broker.Publish(ctx, "chat.broadcast", rawMsg)
}

// GetHistory:
func (l *ChatLogic) GetHistory(userID uuid.UUID) ([]types.Message, error) {
	ctx := context.Background()
	key := "history:" + userID.String()

	rawMsgs, err := l.Broker.GetRecentMessages(ctx, key)
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