package services

import (
	"context"
	"track-my-money/internal/core/domain"

	"github.com/google/uuid"
)

func (s *service) UpdateMovement(ctx context.Context, movementID uuid.UUID, userID uuid.UUID, req domain.UpdateMovementRequest) error {
	return s.repo.UpdateMovement(ctx, movementID, userID, req)
}
