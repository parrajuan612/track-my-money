package postgres

import (
	"context"
	"track-my-money/internal/core/domain"

	"github.com/google/uuid"
)

func (r *analysisRepository) GetExpensesByCategory(ctx context.Context, userID uuid.UUID, month string) ([]domain.CategoryReport, error) {
	var reports []domain.CategoryReport

	// Query con Join para traer el nombre de la categoría
	err := r.db.WithContext(ctx).
		Table("movements.movements m").
		Select("c.name as category_name, SUM(ABS(m.amount)) as amount").
		Joins("JOIN movements.categories c ON m.category_id = c.id").
		Where("m.user_id = ? AND m.type = 'expense' AND to_char(m.date, 'YYYY-MM') = ?", userID, month).
		Group("c.name").
		Order("amount DESC").
		Scan(&reports).Error

	return reports, err
}
