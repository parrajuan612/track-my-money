package services

import (
	"context"
	"track-my-money/internal/core/domain"

	"github.com/google/uuid"
)

func (s *service) GetCategoryDistribution(ctx context.Context, userID uuid.UUID, month string) (domain.CategoryDistributionResponse, error) {

	categories, err := s.analysisRepo.GetExpensesByCategory(ctx, userID, month)
	if err != nil {
		return domain.CategoryDistributionResponse{}, err
	}
	var totalExpenses float64
	for _, cat := range categories {
		totalExpenses += cat.Amount
	}
	if totalExpenses > 0 {
		for i := range categories {
			categories[i].Percentage = (categories[i].Amount / totalExpenses) * 100
		}
	}

	return domain.CategoryDistributionResponse{
		TotalExpenses: totalExpenses,
		Categories:    categories,
	}, nil
}
