package channel

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Repository interface {
	GetChRepo(ctx context.Context, serverID uuid.UUID) ([]Channel, error)
	CreateChRepo(ctx context.Context, serverID uuid.UUID, name string) (Channel, error)
	DeleteChnRepo(ctx context.Context, serverID uuid.UUID, channelID uuid.UUID) (Channel, error)
}

type ServiceCHN struct {
	repo Repository
}

func NewService(repo Repository) *ServiceCHN {
	return &ServiceCHN{
		repo: repo,
	}
}

func (s *ServiceCHN) GetChService(ctx context.Context, serverID uuid.UUID) ([]Channel, error) {

	return s.repo.GetChRepo(ctx, serverID)
}

func (s *ServiceCHN) CreateChService(ctx context.Context, serverID uuid.UUID, name string) (Channel, error) {
	if name == "" {
		return Channel{}, fmt.Errorf("name is required")
	}
	if len(name) > 20 {
		return Channel{}, fmt.Errorf("name must be less than 20 characters")
	}

	return s.repo.CreateChRepo(ctx, serverID, name)
}

func (s *ServiceCHN) DeleteChnService(ctx context.Context, serverID uuid.UUID, channelID uuid.UUID) (Channel, error) {

	return s.repo.DeleteChnRepo(ctx, serverID, channelID)
}
