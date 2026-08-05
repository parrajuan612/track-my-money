package ports

import (
	"context"
	"io"
	"track-my-money/internal/core/domain"

	"github.com/google/uuid"
)

type Service interface {
	ExtractStatementData(file io.Reader, fileName string, password string, bankID string) ([]domain.Movement, string, error)
	SaveMovements(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, req domain.SaveStatementRequest) error
	GetMonthlyComparison(ctx context.Context, userID uuid.UUID, month string) (domain.ComparisonResponse, error)
	GetCategoryDistribution(ctx context.Context, userID uuid.UUID, startDate string, endDate string) (domain.CategoryDistributionResponse, error)
	GetMoneyFlow(ctx context.Context, userID uuid.UUID, req domain.MoneyFlowRequest) ([]domain.TimeSeriesData, error)
	GetMovements(ctx context.Context, filters domain.MovementFilters) ([]domain.MovementTable, int64, error)
	Login(ctx context.Context, email string, password string) (string, *domain.User, error)
	LoginWithGoogle(ctx context.Context, idToken string) (string, *domain.User, error)
	Register(ctx context.Context, name, email, password string) (string, *domain.User, error)
	GetExpenseTrends(ctx context.Context, userID uuid.UUID, startDate string, endDate string, categoryID *int) ([]domain.TrendReport, error)
	UpdateMovement(ctx context.Context, movementID uuid.UUID, userID uuid.UUID, req domain.UpdateMovementRequest) error
	CreateMovement(ctx context.Context, userID uuid.UUID, req domain.CreateMovementRequest) error
	// --- MÓDULO DE CUENTAS BANCARIAS ---
	GetBanks(ctx context.Context) ([]domain.Bank, error)
	GetAccountsByUserID(ctx context.Context, userID string) ([]domain.Account, error)
	CreateAccount(ctx context.Context, account *domain.Account) error
	DeleteMovement(ctx context.Context, id string, userID uuid.UUID) error
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
}
