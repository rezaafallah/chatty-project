package worker

import (
	"context"
	"encoding/json" // اضافه شد
	"my-project/pkg/consts"
	"my-project/pkg/logger"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

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

// ساختار پیام برای Decode کردن JSON
type BroadcastMsg struct {
	ReceiverID string `json:"receiver_id"`
	Payload    string `json:"payload"`
}

func (s *Subscriber) Start(ctx context.Context) {
	pubsub := s.RDB.Subscribe(ctx, consts.TopicChatBroadcast)
	defer pubsub.Close()

	ch := pubsub.Channel()
	s.Log.Info("Gateway Subscriber Started...")

	for msg := range ch {
		// 1. آنمارشال کردن پیام برای پیدا کردن گیرنده
		var bMsg BroadcastMsg
		// فرض: Core پیام را کامل (شامل ID گیرنده) می‌فرستد
		// اگر فرمت دیتای Core چیز دیگری است، اینجا باید هندل شود
		// فعلا فرض می‌کنیم Core کل آبجکت Message را JSON کرده فرستاده
		if err := json.Unmarshal([]byte(msg.Payload), &bMsg); err != nil {
			s.Log.Error("Failed to unmarshal broadcast msg", err)
			continue
		}

		// 2. ارسال به کاربر خاص از طریق Hub
		// نکته: اینجا Payload اصلی را دوباره به کلاینت می‌فرستیم
		s.Hub.BroadcastToUser(bMsg.ReceiverID, []byte(msg.Payload))
		
		s.Log.Debugf("Relayed message to user %s", bMsg.ReceiverID)
	}
}