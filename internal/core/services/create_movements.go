package services

import (
	"context"
	"track-my-money/internal/core/domain"

	"github.com/google/uuid"
)

func (s *service) CreateMovement(ctx context.Context, userID uuid.UUID, req domain.CreateMovementRequest) error {
	return s.repo.CreateMovement(ctx, userID, req)
}
