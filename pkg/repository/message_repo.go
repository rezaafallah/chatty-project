package repository

import (
	"context"
	"database/sql"
	"my-project/types"
)

type MessageRepository interface {
	Save(ctx context.Context, msg types.Message) error
}

type MessageRepositoryPostgres struct {
	Conn *sql.DB
}

func NewMessageRepository(conn *sql.DB) *MessageRepositoryPostgres {
	return &MessageRepositoryPostgres{Conn: conn}
}

func (r *MessageRepositoryPostgres) Save(ctx context.Context, msg types.Message) error {
	_, err := r.Conn.ExecContext(ctx,
		"INSERT INTO messages (id, sender_id, receiver_id, content, created_at) VALUES ($1, $2, $3, $4, $5)",
		msg.ID, msg.SenderID, msg.ReceiverID, msg.Content, msg.CreatedAt)
	return err
}