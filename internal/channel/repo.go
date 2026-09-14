package channel

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewChannelPool(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) GetChRepo(ctx context.Context, serverID uuid.UUID) ([]Channel, error) {
	query := `SELECT c.id, c.server_id, c.name, c.created_at
        FROM channel c
        JOIN servers s ON c.server_id = s.id
        WHERE c.server_id = $1
        ORDER BY c.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, serverID)

	if err != nil {
		return nil, fmt.Errorf("get all apis: %w", err)
	}
	defer rows.Close()

	var channels []Channel

	for rows.Next() {
		var chn Channel

		if err := rows.Scan(
			&chn.ID,
			&chn.Server_id,
			&chn.Name,
			&chn.CreatedAT,
		); err != nil {
			return nil, fmt.Errorf("scan api: %w", err)
		}
		channels = append(channels, chn)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate apis: %w", err)
	}
	return channels, nil
}

func (r *Repo) CreateChRepo(ctx context.Context, serverID uuid.UUID, name string) (Channel, error) {
	query := "INSERT INTO channel (name, server_id) VALUES ($1, $2) RETURNING *"

	var channel Channel

	err := r.db.QueryRow(ctx, query, name, serverID).Scan(
		&channel.ID,
		&channel.Server_id,
		&channel.Name,
		&channel.CreatedAT,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Channel{}, fmt.Errorf("channel not found: %w", err)
		}
		return Channel{}, fmt.Errorf("query error: %w", err)
	}

	return channel, nil
}

func (r *Repo) DeleteChnRepo(ctx context.Context, serverID uuid.UUID, channelID uuid.UUID) (Channel, error) {
	query := `DELETE FROM channel WHERE server_id = $1 AND id = $2 RETURNING *`

	var chn Channel

	err := r.db.QueryRow(ctx, query, serverID, channelID).Scan(
		&chn.ID,
		&chn.Server_id,
		&chn.Name,
		&chn.CreatedAT,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Channel{}, fmt.Errorf("channel not found: %w", err)
		}
		return Channel{}, fmt.Errorf("delete error: %w", err)
	}

	return chn, nil
}
