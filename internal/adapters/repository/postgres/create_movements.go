package postgres

import (
	"context"
	"errors"
	"track-my-money/internal/core/domain"

	"github.com/google/uuid"
)

func (r *postgresRepository) CreateMovement(ctx context.Context, userID uuid.UUID, req domain.CreateMovementRequest) error {
	var accountID uuid.UUID

	// Buscamos la primera cuenta bancaria que le pertenezca a ESTE usuario
	err := r.db.WithContext(ctx).Table("accounts").
		Select("id").
		Where("user_id = ?", userID).
		Limit(1).
		Row().
		Scan(&accountID)

	// Si el usuario es nuevo y aún no ha creado ninguna cuenta, le avisamos
	if err != nil || accountID == uuid.Nil {
		return errors.New("no tienes cuentas bancarias registradas. Por favor crea una primero")
	}

	// Ajustamos el signo del monto (Gastos en negativo, Ingresos en positivo)
	montoFinal := req.Amount
	if req.Type == "expense" && montoFinal > 0 {
		montoFinal = -montoFinal
	} else if req.Type == "income" && montoFinal < 0 {
		montoFinal = -montoFinal
	}

	// Insertamos el registro en la base de datos
	return r.db.WithContext(ctx).Table("movements").Create(map[string]interface{}{
		"user_id":     userID,
		"account_id":  accountID,
		"category_id": req.CategoryID,
		"date":        req.Date,
		"description": req.Description,
		"amount":      montoFinal,
		"type":        req.Type,
		"is_active":   true,
	}).Error
}
