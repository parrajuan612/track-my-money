package services

import (
	"context"
	"track-my-money/internal/core/domain"
)

func (s *service) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	// Simplemente le pedimos los datos al repositorio
	return s.repo.GetUserByID(ctx, id)
}
