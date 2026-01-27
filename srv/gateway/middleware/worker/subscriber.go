package worker

import (
	"context"
	"my-project/pkg/consts"
	"my-project/pkg/logger"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// WebSocketHub
type WebSocketHub interface {
	BroadcastToUser(userID string, message []byte)
}

type Subscriber struct {
	RDB *redis.Client
	Hub WebSocketHub
	Log *logrus.Logger
}

func NewSubscriber(rdb *redis.Client, hub WebSocketHub) *Subscriber {
	return &Subscriber{
		RDB: rdb,
		Hub: hub,
		Log: logger.Setup(), 
	}
}

// Start شروع به گوش دادن به Redis Pub/Sub می‌کند
func (s *Subscriber) Start(ctx context.Context) {
	pubsub := s.RDB.Subscribe(ctx, consts.TopicChatBroadcast)
	defer pubsub.Close()

	ch := pubsub.Channel()
	s.Log.Info("Gateway Subscriber Started: Listening for broadcast messages...")

	for msg := range ch {
		// پیام خام را می‌گیریم (که شامل ReceiverID و Content است)
		// فرض: پیام JSON است. باید Parse شود تا گیرنده مشخص شود.
		// اینجا برای سادگی فرض می‌کنیم Logic ارسال به Hub سپرده شده
		
		// s.Hub.BroadcastToUser(parsedMsg.ReceiverID, []byte(parsedMsg.Payload))
		s.Log.Debugf("Received message from Core: %s", msg.Payload)
	}
}