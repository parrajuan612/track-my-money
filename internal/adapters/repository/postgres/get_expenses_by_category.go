package postgres

import (
	"context"
	"track-my-money/internal/core/domain"

	"github.com/google/uuid"
)

func (r *postgresRepository) GetExpensesByCategory(ctx context.Context, userID uuid.UUID, startDate string, endDate string) ([]domain.CategoryReport, error) {
	var reports []domain.CategoryReport

	// Empezamos la query base
	query := r.db.WithContext(ctx).
		Table("movements m").
		Select("c.name as category_name, SUM(ABS(m.amount)) as amount").
		Joins("JOIN categories c ON m.category_id = c.id").
		Where("m.user_id = ? AND m.type = 'expense'", userID)

	// Agregamos los filtros de fecha solo si nos los envían
	if startDate != "" {
		query = query.Where("m.date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("m.date <= ?", endDate)
	}

	// Terminamos de agrupar y ordenar
	err := query.
		Group("c.name").
		Order("amount DESC").
		Scan(&reports).Error

	return reports, err
}
