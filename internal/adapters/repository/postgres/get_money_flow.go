package postgres

import (
	"context"
	"fmt"
	"time"
	"track-my-money/internal/core/domain"

	"github.com/google/uuid"
)

func (r *postgresRepository) GetMoneyFlow(ctx context.Context, userID uuid.UUID, startDate time.Time, groupBy string, accountID *uuid.UUID, bankID *int) ([]domain.TimeSeriesData, error) {
	var results []domain.TimeSeriesData

	// 1. Iniciamos la base de la consulta
	query := r.db.WithContext(ctx).Table("movements m").
		Select(fmt.Sprintf("TO_CHAR(m.date, '%s') as label, "+
			"SUM(CASE WHEN m.type = 'income' THEN ABS(m.amount) ELSE 0 END) as income, "+
			"SUM(CASE WHEN m.type = 'expense' THEN ABS(m.amount) ELSE 0 END) as expense", groupBy)).
		Where("m.user_id = ? AND m.date >= ?", userID, startDate)

	// 2. Filtro dinámico: Por Cuenta
	if accountID != nil {
		query = query.Where("m.account_id = ?", *accountID)
	}

	// 3. Filtro dinámico: Por Banco (Requiere JOIN)
	if bankID != nil {
		query = query.Joins("JOIN accounts a ON m.account_id = a.id").
			Where("a.bank_id = ?", *bankID)
	}

	// 4. Agrupación y Orden (Agrupamos por la fecha real para que el orden sea correcto)
	err := query.Group("label").
		Order("MIN(m.date) ASC").
		Scan(&results).Error

	return results, err
}
