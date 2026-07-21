package services

import (
	"context"
	"track-my-money/internal/core/domain"
)

// Obtener la lista de bancos
func (s *service) GetBanks(ctx context.Context) ([]domain.Bank, error) {
	return s.repo.GetBanks(ctx)
}

// Obtener las cuentas de un usuario específico
func (s *service) GetAccountsByUserID(ctx context.Context, userID string) ([]domain.Account, error) {
	return s.repo.GetAccountsByUserID(ctx, userID)
}

// Crear una nueva cuenta bancaria
func (s *service) CreateAccount(ctx context.Context, account *domain.Account) error {
	// Aquí podrías agregar reglas de negocio en el futuro, por ejemplo:
	// "Un usuario no puede tener más de 10 cuentas en el plan gratuito"

	return s.repo.CreateAccount(ctx, account)
}
