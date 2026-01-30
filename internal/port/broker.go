package port

import "context"

type MessageBroker interface {
	Publish(ctx context.Context, channel string, msg interface{}) error
	PushQueue(ctx context.Context, queue string, msg []byte) error
	CacheMessage(ctx context.Context, key string, msg []byte) error
	GetRecentMessages(ctx context.Context, key string) ([]string, error)
	SetUserOnline(ctx context.Context, userID string) error
	SetUserOffline(ctx context.Context, userID string) error
	IsUserOnline(ctx context.Context, userID string) bool
}