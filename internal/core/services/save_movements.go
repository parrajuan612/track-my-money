package services

import (
	"context"
	"fmt"
	"time"
	"track-my-money/internal/core/domain"

	"github.com/google/uuid"
)

// Función auxiliar robusta para leer fechas en múltiples formatos
func parseDateRobust(dateStr string) time.Time {
	// 1. Intentar formato estándar de JavaScript/React (ISO 8601: "2025-11-28T00:00:00Z")
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return t
	}
	// 2. Intentar formato SQL clásico ("2025-11-28")
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t
	}
	// 3. Intentar formato latino ("28/11/2025")
	if t, err := time.Parse("02/01/2006", dateStr); err == nil {
		return t
	}

	// Si llega aquí, es un formato súper extraño. Lo imprimimos para debuggear.
	fmt.Println("⚠️ ERROR PARSEANDO FECHA, se recibió:", dateStr)

	// Como último recurso para que la BD no colapse con el año 1, le ponemos la fecha de hoy.
	return time.Now()
}

func (s *service) SaveMovements(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, req domain.SaveStatementRequest) error {

	var domainMovements []domain.Movement
	for _, m := range req.Movements {

		// ¡Usamos nuestra nueva función a prueba de balas!
		parsedDate := parseDateRobust(m.Date)

		domainMovements = append(domainMovements, domain.Movement{
			UserID:      userID,
			AccountID:   accountID,
			CategoryID:  m.CategoryID,
			Date:        parsedDate,
			Description: m.Description,
			Amount:      m.Amount,
			Type:        domain.MovementType(m.Type),
		})
	}

	statement := domain.Statement{
		UserID:      userID,
		AccountID:   accountID,
		BankID:      req.BankID,
		FileName:    req.FileName,
		PeriodMonth: req.PeriodMonth,
	}

	return s.repo.SaveStatementWithMovements(ctx, &statement, domainMovements)
}
