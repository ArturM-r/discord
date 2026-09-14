package channel

import (
	"context"
	"discord/internal/checkmember"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetChRepo(ctx context.Context, serverID uuid.UUID) ([]Channel, error)
	CreateChRepo(ctx context.Context, serverID uuid.UUID, name string) (Channel, error)
	DeleteChnRepo(ctx context.Context, serverID uuid.UUID, channelID uuid.UUID) (Channel, error)
}

type ServiceCHN struct {
	repo Repository
	db   *pgxpool.Pool
}

func NewService(repo Repository, db *pgxpool.Pool) *ServiceCHN {
	return &ServiceCHN{
		repo: repo,
		db:   db,
	}
}

func (s *ServiceCHN) requireMember(ctx context.Context, serverID, userID uuid.UUID) error {
	ok, err := checkmember.IsMember(ctx, s.db, serverID, userID)
	if err != nil {
		return fmt.Errorf("failed to check member: %w", err)
	}
	if !ok {
		return fmt.Errorf("forbidden")
	}
	return nil
}

func (s *ServiceCHN) GetChService(ctx context.Context, serverID uuid.UUID, userID uuid.UUID) ([]Channel, error) {
	if err := s.requireMember(ctx, serverID, userID); err != nil {
		return nil, err
	}
	return s.repo.GetChRepo(ctx, serverID)
}

func (s *ServiceCHN) CreateChService(ctx context.Context, serverID uuid.UUID, userID uuid.UUID, name string) (Channel, error) {
	if name == "" {
		return Channel{}, fmt.Errorf("name is required")
	}
	if len(name) > 20 {
		return Channel{}, fmt.Errorf("name must be less than 20 characters")
	}
	if err := s.requireMember(ctx, serverID, userID); err != nil {
		return Channel{}, err
	}

	return s.repo.CreateChRepo(ctx, serverID, name)
}

func (s *ServiceCHN) DeleteChnService(ctx context.Context, serverID uuid.UUID, channelID uuid.UUID, userID uuid.UUID) (Channel, error) {
	if err := s.requireMember(ctx, serverID, userID); err != nil {
		return Channel{}, err
	}

	return s.repo.DeleteChnRepo(ctx, serverID, channelID)
}
