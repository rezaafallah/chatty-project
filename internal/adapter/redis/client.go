package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
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