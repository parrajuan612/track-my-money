package postgres

import (
	"context"
	"track-my-money/internal/core/domain"
)

func (r *postgresRepository) CreateUser(ctx context.Context, user *domain.User) error {

	err := r.db.WithContext(ctx).
		Table("movements.users").
		Create(user).Error

	return err
}
