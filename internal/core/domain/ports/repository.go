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
	GetExpensesByCategory(ctx context.Context, userID uuid.UUID, startDate string, endDate string) ([]domain.CategoryReport, error)
	GetMoneyFlow(ctx context.Context, userID uuid.UUID, startDate time.Time, groupBy string, accountID *uuid.UUID, bankID *int) ([]domain.TimeSeriesData, error)
}

type MovementRepository interface {
	GetMovements(ctx context.Context, filters domain.MovementFilters) ([]domain.MovementTable, int64, error)
	UpdateMovement(ctx context.Context, movementID uuid.UUID, userID uuid.UUID, req domain.UpdateMovementRequest) error
}

type AppRepository interface {
	StatementRepository
	AnalysisRepository
	MovementRepository
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	GetExpenseTrends(ctx context.Context, userID uuid.UUID, startDate string, endDate string, categoryID *int) ([]domain.TrendReport, error)
	UpdateMovement(ctx context.Context, movementID uuid.UUID, userID uuid.UUID, req domain.UpdateMovementRequest) error
	CreateMovement(ctx context.Context, userID uuid.UUID, req domain.CreateMovementRequest) error
	// --- MÓDULO DE CUENTAS BANCARIAS ---
	GetBanks(ctx context.Context) ([]domain.Bank, error)
	DeleteMovement(ctx context.Context, id string, userID uuid.UUID) error
	GetAccountsByUserID(ctx context.Context, userID string) ([]domain.Account, error)
	CreateAccount(ctx context.Context, account *domain.Account) error
}
