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
	DeleteMember(ctx context.Context, serverID uuid.UUID, userID uuid.UUID) (Member, error)
}

type Service struct {
	repo        Repository
	db          *pgxpool.Pool
	memberCache *checkmember.MemberCache
}

func NewServiceMbr(repo Repository, db *pgxpool.Pool, memberCache *checkmember.MemberCache) *Service {
	return &Service{
		repo:        repo,
		db:          db,
		memberCache: memberCache,
	}
}

func (s *Service) CreateMember(ctx context.Context, userID uuid.UUID, serverID uuid.UUID) (Member, error) {
	member, err := s.repo.CreateMember(ctx, userID, serverID)
	if err != nil {
		return Member{}, err
	}
	s.memberCache.Add(serverID, userID)
	return member, nil
}

func (s *Service) DeleteMember(ctx context.Context, userID uuid.UUID, targetID uuid.UUID, serverID uuid.UUID) (Member, error) {
	role, err := checkmember.GetMemberRole(ctx, s.db, serverID, userID)
	if err != nil {
		return Member{}, fmt.Errorf("failed to get role: %w", err)
	}

	if userID == targetID {
		member, err := s.repo.DeleteMember(ctx, serverID, userID)
		if err != nil {
			return Member{}, err
		}
		s.memberCache.RemoveMember(serverID, userID)
		return member, nil
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

	member, err := s.repo.DeleteMember(ctx, serverID, targetID)
	if err != nil {
		return Member{}, err
	}
	s.memberCache.RemoveMember(serverID, targetID)
	return member, nil
}
