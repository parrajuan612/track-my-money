package services

import (
	"context"
	"fmt"
	"time"
	"track-my-money/internal/core/domain"

	"github.com/google/uuid"
)

func (s *service) GetMonthlyComparison(ctx context.Context, userID uuid.UUID, currentMonth string) (domain.ComparisonResponse, error) {
	layout := "2006-01"
	t, err := time.Parse(layout, currentMonth)
	if err != nil {
		return domain.ComparisonResponse{}, fmt.Errorf("formato de mes inválido: %w", err)
	}
	prevMonthStr := t.AddDate(0, -1, 0).Format(layout)
	currInc, currExp, err := s.repo.GetMonthlyTotals(ctx, userID, currentMonth)
	if err != nil {
		return domain.ComparisonResponse{}, err
	}
	prevInc, prevExp, err := s.repo.GetMonthlyTotals(ctx, userID, prevMonthStr)
	if err != nil {
		return domain.ComparisonResponse{}, err
	}
	incVar := s.calculatePercentageVariation(currInc, prevInc)
	expVar := s.calculatePercentageVariation(currExp, prevExp)
	return domain.ComparisonResponse{
		CurrentMonth: domain.MonthTotals{
			Income:  currInc,
			Expense: currExp,
		},
		PreviousMonth: domain.MonthTotals{
			Income:  prevInc,
			Expense: prevExp,
		},
		Variations: domain.Variations{
			IncomePercentage:  incVar,
			ExpensePercentage: expVar,
		},
	}, nil
}
func (s *service) calculatePercentageVariation(current, previous float64) float64 {
	if previous == 0 {
		if current > 0 {
			return 100.0
		}
		return 0.0
	}
	return ((current - previous) / previous) * 100
}
