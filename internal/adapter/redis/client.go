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