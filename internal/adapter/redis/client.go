package redis

import (
	"context"
	"time"
	"github.com/redis/go-redis/v9"
	"my-project/pkg/repository"
)

type Client struct {
	RDB *redis.Client
}

func New(addr string) *Client {
	return &Client{
		RDB: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

func (c *Client) Publish(ctx context.Context, channel string, msg interface{}) error {
	return c.RDB.Publish(ctx, channel, msg).Err()
}

func (c *Client) PushQueue(ctx context.Context, queue string, msg []byte) error {
	return c.RDB.RPush(ctx, queue, msg).Err()
}

// --- CacheMessage: ---

func (c *Client) CacheMessage(ctx context.Context, key string, msg []byte) error {
	pipe := c.RDB.Pipeline()
	pipe.LPush(ctx, key, msg)
	pipe.LTrim(ctx, key, 0, 49)   // 50 msg
	_, err := pipe.Exec(ctx)
	return err
}

// GetRecentMessages
func (c *Client) GetRecentMessages(ctx context.Context, key string) ([]string, error) {
	return c.RDB.LRange(ctx, key, 0, -1).Result()
}

// Presence Methods

func (c *Client) SetUserOnline(ctx context.Context, userID string) error {
	key := repository.OnlineKey(userID)
	return c.RDB.Set(ctx, key, "1", 60*time.Second).Err()
}

func (c *Client) SetUserOffline(ctx context.Context, userID string) error {
	key := repository.OnlineKey(userID)
	return c.RDB.Del(ctx, key).Err()
}

func (c *Client) IsUserOnline(ctx context.Context, userID string) bool {
	key := repository.OnlineKey(userID)
	exists, _ := c.RDB.Exists(ctx, key).Result()
	return exists > 0
}