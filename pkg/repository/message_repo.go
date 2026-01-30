package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"my-project/pkg/uid"
	"my-project/types"
)

type MessageRepository interface {
	Save(ctx context.Context, msg types.Message) error
	GetHistory(ctx context.Context, u1, u2 uid.ID) ([]types.Message, error)
	GetByID(ctx context.Context, id uid.ID) (*types.Message, error)
	UpdateContent(ctx context.Context, id uid.ID, content string) error
	SoftDelete(ctx context.Context, id uid.ID) error
}

type messageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Save(ctx context.Context, msg types.Message) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO messages (id, sender_id, receiver_id, content, created_at, edited_at, deleted_at) VALUES ($1, $2, $3, $4, $5, 0, 0)",
		msg.ID, msg.SenderID, msg.ReceiverID, msg.Content, msg.CreatedAt,
	)
	return err
}

func (r *messageRepository) GetHistory(ctx context.Context, u1, u2 uid.ID) ([]types.Message, error) {
	// Only fetch messages that are NOT deleted (deleted_at = 0)
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, sender_id, receiver_id, content, created_at, edited_at FROM messages WHERE ((sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)) AND deleted_at = 0 ORDER BY created_at ASC",
		u1, u2,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []types.Message
	for rows.Next() {
		var m types.Message
		if err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.Content, &m.CreatedAt, &m.EditedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}

// GetByID fetches a single message to check ownership
func (r *messageRepository) GetByID(ctx context.Context, id uid.ID) (*types.Message, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, sender_id, receiver_id, content, created_at, deleted_at FROM messages WHERE id = $1", id)
	
	var m types.Message
	err := row.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.Content, &m.CreatedAt, &m.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &m, nil
}

func (r *messageRepository) UpdateContent(ctx context.Context, id uid.ID, content string) error {
	// Update content and set edited_at timestamp
	_, err := r.db.ExecContext(ctx, "UPDATE messages SET content = $1, edited_at = $2 WHERE id = $3", content, time.Now().Unix(), id)
	return err
}

func (r *messageRepository) SoftDelete(ctx context.Context, id uid.ID) error {
	// Set deleted_at timestamp instead of removing the row
	_, err := r.db.ExecContext(ctx, "UPDATE messages SET deleted_at = $1 WHERE id = $2", time.Now().Unix(), id)
	return err
}