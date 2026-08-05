package postgres

import (
	"context"
	"track-my-money/internal/core/domain"
)

// Obtener todos los bancos disponibles
func (r *postgresRepository) GetBanks(ctx context.Context) ([]domain.Bank, error) {
	var banks []domain.Bank
	err := r.db.WithContext(ctx).Order("name ASC").Find(&banks).Error
	return banks, err
}

// Obtener las cuentas creadas por el usuario autenticado
func (r *postgresRepository) GetAccountsByUserID(ctx context.Context, userID string) ([]domain.Account, error) {
	var accounts []domain.Account
	err := r.db.WithContext(ctx).
		Preload("Bank"). // Carga automáticamente los datos del banco asociado
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&accounts).Error
	return accounts, err
}

// Crear una nueva cuenta asignada al usuario
func (r *postgresRepository) CreateAccount(ctx context.Context, account *domain.Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}
func (r *postgresRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	// GORM buscará al usuario por su ID
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
