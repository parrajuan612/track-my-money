package services

import (
	"context"

	"github.com/google/uuid"
)

func (s *service) DeleteMovement(ctx context.Context, id string, userID uuid.UUID) error {
	return s.repo.DeleteMovement(ctx, id, userID)
}
