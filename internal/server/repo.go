package server

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

func NewMessagePool(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) CreateSRV(ctx context.Context, userID uuid.UUID, name string) (Server, error) {
	tx, err := r.db.Begin(ctx)

	if err != nil {
		return Server{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO server (name, owner_id) VALUES ($1, $2) RETURNING *`

	var srv Server

	err = tx.QueryRow(ctx, query, name, userID).Scan(
		&srv.ID,
		&srv.Name,
		&srv.OwnerID,
	)
	if err != nil {
		return Server{}, fmt.Errorf("failed to create server: %w", err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO member (server_id, user_id, role) VALUES ($1, $2, $3)`, srv.ID, srv.OwnerID, "owner")

	if err != nil {
		return Server{}, fmt.Errorf("failed to insert:%w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return Server{}, fmt.Errorf("failed to commit: %w", err)
	}

	return srv, nil
}

func (r *Repo) GetSrvUser(ctx context.Context, userID uuid.UUID) ([]Server, error) {
	query := `SELECT 
    s.id, 
    s.name, 
    s.owner_id,
	s.created_at
	FROM servers s
	JOIN members m ON s.id = m.server_id
	WHERE m.user_id = $1
	ORDER BY s.created_at DESC
	LIMIT 20
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get all apis: %w", err)
	}
	defer rows.Close()

	servers := make([]Server, 0)

	for rows.Next() {
		var srv Server

		if err := rows.Scan(
			&srv.ID,
			&srv.Name,
			&srv.OwnerID,
			&srv.CreatedAT,
		); err != nil {
			return nil, fmt.Errorf("scan api: %w", err)
		}
		servers = append(servers, srv)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate apis: %w", err)
	}
	return servers, nil
}

func (r *Repo) GetInfo(ctx context.Context, id uuid.UUID) (Server, error) {
	query := "SELECT id, name, owner_id FROM servers WHERE id = $1"

	var srv Server

	err := r.db.QueryRow(ctx, query, id).Scan(
		&srv.ID,
		&srv.Name,
		&srv.OwnerID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Server{}, fmt.Errorf("server not found: %w", err)
		}
		return Server{}, fmt.Errorf("query error: %w", err)
	}

	return srv, nil
}

func (r *Repo) DeleteServer(ctx context.Context, id uuid.UUID) (Server, error) {
	query := "DELETE FROM servers WHERE id = $1 RETURNING id, name, owner_id"

	var srv Server

	err := r.db.QueryRow(ctx, query, id).Scan(
		&srv.ID,
		&srv.Name,
		&srv.OwnerID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Server{}, fmt.Errorf("server not found: %w", err)
		}
		return Server{}, fmt.Errorf("delete error: %w", err)
	}

	return srv, nil
}
