package logic

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"my-project/pkg/broker"
	"my-project/pkg/consts"
	"my-project/pkg/logger"
	"my-project/pkg/repository"
	"my-project/pkg/uid"
	"my-project/types"
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

// ProcessIncomingMessage handles saving and broadcasting a NEW message.
func (l *ChatLogic) ProcessIncomingMessage(ctx context.Context, rawMsg []byte) error {
	var msg types.Message
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		l.Log.WithError(err).Error("Unmarshal failed in logic")
		return err
	}

	if msg.Content == "" {
		return errors.New(consts.ErrEmptyContent)
	}

	// Generate a new ID if one wasn't provided (usually it isn't)
	if msg.ID == uid.Nil {
		msg.ID = uid.New()
	}

	// 1. Save to Database
	err := l.Repo.Save(ctx, msg)
	if err != nil {
		l.Log.WithError(err).Error("Database Save failed")
		return err
	}

	// 2. Cache in Redis (for history)
	updatedMsgBytes, _ := json.Marshal(msg)
	senderKey := repository.HistoryKey(msg.SenderID.String())
	receiverKey := repository.HistoryKey(msg.ReceiverID.String())

	// Log warnings on cache failure, but don't stop the flow
	if err := l.Broker.CacheMessage(ctx, senderKey, updatedMsgBytes); err != nil {
		l.Log.WithError(err).Warn("Failed to cache message for sender")
	}
	if err := l.Broker.CacheMessage(ctx, receiverKey, updatedMsgBytes); err != nil {
		l.Log.WithError(err).Warn("Failed to cache message for receiver")
	}

	// 3. Broadcast Event (Wrap in Job/Envelope)
	// We wrap the message in a 'Job' struct so the consumer knows this is a "New Message" event.
	job := types.Job{
		Type:    consts.EventMessageNew,
		Payload: updatedMsgBytes,
	}
	
	jobBytes, _ := json.Marshal(job)

	// Publish to the global broadcast topic
	return l.Broker.Publish(ctx, consts.TopicChatBroadcast, jobBytes)
}

// EditMessage handles updating the content of an existing message.
func (l *ChatLogic) EditMessage(ctx context.Context, req types.EditMessageReq) error {
	// 1. Retrieve the original message to check ownership
	msg, err := l.Repo.GetByID(ctx, req.MessageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.New("message not found")
	}

	// 2. Ownership & Validity Checks
	if msg.SenderID != req.UserID {
		return errors.New("unauthorized: you can only edit your own messages")
	}

	if msg.DeletedAt > 0 {
		return errors.New("cannot edit a deleted message")
	}

	// 3. Update Database
	if err := l.Repo.UpdateContent(ctx, req.MessageID, req.Content); err != nil {
		l.Log.WithError(err).Error("Failed to update message content in DB")
		return err
	}

	// 4. Broadcast Event (Edit Message)
	// We construct a temporary message struct with the updated data for the client
	msg.Content = req.Content
	msg.EditedAt = time.Now().Unix()

	payloadBytes, _ := json.Marshal(msg)
	
	job := types.Job{
		Type:    consts.EventMessageEdit,
		Payload: payloadBytes,
	}
	jobBytes, _ := json.Marshal(job)

	return l.Broker.Publish(ctx, consts.TopicChatBroadcast, jobBytes)
}

// DeleteMessage handles soft-deleting a message.
func (l *ChatLogic) DeleteMessage(ctx context.Context, req types.DeleteMessageReq) error {
	// 1. Retrieve the original message
	msg, err := l.Repo.GetByID(ctx, req.MessageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.New("message not found")
	}

	// 2. Ownership Check
	if msg.SenderID != req.UserID {
		return errors.New("unauthorized: you can only delete your own messages")
	}

	// 3. Soft Delete in Database
	if err := l.Repo.SoftDelete(ctx, req.MessageID); err != nil {
		l.Log.WithError(err).Error("Failed to soft delete message in DB")
		return err
	}

	// 4. Broadcast Event (Delete Message)
	// Clients just need the ID to remove it from their UI
	deletePayload := map[string]interface{}{
		"id": req.MessageID,
	}
	payloadBytes, _ := json.Marshal(deletePayload)

	job := types.Job{
		Type:    consts.EventMessageDelete,
		Payload: payloadBytes,
	}
	jobBytes, _ := json.Marshal(job)

	return l.Broker.Publish(ctx, consts.TopicChatBroadcast, jobBytes)
}

// GetHistory fetches recent messages from Redis (or DB if we implemented fallback).
func (l *ChatLogic) GetHistory(userID uid.ID) ([]types.Message, error) {
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