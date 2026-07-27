package postgres

import (
	"context"
	"errors"
	"track-my-money/internal/core/domain"

	"github.com/google/uuid"
)

func (r *postgresRepository) UpdateMovement(ctx context.Context, movementID uuid.UUID, userID uuid.UUID, req domain.UpdateMovementRequest) error {

	// Como tus gastos en la BD podrían guardarse en negativo, aseguramos el signo
	montoFinal := req.Amount
	if req.Type == "expense" && montoFinal > 0 {
		montoFinal = -montoFinal
	} else if req.Type == "income" && montoFinal < 0 {
		montoFinal = -montoFinal // Lo volvemos positivo
	}

	result := r.db.WithContext(ctx).
		Table("movements").          // Nombre exacto de tu tabla
		Where("id = ?", movementID). // Quitamos el filtro de user_id temporalmente
		Updates(map[string]interface{}{
			"description": req.Description,
			"amount":      montoFinal, // Usamos el monto con el signo correcto
			"type":        req.Type,
			"category_id": req.CategoryID,
		})

	if result.Error != nil {
		return result.Error
	}

	// ¡ESTA ES LA MAGIA! Si GORM actualiza 0 filas, ahora sí lanzará un error
	if result.RowsAffected == 0 {
		return errors.New("no se encontró el movimiento para actualizar o no hubo cambios")
	}

	return nil
}
