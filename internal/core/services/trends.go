package services

import (
	"context"
	"fmt"
	"track-my-money/internal/core/domain"

	"github.com/google/uuid"
)

func (s *service) GetExpenseTrends(ctx context.Context, userID uuid.UUID, startDate string, endDate string, categoryID *int) ([]domain.TrendReport, error) {

	trends, err := s.repo.GetExpenseTrends(ctx, userID, startDate, endDate, categoryID)
	if err != nil {
		return []domain.TrendReport{}, nil
	}
	fmt.Println(trends)
	return trends, nil
}
