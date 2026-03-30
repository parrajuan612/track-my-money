package postgres

import (
	"context"
	"track-my-money/internal/core/domain/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type analysisRepository struct {
	db *gorm.DB
}

func NewAnalysisRepository(db *gorm.DB) ports.AnalysisRepository {
	return &analysisRepository{db: db}
}

func (r *analysisRepository) GetMonthlyTotals(ctx context.Context, userID uuid.UUID, month string) (float64, float64, error) {
	var results []struct {
		Type   string
		Amount float64
	}
	err := r.db.WithContext(ctx).
		Table("movements.movements").
		Select("type, SUM(ABS(amount)) as amount").
		Where("user_id = ? AND to_char(date, 'YYYY-MM') = ?", userID, month).
		Group("type").
		Scan(&results).Error

	if err != nil {
		return 0, 0, err
	}

	var income, expense float64
	for _, res := range results {
		if res.Type == "income" {
			income = res.Amount
		} else if res.Type == "expense" {
			expense = res.Amount
		}
	}

	return income, expense, nil
}
