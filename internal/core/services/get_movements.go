package services

import (
	"context"
	"track-my-money/internal/core/domain"
)

func (s *service) GetMovements(ctx context.Context, filters domain.MovementFilters) ([]domain.MovementTable, int64, error) {

	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PageSize <= 0 {
		filters.PageSize = 10
	}
	if filters.SortBy == "" {
		filters.SortBy = "date"
	}
	if filters.SortOrder == "" {
		filters.SortOrder = "DESC"
	}

	return s.repo.GetMovements(ctx, filters)
}
