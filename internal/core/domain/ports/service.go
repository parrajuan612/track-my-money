package ports

import (
	"context"
	"io"
	"track-my-money/internal/core/domain"
)

type Service interface {
	ExtractStatementData(file io.Reader, fileName string, password string, bankID string) ([]domain.Movement, string, error)
	SaveMovements(c context.Context, req domain.SaveStatementRequest) error
}
