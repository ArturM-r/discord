package message

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewMessagePool(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

func (m *Repo) Get(ctx context.Context, channelID uuid.UUID, userID uuid.UUID) ([]Message, error) {
	query := `
	SELECT m.id, m.channel_id, u.username, m.user_id, m.content, m.created_at 
	FROM messages m
	JOIN users u ON m.user_id = u.id
	JOIN channels c ON m.channel_id = c.id
	JOIN members mb ON c.server_id = mb.server_id AND mb.user_id = $2
	WHERE m.channel_id = $1
	ORDER BY m.created_at DESC
	LIMIT 50
	`

	rows, err := m.db.Query(ctx, query, channelID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []Message

	for rows.Next() {
		var msg Message

		if err := rows.Scan(
			&msg.ID,
			&msg.ChannelID,
			&msg.UserName,
			&msg.UserID,
			&msg.Content,
			&msg.CreatedAT,
		); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate messages: %w", err)
	}
	return messages, nil
}

func (m *Repo) WriteMessage(ctx context.Context, channelID uuid.UUID, userID uuid.UUID, content string) (Message, error) {
	query := `
	WITH inserted AS (
		INSERT INTO messages (channel_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, channel_id, user_id, content, created_at
	)
	SELECT i.id, i.channel_id, i.user_id, u.username, i.content, i.created_at
	FROM inserted i
	JOIN users u ON i.user_id = u.id
	`

	var msg Message
	err := m.db.QueryRow(ctx, query, channelID, userID, content).Scan(
		&msg.ID,
		&msg.ChannelID,
		&msg.UserID,
		&msg.UserName,
		&msg.Content,
		&msg.CreatedAT,
	)
	if err != nil {
		return Message{}, fmt.Errorf("failed to write message: %w", err)
	}
	return msg, nil
}
