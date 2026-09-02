package server

import (
	"context"
	"fmt"
	"strings"

	"discord/internal/checkmember"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateSRV(ctx context.Context, userID uuid.UUID, name string) (Server, error)
	GetSrvUser(ctx context.Context, userID uuid.UUID) ([]Server, error)
	GetInfo(ctx context.Context, id uuid.UUID) (Server, error)
	DeleteServer(ctx context.Context, id uuid.UUID) (Server, error)
}

type ServiceSrv struct {
	repo        Repository
	db          *pgxpool.Pool
	memberCache *checkmember.MemberCache
}

func NewService(repo Repository, db *pgxpool.Pool, memberCache *checkmember.MemberCache) *ServiceSrv {
	return &ServiceSrv{
		repo:        repo,
		db:          db,
		memberCache: memberCache,
	}
}

func (s *ServiceSrv) CreateSrvService(ctx context.Context, userID uuid.UUID, name string) (Server, error) {
	if len(strings.TrimSpace(name)) == 0 {
		return Server{}, fmt.Errorf("name cant be empty")
	}
	if len(name) > 100 {
		return Server{}, fmt.Errorf("name too long")
	}
	server, err := s.repo.CreateSRV(ctx, userID, name)
	if err != nil {
		return Server{}, err
	}
	s.memberCache.Add(server.ID, userID)
	return server, nil
}

func (s *ServiceSrv) GetSrvUserService(ctx context.Context, userID uuid.UUID) ([]Server, error) {

	return s.repo.GetSrvUser(ctx, userID)
}

func (s *ServiceSrv) GetInfoService(ctx context.Context, id uuid.UUID, userID uuid.UUID) (Server, error) {
	ok, err := checkmember.IsMember(ctx, s.db, id, userID)
	if err != nil {
		return Server{}, fmt.Errorf("failed to check member: %w", err)
	}
	if !ok {
		return Server{}, fmt.Errorf("forbidden")
	}
	return s.repo.GetInfo(ctx, id)
}

func (s *ServiceSrv) DeleteSrvService(ctx context.Context, id uuid.UUID, userID uuid.UUID) (Server, error) {
	role, err := checkmember.GetMemberRole(ctx, s.db, id, userID)
	if err != nil {
		return Server{}, fmt.Errorf("forbidden")
	}

	if role != "owner" {
		return Server{}, fmt.Errorf("forbidden")
	}

	server, err := s.repo.DeleteServer(ctx, id)
	if err != nil {
		return Server{}, err
	}
	s.memberCache.RemoveServer(id)
	return server, nil
}
