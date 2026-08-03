package members

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

func NewRepository(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) CreateMember(ctx context.Context, userID uuid.UUID, serverID uuid.UUID) (Member, error) {
	query := `
	INSERT INTO member (user_id, server_id, role) VALUES ($1, $2, $3) Returning *
	`
	var mbr Member

	err := r.db.QueryRow(ctx, query, userID, serverID).Scan(
		&mbr.ID,
		&mbr.ServerID,
		&mbr.UserID,
		&mbr.Role,
	)
	if err != nil {
		return Member{}, fmt.Errorf("query error: %w", err)
	}

	return mbr, nil
}

func (r *Repo) DeleteMember(ctx context.Context, serverID uuid.UUID, userID uuid.UUID) (Member, error) {
	query := `DELETE FROM member WHERE server_id = $1 AND user_id = $2 RETURNING *`

	var mbr Member

	err := r.db.QueryRow(ctx, query, serverID, userID).Scan(
		&mbr.ID,
		&mbr.ServerID,
		&mbr.UserID,
		&mbr.Role,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Member{}, fmt.Errorf("member not found")
		}
		return Member{}, fmt.Errorf("query error: %w", err)
	}
	return mbr, nil
}
