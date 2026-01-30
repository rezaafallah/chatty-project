package worker

import (
	"context"
	"encoding/json"
	"my-project/pkg/consts"
	"my-project/pkg/logger"
	"my-project/types"
	"github.com/redis/go-redis/v9"
)

type WebSocketHub interface {
	BroadcastToUser(userID string, message []byte)
}

type Subscriber struct {
	RDB *redis.Client
	Hub WebSocketHub
	Log logger.Logger
}

func NewSubscriber(rdb *redis.Client, hub WebSocketHub) *Subscriber {
	return &Subscriber{
		RDB: rdb,
		Hub: hub,
		Log: logger.Setup(), 
	}
}

func (s *Subscriber) Start(ctx context.Context) {
	pubsub := s.RDB.Subscribe(ctx, consts.TopicChatBroadcast)
	defer pubsub.Close()

	ch := pubsub.Channel()
	s.Log.Info("Gateway Subscriber Started...")

	for msg := range ch {
		// 1. Unmarshal into Domain Message to check Receiver
		var domainMsg types.Message
		
		if err := json.Unmarshal([]byte(msg.Payload), &domainMsg); err != nil {
			s.Log.Error("Failed to unmarshal broadcast msg", err)
			continue
		}

		// 2. Send to specific user via Hub
		// We send the full payload back to the client so they see sender_id etc.
		s.Hub.BroadcastToUser(domainMsg.ReceiverID.String(), []byte(msg.Payload))
		
		s.Log.Debugf("Relayed message to user %s", domainMsg.ReceiverID)
	}
}