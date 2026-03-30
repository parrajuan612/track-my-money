package ports

import (
	"context"
	"io"
	"track-my-money/internal/core/domain"

	"github.com/google/uuid"
)

type Service interface {
	ExtractStatementData(file io.Reader, fileName string, password string, bankID string) ([]domain.Movement, string, error)
	SaveMovements(c context.Context, req domain.SaveStatementRequest) error
	GetMonthlyComparison(ctx context.Context, userID uuid.UUID, month string) (domain.ComparisonResponse, error)
	GetCategoryDistribution(ctx context.Context, userID uuid.UUID, month string) (domain.CategoryDistributionResponse, error)
	GetMoneyFlow(ctx context.Context, userID uuid.UUID, req domain.MoneyFlowRequest) ([]domain.TimeSeriesData, error)
}
