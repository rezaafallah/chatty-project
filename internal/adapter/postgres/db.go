package core

import (
	"context"
	"encoding/json"
	"my-project/internal/adapter/postgres"
	"my-project/internal/adapter/redis"
	"my-project/types"
	"github.com/sirupsen/logrus"
)

// Logic
type Logic struct {
	DB    *postgres.DB
	Redis *redis.Client
	Log   *logrus.Logger
}

func NewLogic(db *postgres.DB, r *redis.Client, l *logrus.Logger) *Logic {
	return &Logic{DB: db, Redis: r, Log: l}
}

// ProcessMessage 
func (l *Logic) ProcessMessage(ctx context.Context, rawMsg []byte) error {
	var msg types.Message
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return err
	}

	// 1. Logic: Save to DB (Persistent)
	// Manual SQL Insert
	_, err := l.DB.Conn.ExecContext(ctx, "INSERT INTO messages ...", msg.ID, msg.Payload)
	if err != nil {
		l.Log.Error("DB Save Failed", err)
		return err
	}

	// 2. Logic: Publish to PubSub (Realtime delivery)
	// Gateway will listen to this channel
	return l.Redis.Publish(ctx, "chat.broadcast", rawMsg)
}