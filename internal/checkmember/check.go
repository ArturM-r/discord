package checkmember

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func IsMember(ctx context.Context, db *pgxpool.Pool, userID uuid.UUID, serverID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(
    SELECT 1 FROM members 
    WHERE server_id = $1 AND user_id = $2
	)`

	var exist bool

	err := db.QueryRow(ctx, query, serverID, userID).Scan(&exist)

	if err != nil {
		return false, fmt.Errorf("query error: %w", err)
	}

	return exist, nil
}

func GetMemberRole(ctx context.Context, db *pgxpool.Pool, serverID, userID uuid.UUID) (string, error) {
	var role string
	err := db.QueryRow(ctx,
		"SELECT role FROM members WHERE server_id = $1 AND user_id = $2",
		serverID, userID,
	).Scan(&role)
	if err != nil {
		return "", fmt.Errorf("member not found: %w", err)
	}
	return role, nil
}
