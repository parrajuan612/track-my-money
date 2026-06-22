package postgres

import (
	"context"
	"track-my-money/internal/core/domain"

	"github.com/google/uuid"
)

func (r *postgresRepository) GetExpenseTrends(ctx context.Context, userID uuid.UUID, startDate string, endDate string, categoryID *int) ([]domain.TrendReport, error) {
	var trends []domain.TrendReport

	query := r.db.WithContext(ctx).
		Table("movements.movements m").
		Select("to_char(m.date, 'YYYY-MM-DD') as name, SUM(ABS(m.amount)) as amount").
		Where("m.user_id = ? AND m.type = 'expense'", userID)

	if startDate != "" {
		query = query.Where("m.date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("m.date <= ?", endDate)
	}
	if categoryID != nil {
		query = query.Where("m.category_id = ?", *categoryID)
	}

	err := query.
		Group("to_char(m.date, 'YYYY-MM-DD')").
		Order("name ASC").
		Scan(&trends).Error

	return trends, err
}
