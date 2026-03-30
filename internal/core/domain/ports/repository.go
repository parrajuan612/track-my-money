package ports

import (
	"context"
	"time"
	"track-my-money/internal/core/domain"

	"github.com/google/uuid"
)

type Repository interface {
	SaveStatementWithMovements(ctx context.Context, stmt *domain.Statement, movs []domain.Movement) error
}
type AnalysisRepository interface {
	GetMonthlyTotals(ctx context.Context, userID uuid.UUID, month string) (float64, float64, error)
	GetExpensesByCategory(ctx context.Context, userID uuid.UUID, month string) ([]domain.CategoryReport, error)
	GetMoneyFlow(ctx context.Context, userID uuid.UUID, startDate time.Time, groupBy string, accountID *uuid.UUID, bankID *int) ([]domain.TimeSeriesData, error)
}
