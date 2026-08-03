package members

import (
	"context"
	"discord/internal/checkmember"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateMember(ctx context.Context, userID uuid.UUID, serverID uuid.UUID) (Member, error)
	DeleteMember(ctx context.Context, userID uuid.UUID, serverID uuid.UUID) (Member, error)
}

type Service struct {
	repo Repository
	db   *pgxpool.Pool
}

func NewServiceMbr(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateMember(ctx context.Context, userID uuid.UUID, serverID uuid.UUID) (Member, error) {
	role, err := checkmember.GetMemberRole(ctx, s.db, serverID, userID)

	if err != nil {
		return Member{}, fmt.Errorf("failed to get role: %w", err)
	}
	switch role {
	case "owner", "admin", "member":
		return s.repo.CreateMember(ctx, userID, serverID)
	default:
		return Member{}, fmt.Errorf("forbidden")
	}
}

func (s *Service) DeleteMember(ctx context.Context, userID uuid.UUID, targetID uuid.UUID, serverID uuid.UUID) (Member, error) {
	role, err := checkmember.GetMemberRole(ctx, s.db, serverID, userID)
	if err != nil {
		return Member{}, fmt.Errorf("failed to get role: %w", err)
	}

	if userID == targetID {
		return s.repo.DeleteMember(ctx, serverID, userID)
	}

	if role != "owner" && role != "admin" {
		return Member{}, fmt.Errorf("forbidden")
	}

	targetRole, err := checkmember.GetMemberRole(ctx, s.db, serverID, targetID)
	if err != nil {
		return Member{}, fmt.Errorf("target not found")
	}
	if targetRole == "owner" {
		return Member{}, fmt.Errorf("cannot kick owner")
	}

	return s.repo.DeleteMember(ctx, serverID, targetID)
}
