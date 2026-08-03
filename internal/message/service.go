package message

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Repository interface {
	Get(ctx context.Context, channelID uuid.UUID, userID uuid.UUID) ([]Message, error)
}

type ServiceMSG struct {
	repo Repository
}

func NewService(repo Repository) *ServiceMSG {
	return &ServiceMSG{
		repo: repo,
	}
}

func (m *ServiceMSG) GetByChannel(ctx context.Context, channelID string, userID uuid.UUID) ([]Message, error) {
	channeUUID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, fmt.Errorf("type not valid: %w", err)
	}

	return m.repo.Get(ctx, channeUUID, userID)
}
