package services

import (
	"context"
	"fmt"
	"time"
	"track-my-money/internal/core/domain"

	"github.com/google/uuid"
)

func (s *service) SaveMovements(ctx context.Context, req domain.SaveStatementRequest) error {

	userUUID, _ := uuid.Parse("296f368f-f7b4-4388-8934-209e146de03c")
	accountUUID, _ := uuid.Parse("3bf374ea-db66-4f47-ab8a-7d0156c4440f")
	fmt.Println(accountUUID)
	var domainMovements []domain.Movement
	for _, m := range req.Movements {
		parsedDate, _ := time.Parse("2006-01-02", m.Date)

		domainMovements = append(domainMovements, domain.Movement{
			UserID:      userUUID,
			AccountID:   accountUUID,
			CategoryID:  m.CategoryID,
			Date:        parsedDate,
			Description: m.Description,
			Amount:      m.Amount,
			Type:        domain.MovementType(m.Type),
		})
	}

	statement := domain.Statement{
		UserID:      userUUID,
		AccountID:   accountUUID,
		BankID:      1,
		FileName:    req.FileName,
		PeriodMonth: req.PeriodMonth,
	}

	return s.repository.SaveStatementWithMovements(ctx, &statement, domainMovements)
}
