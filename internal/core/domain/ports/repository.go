package ports

import (
	"context"
	"time"
	"track-my-money/internal/core/domain"

	"github.com/google/uuid"
)

type StatementRepository interface {
	SaveStatementWithMovements(ctx context.Context, stmt *domain.Statement, movs []domain.Movement) error
}

type AnalysisRepository interface {
	GetMonthlyTotals(ctx context.Context, userID uuid.UUID, month string) (float64, float64, error)
	GetExpensesByCategory(ctx context.Context, userID uuid.UUID, month string) ([]domain.CategoryReport, error)
	GetMoneyFlow(ctx context.Context, userID uuid.UUID, startDate time.Time, groupBy string, accountID *uuid.UUID, bankID *int) ([]domain.TimeSeriesData, error)
}

type MovementRepository interface {
	GetMovements(ctx context.Context, filters domain.MovementFilters) ([]domain.MovementTable, int64, error)
}

type AppRepository interface {
	StatementRepository
	AnalysisRepository
	MovementRepository
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
}
