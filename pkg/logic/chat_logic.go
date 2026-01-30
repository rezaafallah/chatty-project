package logic

import (
	"context"
	"encoding/json"
	"my-project/pkg/broker"
	"my-project/pkg/repository"
	"my-project/types"
	"errors"
	"github.com/google/uuid"
	"my-project/pkg/logger"
	"my-project/pkg/consts"
)

	type ChatLogic struct {
		Repo   repository.MessageRepository
		Broker broker.MessageBroker
		Log    logger.Logger
	}

	func NewChatLogic(repo repository.MessageRepository, broker broker.MessageBroker, log logger.Logger) *ChatLogic {
		return &ChatLogic{
			Repo:   repo,
			Broker: broker,
			Log:    log,
		}
	}

	func (l *ChatLogic) ProcessIncomingMessage(ctx context.Context, rawMsg []byte) error {
	var msg types.Message
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		l.Log.WithError(err).Error("Unmarshal failed in logic") // Log error
		return err
	}

	if msg.Content == "" {
		return errors.New(consts.ErrEmptyContent)
	}

	if msg.ID == uuid.Nil {
		msg.ID = uuid.New()
	}

	err := l.Repo.Save(ctx, msg)
	if err != nil {
		l.Log.WithError(err).Error("Database Save failed")
		return err
	}

	// Cache & Broadcast
	updatedMsgBytes, _ := json.Marshal(msg)

	senderKey := repository.HistoryKey(msg.SenderID.String())
	receiverKey := repository.HistoryKey(msg.ReceiverID.String())
	
	// Don't ignore errors (Log them as Warnings)
	if err := l.Broker.CacheMessage(ctx, senderKey, updatedMsgBytes); err != nil {
		l.Log.WithError(err).Warn("Failed to cache message for sender")
	}
	if err := l.Broker.CacheMessage(ctx, receiverKey, updatedMsgBytes); err != nil {
		l.Log.WithError(err).Warn("Failed to cache message for receiver")
	}

	return l.Broker.Publish(ctx, consts.TopicChatBroadcast, updatedMsgBytes)
}

	// GetHistory:
	func (l *ChatLogic) GetHistory(userID uuid.UUID) ([]types.Message, error) {
	ctx := context.Background()
	key := repository.HistoryKey(userID.String())

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