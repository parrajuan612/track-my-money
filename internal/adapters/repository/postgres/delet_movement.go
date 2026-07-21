package postgres

import (
	"context"

	"github.com/google/uuid"
)

func (r *postgresRepository) DeleteMovement(ctx context.Context, id string, userID uuid.UUID) error {
	// Usamos Exec con SQL directo para mayor seguridad y precisión,
	// asegurando que solo el dueño del movimiento pueda borrarlo.
	err := r.db.WithContext(ctx).Exec(
		"DELETE FROM movements.movements WHERE id = ? AND user_id = ?",
		id, userID,
	).Error

	return err
}
