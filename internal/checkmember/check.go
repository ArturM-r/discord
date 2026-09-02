package checkmember

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MemberCache struct {
	servers map[uuid.UUID]map[uuid.UUID]struct{}
	mu      sync.RWMutex
}

func NewMemberCache(m map[uuid.UUID]map[uuid.UUID]struct{}) *MemberCache {
	return &MemberCache{
		servers: m,
	}
}

func Unload(ctx context.Context, db *pgxpool.Pool) (map[uuid.UUID]map[uuid.UUID]struct{}, error) {
	query := "SELECT user_id, server_id FROM member"

	rows, err := db.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("failed to get: %w", err)
	}
	defer rows.Close()

	members := make(map[uuid.UUID]map[uuid.UUID]struct{})

	for rows.Next() {
		var userID, serverID uuid.UUID
		if err := rows.Scan(&userID, &serverID); err != nil {
			return nil, fmt.Errorf("failed to find user:%w", err)
		}

		if members[serverID] == nil {
			members[serverID] = make(map[uuid.UUID]struct{})
		}
		members[serverID][userID] = struct{}{}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return members, nil
}

func (c *MemberCache) IsMemberCache(serverID, userID uuid.UUID) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.servers[serverID][userID]
	return ok
}

func (c *MemberCache) Add(serverID, userID uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.servers[serverID] == nil {
		c.servers[serverID] = make(map[uuid.UUID]struct{})
	}
	c.servers[serverID][userID] = struct{}{}
}

func (c *MemberCache) RemoveMember(serverID, userID uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.servers[serverID], userID)
	if len(c.servers[serverID]) == 0 {
		delete(c.servers, serverID)
	}
}

func (c *MemberCache) RemoveServer(serverID uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.servers, serverID)
}

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
